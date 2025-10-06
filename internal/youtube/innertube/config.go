package innertube

import _ "embed"

const (
	REFERER_YOUTUBE = "https://www.youtube.com/"
	USER_AGENT_WEB  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.157 Safari/537.36"
)

//go:embed innertube_api_key.txt
var webApiKey string

var config = Config{
	Host:    "youtubei.googleapis.com",
	BaseURL: "https://youtubei.googleapis.com/youtubei/v1/",
	Clients: []ClientContext{
		{ClientID: 1, ClientName: "WEB", ClientVersion: "2.20251002.00.00", UserAgent: USER_AGENT_WEB, Referer: REFERER_YOUTUBE, APIKey: webApiKey},
	},
}
