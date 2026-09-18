package models

var Request struct {
	URL string `json:"url"` // tellsw the decoder to look for a field called url in the request body and store it in the URL field of the request struct
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}
