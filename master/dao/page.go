package dao

type Page struct {
	Total int

	Page  int
	Limit int

	Data interface{}
}
