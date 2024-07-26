package handler

type SaveURLResponse struct {
	Status string `json:"status"`
}

type GetURLResponse struct {
	URL string `json:"url"`
}
