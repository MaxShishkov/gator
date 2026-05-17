package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/MaxShishkov/gator/internal/config"
	"github.com/MaxShishkov/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	dbQueries := database.New(db)

	programState := state{db: dbQueries, config: &cfg}
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
	cmds.register("register", handlerRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", hanndleGetUsers)
	return &cmds
}
