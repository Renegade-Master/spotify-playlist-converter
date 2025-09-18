package main

import (
	"context"
	"flag"
	"log"

	"github.com/Renegade-Master/spotify-playlist-converter/internal/spotify"
)

func main() {
	log.Printf("PlaylistConverter")

	userID := flag.String("user", "", "the Spotify user ID to look up")
	flag.Parse()

	ctx := context.Background()

	user := spotify.NewSpotifyUser(ctx, *userID)

	log.Printf("User ID: [%s]", user.ID)
	log.Printf("Display name: [%s]", user.DisplayName)
	log.Printf("Spotify URI: [%s]", string(user.URI))
	log.Printf("Endpoint: [%s]", user.Endpoint)

	spotify.Authenticate()
}
