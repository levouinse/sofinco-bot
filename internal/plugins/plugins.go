package plugins

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/config"
	"github.com/levouinse/sofinco-bot/internal/database"
)

type Context struct {
	API     *tgbotapi.BotAPI
	DB      *database.Database
	Config  *config.Config
	Message *tgbotapi.Message
	User    *database.User
	Args    []string
	Command string
}

type Plugin interface {
	Commands() []string
	Tags() []string
	Help() string
	Execute(ctx *Context) error
	RequireLimit() bool
	RequirePremium() bool
	RequireGroup() bool
	RequireAdmin() bool
}

type BasePlugin struct {
	commands       []string
	tags           []string
	help           string
	requireLimit   bool
	requirePremium bool
	requireGroup   bool
	requireAdmin   bool
}

func (p *BasePlugin) Commands() []string      { return p.commands }
func (p *BasePlugin) Tags() []string          { return p.tags }
func (p *BasePlugin) Help() string            { return p.help }
func (p *BasePlugin) RequireLimit() bool      { return p.requireLimit }
func (p *BasePlugin) RequirePremium() bool    { return p.requirePremium }
func (p *BasePlugin) RequireGroup() bool      { return p.requireGroup }
func (p *BasePlugin) RequireAdmin() bool      { return p.requireAdmin }
func (p *BasePlugin) Execute(ctx *Context) error { return nil }

var Registry = make(map[string]Plugin)

func Register(plugin Plugin) {
	for _, cmd := range plugin.Commands() {
		Registry[cmd] = plugin
	}
}
