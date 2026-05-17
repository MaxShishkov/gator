package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/MaxShishkov/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const errDupUser = "23505"

func hanndleGetUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("error fetching users: %w", err)
	}

	currUser := s.config.CurrentUserName

	for _, user := range users {
		if user.Name == currUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handleReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("error resetting database: %w", err)
	}

	fmt.Println("Database reset successfully")
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username is required")
	}

	ctx := context.Background()
	user, err := s.db.GetUser(ctx, cmd.args[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("error fetching user: %v does not exist\n", cmd.args[0])
			os.Exit(1)
		}
		return fmt.Errorf("error fetching user: %w\n", err)
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("error setting user: %w\n", err)
	}

	fmt.Printf("User set to %s\n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.name)
	}

	ctx := context.Background()
	userArgs := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	user, err := s.db.CreateUser(ctx, userArgs)
	if pqErr, ok := errors.AsType[*pq.Error](err); ok {
		if pqErr.Code == errDupUser {
			fmt.Println("User already exists")
			os.Exit(1)
		}
	} else if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("error setting user: %w\n", err)
	}

	fmt.Println("User created:")
	printUser(user)
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:		%s\n", user.ID)
	fmt.Printf(" * Name:	%s\n", user.Name)
}
