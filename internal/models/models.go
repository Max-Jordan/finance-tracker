// Package models describe entities in the app
package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Type string

const (
	Expense = "expense"
	Income  = "income"
)

type Category struct {
	ID   int64  `json:"ID"`
	Name string `json:"category_name"`
}

type Record struct {
	ID       int64
	Type     Type
	Category Category
	Amount   int64
	Date     string
}

type RecordDTO struct {
	ID       int64    `json:"ID"`
	Type     Type     `json:"type"`
	Category Category `json:"category"`
	Amount   string   `json:"amount"`
	Date     string   `json:"date"`
}

func NewRecordDTO(r Record) RecordDTO {
	return RecordDTO{ID: r.ID, Type: r.Type, Category: r.Category, Amount: parseIntAmount(r.Amount), Date: r.Date}
}

func parseIntAmount(amount int64) string {
	if amount >= 100 {
	unit := amount / 100
	sUnit := strconv.Itoa(int(unit))
	stringAmount := strconv.Itoa(int(amount))
	res := strings.SplitAfterN(stringAmount, sUnit, 2)
	return strings.Join(res, ",")
	} else if amount >= 10 {
		sSubunit := strconv.Itoa(int(amount))
		return "0," + sSubunit
	} else {
		return "0," + fmt.Sprintf("%02d", amount)
	}
}
