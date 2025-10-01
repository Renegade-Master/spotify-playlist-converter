package main

import (
	"log"

	"github.com/Renegade-Master/spotify-playlist-converter/internal/spotify"
	"github.com/Renegade-Master/spotify-playlist-converter/internal/youtube"
)

func main() {
	log.Printf("PlaylistConverter")

	spotifyClient := spotify.NewSpotify()
	spotifyClient.ListInfo()
	spotifyClient.ListPlaylists()

	newTube := youtube.NewYouTube()
	newTube.ListPlaylists()
}
