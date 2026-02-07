package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

type TikTokPlugin struct {
	plugins.BasePlugin
}

func init() {
	plugins.Register(&TikTokPlugin{})
}

func (p *TikTokPlugin) Commands() []string { return []string{"tiktok", "tt", "ttdl"} }
func (p *TikTokPlugin) Tags() []string     { return []string{"downloader"} }
func (p *TikTokPlugin) Help() string       { return "Download TikTok video" }
func (p *TikTokPlugin) RequireLimit() bool { return true }

func (p *TikTokPlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "Masukkan link TikTok!\n\nContoh:\n/tiktok https://vt.tiktok.com/xxx")
		ctx.API.Send(msg)
		return nil
	}

	link := ctx.Args[0]
	waitMsg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "⏳ Downloading...")
	sent, _ := ctx.API.Send(waitMsg)

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/download/tiktok?url=%s&apikey=%s",
		url.QueryEscape(link), ctx.Config.APIKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("gagal download: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Status bool `json:"status"`
		Result struct {
			Title       string `json:"title"`
			VideoNoWM   string `json:"video"`
			Music       string `json:"music"`
			Author      string `json:"author"`
			Views       int64  `json:"views"`
			Likes       int64  `json:"likes"`
			Comments    int64  `json:"comments"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil || !result.Status {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("gagal mendapatkan video")
	}

	ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))

	caption := fmt.Sprintf("📹 *TikTok Downloader*\n\n"+
		"👤 Author: %s\n"+
		"📝 Title: %s\n"+
		"👁 Views: %d\n"+
		"❤️ Likes: %d\n"+
		"💬 Comments: %d",
		result.Result.Author, result.Result.Title,
		result.Result.Views, result.Result.Likes, result.Result.Comments)

	video := tgbotapi.NewVideo(ctx.Message.Chat.ID, tgbotapi.FileURL(result.Result.VideoNoWM))
	video.Caption = caption
	video.ParseMode = "Markdown"
	ctx.API.Send(video)

	if !ctx.User.Premium {
		ctx.User.Limit--
		ctx.DB.SaveUser(ctx.User)
	}

	return nil
}
