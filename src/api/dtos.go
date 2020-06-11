package api

type ShortURL struct {
	IPAddress string `json:"ip_address"`
	Views     int64  `json:"views"`
	Code      string `json:"code"`
	CreatedAt int64  `json:"created_at"`
	Original  string `json:"original"`
	URL       string `json:"url"`
}

type ShortenURLForm struct {
	URL string `json:"url"`
}
