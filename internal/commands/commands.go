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
		return errors.New("please specify a username: login <username>")
	}

	username := cmd.Args[0]
	err := s.Cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("user set to:", username)
	return nil
}

type CliCommand struct {
	name        string
	description string
	callback    func(cfg *config.Config, args ...string) error
}

//func GetCommands() map[string]CliCommand {
//	return map[string]CliCommand{
//		"help": {
//			name:        "help",
//			description: "displays a help message",
//			callback:    CommandHelp,
//		},
//		"login": {
//			name:        "login",
//			description: "login to the system",
//			callback:    LoginHelp,
//		},
//	}
//}
//
//func CommandHelp(*config.Config, ...string) error {
//	for _, cmd := range GetCommands() {
//		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
//	}
//	return nil
//}
//
//func LoginHelp(*config.Config, ...string) error {
//	if len(os.Args) != 1 {
//		fmt.Println("please specify a username: login <username>")
//	}
//	return nil
//}
