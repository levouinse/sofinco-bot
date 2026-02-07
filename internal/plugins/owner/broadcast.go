package owner

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

type BroadcastPlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &BroadcastPlugin{}
	plugins.Register(p)
}

func (p *BroadcastPlugin) Commands() []string { return []string{"broadcast", "bc"} }
func (p *BroadcastPlugin) Tags() []string     { return []string{"owner"} }
func (p *BroadcastPlugin) Help() string       { return "Broadcast message to all users (owner only)" }
func (p *BroadcastPlugin) RequireLimit() bool { return false }

func (p *BroadcastPlugin) Execute(ctx *plugins.Context) error {
	// Check if user is owner
	isOwner := false
	for _, ownerID := range ctx.Config.OwnerIDs {
		if ctx.Message.From.ID == ownerID {
			isOwner = true
			break
		}
	}

	if !isOwner {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "❌ Command ini hanya untuk owner!"))
		return nil
	}

	if len(ctx.Args) == 0 {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Usage: /broadcast <message>"))
		return nil
	}

	message := strings.Join(ctx.Args, " ")
	
	// Get all users
	users := ctx.DB.GetAllUsers()
	
	statusMsg := tgbotapi.NewMessage(ctx.Message.Chat.ID, 
		fmt.Sprintf("📢 Broadcasting to %d users...", len(users)))
	sent, _ := ctx.API.Send(statusMsg)

	success := 0
	failed := 0

	for _, user := range users {
		msg := tgbotapi.NewMessage(user.ID, 
			fmt.Sprintf("📢 *Broadcast Message*\n\n%s", message))
		msg.ParseMode = "Markdown"
		
		if _, err := ctx.API.Send(msg); err != nil {
			failed++
		} else {
			success++
		}
	}

	result := fmt.Sprintf("✅ Broadcast selesai!\n\n"+
		"✓ Berhasil: %d\n"+
		"✗ Gagal: %d\n"+
		"📊 Total: %d",
		success, failed, len(users))

	edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, result)
	ctx.API.Send(edit)

	return nil
}
