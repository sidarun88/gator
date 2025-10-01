package main

import (
	"context"
	"errors"
	"fmt"
	"gator/internal/database"
	"gator/internal/rss"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func scrapeFeeds(s *state) error {
	toFetchFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Fetching feed for %s...\n", toFetchFeed.Name)
	feed, err := rss.FetchFeed(context.Background(), toFetchFeed.Url)
	if err != nil {
		return err
	}

	markFetchedParams := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		ID:        toFetchFeed.ID,
	}
	err = s.db.MarkFeedFetched(context.Background(), markFetchedParams)
	if err != nil {
		return err
	}

	fmt.Printf("Saving feed posts for %s...\n", toFetchFeed.Name)
	for _, t := range feed.Channel.Items {
		pubDate, err := time.Parse(time.RFC1123Z, t.PubDate)
		if err != nil {
			fmt.Printf("Error parsing pub date: %s\n", err)
			pubDate = time.Now()
		}

		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       t.Title,
			Description: t.Description,
			Url:         t.Link,
			PublishedAt: pubDate,
			FeedID:      toFetchFeed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			var pqErr *pq.Error
			isDupUrlErr := errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation"
			if !isDupUrlErr {
				fmt.Printf("Error creating post: %s\n", err)
			}
		}
	}
	fmt.Println("Saved feed posts")
	fmt.Println()

	return nil
}

func handlerAggregation(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("expected 1 argument, got %d args: %v", len(cmd.args)-1, cmd.args[1:])
	}

	dur, err := time.ParseDuration(cmd.args[1])
	if err != nil {
		return err
	}

	fmt.Printf("Collectind feeds every %s\n", cmd.args[1])
	ticker := time.NewTicker(dur)
	for ; ; <-ticker.C {
		_ = scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 3 {
		return fmt.Errorf("expected 2 arguments, got %d args: %v", len(cmd.args)-1, cmd.args[1:])
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[1],
		Url:       cmd.args[2],
		UserID:    user.ID,
	}
	dbFeed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    dbFeed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return err
	}

	fmt.Println(dbFeed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected no arguments, got %d args: %v", len(cmd.args)-1, cmd.args[1:])
	}

	feeds, err := s.db.Feeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("- Name: %s\n", feed.Name)
		fmt.Printf("- URL: %s\n", feed.Url)
		fmt.Printf("- User: %s\n", feed.UserName)
	}

	return nil
}
