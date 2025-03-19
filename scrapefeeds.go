package main

import (
	"context"
	"fmt"
	"log"
)

func scrapeFeeds(s *state) {
	next, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Fatal("Next feed to fetch could not be determined.")
	}

	err = s.db.MarkFeedFetched(context.Background(), next.ID)
	if err != nil {
		log.Fatal("Feed could not be marked.")
	}

	feed, err := fetchFeed(context.Background(), next.Url)
	if err != nil {
		log.Fatal("Feed could not be fetched.")
	}

	fmt.Printf("%v:\n", feed.Channel.Title)
	for _, item := range feed.Channel.Item {
		fmt.Println(item)
	}
}