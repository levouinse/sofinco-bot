package owner

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/levouinse/sofinco-bot/internal/plugins"
)

type ExecPlugin struct {
	plugins.BasePlugin
}

func init() {
	p := &ExecPlugin{}
	plugins.Register(p)
}

func (p *ExecPlugin) Commands() []string { return []string{"exec", "$"} }
func (p *ExecPlugin) Tags() []string     { return []string{"owner"} }
func (p *ExecPlugin) Help() string       { return "Execute shell command (owner only)" }
func (p *ExecPlugin) RequireLimit() bool { return false }

func (p *ExecPlugin) Execute(ctx *plugins.Context) error {
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
		ctx.API.Send(tgbotapi.NewMessage(ctx.Message.Chat.ID, "Usage: /exec <command>"))
		return nil
	}

	command := strings.Join(ctx.Args, " ")
	
	msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, "⏳ Executing...")
	sent, _ := ctx.API.Send(msg)

	// Execute command with timeout
	cmdExec := exec.Command("bash", "-c", command)
	
	done := make(chan error, 1)
	var output []byte
	var err error

	go func() {
		output, err = cmdExec.CombinedOutput()
		done <- err
	}()

	select {
	case <-time.After(30 * time.Second):
		cmdExec.Process.Kill()
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, "❌ Command timeout (30s)")
		ctx.API.Send(edit)
		return nil
	case err := <-done:
		result := string(output)
		if result == "" {
			result = "✅ Command executed successfully (no output)"
		}
		
		if err != nil {
			result = fmt.Sprintf("❌ Error: %v\n\nOutput:\n%s", err, result)
		}

		// Limit output length
		if len(result) > 4000 {
			result = result[:4000] + "\n\n... (truncated)"
		}

		response := fmt.Sprintf("```bash\n$ %s\n\n%s\n```", command, result)
		edit := tgbotapi.NewEditMessageText(ctx.Message.Chat.ID, sent.MessageID, response)
		edit.ParseMode = "Markdown"
		ctx.API.Send(edit)
	}

	return nil
}
