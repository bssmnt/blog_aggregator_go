package commands

import (
	"blog_aggregator_go/internal/config"
	"blog_aggregator_go/internal/database"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CommandNames map[string]func(*State, Command) error
}

func NewCommands() *Commands {
	return &Commands{
		CommandNames: make(map[string]func(*State, Command) error),
	}
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

	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user %s not found", username)
	}

	err = s.Cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("user set to:", username)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {

	if len(cmd.Args) != 1 {
		return errors.New("please specify a username: register <username>")
	}
	username := cmd.Args[0]
	ctx := context.Background()
	params := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: username}

	user, err := s.Db.CreateUser(ctx, params)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return fmt.Errorf("user %s already exists", username)
			}
		}
		return err
	}

	err = s.Cfg.SetUser(username)
	if err != nil {
		return err
	}

	err = config.Save(*s.Cfg)
	if err != nil {
		return err
	}

	fmt.Println("user set to:", user)

	return nil
}

func Reset(s *State, cmd Command) error {
	err := s.Db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func Users(s *State, cmd Command) error {
	allUsers, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(allUsers) == 0 {
		fmt.Println("no users found")
	}

	currentUser := s.Cfg.CurrentUserName
	for _, user := range allUsers {
		if user == currentUser {
			fmt.Printf("* %s (current)\n", currentUser)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}
