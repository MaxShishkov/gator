package main

import (
	"context"
	"fmt"
)

func handleAgg(s *state, cmd command) error {
	var url string
	if len(cmd.args) == 0 {
		url = "https://www.wagslane.dev/index.xml"
		//return fmt.Errorf("feed URL is required")
	} else {
		url = cmd.args[0]
	}

	ctx := context.Background()
	feed, err := fetchFeed(ctx, url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	printFeed(feed)
	return nil
}

func printFeed(feed *RSSFeed) {
	fmt.Printf("Feed: %s\n", feed.Channel.Title)
	fmt.Printf("Description: %s\n", feed.Channel.Description)
	fmt.Printf("Link: %s\n", feed.Channel.Link)
	fmt.Println("Items:")
	for _, item := range feed.Channel.Item {
		fmt.Printf("- %s\n  Link: %s\n  Published: %s\n  Description: %s\n", item.Title, item.Link, item.PubDate, item.Description)
	}
}
