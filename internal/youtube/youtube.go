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
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Renegade-Master/spotify-playlist-converter/internal/util"
	"google.golang.org/api/youtube/v3"
)

type YouTube struct {
	client    *youtube.Service
	rawClient *http.Client
	Credits   int
}

func NewYouTube() *YouTube {
	youtubeService, youtubeClient := createYouTubeService()
	return &YouTube{client: youtubeService, rawClient: youtubeClient}
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

	yt.addAllIdsToPlaylist(playlistId, trackIds...)
	return nil
}

func (yt *YouTube) addAllIdsToPlaylist(playlistId string, trackIds ...string) error {
	// Check if we have SAPISID cookie - if not, fall back to official API
	var sapisid string
	if yt.rawClient != nil && yt.rawClient.Jar != nil {
		youtubeURL, _ := url.Parse("https://www.youtube.com")
		cookies := yt.rawClient.Jar.Cookies(youtubeURL)

		for _, cookie := range cookies {
			if cookie.Name == "SAPISID" {
				sapisid = cookie.Value
				break
			}
		}
	}

	// If no SAPISID cookie, use official API
	if sapisid == "" {
		log.Println("No SAPISID cookie found, using official API")
		return yt.addIdsToPlaylistOfficial(playlistId, trackIds...)
	}

	// Try unofficial API
	err := yt.addIdsToPlaylistUnofficial(playlistId, trackIds...)
	if err != nil {
		log.Printf("Unofficial API failed: %v, falling back to official API", err)
		return yt.addIdsToPlaylistOfficial(playlistId, trackIds...)
	}

	return nil
}

func (yt *YouTube) addIdsToPlaylistOfficial(playlistId string, trackIds ...string) error {
	log.Printf("Adding %d tracks to playlist %s using official API", len(trackIds), playlistId)

	successCount := 0
	for i, videoId := range trackIds {
		playlistItem := &youtube.PlaylistItem{
			Snippet: &youtube.PlaylistItemSnippet{
				PlaylistId: playlistId,
				ResourceId: &youtube.ResourceId{
					Kind:    "youtube#video",
					VideoId: videoId,
				},
			},
		}

		call := yt.client.PlaylistItems.Insert([]string{"snippet"}, playlistItem)
		_, err := call.Do()
		if err != nil {
			log.Printf("Error adding video %s to playlist: %v", videoId, err)
			continue
		}

		yt.Credits += 50
		successCount++

		if (i+1)%10 == 0 {
			log.Printf("Progress: %d/%d videos added", i+1, len(trackIds))
		}
	}

	log.Printf("Successfully added %d/%d tracks to playlist", successCount, len(trackIds))

	if successCount == 0 {
		return fmt.Errorf("failed to add any tracks to playlist")
	}

	return nil
}

func (yt *YouTube) addIdsToPlaylistUnofficial(playlistId string, trackIds ...string) error {
	apiKey := "AAA"
	apiURL := fmt.Sprintf("https://www.youtube.com/youtubei/v1/browse/edit_playlist?key=%s", apiKey)

	var actions []interface{}
	for _, id := range trackIds {
		actions = append(actions, map[string]interface{}{
			"addedVideoId": id,
			"action":       "ACTION_ADD_VIDEO",
		})
	}

	// Get SAPISID cookie for authentication
	var sapisid string
	if yt.rawClient != nil && yt.rawClient.Jar != nil {
		youtubeURL, _ := url.Parse("https://www.youtube.com")
		cookies := yt.rawClient.Jar.Cookies(youtubeURL)

		for _, cookie := range cookies {
			if cookie.Name == "SAPISID" {
				sapisid = cookie.Value
				break
			}
		}
	}

	if sapisid == "" {
		return fmt.Errorf("SAPISID cookie not found")
	}

	// Generate SAPISIDHASH
	timestamp := time.Now().Unix()
	origin := "https://www.youtube.com"
	hashInput := fmt.Sprintf("%d %s %s", timestamp, sapisid, origin)
	hash := sha1.Sum([]byte(hashInput))
	sapisidHash := fmt.Sprintf("SAPISIDHASH %d_%x", timestamp, hash)

	bodyMap := map[string]interface{}{
		"context": map[string]interface{}{
			"client": map[string]interface{}{
				"clientName":    "WEB",
				"clientVersion": "2.20250102.01.00",
			},
		},
		"actions":    actions,
		"playlistId": playlistId,
	}

	requestBody, err := json.Marshal(bodyMap)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://www.youtube.com")
	req.Header.Set("Authorization", sapisidHash)

	// Add all cookies
	if yt.rawClient != nil && yt.rawClient.Jar != nil {
		for _, cookie := range yt.rawClient.Jar.Cookies(req.URL) {
			req.AddCookie(cookie)
		}
	}

	// Execute request
	resp, err := yt.rawClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, _ := io.ReadAll(resp.Body)

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unofficial API Response Status: %d", resp.StatusCode)
		log.Printf("Unofficial API Response Body: %s", string(body))
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	log.Printf("Successfully added %d tracks to playlist %s using unofficial API", len(trackIds), playlistId)
	return nil
}
