package tools

import (
	"fmt"
	"runtime"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

type StatsPlugin struct {
	plugins.BasePlugin
	startTime time.Time
}

func init() {
	p := &StatsPlugin{
		startTime: time.Now(),
	}
	plugins.Register(p)
}

func (p *StatsPlugin) Commands() []string { return []string{"stats", "statistics"} }
func (p *StatsPlugin) Tags() []string     { return []string{"info"} }
func (p *StatsPlugin) Help() string       { return "Show bot statistics" }
func (p *StatsPlugin) RequireLimit() bool { return false }

func (p *StatsPlugin) Execute(ctx *plugins.Context) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(p.startTime).Round(time.Second)
	
	// Get database stats
	totalUsers := ctx.DB.GetTotalUsers()
	
	msg := fmt.Sprintf("📊 *Bot Statistics*\n\n"+
		"👥 Total Users: %d\n"+
		"⏰ Uptime: %s\n"+
		"💾 Memory: %.2f MB\n"+
		"🔧 Goroutines: %d\n"+
		"🤖 Go Version: %s",
		totalUsers,
		uptime,
		float64(m.Alloc)/1024/1024,
		runtime.NumGoroutine(),
		runtime.Version())

	reply := tgbotapi.NewMessage(ctx.Message.Chat.ID, msg)
	reply.ParseMode = "Markdown"
	ctx.API.Send(reply)

	return nil
}
