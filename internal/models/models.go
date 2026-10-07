// Package models describe entities in the app
package models

type Type string

const (
	Expense = "expense"
	Income  = "income"
)

type Category struct {
	ID   int64   `json:"ID"`
	Name string `json:"category_name"`
}

type Record struct {
	ID       int64    `json:"ID"`
	Type     Type     `json:"type"`
	Category Category `json:"category"`
	Amount   int64    `json:"amount"`
	Date     string   `json:"date"`
}
