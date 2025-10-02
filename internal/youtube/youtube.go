/*
 * Copyright (c) 2025.
 *
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

type ChannelResource string

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

func (yt YouTube) FindTrack(query string) {
	log.Printf("Searching for: [%s]\n", query)

	call := yt.client.Search.List([]string{"snippet"})
	call = call.Q(query)
	call = call.MaxResults(5)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error retrieving track: %s", err)
	}

	log.Println("Your YouTube Music Search Results:")
	log.Println("========================================")

	if len(response.Items) == 0 {
		log.Println("No tracks found.")
	} else {
		for i, track := range response.Items {
			log.Printf("%d. %s\n", i+1, track.Snippet.Title)
			log.Printf("   ID: %s\n", track.Id.VideoId)
			log.Printf("   Description: %s\n", track.Snippet.Description)
			log.Println()
		}
	}
}
