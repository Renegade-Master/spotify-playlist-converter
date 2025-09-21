package youtubeMusic

import (
	"log"

	"github.com/prettyirrelevant/ytmusicapi"
)

type YoutubeMusic struct{}

func NewYoutubeMusic() YoutubeMusic {
	return YoutubeMusic{}
}

func Playlists() {
	ytmusicapi.Setup()

	playlist, err := ytmusicapi.CreatePlaylist("test", "test", ytmusicapi.PRIVATE, "", []string{})
	if err != nil {
		log.Fatalf("Error creating playlist: [%s]", err)
	}

	log.Printf("Playlist ID: [%s]", playlist)
}
