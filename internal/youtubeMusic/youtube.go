package youtubeMusic

import "github.com/prettyirrelevant/ytmusicapi"

func Playlists() {
	ytmusicapi.GetPlaylist("", 1)
}
