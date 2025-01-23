package main

import (
	"blog_aggregator_go/internal/commands"
	"blog_aggregator_go/internal/config"
	"blog_aggregator_go/internal/database"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	cmds := &commands.Commands{
		CommandNames: make(map[string]func(*commands.State, commands.Command) error),
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	state := &commands.State{
		Cfg: &cfg,
		Db:  dbQueries,
	}

	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)

	if len(os.Args) < 2 {
		log.Fatal("please provide a command")
	}

	cmd := commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := cmds.Run(state, cmd); err != nil {
		log.Fatal(err)
	}

}
