package domain

type URL struct {
	FullURL      string
	ShortenedURL string
}

type URLShortenRepository interface {
	GetByFull(fullURL string) (URL, error)
	GetByShortened(shortenedURL string) (URL, error)
}
