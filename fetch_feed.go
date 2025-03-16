package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	blank := RSSFeed{}
	ptr := &blank
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return ptr, err
	}

	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ptr, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return ptr, err
	}

	err = xml.Unmarshal(dat, ptr)
	ptr.Channel.Title = html.UnescapeString(ptr.Channel.Title)
	ptr.Channel.Description = html.UnescapeString(ptr.Channel.Description)
	for i, item := range ptr.Channel.Item {
		ptr.Channel.Item[i].Description = html.UnescapeString(item.Description)
		ptr.Channel.Item[i].Title = html.UnescapeString(item.Title)
	}

	return ptr, err
}