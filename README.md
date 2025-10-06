# Spotify Playlist Converter

[![Quality gate](https://sonarcloud.io/api/project_badges/quality_gate?project=Renegade-Master_spotify-playlist-converter)](https://sonarcloud.io/summary/new_code?id=Renegade-Master_spotify-playlist-converter)

## Building

### Setup

To build the application, there are some files which must be created:

```shell
touch internal/spotify/spotify_client_id.txt \
  internal/spotify/spotify_client_secret.txt \
  internal/youtube/google_client_secret.json
```

The `spotify_client_id.txt` and `spotify_client_secret.txt` should be filled with the Spotify Client ID and Secret from
the spotify-playlist-converter Spotify Application.

The `google_client_secret.json` should be filled with the Client Secret JSON file from the Google Desktop Application. 

### Build

To build the application, run the following command:

```shell
go build -o out/ -tags generate ./...
```

To build for other platforms or architectures, use the following command:

```shell
GOOS=android GOARCH=arm64 go build -o out/ -tags generate ./...
```

## Usage

Run the following command on the binary:

```shell
$ playlistConverter
```

A webpage will be opened in your browser, or you will be shown a URL to open manually, for both Spotify and YouTube.

## References

Spotify API:

* https://github.com/zmb3/spotify
* https://arkoes.medium.com/using-spotify-web-api-in-go-fa10373d5efb

YouTube QuickStart

* https://developers.google.com/youtube/v3/quickstart/go

Google API Samples:

* https://github.com/youtube/api-samples

YouTube Videos:

* https://www.youtube.com/watch?v=0aR9xvrRP2g
* https://www.youtube.com/watch?v=QY8dhl1EQfI
* https://www.youtube.com/watch?v=Q49gGXCCY_4&list=PL_cUvD4qzbkx_-4roA33KKz37BDX_JAZb&index=2

InnerTube API:

* https://github.com/wslyyy/youtube-go
