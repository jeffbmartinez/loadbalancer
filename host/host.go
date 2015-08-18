package host

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jeffbmartinez/log"
)

type Host struct {
	Hostname string `json:"hostname"`
	Weight   int    `json:"weight"`
}

func (h Host) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	newURL, err := url.Parse(h.Hostname + request.URL.String())
	if err != nil {
		log.Printf("Couldn't parse url (%v): %v\n", newURL, err)
		return
	}

	newRequest := &http.Request{
		URL:    newURL,
		Header: request.Header,
	}

	intermediateResponse, err := http.DefaultClient.Do(newRequest)
	if err != nil {
		log.Errorf("Had a problem with a response: %v\n", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := ioutil.ReadAll(intermediateResponse.Body)
	if err != nil {
		log.Errorf("Had a problem reading response body: %v\n", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// There can be multiple values for a given header key. Here I am
	// clearing any values that may pre-exist in the header and replacing
	// them with the values from the response.
	for key, values := range intermediateResponse.Header {
		response.Header().Del(key)

		for _, value := range values {
			response.Header().Add(key, value)
		}
	}

	response.WriteHeader(intermediateResponse.StatusCode)
	response.Write(responseBody)
}
