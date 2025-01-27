package main

import (
	"blog_aggregator_go/internal/commands"
	"blog_aggregator_go/internal/config"
	"blog_aggregator_go/internal/database"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	dbQueries, err := database.InitDB(dbURL)
	if err != nil {
		log.Fatal(err)
	}

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
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerUsers)
	cmds.Register("agg", commands.HandlerAgg)
	cmds.Register("addfeed", commands.HandlerAddFeed)
	cmds.Register("feeds", commands.HandlerFeeds)
	cmds.Register("follow", commands.Follow)
	cmds.Register("following", commands.Following)

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
