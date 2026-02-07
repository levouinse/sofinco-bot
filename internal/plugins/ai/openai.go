package ai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

type AIPlugin struct {
	plugins.BasePlugin
}

func init() {
	plugins.Register(&AIPlugin{})
}

func (p *AIPlugin) Commands() []string { return []string{"ai", "ask", "chatgpt"} }
func (p *AIPlugin) Tags() []string     { return []string{"ai"} }
func (p *AIPlugin) Help() string       { return "Chat with AI" }
func (p *AIPlugin) RequireLimit() bool { return true }

func (p *AIPlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "Masukkan pertanyaan!\n\nContoh:\n/ai apa itu golang?")
		ctx.API.Send(msg)
		return nil
	}

	question := strings.Join(ctx.Args, " ")
	
	waitMsg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "🤔 Thinking...")
	sent, _ := ctx.API.Send(waitMsg)

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/search/openai-chat?text=%s&apikey=%s",
		url.QueryEscape(question), ctx.Config.APIKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("gagal menghubungi AI: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &result); err != nil || !result.Status {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("gagal mendapatkan jawaban")
	}

	ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))

	reply := tgbotapi.NewMessage(ctx.Message.Chat.ID, result.Message)
	reply.ReplyToMessageID = ctx.Message.MessageID
	ctx.API.Send(reply)

	if !ctx.User.Premium {
		ctx.User.Limit--
		ctx.DB.SaveUser(ctx.User)
	}

	return nil
}
