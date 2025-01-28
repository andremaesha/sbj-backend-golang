package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Curl[T any](method, url string, headers map[string]string, requestBody any, responseType T) (T, error) {
	// Encode the request body if it's not nil
	var requestBodyBytes []byte
	if requestBody != nil {
		var err error
		requestBodyBytes, err = json.Marshal(requestBody)
		if err != nil {
			return responseType, fmt.Errorf("failed to encode request body: %v", err)
		}
	}

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return responseType, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return responseType, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return responseType, fmt.Errorf("failed to read response body: %v", err)
	}

	// Decode the response body into the provided responseType
	err = json.Unmarshal(body, &responseType)
	if err != nil {
		return responseType, fmt.Errorf("failed to decode response body: %v", err)
	}

	return responseType, nil
}
