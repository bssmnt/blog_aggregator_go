package main

import (
	"github.com/bssmnt/blog_aggregator_go/internal/commands"
	"github.com/bssmnt/blog_aggregator_go/internal/config"
	"log"
	"os"
)

func main() {
	cmds := &commands.Commands{
		CommandNames: make(map[string]func(*commands.State, commands.Command) error),
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	state := &commands.State{
		Cfg: &cfg,
	}

	cmds.Register("login", commands.HandlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments provided")
	}

	cmd := commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := cmds.Run(state, cmd); err != nil {
		log.Fatal(err)
	}
}
