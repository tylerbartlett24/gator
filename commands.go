package main

type command struct {
	Name      string
	Arguments []string
}

type commands struct {
	commandList map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandList[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	err := c.commandList["login"](s, cmd)
	return err
}