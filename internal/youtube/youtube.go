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
	"net/http"

	"github.com/Renegade-Master/spotify-playlist-converter/internal/util"
	"github.com/Renegade-Master/spotify-playlist-converter/internal/youtube/innertube"
	"google.golang.org/api/youtube/v3"
)

type YouTube struct {
	client    *youtube.Service
	intClient *innertube.InnerTube
	Credits   int
}

func NewYouTube() *YouTube {
	youtubeService := createYouTubeService()

	httpClient := &http.Client{}
	innerTubeService, _ := innertube.NewInnerTube(httpClient, "WEB", "2.20230728.00.00", "", "", "", nil, true)

	return &YouTube{
		client:    youtubeService,
		intClient: innerTubeService,
	}
}

func (yt *YouTube) ListChannels() {
	call := yt.client.Channels.List([]string{"snippet", "contentDetails"}).
		Mine(true).
		MaxResults(50)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error retrieving channels: %v", err)
	}
	yt.Credits += 1

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

func (yt *YouTube) ListPlaylists() {
	playlists := yt.GetPlaylists()

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

func (yt *YouTube) GetPlaylists() []*youtube.Playlist {
	var playlists []*youtube.Playlist
	var nextPageToken string

	for {
		call := yt.client.Playlists.List([]string{"snippet", "contentDetails"}).
			Mine(true).
			MaxResults(50).
			PageToken(nextPageToken)

		response, err := call.Do()
		if err != nil {
			log.Fatalf("Error retrieving Playlists: %v", err)
		}
		yt.Credits += 1

		if len(response.Items) == 0 {
			log.Println("No Playlists found.")
		} else {
			playlists = append(playlists, response.Items...)
		}

		if response.NextPageToken == "" {
			break
		}

		nextPageToken = response.NextPageToken
	}

	return playlists
}

func (yt *YouTube) GetPlaylist(playlistId string) *youtube.Playlist {
	playlists := yt.GetPlaylists()

	if len(playlists) == 0 {
		log.Println("No Playlists found.")
	} else {
		for _, playlist := range playlists {
			if playlist.Id == playlistId {
				return playlist
			}
		}
	}

	return nil
}

func (yt *YouTube) GetPlaylistItems(playlistId string) []*youtube.PlaylistItem {
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
		yt.Credits += 1

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

// GetTrackUnofficial is a method of Searching YouTube without using Credits
func (yt *YouTube) GetTrackUnofficial(query string, maxResults int64) string {
	paramsTypeVideo := "EgIQAQ%3D%3D"

	data, err := yt.intClient.Search(&query, &paramsTypeVideo, nil)
	if err != nil {
		log.Fatalf("Error retrieving track: %s", err)
	}

	contents := data["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"].(map[string]interface{})["primaryContents"].(map[string]interface{})["sectionListRenderer"].(map[string]interface{})["contents"]
	resultsContent := contents.([]interface{})[0].(map[string]interface{})["itemSectionRenderer"].(map[string]interface{})["contents"]
	results := resultsContent.([]interface{})

	// Iterate through the first maxResults results
	var weightedTracks []WeightedSimpleSearchResult
	for i := 0; i < int(maxResults); i++ {
		result := results[i].(map[string]interface{})["videoRenderer"].(map[string]interface{})
		title := result["title"].(map[string]interface{})["runs"].([]interface{})[0].(map[string]interface{})["text"].(string)
		videoId := result["videoId"].(string)

		log.Printf("Found Track: [%s] [%s]\n", title, videoId)

		distance := util.LevenshteinDistance(query, title)
		weightedTracks = append(weightedTracks, WeightedSimpleSearchResult{Result: title, Weight: distance, Id: videoId})
	}

	distance := func(track1, track2 *WeightedSimpleSearchResult) bool {
		return track1.Weight < track2.Weight
	}

	BySimpleSearchResult(distance).SortSimpleSearchResult(weightedTracks)

	// Return the top (i.e. most similar) Result
	return weightedTracks[0].Id
}

func (yt *YouTube) GetTrack(query string, maxResults int64) *youtube.SearchResult {
	log.Printf("Searching for: [%s]\n", query)

	call := yt.client.Search.List([]string{"snippet"}).
		Q(query).
		MaxResults(maxResults)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error retrieving track: %s", err)
	}
	yt.Credits += 100

	if len(response.Items) == 0 {
		log.Println("No tracks found.")
		return &youtube.SearchResult{}
	} else {
		var weightedTracks []WeightedSearchResult
		for _, track := range response.Items {
			log.Printf("Found Track: [%s]\n", track.Snippet.Title)
			youTubeTitle := track.Snippet.Title

			distance := util.LevenshteinDistance(query, youTubeTitle)
			weightedTracks = append(weightedTracks, WeightedSearchResult{Result: track, Weight: distance})
		}

		distance := func(track1, track2 *WeightedSearchResult) bool {
			return track1.Weight < track2.Weight
		}

		BySearchResult(distance).SortSearchResult(weightedTracks)

		// Return the top (i.e. most similar) Result
		return weightedTracks[0].Result
	}
}

// CreatePlaylist will create a YouTube Playlist if it does not already exist.
// Returns the Playlist ID of the new Playlist, or the existing Playlist by the
// same name, as well as a Boolean to indicate if this is a new Playlist.
func (yt *YouTube) CreatePlaylist(name string) (string, bool) {
	playlists := yt.GetPlaylists()
	for _, playlist := range playlists {
		if playlist.Snippet.Title == name {
			log.Printf("Playlist [%s] already exists\n", name)
			return playlist.Id, false
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
	yt.Credits += 50

	log.Printf("Created Playlist: [%s], ID: [%s]\n", response.Snippet.Title, response.Id)
	return response.Id, true
}

func (yt *YouTube) AddToPlaylist(playlistId string, trackIds ...string) error {
	playlistItems := yt.GetPlaylistItems(playlistId)

	// Check if any Tracks already exist in the Playlist
	var badTrackIds []int
	for _, playlistItem := range playlistItems {
		for idx, trackId := range trackIds {
			if playlistItem.Snippet.ResourceId.VideoId == trackId {
				log.Printf("Track [%s] already exists in Playlist [%s]\n", trackId, playlistId)
				badTrackIds = append(badTrackIds, idx)
			}
		}
	}

	// Remove Tracks that already exist
	if len(badTrackIds) == len(trackIds) {
		log.Printf("All Tracks are already present in the Playlist")
		return nil
	}

	itemsRemoved := 0
	log.Printf("Removing [%d] Tracks that already exist in Playlist [%s]\n", len(badTrackIds), playlistId)
	for _, idx := range badTrackIds {
		trackIds = util.RemoveIndexString(trackIds, idx-itemsRemoved)
		itemsRemoved++
	}

	// Add all found Tracks to the Playlist
	for _, trackId := range trackIds {
		// ToDo: It would be nice if it was possible to add all Tracks in one call. May be possible using raw HTTP Requests instead of the library
		call := yt.client.PlaylistItems.Insert([]string{"snippet"}, &youtube.PlaylistItem{
			Snippet: &youtube.PlaylistItemSnippet{
				PlaylistId: playlistId,
				ResourceId: &youtube.ResourceId{
					Kind:    "youtube#video",
					VideoId: trackId,
				},
			},
		})

		_, err := call.Do()
		if err != nil {
			log.Printf("Error adding Track ID [%s] to Playlist [%s]: [%v]", trackId, playlistId, err)
			return err
		}
		yt.Credits += 50

		log.Printf("Added Track ID [%s] to Playlist [%s]\n", trackId, playlistId)
	}

	return nil
}
