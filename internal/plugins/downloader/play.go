package downloader

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

type PlayPlugin struct {
	plugins.BasePlugin
}

type YTSearchResult struct {
	Videos []struct {
		VideoID   string `json:"videoId"`
		Title     string `json:"title"`
		URL       string `json:"url"`
		Image     string `json:"image"`
		Timestamp string `json:"timestamp"`
		Views     int64  `json:"views"`
		Author    struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"author"`
	} `json:"videos"`
}

func init() {
	plugins.Register(&PlayPlugin{
		BasePlugin: plugins.BasePlugin{},
	})
}

func (p *PlayPlugin) Commands() []string { return []string{"play", "ds", "song"} }
func (p *PlayPlugin) Tags() []string     { return []string{"downloader"} }
func (p *PlayPlugin) Help() string       { return "Download audio from YouTube" }
func (p *PlayPlugin) RequireLimit() bool { return true }

func (p *PlayPlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Masukkan judul/link YouTube!\n\nContoh:\n/play taylor swift"))
		return nil
	}

	query := strings.Join(ctx.Args, " ")
	
	waitMsg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "Tunggu sebentar...")
	sent, _ := ctx.API.Send(waitMsg)

	searchURL := fmt.Sprintf("https://api.betabotz.eu.org/api/search/yts?query=%s&apikey=%s",
		url.QueryEscape(query), ctx.Config.APIKey)

	resp, err := http.Get(searchURL)
	if err != nil {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("gagal mencari video: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var searchResult YTSearchResult
	if err := json.Unmarshal(body, &searchResult); err != nil || len(searchResult.Videos) == 0 {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("video tidak ditemukan")
	}

	vid := searchResult.Videos[0]

	downloadURL := fmt.Sprintf("https://api.betabotz.eu.org/api/download/yt?url=%s&apikey=%s",
		url.QueryEscape(vid.URL), ctx.Config.APIKey)

	resp2, err := http.Get(downloadURL)
	if err != nil {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("gagal download: %w", err)
	}
	defer resp2.Body.Close()

	body2, _ := io.ReadAll(resp2.Body)
	var result struct {
		Status bool `json:"status"`
		Result struct {
			Title     string `json:"title"`
			Duration  string `json:"duration"`
			Views     string `json:"views"`
			Author    string `json:"author"`
			Thumbnail string `json:"thumbnail"`
			MP3       string `json:"mp3"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body2, &result); err != nil || !result.Status {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("gagal mendapatkan link audio")
	}

	ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))

	caption := fmt.Sprintf("∘ Title : %s\n∘ ID : %s\n∘ Duration : %s\n∘ Viewers : %s\n∘ Author : %s\n∘ Url : %s",
		vid.Title, vid.VideoID, vid.Timestamp, result.Result.Views, vid.Author.Name, vid.URL)

	photo := tgbotapi.NewPhoto(ctx.Message.Chat.ID, tgbotapi.FileURL(vid.Image))
	photo.Caption = caption
	ctx.API.Send(photo)

	audio := tgbotapi.NewAudio(ctx.Message.Chat.ID, tgbotapi.FileURL(result.Result.MP3))
	audio.Title = vid.Title
	audio.Performer = vid.Author.Name
	ctx.API.Send(audio)

	if !ctx.User.Premium {
		ctx.User.Limit--
		ctx.DB.SaveUser(ctx.User)
	}

	return nil
}
