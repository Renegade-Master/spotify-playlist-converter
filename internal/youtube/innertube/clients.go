/*
 *    Copyright (c) 2024 wslyyy
 *
 *    Permission is hereby granted, free of charge, to any person obtaining a copy
 *    of this software and associated documentation files (the "Software"), to deal
 *    in the Software without restriction, including without limitation the rights
 *    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *    copies of the Software, and to permit persons to whom the Software is
 *    furnished to do so, subject to the following conditions:
 *
 *    The above copyright notice and this permission notice shall be included in all
 *    copies or substantial portions of the Software.
 *
 *    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *    SOFTWARE.
 */

package innertube

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

var ch = make(chan string)

// InnerTube struct
type InnerTube struct {
	Adaptor Adaptor
}

// Adaptor interface
type Adaptor interface {
	Dispatch(endpoint string, params map[string]string, body map[string]interface{}) (map[string]interface{}, error)
}

// NewInnerTube creates a new InnerTube instance
func NewInnerTube() (*InnerTube, error) {
	log.Println("Start Cookie Extract Attempt")
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://music.youtube.com/account/login/", nil)
	if err != nil {
		log.Fatalf("Error creating request: [%v]", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error retrieving response: [%v]", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Error retrieving response: [%v] %s", resp.StatusCode, string(body))
	}
	log.Println("Cookie Extract Attempt Complete")

	ytContext := GetContext("WEB")
	return &InnerTube{
		Adaptor: NewInnerTubeAdaptor(ytContext, &http.Client{Jar: jar}),
	}, nil
}

// Call method to make requests
func (it *InnerTube) Call(endpoint string, params map[string]string, body map[string]interface{}) (map[string]interface{}, error) {
	response, err := it.Adaptor.Dispatch(endpoint, params, body)
	if err != nil {
		return nil, err
	}

	delete(response, "responseContext")
	return response, nil
}

func (it *InnerTube) Search(query *string, continuation *string) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"query":        query,
		"params":       "EgIQAQ%3D%3D",
		"continuation": continuation,
	}

	//log.Println("body: ", body)
	//log.Println("Filter(body): ", Filter(body))
	return it.Call("SEARCH", nil, Filter(body))
}

func (it *InnerTube) AddToPlaylist(playlistId *string, trackIds ...string) (map[string]interface{}, error) {
	var actions []map[string]interface{}

	for _, trackId := range trackIds {
		actions = append(actions, map[string]interface{}{
			"addedVideoId": trackId,
			"action":       "ACTION_ADD_VIDEO",
		})
	}

	body := map[string]interface{}{
		"actions":      actions,
		"playlistId":   playlistId,
		"params":       "EgIQAw%3D%3D",
		"continuation": nil,
	}

	log.Println("AddToPlaylist body: ", body)
	log.Println("Filter(body): ", Filter(body))

	return it.Call("BROWSE/EDIT_PLAYLIST", nil, Filter(body))
}

// getTokenFromWeb requests a token from the web, then returns the retrieved token.
func getCookiesFromWeb(config *oauth2.Config) *oauth2.Token {
	http.HandleFunc("/", completeAuth)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for favicon.ico")
	})

	go func() {
		log.Println("Starting HTTP server")
		err := http.ListenAndServe(":8002", nil)
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
}
