package main

import (
	"github.com/Renegade-Master/spotify-playlist-converter/internal/spotify"
	"github.com/Renegade-Master/spotify-playlist-converter/internal/youtube"
)

func main() {
	spotifyClient := spotify.NewSpotify()
	spotifyClient.ListInfo()
	spotifyClient.ListPlaylists()

	newTube := youtube.NewYouTube()
	newTube.ListChannels()
	newTube.ListPlaylists()
}
