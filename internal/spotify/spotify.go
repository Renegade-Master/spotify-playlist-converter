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
		log.Printf("%d. %s\n", idx+1, playlist.Name)
		log.Printf("   Description: %s\n", playlist.Description)
		log.Printf("   ID: %s\n", playlist.ID)
		log.Printf("   URI: %s\n", playlist.URI)
		log.Println()
	}
}

func (s Spotify) GetPlaylists() []spotify.SimplePlaylist {
	playlists, _ := s.client.GetPlaylistsForUser(context.Background(), s.privateClient.ID)

	return playlists.Playlists
}

func (s Spotify) GetPlaylist(playlistId spotify.ID) *spotify.FullPlaylist {
	playlist, _ := s.client.GetPlaylist(context.Background(), playlistId)

	return playlist
}

func (s Spotify) ListPlaylist(playlistId spotify.ID) {
	playlist, err := s.client.GetPlaylistItems(context.Background(), playlistId)
	if err != nil {
		log.Fatalf("Error retrieving playlist: [%s]", err)
	}

	log.Printf("Found [%d] Tracks", len(playlist.Items))
	for idx, track := range playlist.Items {
		log.Printf("%d. %s\n", idx+1, track.Track.Track.Name)
		log.Printf("   Album: %s\n", track.Track.Track.Album.Name)
		log.Printf("   Artists: %s\n", track.Track.Track.Artists)
		log.Printf("   Duration: %d\n", track.Track.Track.Duration)
		log.Printf("   ID: %s\n", track.Track.Track.ID)
		log.Printf("   URI: %s\n", track.Track.Track.URI)
		log.Println()
	}
}
