package innertube

import (
	"net/http"
)

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
	context := GetContext("WEB")

	return &InnerTube{
		Adaptor: NewInnerTubeAdaptor(context, &http.Client{}),
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

func (it *InnerTube) Search(query *string, params *string, continuation *string) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"query":        query,
		"params":       params,
		"continuation": continuation,
	}
	//log.Println("body: ", body)
	//log.Println("Filter(body): ", Filter(body))
	return it.Call("SEARCH", nil, Filter(body))
}
