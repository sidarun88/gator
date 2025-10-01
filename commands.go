package main

import "fmt"

type command struct {
	name string
	args []string
}

type commands struct {
	commandsMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("expected atleast 1 argument, got %d args: %v", len(cmd.args)-1, cmd.args[1:])
	}

	cmdName := cmd.args[0]
	f, ok := c.commandsMap[cmdName]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmdName)
	}

	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandsMap[name] = f
}
