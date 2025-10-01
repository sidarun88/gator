package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sidarun88/gator/internal/database"

	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("expected 1 argument, got %d args: %v", len(cmd.args)-1, cmd.args[1:])
	}

	username := cmd.args[1]
	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("user %s already exists", username)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	dbUser, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(dbUser.Name)
	if err != nil {
		return err
	}

	fmt.Printf("user %s created\n", username)
	fmt.Printf("Data %v\n", dbUser)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("expected 1 argument, got %d args: %v", len(cmd.args)-1, cmd.args[1:])
	}

	username := cmd.args[1]
	dbUser, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(dbUser.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s user logged in successfully\n", dbUser.Name)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected no arguments, got %d args: %v", len(cmd.args)-1, cmd.args[1:])
	}

	dbUsers, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, u := range dbUsers {
		if u.Name == s.cfg.Username {
			fmt.Printf("* %s (current)\n", u.Name)
		} else {
			fmt.Printf("* %s\n", u.Name)
		}
	}

	return nil
}
