package store

import "fmt"

type (
	CommandKey string
	Commands   struct {
		list map[CommandKey]*func()
	}
)

func NewCommands() Commands {
	return Commands{
		list: make(map[CommandKey]*func()),
	}
}

func (c *Commands) Assign(cmds map[CommandKey]*func()) {
	for k, v := range cmds {
		c.list[k] = v
	}
}

func (c *Commands) Execute(key CommandKey) error {
	f, ok := c.list[key]
	if !ok {
		return fmt.Errorf("no command with key %v", key)
	}
	cmd := *f
	cmd()
	return nil
}
