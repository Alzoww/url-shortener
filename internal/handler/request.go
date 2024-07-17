package handler

import "errors"

type SaveURLRequest struct {
	URL   string `json:"url"`
	Alias string `json:"alias"`
}

func (r *SaveURLRequest) validate() error {
	if r.URL == "" {
		return errors.New("url is required")
	}

	if r.Alias == "" {
		return errors.New("alias is required")
	}

	return nil
}
