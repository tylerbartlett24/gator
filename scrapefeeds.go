package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tylerbartlett24/gator/internal/database"
)

func scrapeFeeds(s *state) {
	ctx := context.Background()
	next, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Fatal("Next feed to fetch could not be determined.")
	}

	err = s.db.MarkFeedFetched(ctx, next.ID)
	if err != nil {
		log.Fatal("Feed could not be marked.")
	}

	feed, err := fetchFeed(ctx, next.Url)
	if err != nil {
		log.Fatal("Feed could not be fetched.")
	}

	currentTime := time.Now()
	updateTime := sql.NullTime{
		Time:  currentTime,
		Valid: true,
	}
	for _, item := range feed.Channel.Item {
		description := sql.NullString{
			String: item.Description,
			Valid: true,
		}
		pubAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			pubAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		params := database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: currentTime,
			UpdatedAt: updateTime,
			Title: item.Title,
			Url: item.Link,
			Description: description,
			PublishedAt: pubAt,
			FeedID: next.ID,

		}
		err = s.db.CreatePost(ctx, params)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
}