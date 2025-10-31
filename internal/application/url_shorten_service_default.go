package application

import (
	"github.com/Vla8islav/urlshortener/internal/domain"
)

type urlShortenServiceDefault struct {
	URLRepo domain.URLShortenRepository
}

func NewURLShortenService(repo domain.URLShortenRepository) domain.URLShortenService {
	return &urlShortenServiceDefault{URLRepo: repo}
}

func (u *urlShortenServiceDefault) GetByFull(longURL string) (domain.URL, error) {
	return u.URLRepo.GetByFull(longURL)
}

func (u *urlShortenServiceDefault) GetByShortened(shortURL string) (domain.URL, error) {
	return u.URLRepo.GetByShortened(shortURL)
}
