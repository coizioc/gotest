package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CmdFunction is the type that defines functions that are executing using the command handler.
type CmdFunction func(ctx *Ctx, args []string) error

// CmdHandler contains the session, prefix, and command functions
type CmdHandler struct {
	s        *discordgo.Session
	prefix   string
	commands map[string]CmdFunction
}

// NewCmdHandler creates a new CmdHandler with session and prefix as parameters
func NewCmdHandler(s *discordgo.Session, prefix string) *CmdHandler {
	return &CmdHandler{s, prefix, map[string]CmdFunction{}}
}

// NewCmd adds a new command to the CmdHandler.
func (c *CmdHandler) NewCmd(name string, cmd CmdFunction) {
	c.commands[name] = cmd
}

// Handle handles a command from a new message using the CmdHandler.
func (c *CmdHandler) Handle(m *discordgo.MessageCreate) {
	// Generate Ctx object.
	ctx, err := GetCtx(c.s, m)
	if err != nil {
		return
	}

	// Check that the message did not come from the bot itself.
	if m.Author.ID == c.s.State.User.ID {
		return
	}

	// Check if message has the bot prefix.
	content := m.Content
	if !strings.HasPrefix(content, c.prefix) {
		return
	}

	// Remove prefix from content and create args from content.
	content = content[len(c.prefix):len(content)]
	args := strings.Split(content, " ")

	// If args[0] is a command in CmdHandler, run the command.
	if cmdf, ok := c.commands[args[0]]; ok {
		err := cmdf(ctx, args)
		if err != nil {
			fmt.Println("error running command: %s", err)
		}
	}
}
