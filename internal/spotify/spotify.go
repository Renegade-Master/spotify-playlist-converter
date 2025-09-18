package spotify

import (
	"net/http"
	"os"

	"github.com/pkg/browser"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func auth1(state string) {
	redirectUri := os.Getenv("SPOTIFY_REDIRECT_URI")

	// the redirect URL must be an exact match of a URL you've registered for your application
	// scopes determine which permissions the user is prompted to authorize
	auth := spotifyauth.New(spotifyauth.WithRedirectURL(redirectUri), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))

	// get the user to this URL - how you do that is up to you
	// you should specify a unique state string to identify the session
	url := auth.AuthURL(state)
	browser.OpenURL(url)
}

// the user will eventually be redirected back to your redirect URL
// typically you'll have a handler set up like the following:
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// use the same state string here that you used to generate the URL
	token, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	// create a client using the specified token
	client := spotify.New(auth.Client(r.Context(), token))

	// the client can now be used to make authenticated requests
}
