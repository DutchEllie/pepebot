package mux

import (
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Execute(s *discordgo.Session, m *discordgo.MessageCreate)
}

func NotFound(s *discordgo.Session, m *discordgo.MessageCreate) {
	return
}

func NotFoundHandler() Command { return HandlerFunc(NotFound) }

type HandlerFunc func(s *discordgo.Session, m *discordgo.MessageCreate)

func (f HandlerFunc) Execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	f(s, m)
}

/* The CommandMux struct is a type of mux for Discord commands. It's modelled after the net/http ServeMux */
type CommandMux struct {
	mu     sync.RWMutex
	m      map[string]muxEntry
	Prefix string
}

func NewCommandMux() *CommandMux { return new(CommandMux) }

func (c *CommandMux) removeFirst(command string) string {
	split := strings.SplitN(strings.TrimSpace(command), " ", 2)
	if len(split) > 1 {
		return split[1]
	}
	return ""
}

func (c *CommandMux) firstCommand(command string) string {
	split := strings.SplitN(strings.TrimSpace(command), " ", 2)
	if len(split) > 0 {
		return split[0]
	}
	return ""
}

func (c *CommandMux) Handler(m *discordgo.MessageCreate) (cmd Command, pattern string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if strings.HasPrefix(m.Content, c.Prefix) {
		/* Special case for this bot alone. It has a command that is only it's prefix
		So we check if the whole message is only the prefix before proceding.
		So please don't forget to add the command, since it's totally hardcoded here. */
		if strings.TrimSpace(m.Content) == c.Prefix {
			return c.m[c.Prefix].h, c.m[c.Prefix].pattern
		}

		m := c.removeFirst(m.Content) /* Here the prefix is removed, so we're left with only the first keyword */
		cmd, ok := c.m[c.firstCommand(m)]
		if ok {
			return cmd.h, cmd.pattern
		}
	}

	/* Here is where I might add the whole checking for bad words part */
	return nil, ""
}

func (c *CommandMux) Execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	h, _ := c.Handler(m)
	if h == nil {
		//log.Printf("There exists no handler for %s\n", m.Content)
		return
	}
	h.Execute(s, m)
}

func (c *CommandMux) Handle(pattern string, handler Command) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if pattern == "" {
		panic("commandmux: invalid pattern")
	}
	if handler == nil {
		panic("commandmux: nil handler")
	}
	if _, exist := c.m[pattern]; exist {
		panic("commandmux: multiple registrations for " + pattern)
	}

	if c.m == nil {
		c.m = make(map[string]muxEntry)
	}
	e := muxEntry{h: handler, pattern: pattern}
	c.m[pattern] = e
}

func (c *CommandMux) HandleFunc(pattern string, handler func(s *discordgo.Session, m *discordgo.MessageCreate)) {
	if handler == nil {
		panic("commandmux: nil handler")
	}
	c.Handle(pattern, HandlerFunc(handler))
}

/* The muxEntry struct contains the actual Command implementation as well as the pattern (discord command)
it will be matched against */
type muxEntry struct {
	h       Command
	pattern string
}
