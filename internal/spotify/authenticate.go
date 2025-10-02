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
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/browser"
	"github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"
)

const redirectURI = "http://127.0.0.1:8000/callback"

//go:embed spotify_client_id.txt
var clientId string

//go:embed spotify_client_secret.txt
var clientSecret string

var (
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
		),
		spotifyauth.WithClientID(strings.TrimSuffix(clientId, "\n")),
		spotifyauth.WithClientSecret(strings.TrimSuffix(clientSecret, "\n")))
	ch    = make(chan *spotify.Client)
	state = uuid.New().String()
)

func getSpotifyClient() *spotify.Client {
	if clientId == "" || clientSecret == "" {
		log.Fatal("SPOTIFY_ID or SPOTIFY_SECRET are blank. This binary was compiled incorrectly.")
	}

	http.HandleFunc("/callback", completeAuth)

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

func getSpotifyPrivateUser(ctx context.Context, client spotify.Client) *spotify.PrivateUser {
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

	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
