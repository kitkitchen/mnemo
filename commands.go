package mnemo

import "fmt"

type (
	// CommandKey is a unique identifier for a command.
	CommandKey string
	// Commands is a collection of commands.
	Commands struct {
		list map[CommandKey]func()
	}
)

// NewCommands creates a new collection of commands.
func NewCommands() Commands {
	return Commands{
		list: make(map[CommandKey]func()),
	}
}

// Assign assigns a map of commands to the collection.
func (c *Commands) Assign(cmds map[CommandKey]func()) {
	for k, v := range cmds {
		c.list[k] = v
	}
}

// Execute executes a command and returns an error if the command does not exist.
func (c *Commands) Execute(key CommandKey) error {
	f, ok := c.list[key]
	if !ok {
		return fmt.Errorf("no command with key %v", key)
	}
	cmd := f
	cmd()
	return nil
}
