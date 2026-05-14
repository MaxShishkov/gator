package main

import (
	"log"
	"os"

	"github.com/MaxShishkov/gator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	programState := state{config: &cfg}
	cmds := registerCommands()
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("no command provided")
	}
	cmd := command{args[1], args[2:]}
	err = cmds.run(&programState, cmd)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}

}

func registerCommands() *commands {
	cmds := commands{handlers: make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogin)
	return &cmds
}
