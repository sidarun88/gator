package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/sidarun88/gator/internal/config"
	"github.com/sidarun88/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	gatorConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", gatorConfig.DbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	gatorState := &state{
		db:  dbQueries,
		cfg: gatorConfig,
	}
	gatorCommands := &commands{
		commandsMap: make(map[string]func(*state, command) error),
	}
	gatorCommands.register("login", handlerLogin)
	gatorCommands.register("register", handlerRegister)
	gatorCommands.register("reset", handlerReset)
	gatorCommands.register("users", handlerUsers)
	gatorCommands.register("agg", handlerAggregation)
	gatorCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	gatorCommands.register("feeds", handlerFeeds)
	gatorCommands.register("follow", middlewareLoggedIn(handlerFollow))
	gatorCommands.register("following", middlewareLoggedIn(handlerFollowing))
	gatorCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	gatorCommands.register("browse", middlewareLoggedIn(handlerPosts))

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("expected atleast 1 argument, got %d args: %v", len(args)-1, args[1:])
	}

	cmd := command{
		name: args[1],
		args: args[1:],
	}
	err = gatorCommands.run(gatorState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
