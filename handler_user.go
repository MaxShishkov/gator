package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username is required")
	}

	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}

	fmt.Printf("User set to %s\n", cmd.args[0])
	return nil
}
