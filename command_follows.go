package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("expected 1 argument, got %d args: %v", len(cmd.args)-1, cmd.args[1:])
	}

	dbFeed, err := s.db.Feed(context.Background(), cmd.args[1])
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
	dbFollow, err := s.db.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return err
	}

	fmt.Printf("%s user is now following %s feed", dbFollow.UserName, dbFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected no arguments, got %d args: %v", len(cmd.args)-1, cmd.args[1:])
	}

	dbUserFollowing, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(dbUserFollowing) == 0 {
		fmt.Printf("%s user is not following any feeds", user.Name)
		return nil
	}

	fmt.Printf("Feeds being followed by current user %s:\n", user.Name)
	for _, follow := range dbUserFollowing {
		fmt.Printf(" - %s\n", follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("expected 1 argument, got %d args: %v", len(cmd.args)-1, cmd.args[1:])
	}

	dbFeed, err := s.db.Feed(context.Background(), cmd.args[1])
	if err != nil {
		return err
	}

	removeFollowParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: dbFeed.ID,
	}
	err = s.db.DeleteFeedFollow(context.Background(), removeFollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("User %s unfollowed %s feed", user.Name, dbFeed.Name)
	return nil
}
