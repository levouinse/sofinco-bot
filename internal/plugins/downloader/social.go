package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

// Instagram Downloader
type InstagramPlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &InstagramPlugin{}
	plugins.Register(p)
}

func (p *InstagramPlugin) Commands() []string { return []string{"instagram", "ig", "igdl"} }
func (p *InstagramPlugin) Tags() []string     { return []string{"downloader"} }
func (p *InstagramPlugin) Help() string       { return "Download Instagram video/photo" }
func (p *InstagramPlugin) RequireLimit() bool { return true }

func (p *InstagramPlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Usage: /instagram <url>"))
		return nil
	}

	url := ctx.Args[0]
	msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "⏳ Downloading...")
	sent, _ := ctx.API.Send(msg)

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/download/instagram?url=%s&apikey=%s", url, ctx.Config.APIKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to download")
		ctx.API.Send(edit)
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Result  []struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to parse response")
		ctx.API.Send(edit)
		return err
	}

	if !apiResponse.Status || len(apiResponse.Result) == 0 {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ No media found")
		ctx.API.Send(edit)
		return nil
	}

	ctx.API.Request(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))

	for _, media := range apiResponse.Result {
		if media.Type == "video" {
			video := tgbotapi.NewVideo(ctx.Message.Chat.ID, tgbotapi.FileURL(media.URL))
			video.Caption = "📥 Instagram Video"
			ctx.API.Send(video)
		} else {
			photo := tgbotapi.NewPhoto(ctx.Message.Chat.ID, tgbotapi.FileURL(media.URL))
			photo.Caption = "📥 Instagram Photo"
			ctx.API.Send(photo)
		}
	}

	return nil
}

// Facebook Downloader
type FacebookPlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &FacebookPlugin{}
	plugins.Register(p)
}

func (p *FacebookPlugin) Commands() []string { return []string{"facebook", "fb", "fbdl"} }
func (p *FacebookPlugin) Tags() []string     { return []string{"downloader"} }
func (p *FacebookPlugin) Help() string       { return "Download Facebook video" }
func (p *FacebookPlugin) RequireLimit() bool { return true }

func (p *FacebookPlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Usage: /facebook <url>"))
		return nil
	}

	url := ctx.Args[0]
	msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "⏳ Downloading...")
	sent, _ := ctx.API.Send(msg)

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/download/facebook?url=%s&apikey=%s", url, ctx.Config.APIKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to download")
		ctx.API.Send(edit)
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Result  struct {
			Title string `json:"title"`
			HD    string `json:"hd"`
			SD    string `json:"sd"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to parse response")
		ctx.API.Send(edit)
		return err
	}

	if !apiResponse.Status {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to download")
		ctx.API.Send(edit)
		return nil
	}

	ctx.API.Request(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))

	videoURL := apiResponse.Result.HD
	if videoURL == "" {
		videoURL = apiResponse.Result.SD
	}

	video := tgbotapi.NewVideo(ctx.Message.Chat.ID, tgbotapi.FileURL(videoURL))
	video.Caption = fmt.Sprintf("📥 *Facebook Video*\n\n%s", apiResponse.Result.Title)
	video.ParseMode = "Markdown"
	ctx.API.Send(video)

	return nil
}

// Twitter Downloader
type TwitterPlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &TwitterPlugin{}
	plugins.Register(p)
}

func (p *TwitterPlugin) Commands() []string { return []string{"twitter", "tw", "twdl", "x"} }
func (p *TwitterPlugin) Tags() []string     { return []string{"downloader"} }
func (p *TwitterPlugin) Help() string       { return "Download Twitter/X video" }
func (p *TwitterPlugin) RequireLimit() bool { return true }

func (p *TwitterPlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Usage: /twitter <url>"))
		return nil
	}

	url := ctx.Args[0]
	msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "⏳ Downloading...")
	sent, _ := ctx.API.Send(msg)

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/download/twitter?url=%s&apikey=%s", url, ctx.Config.APIKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to download")
		ctx.API.Send(edit)
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Result  struct {
			Desc  string `json:"desc"`
			Video string `json:"video"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to parse response")
		ctx.API.Send(edit)
		return err
	}

	if !apiResponse.Status {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to download")
		ctx.API.Send(edit)
		return nil
	}

	ctx.API.Request(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))

	video := tgbotapi.NewVideo(ctx.Message.Chat.ID, tgbotapi.FileURL(apiResponse.Result.Video))
	video.Caption = fmt.Sprintf("📥 *Twitter Video*\n\n%s", apiResponse.Result.Desc)
	video.ParseMode = "Markdown"
	ctx.API.Send(video)

	return nil
}

// Spotify Downloader
type SpotifyPlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &SpotifyPlugin{}
	plugins.Register(p)
}

func (p *SpotifyPlugin) Commands() []string { return []string{"spotify", "spotifydl"} }
func (p *SpotifyPlugin) Tags() []string     { return []string{"downloader"} }
func (p *SpotifyPlugin) Help() string       { return "Download Spotify track" }
func (p *SpotifyPlugin) RequireLimit() bool { return true }

func (p *SpotifyPlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Usage: /spotify <url>"))
		return nil
	}

	url := ctx.Args[0]
	msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "⏳ Downloading...")
	sent, _ := ctx.API.Send(msg)

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/download/spotify?url=%s&apikey=%s", url, ctx.Config.APIKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to download")
		ctx.API.Send(edit)
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Result  struct {
			Title     string `json:"title"`
			Artist    string `json:"artist"`
			Thumbnail string `json:"thumbnail"`
			Download  string `json:"download"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to parse response")
		ctx.API.Send(edit)
		return err
	}

	if !apiResponse.Status {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to download")
		ctx.API.Send(edit)
		return nil
	}

	ctx.API.Request(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))

	audio := tgbotapi.NewAudio(ctx.Message.Chat.ID, tgbotapi.FileURL(apiResponse.Result.Download))
	audio.Title = apiResponse.Result.Title
	audio.Performer = apiResponse.Result.Artist
	audio.Caption = fmt.Sprintf("🎵 *%s*\n👤 %s", apiResponse.Result.Title, apiResponse.Result.Artist)
	audio.ParseMode = "Markdown"
	ctx.API.Send(audio)

	return nil
}

// MediaFire Downloader
type MediaFirePlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &MediaFirePlugin{}
	plugins.Register(p)
}

func (p *MediaFirePlugin) Commands() []string { return []string{"mediafire", "mf"} }
func (p *MediaFirePlugin) Tags() []string     { return []string{"downloader"} }
func (p *MediaFirePlugin) Help() string       { return "Download from MediaFire" }
func (p *MediaFirePlugin) RequireLimit() bool { return true }

func (p *MediaFirePlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Usage: /mediafire <url>"))
		return nil
	}

	url := ctx.Args[0]
	msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "⏳ Getting file info...")
	sent, _ := ctx.API.Send(msg)

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/download/mediafire?url=%s&apikey=%s", url, ctx.Config.APIKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to get file")
		ctx.API.Send(edit)
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Result  struct {
			Filename string `json:"filename"`
			Filesize string `json:"filesize"`
			Link     string `json:"link"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to parse response")
		ctx.API.Send(edit)
		return err
	}

	if !apiResponse.Status {
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Failed to get file")
		ctx.API.Send(edit)
		return nil
	}

	result := fmt.Sprintf("📁 *MediaFire*\n\n"+
		"📄 File: %s\n"+
		"💾 Size: %s\n\n"+
		"🔗 [Download Link](%s)",
		apiResponse.Result.Filename,
		apiResponse.Result.Filesize,
		apiResponse.Result.Link)

	edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, result)
	edit.ParseMode = "Markdown"
	ctx.API.Send(edit)

	return nil
}
