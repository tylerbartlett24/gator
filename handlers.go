package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tylerbartlett24/gator/internal/database"
)

func HandlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("no username given")
	}
	_, err := s.db.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		return errors.New("user not in database")
	}
	err = s.cfg.SetUser(cmd.Arguments[0])
	if err != nil {
		return err
	}
	
	fmt.Printf("Current user has been set to %v.\n", cmd.Arguments[0])
	return nil
}

func HandlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("no username given")
	}
	currentTime := time.Now()
	username := cmd.Arguments[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("user %s already exists", username)
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      username,
	}
	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil { 
		return err
	}

	err = s.cfg.SetUser(username)
	if err != nil { 
		return err
	}

	fmt.Printf("User created.\nInfo: %+v\n", user)
	return nil	
}