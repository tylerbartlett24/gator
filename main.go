package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/tylerbartlett24/gator/internal/config"
	"github.com/tylerbartlett24/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

const dbURL = "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

func main() {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database")
	}
	dbQueries := database.New(db)
	configuration, err := config.Read()
	if err != nil {
		log.Fatalf("could not read config file: %v", err)
	}
	configPtr := &configuration
	sysState := state{
		db: dbQueries,
		cfg: configPtr,
	}
	statePtr :=  &sysState
	validCommands := commands{
		commandList: make(map[string]func (*state, command) error),
	}
	
	validCommands.register("login", HandlerLogin)
	validCommands.register("register", HandlerRegister)
	validCommands.register("reset", HandlerReset)
	validCommands.register("users", HandlerUsers)
	validCommands.register("agg", HandlerAgg)
	validCommands.register("addfeed", middlewareLoggedIn(HandlerAddFeed))
	validCommands.register("feeds", HandlerFeeds)
	validCommands.register("follow", middlewareLoggedIn(HandlerFollow))
	validCommands.register("following", middlewareLoggedIn(HandlerFollowing))
	validCommands.register("unfollow", middlewareLoggedIn(HandlerUnfollow))

	args := os.Args
	if len(args) < 2 {
		log.Fatal("no command supplied")
	}
	input := command{
		Name: args[1],
		Arguments: args[2:],
	}
	err = validCommands.run(statePtr, input)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}