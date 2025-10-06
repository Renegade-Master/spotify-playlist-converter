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
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

//go:embed google_client_secret.json
var googleAuthFile string

var ch = make(chan string)

func createYouTubeService() (*youtube.Service, *http.Client) {
	ctx := context.Background()

	if googleAuthFile == "" {
		log.Fatal("Google Client Secret is blank. This binary was compiled incorrectly.")
	}

	// Read credentials from file
	b := []byte(googleAuthFile)

	// Configure OAuth2 with required scopes
	config, err := google.ConfigFromJSON(b, youtube.YoutubeForceSslScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	// Create YouTube service
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Error creating YouTube service: %v", err)
	}

	return service, client
}

// getClient retrieves a token, saves the token, and returns the configured client.
func getClient(config *oauth2.Config) *http.Client {
	tok := getTokenFromWeb(config)
	return config.Client(context.Background(), tok)
}

// getTokenFromWeb requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	http.HandleFunc("/", completeAuth)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for favicon.ico")
	})

	go func() {
		log.Println("Starting HTTP server")
		err := http.ListenAndServe(":8001", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser:\n%v\n", authURL)

	if err := browser.OpenURL(authURL); err != nil {
		log.Printf("Could not open browser automatically. Please visit the URL above to log in: Error: [%v]", err)
	}

	authCode := <-ch

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request for:", r.URL.String())

	values := r.URL.Query()
	if e := values.Get("error"); e != "" {
		log.Fatalf(`return nil, errors.New("spotify: auth failed - " + e)`)
	}
	code := values.Get("code")
	if code == "" {
		log.Fatalf(`return nil, errors.New("spotify: didn't get access code")`)
	}

	r.Context().Done()

	ch <- code
}

// tokenFromFile retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	log.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
