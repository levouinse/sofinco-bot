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

type GHStalkPlugin struct {
	plugins.BasePlugin
}

func init() {
	plugins.Register(&GHStalkPlugin{})
}

func (p *GHStalkPlugin) Commands() []string { return []string{"ghstalk", "github", "gh"} }
func (p *GHStalkPlugin) Tags() []string     { return []string{"stalker"} }
func (p *GHStalkPlugin) Help() string       { return "Stalk GitHub profile" }
func (p *GHStalkPlugin) RequireLimit() bool { return false }

func (p *GHStalkPlugin) Execute(ctx *plugins.Context) error {
	if len(ctx.Args) == 0 {
		msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "Masukkan username GitHub!\n\nContoh:\n/ghstalk username")
		ctx.API.Send(msg)
		return nil
	}

	username := strings.TrimPrefix(ctx.Args[0], "@")
	
	waitMsg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "🔍 Stalking GitHub...")
	sent, _ := ctx.API.Send(waitMsg)

	apiURL := fmt.Sprintf("https://api.github.com/users/%s", url.QueryEscape(username))

	resp, err := http.Get(apiURL)
	if err != nil {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("gagal stalk: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("user tidak ditemukan")
	}

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Login      string `json:"login"`
		Name       string `json:"name"`
		Bio        string `json:"bio"`
		Company    string `json:"company"`
		Location   string `json:"location"`
		Email      string `json:"email"`
		Blog       string `json:"blog"`
		Followers  int64  `json:"followers"`
		Following  int64  `json:"following"`
		PublicRepos int64 `json:"public_repos"`
		PublicGists int64 `json:"public_gists"`
		AvatarURL  string `json:"avatar_url"`
		HTMLURL    string `json:"html_url"`
		CreatedAt  string `json:"created_at"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))
		return fmt.Errorf("gagal parse data")
	}

	ctx.API.Send(tgbotapi.NewDeleteMessage(ctx.Message.Chat.ID, sent.MessageID))

	bio := result.Bio
	if bio == "" {
		bio = "No bio"
	}

	company := result.Company
	if company == "" {
		company = "-"
	}

	location := result.Location
	if location == "" {
		location = "-"
	}

	caption := fmt.Sprintf("🐙 *GitHub Profile*\n\n"+
		"👤 Username: %s\n"+
		"📝 Name: %s\n"+
		"📄 Bio: %s\n"+
		"🏢 Company: %s\n"+
		"📍 Location: %s\n\n"+
		"👥 Followers: %d\n"+
		"➡️ Following: %d\n"+
		"📦 Repos: %d\n"+
		"📝 Gists: %d\n\n"+
		"🔗 Profile: %s",
		result.Login,
		result.Name,
		bio,
		company,
		location,
		result.Followers,
		result.Following,
		result.PublicRepos,
		result.PublicGists,
		result.HTMLURL)

	photo := tgbotapi.NewPhoto(ctx.Message.Chat.ID, tgbotapi.FileURL(result.AvatarURL))
	photo.Caption = caption
	photo.ParseMode = "Markdown"
	ctx.API.Send(photo)

	return nil
}
