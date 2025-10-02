package main

import (
	"log"

	"github.com/Renegade-Master/spotify-playlist-converter/internal/spotify"
	"github.com/Renegade-Master/spotify-playlist-converter/internal/youtube"
)

func main() {
	spotifyClient := spotify.NewSpotify()
	//spotifyClient.ListInfo()
	//spotifyClient.ListPlaylists()
	//spotifyClient.ListPlaylist(spotifyClient.GetPlaylists()[0].ID)

	newTube := youtube.NewYouTube()
	//newTube.ListChannels()
	//newTube.ListPlaylists()

	spotifyPlaylistId := spotifyClient.GetPlaylists()[0].ID
	spotifyPlaylist := spotifyClient.GetPlaylist(spotifyPlaylistId)
	spotifyTrack := spotifyPlaylist.Tracks.Tracks[0].Track

	log.Printf("Searching YouTube for Track [%s] from Playlist [%s]\n", spotifyTrack.Name, spotifyPlaylist.Name)
	newTube.FindTrack(spotifyTrack.Name)
}
