package main

import (
	"context"
	"log"

	"github.com/Renegade-Master/spotify-playlist-converter/internal/spotify"
)

func main() {
	log.Printf("PlaylistConverter")

	ctx := context.Background()

	spotifyClient := spotify.GetSpotifyClient()

	spotifyPrivateUser := spotify.GetSpotifyPrivateUser(ctx, *spotifyClient)
	log.Printf("User ID: [%s]", spotifyPrivateUser.ID)
	log.Printf("Display name: [%s]", spotifyPrivateUser.DisplayName)
	log.Printf("Spotify URI: [%s]", string(spotifyPrivateUser.URI))
	log.Printf("Endpoint: [%s]", spotifyPrivateUser.Endpoint)

	playlists, _ := spotifyClient.GetPlaylistsForUser(ctx, spotifyPrivateUser.ID)

	log.Printf("Found [%d] playlists", len(playlists.Playlists))
	for idx, playlist := range playlists.Playlists {
		log.Printf("Playlist [%d]: [%s]", idx, playlist.Name)
	}
}
