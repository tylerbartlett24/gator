package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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
	if len(cmd.Arguments) != 1 {
		return errors.New("usage: agg <time between requests>")
	}

	duration, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return err
	}
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func HandlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("usage: addfeed <feed name> <feed url>")
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

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    user.ID,
		FeedID:    params.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	} else {
		fmt.Printf("%v is now following %v.\n", s.cfg.Username, feed.Name)
	}

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


func HandlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("Name: %v, URL: %v, Creator: %v\n", feed.Name, feed.Url, 
		feed.Name_2.String)
	}
	return err
}

func HandlerFollow(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	feed, err := s.db.GetFeed(ctx, cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("could not locate feeds: %w", err)
	}

	currentTime := time.Now()
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(ctx, params)
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	}

	fmt.Printf("%v is now following %v\n", user.Name, feed.Name)
	return err
}

func HandlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not find follows: %w", err)
	}

	if len(follows) == 0 {
		fmt.Printf("%v does not follow any feeds.\n", user.Name)
		return err
	}
	fmt.Printf("%v is following:\n", follows[0].Username)
	for _, follow := range follows {
		fmt.Println(follow.Name)
	}
	return err
}

func HandlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("usage: unfollow <feed url>")
	}

	ctx := context.Background()
	feed, err := s.db.GetFeed(ctx, cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("could not find feed: %w", err)
	}

	params := database.DeleteFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFollow(ctx, params)
	if err == nil {
		fmt.Printf("%v no longer follows %v.\n", user.Name, feed.Name)
	}
	return err
}

func HandlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Arguments) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Arguments[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}
