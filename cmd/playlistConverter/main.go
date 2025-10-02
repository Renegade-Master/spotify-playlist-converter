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
