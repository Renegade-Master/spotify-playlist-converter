package main

import (
	"context"
	"log"
	"os"

	"github.com/Renegade-Master/spotify-playlist-converter/internal/youtubeMusic"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	log.Printf("PlaylistConverter")

	ctx := context.Background()

	//spotifyClient := spotify.GetSpotifyClient()

	//spotifyPrivateUser := spotify.GetSpotifyPrivateUser(ctx, *spotifyClient)
	//log.Printf("User ID: [%s]", spotifyPrivateUser.ID)
	//log.Printf("Display name: [%s]", spotifyPrivateUser.DisplayName)
	//log.Printf("Spotify URI: [%s]", string(spotifyPrivateUser.URI))
	//log.Printf("Endpoint: [%s]", spotifyPrivateUser.Endpoint)

	//playlists, _ := spotifyClient.GetPlaylistsForUser(ctx, spotifyPrivateUser.ID)

	//log.Printf("Found [%d] playlists", len(playlists.Playlists))
	//for idx, playlist := range playlists.Playlists {
	//	log.Printf("Playlist [%d]: [%s]", idx, playlist.Name)
	//}

	//youtubeClient, err := youtube.NewService(ctx, option.WithAPIKey(""))
	//if err != nil {
	//	log.Fatalf("Unable to retrieve YouTube client: %s", err)
	//}

	//youtubeClient.PlaylistItems.List([]string{})

	ytClientId := os.Getenv("YOUTUBE_CLIENT_ID")
	ytClientSecret := os.Getenv("YOUTUBE_CLIENT_SECRET")

	youtubeClient := youtubeMusic.GetClient(ctx, ytClientId, ytClientSecret, youtube.YoutubeReadonlyScope)
	youtubeService, err := youtube.NewService(ctx, option.WithHTTPClient(youtubeClient))
	if err != nil {
		log.Fatalf("Unable to retrieve YouTube client: [%s]", err)
	}

	request := youtubeService.Playlists.List([]string{"snippet,contentDetails"})
	response, err := request.Do()
	if err != nil {
		log.Fatalf("Unable to retrieve playlist: [%s]", err)
	}
	for _, playlist := range response.Items {
		log.Printf("Playlist ID: [%s]", playlist.Id)
	}

	youtubeMusic.Playlists()
}
