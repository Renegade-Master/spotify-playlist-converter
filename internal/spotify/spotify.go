package spotify

import (
	"context"
	"log"

	"github.com/zmb3/spotify/v2"
)

type Spotify struct {
	client        *spotify.Client
	privateClient *spotify.PrivateUser
}

func NewSpotify() Spotify {
	spotifyClient := createSpotifyService()
	spotifyPrivateUser := getSpotifyPrivateUser(context.Background(), *spotifyClient)
	return Spotify{client: spotifyClient, privateClient: spotifyPrivateUser}
}

func createSpotifyService() *spotify.Client {
	return getSpotifyClient()
}

func (s Spotify) ListInfo() {
	log.Printf("User ID: [%s]", s.privateClient.ID)
	log.Printf("Display name: [%s]", s.privateClient.DisplayName)
	log.Printf("Spotify URI: [%s]", string(s.privateClient.URI))
	log.Printf("Endpoint: [%s]", s.privateClient.Endpoint)
}

func (s Spotify) ListPlaylists() {
	playlists, _ := s.client.GetPlaylistsForUser(context.Background(), s.privateClient.ID)

	log.Printf("Found [%d] playlists", len(playlists.Playlists))
	for idx, playlist := range playlists.Playlists {
		log.Printf("Playlist [%d]: [%s]", idx, playlist.Name)
	}
}
