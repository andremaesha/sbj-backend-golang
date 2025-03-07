package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Curl Mengirim request JSON atau file (multipart/form-data)
func Curl[T any](method, url string, headers map[string]string, requestBody any, files map[string]string, responseType T) (T, error) {
	var requestBodyBytes *bytes.Buffer
	var contentType string

	// Jika ada file, gunakan multipart/form-data
	if len(files) > 0 {
		requestBodyBytes = &bytes.Buffer{}
		writer := multipart.NewWriter(requestBodyBytes)

		file, err := os.Open(files["file"])
		if err != nil {
			return responseType, fmt.Errorf("failed to open file: %v", err)
		}

		part, err := writer.CreateFormFile("file", filepath.Base(files["file"]))
		if err != nil {
			err = file.Close()
			if err != nil {
				panic(err)
			}
			return responseType, fmt.Errorf("failed to create form file: %v", err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return responseType, fmt.Errorf("failed to copy file data: %v", err)
		}
		err = file.Close()
		if err != nil {
			panic(err)
		}

		_ = writer.WriteField("api_key", files["api_key"])
		_ = writer.WriteField("timestamp", files["timestamp"])
		_ = writer.WriteField("signature", files["signature"])
		_ = writer.WriteField("public_id", files["public_id"])
		_ = writer.WriteField("folder", files["folder"])
		contentType = writer.FormDataContentType()
		err = writer.Close()
		if err != nil {
			panic(err)
		}
	} else {
		// Jika tanpa file, gunakan JSON
		requestBodyBytes = &bytes.Buffer{}
		if requestBody != nil {
			jsonData, err := json.Marshal(requestBody)
			if err != nil {
				return responseType, fmt.Errorf("failed to encode request body: %v", err)
			}
			requestBodyBytes.Write(jsonData)
		}
		contentType = "application/json"
	}

	// Buat HTTP request
	req, err := http.NewRequest(method, url, requestBodyBytes)
	if err != nil {
		return responseType, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", contentType)

	// Kirim request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return responseType, fmt.Errorf("failed to make request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	// Baca response body
	body, err := io.ReadAll(resp.Body)
	println(string(body))
	if err != nil {
		return responseType, fmt.Errorf("failed to read response body: %v", err)
	}

	// Decode response body ke responseType
	err = json.Unmarshal(body, &responseType)
	if err != nil {
		return responseType, fmt.Errorf("failed to decode response body: %v", err)
	}

	return responseType, nil
}
