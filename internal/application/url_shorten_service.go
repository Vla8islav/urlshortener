package application

import (
	"github.com/Vla8islav/urlshortener/internal/domain"
)

type URLShortenService struct {
	URLRepo domain.UrlShortenRepository
}

func NewURLShortenService(repo domain.UrlShortenRepository) *URLShortenService {
	return &URLShortenService{URLRepo: repo}
}

func (u *URLShortenService) GetShortURL(longURL string) (domain.URL, error) {
	return u.URLRepo.GetByFull(longURL)
}

func (u *URLShortenService) GetLongURL(shortURL string) (domain.URL, error) {
	return u.URLRepo.GetByShortened(shortURL)
}
