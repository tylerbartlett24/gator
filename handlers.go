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

func HandlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	return err
}

func HandlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	currentUser := s.cfg.Username
	for _, user := range users {
		output := "* " + user
		if user == currentUser {
			output += " (current)"
		}
		fmt.Println(output)
	}
	return err
}

func HandlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("could not fetch feed: %w", err)
	}

	fmt.Println(feed)
	return err
}

func HandlerAddFeed(s *state, cmd command) error {
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.Username)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")
	return err
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}

func HandlerFeeds (s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Println(feed)
	}
	return err
}
