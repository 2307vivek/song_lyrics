package types

type Song struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Artist Artist `json:"aritst"`
}