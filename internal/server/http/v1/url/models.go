package urlroute

type CreateURLAliasRequest struct {
	OriginalURL string `json:"original_url"`
}

type CreateURLAliasResponse struct {
	Alias string `json:"alias"`
}

type GetOriginalByAliasResponse struct {
	OriginalURL string `json:"original_url"`
}
