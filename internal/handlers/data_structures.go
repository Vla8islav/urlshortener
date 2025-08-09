package handlers

type ShortenRequestPayload struct {
	FullURL string `json:"url" validate:"required,url"`
}

type ShortenResponsePayload struct {
	ShortenedURL string `json:"result" validate:"required,url"`
}
