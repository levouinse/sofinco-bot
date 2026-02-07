package maker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

type JadiAnimePlugin struct {
	plugins.BasePlugin
}

func init() {
	plugins.Register(&JadiAnimePlugin{})
}

func (p *JadiAnimePlugin) Commands() []string { return []string{"jadianime", "toanime", "anime"} }
func (p *JadiAnimePlugin) Tags() []string     { return []string{"maker"} }
func (p *JadiAnimePlugin) Help() string       { return "Convert photo to anime style" }
func (p *JadiAnimePlugin) RequireLimit() bool { return true }

func (p *JadiAnimePlugin) Execute(ctx *plugins.Context) error {
	if ctx.Message.ReplyToMessage == nil || ctx.Message.ReplyToMessage.Photo == nil {
		msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "Reply ke foto yang ingin dijadikan anime!")
		ctx.API.Send(msg)
		return nil
	}

	waitMsg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "🎨 Converting to anime...")
	sent, _ := ctx.API.Send(waitMsg)

	photos := ctx.Message.ReplyToMessage.Photo
	fileID := photos[len(photos)-1].FileID

	file, err := ctx.API.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("gagal mendapatkan file: %w", err)
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", ctx.Config.BotToken, file.FilePath)

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/maker/jadianime?url=%s&apikey=%s",
		url.QueryEscape(fileURL), ctx.Config.APIKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("gagal convert: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Status bool   `json:"status"`
		Result string `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil || !result.Status {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("gagal convert image")
	}

	ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))

	photo := tgbotapi.NewPhoto(ctx.Message.Chat.ID, tgbotapi.FileURL(result.Result))
	photo.Caption = "🎨 Anime Style"
	ctx.API.Send(photo)

	if !ctx.User.Premium {
		ctx.User.Limit--
		ctx.DB.SaveUser(ctx.User)
	}

	return nil
}
