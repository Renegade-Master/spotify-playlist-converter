package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	log.Printf("PlaylistConverter")

	userID := flag.String("user", "", "the Spotify user ID to look up")
	flag.Parse()

	ctx := context.Background()

	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	user, _ := client.GetUsersPublicProfile(ctx, spotify.ID(*userID))

	log.Println("User ID:", user.ID)
	log.Println("Display name:", user.DisplayName)
	log.Println("Spotify URI:", string(user.URI))
	log.Println("Endpoint:", user.Endpoint)
	log.Println("Followers:", user.Followers.Count)
}
