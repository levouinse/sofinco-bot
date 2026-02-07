package stalker

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

type IGStalkPlugin struct {
	plugins.BasePlugin
}

func init() {
	plugins.Register(&IGStalkPlugin{})
}

func (p *IGStalkPlugin) Commands() []string { return []string{"igstalk", "ig", "instagram"} }
func (p *IGStalkPlugin) Tags() []string     { return []string{"stalker"} }
func (p *IGStalkPlugin) Help() string       { return "Stalk Instagram profile" }
func (p *IGStalkPlugin) RequireLimit() bool { return true }

func (p *IGStalkPlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "Masukkan username Instagram!\n\nContoh:\n/igstalk username")
		ctx.API.Send(msg)
		return nil
	}

	username := strings.TrimPrefix(ctx.Args[0], "@")
	
	waitMsg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "🔍 Stalking...")
	sent, _ := ctx.API.Send(waitMsg)

	apiURL := fmt.Sprintf("https://api.betabotz.eu.org/api/stalk/ig?username=%s&apikey=%s",
		url.QueryEscape(username), ctx.Config.APIKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("gagal stalk: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Status bool `json:"status"`
		Result struct {
			Username  string `json:"username"`
			FullName  string `json:"fullName"`
			Bio       string `json:"biography"`
			Followers int64  `json:"followers"`
			Following int64  `json:"following"`
			Posts     int64  `json:"posts"`
			IsPrivate bool   `json:"isPrivate"`
			IsVerified bool  `json:"isVerified"`
			ProfilePic string `json:"profilePicUrl"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil || !result.Status {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("user tidak ditemukan")
	}

	ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))

	verified := ""
	if result.Result.IsVerified {
		verified = "✓"
	}
	
	private := "Public"
	if result.Result.IsPrivate {
		private = "Private"
	}

	caption := fmt.Sprintf("📸 *Instagram Profile*\n\n"+
		"👤 Username: @%s %s\n"+
		"📝 Name: %s\n"+
		"📄 Bio: %s\n\n"+
		"👥 Followers: %d\n"+
		"➡️ Following: %d\n"+
		"📷 Posts: %d\n"+
		"🔒 Status: %s",
		result.Result.Username, verified,
		result.Result.FullName,
		result.Result.Bio,
		result.Result.Followers,
		result.Result.Following,
		result.Result.Posts,
		private)

	photo := tgbotapi.NewPhoto(ctx.Message.Chat.ID, tgbotapi.FileURL(result.Result.ProfilePic))
	photo.Caption = caption
	photo.ParseMode = "Markdown"
	ctx.API.Send(photo)

	if !ctx.User.Premium {
		ctx.User.Limit--
		ctx.DB.SaveUser(ctx.User)
	}

	return nil
}
