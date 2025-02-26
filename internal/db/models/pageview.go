package models

type PageView struct {
	Id   string      `pg:"type:uuid" pg:",pk" json:"id"`
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}
