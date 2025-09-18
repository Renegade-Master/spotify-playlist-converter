package spotify

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/browser"
	"github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://127.0.0.1:8000/callback"

var (
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopePlaylistReadPrivate,
		spotifyauth.ScopePlaylistReadCollaborative,
	))
	ch    = make(chan *spotify.Client)
	state = uuid.New().String()
)

func GetSpotifyClient() *spotify.Client {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go func() {
		log.Println("Starting HTTP server")
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	log.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	if err := browser.OpenURL(url); err != nil {
		log.Printf("Could not open browser automatically. Please visit the URL above to log in: Error: [%v]", err)
	}

	// wait for auth to complete
	client := <-ch

	return client
}

func GetSpotifyPrivateUser(ctx context.Context, client spotify.Client) *spotify.PrivateUser {
	user, err := client.CurrentUser(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return user
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request for:", r.URL.String())

	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
