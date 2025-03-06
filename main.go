package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tylerbartlett24/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	configuration, err := config.Read()
	if err != nil {
		log.Fatalf("could not read config file: %v", err)
	}
	configPtr := &configuration
	sysState := state{
		cfg: configPtr,
	}
	statePtr :=  &sysState
	validCommands := commands{
		commandList: make(map[string]func (*state, command) error),
	}
	validCommands.register("login", HandlerLogin)
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
