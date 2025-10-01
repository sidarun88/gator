package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected no arguments, got %d args: %v", len(cmd.args)-1, cmd.args[1:])
	}

	err := s.db.ResetDB(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Reset DB successfully")
	return nil
}
