package owner

import (
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

type AddPremiumPlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &AddPremiumPlugin{}
	plugins.Register(p)
}

func (p *AddPremiumPlugin) Commands() []string { return []string{"addprem", "addpremium"} }
func (p *AddPremiumPlugin) Tags() []string     { return []string{"owner"} }
func (p *AddPremiumPlugin) Help() string       { return "Add premium user (owner only)" }
func (p *AddPremiumPlugin) RequireLimit() bool { return false }

func (p *AddPremiumPlugin) Execute(ctx *plugins.Context) error {
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

	if len(ctx.Args) < 1 {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, 
			"Usage: /addprem <user_id> [days]\n\nExample:\n/addprem 123456789 30"))
		return nil
	}

	userID, err := strconv.ParseInt(ctx.Args[0], 10, 64)
	if err != nil {
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "❌ Invalid user ID"))
		return nil
	}

	days := 30 // default 30 days
	if len(ctx.Args) > 1 {
		if d, err := strconv.Atoi(ctx.Args[1]); err == nil {
			days = d
		}
	}

	// Get or create user
	user, err := ctx.DB.GetOrCreateUser(userID, "", "")
	if err != nil {
		return fmt.Errorf("gagal mendapatkan user: %w", err)
	}

	// Set premium
	user.Premium = true
	user.PremiumUntil = time.Now().Add(time.Duration(days) * 24 * time.Hour)
	
	if err := ctx.DB.SaveUser(user); err != nil {
		return fmt.Errorf("gagal menyimpan user: %w", err)
	}

	msg := fmt.Sprintf("✅ Premium berhasil ditambahkan!\n\n"+
		"👤 User ID: %d\n"+
		"⏰ Durasi: %d hari\n"+
		"📅 Expired: %s",
		userID, days, user.PremiumUntil.Format("2006-01-02 15:04:05"))

	ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, msg))

	// Notify user
	notif := tgbotapi.NewMessage(userID, 
		fmt.Sprintf("🎉 *Selamat!*\n\n"+
			"Kamu telah mendapatkan akses Premium selama %d hari!\n\n"+
			"✨ Benefit:\n"+
			"• Unlimited limit\n"+
			"• Priority support\n"+
			"• Access to premium features\n\n"+
			"Expired: %s",
			days, user.PremiumUntil.Format("2006-01-02")))
	notif.ParseMode = "Markdown"
	ctx.API.Send(notif)

	return nil
}
