package urlroute

type CreateShortURLRequest struct {
	OriginalURL string `json:"original_url"`
}

type CreateShortURLResponse struct {
	ShortURL string `json:"short_url"`
}
