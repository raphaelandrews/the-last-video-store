package main

import (
	"fmt"
	"os"

	"github.com/thelastvideostore/internal/store"
)

func main() {
	cfgPath := "thelastvideostore.db"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	os.Remove(cfgPath)

	s, err := store.Open(cfgPath)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	seedUsers(s)
	seedMovies(s)
	seedMerch(s)
	seedSnackBar(s)
	seedGames(s)
	fmt.Printf("Seeded %d movies/series/games, %d users, 26 merch items, and snack bar.\n", countMovies(s), 12)
}
