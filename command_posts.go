package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sidarun88/gator/internal/database"
)

func handlerPosts(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 && len(cmd.args) != 1 {
		return fmt.Errorf("expected 0 or 1 arguments, got %d args: %v", len(cmd.args)-1, cmd.args[1:])
	}

	var limit int32 = 2
	if len(cmd.args) == 2 {
		parsedLimit, err := strconv.Atoi(cmd.args[1])
		if err != nil {
			return err
		}
		limit = int32(parsedLimit)
	}

	postParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}
	dbPosts, err := s.db.GetPostsForUser(context.Background(), postParams)
	if err != nil {
		return err
	}

	for _, post := range dbPosts {
		fmt.Printf("Post ID: %v\n", post.ID)
		fmt.Printf("  - Title: %v\n", post.Title)
		fmt.Printf("  - Published Date: %v\n", post.PublishedAt)
	}

	return nil
}
