# Spotify Playlist Converter

## Usage

Run the following command on the binary:

```shell
$ playlistConverter
```

A webpage will be opened in your browser, or you will be shown a URL to open manually.

## Notes

YTMusicAPI is not working as smoothly as I would like.

When retrieving your Authentication Headers, you must apply the following changes:
1. Convert to JSON
2. Change 'Cookies' to 'cookies'

## References

Spotify API:
* https://github.com/zmb3/spotify
* https://arkoes.medium.com/using-spotify-web-api-in-go-fa10373d5efb

YouTube Music API:
* https://pkg.go.dev/github.com/prettyirrelevant/ytmusicapi
* https://ytmusicapi.readthedocs.io/en/stable/setup/browser.html#copy-authentication-headers
* https://github.com/sigma67/ytmusicapi

Google API Samples:
* https://github.com/youtube/api-samples
