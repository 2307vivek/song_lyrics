package types

type EsResponse struct {
	Hits Hits `json:"hits"`
}

type Hits struct {
	Total Total `json:"total"`
}

type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}
