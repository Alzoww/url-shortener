package service

import "github.com/Alzoww/url-shortener/internal/storage"

type URLService struct {
	storage storage.Interface
}

func NewURLService(storage storage.Interface) *URLService {
	return &URLService{
		storage: storage,
	}
}
