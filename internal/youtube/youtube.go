/*
 *    Copyright (c) 2025 [renegade@renegade-master.com]
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

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

func (yt YouTube) ListChannels() {
	// List user's playlists
	call := yt.client.Channels.List([]string{"snippet", "contentDetails"})
	call = call.Mine(true)
	call = call.MaxResults(50)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error retrieving channels: %v", err)
	}

	log.Println("Your YouTube Music Channels:")
	log.Println("========================================")

	if len(response.Items) == 0 {
		log.Println("No channels found.")
	} else {
		for i, channel := range response.Items {
			log.Printf("%d. %s\n", i+1, channel.Snippet.Title)
			log.Printf("   ID: %s\n", channel.Id)
			log.Printf("   Title: %s\n", channel.Snippet.Title)
			log.Printf("   Description: %s\n", channel.Snippet.Description)
			log.Printf("   Custom URL: %s\n", channel.Snippet.CustomUrl)
			log.Println()
		}
	}
}

func (yt YouTube) ListPlaylists() {
	playlists := yt.GetPlaylists()

	log.Println("Your YouTube Music Playlists:")
	log.Println("========================================")

	if len(playlists) == 0 {
		log.Println("No playlists found.")
	} else {
		for i, playlist := range playlists {
			log.Printf("%d. %s\n", i+1, playlist.Snippet.Title)
			log.Printf("   ID: %s\n", playlist.Id)
			log.Printf("   Description: %s\n", playlist.Snippet.Description)
			log.Printf("   Videos: %d\n", playlist.ContentDetails.ItemCount)
			log.Println()
		}
	}
}

func (yt YouTube) GetPlaylists() []*youtube.Playlist {
	call := yt.client.Playlists.List([]string{"snippet", "contentDetails"})
	call = call.Mine(true)
	call = call.MaxResults(50)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error retrieving Playlists: %v", err)
	}

	if len(response.Items) == 0 {
		log.Println("No playlists found.")
		return []*youtube.Playlist{}
	} else {
		return response.Items
	}
}

func (yt YouTube) GetPlaylist(playlistId string) *youtube.Playlist {
	call := yt.client.Playlists.List([]string{"snippet", "contentDetails"})
	call = call.Mine(true)
	call = call.MaxResults(50)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error retrieving Playlists: %v", err)
	}

	if len(response.Items) == 0 {
		log.Println("No playlists found.")
	} else {
		for _, playlist := range response.Items {
			if playlist.Id == playlistId {
				return playlist
			}
		}
	}

	return nil
}

func (yt YouTube) GetPlaylistItems(playlistId string) []*youtube.PlaylistItem {
	var playlistItems []*youtube.PlaylistItem
	var nextPageToken string

	for {
		call := yt.client.PlaylistItems.List([]string{"snippet", "contentDetails"}).
			PlaylistId(playlistId).
			MaxResults(50).
			PageToken(nextPageToken)

		response, err := call.Do()
		if err != nil {
			log.Fatalf("Error retrieving Tracks: [%v]", err)
		}

		if len(response.Items) == 0 {
			log.Println("No Tracks found.")
		} else {
			playlistItems = append(playlistItems, response.Items...)
		}

		if response.NextPageToken == "" {
			break
		}

		nextPageToken = response.NextPageToken
	}

	return playlistItems
}

func (yt YouTube) ListTracks(query string, maxResults int64) {
	track := yt.GetTracks(query, maxResults)

	log.Println("Your YouTube Music Search Results:")
	log.Println("========================================")

	if len(track) == 0 {
		log.Println("No tracks found.")
	} else {
		for i, track := range track {
			log.Printf("%d. %s\n", i+1, track.Snippet.Title)
			log.Printf("   ID: %s\n", track.Id.VideoId)
			log.Printf("   Description: %s\n", track.Snippet.Description)
			log.Println()
		}
	}
}

func (yt YouTube) GetTracks(query string, maxResults int64) []*youtube.SearchResult {
	log.Printf("Searching for: [%s]\n", query)

	call := yt.client.Search.List([]string{"snippet"})
	call = call.Q(query)
	call = call.MaxResults(maxResults)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error retrieving track: %s", err)
	}

	if len(response.Items) == 0 {
		log.Println("No tracks found.")
		return []*youtube.SearchResult{}
	} else {
		return response.Items
	}
}

func (yt YouTube) CreatePlaylist(name string) string {
	playlists := yt.GetPlaylists()
	for _, playlist := range playlists {
		if playlist.Snippet.Title == name {
			log.Printf("Playlist [%s] already exists\n", name)
			return playlist.Id
		}
	}

	playlist := &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{
			Title:       name,
			Description: "Playlist created by Spotify Playlist Converter",
			Tags:        []string{"spotify-playlist-converter"},
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: "private",
		},
	}

	call := yt.client.Playlists.Insert([]string{"snippet", "status"}, playlist)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error creating Playlist: %v", err)
	}

	log.Printf("Created Playlist: [%s], ID: [%s]\n", response.Snippet.Title, response.Id)
	return response.Id
}

func (yt YouTube) AddToPlaylist(playlistId string, trackId string) string {
	playlistItems := yt.GetPlaylistItems(playlistId)

	for _, playlistItem := range playlistItems {
		if playlistItem.Snippet.ResourceId.VideoId == trackId {
			log.Printf("Track [%s] already exists in Playlist [%s]\n", trackId, playlistId)
			return playlistItem.Id
		}
	}

	call := yt.client.PlaylistItems.Insert([]string{"snippet"}, &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: playlistId,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: trackId,
			},
		},
	})

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error adding track to playlist: %v", err)
	}

	return response.Id
}
