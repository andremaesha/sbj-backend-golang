package domain

type RequestFormDataCloudinary struct {
	File      string `json:"file"`
	ApiKey    string `json:"api_key"`
	Timestamp string `json:"timestamp"`
	Signature string `json:"signature"`
	PublicId  string `json:"public_id"`
	Folder    string `json:"folder"`
}

type ResponseCloudinary struct {
	AssetId   string                   `json:"asset_id"`
	PublicId  string                   `json:"public_id"`
	SecureUrl string                   `json:"secure_url,omitempty"`
	Error     *ResponseCloudinaryError `json:"error,omitempty"`
}

type ResponseCloudinaryError struct {
	Message string `json:"message"`
}
