package youtube

import (
	"log"

	"google.golang.org/api/youtube/v3"
)

type YouTube struct {
	client *youtube.Service
}

func NewYouTube() YouTube {
	youtubeService := createYouTubeService()
	return YouTube{client: youtubeService}
}

func (yt YouTube) ListPlaylists() {
	// List user's playlists
	call := yt.client.Playlists.List([]string{"snippet", "contentDetails"})
	call = call.Mine(true)
	call = call.MaxResults(50)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error retrieving playlists: %v", err)
	}

	log.Println("Your YouTube Music Playlists:")
	log.Println("========================================")

	if len(response.Items) == 0 {
		log.Println("No playlists found.")
	} else {
		for i, playlist := range response.Items {
			log.Printf("%d. %s\n", i+1, playlist.Snippet.Title)
			log.Printf("   ID: %s\n", playlist.Id)
			log.Printf("   Description: %s\n", playlist.Snippet.Description)
			log.Printf("   Videos: %d\n", playlist.ContentDetails.ItemCount)
			log.Println()
		}
	}
}
