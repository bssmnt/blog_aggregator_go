package commands

import (
	"blog_aggregator_go/internal/config"
	"errors"
	"fmt"
)

type State struct {
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CommandNames map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandNames[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	if handler, exists := c.CommandNames[cmd.Name]; exists {
		return handler(s, cmd)
	} else {
		return errors.New("command not found")
	}
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("please specify a username: gator login <username>")
	}

	username := cmd.Args[0]
	err := s.Cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("User set to:", username)
	return nil
}
