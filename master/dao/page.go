package dao

type Page struct {
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Data  interface{} `json:"data"`
}
