package innertube

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type InnerTubeAdaptor struct {
	context ClientContext
	session *http.Client
}

func NewInnerTubeAdaptor(context ClientContext, session *http.Client) *InnerTubeAdaptor {
	if session == nil {
		session = &http.Client{}
	}
	return &InnerTubeAdaptor{
		context: context,
		session: session,
	}
}

func (ita *InnerTubeAdaptor) buildRequest(endpoint string, params map[string]string, body map[string]interface{}) (*http.Request, error) {
	body = Contextualise(ita.context, body)
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", config.BaseURL+strings.ToLower(endpoint), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header = ita.context.Headers()

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (ita *InnerTubeAdaptor) request(endpoint string, params map[string]string, body map[string]interface{}) (*http.Response, error) {
	req, err := ita.buildRequest(endpoint, params, body)

	log.Println("Method: ", req.Method)
	log.Println("URL: ", req.URL)
	log.Println("Request Headers: ")
	for k, v := range req.Header {
		log.Println(fmt.Sprintf("%s: ", k), v)
	}
	log.Println("Header: ", req.Header)
	log.Println("Cookies: ", req.Cookies())
	log.Println("UserAgent: ", req.UserAgent())
	log.Println("Form: ", req.Form)
	log.Println("PostForm", req.PostForm)
	log.Println("Body: ", req.Body)

	if err != nil {
		return nil, err
	}
	return ita.session.Do(req)
}

func (ita *InnerTubeAdaptor) Dispatch(endpoint string, params map[string]string, body map[string]interface{}) (map[string]interface{}, error) {
	resp, err := ita.request(endpoint, params, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with status code [%d]. Status [%s]", resp.StatusCode, resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "" && !isJSONContentType(contentType) {
		return nil, fmt.Errorf("expected JSON response, got %q", contentType)
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var reader io.Reader
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzr, err := gzip.NewReader(bytes.NewReader(bodyResp))
		if err != nil {
			return nil, err
		}
		defer gzr.Close()
		reader = gzr
	} else {
		reader = bytes.NewReader(bodyResp)
	}

	var responseData map[string]interface{}
	if err := json.NewDecoder(reader).Decode(&responseData); err != nil {
		return nil, err
	}

	// responseData["responseContext"]
	if responseContext, ok := responseData["responseContext"].(map[string]interface{}); ok {
		if visitorData, ok := responseContext["visitorData"].(string); ok {
			ita.context.XGoogVisitorId = visitorData
		}
	}

	if errorData, ok := responseData["error"]; ok {
		return nil, fmt.Errorf("request error: %v", errorData)
	}

	return responseData, nil
}

func isJSONContentType(contentType string) bool {
	return contentType == "application/json" || contentType == "application/json; charset=utf-8" || contentType == "application/json; charset=UTF-8"
}
