// Package internal describe the app
package internal

import (
	"database/sql"
	"errors"
	"finance-tracker/internal/models"
	"finance-tracker/internal/storage"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type App struct {
	categoryStorage storage.CategoryRepository
	recordStorage   storage.RecordRepository
}

func NewApp(catStorage storage.CategoryRepository, recStorage storage.RecordRepository) *App {
	return &App{categoryStorage: catStorage, recordStorage: recStorage}
}

func (a *App) AddExpense(catName string, amount string) error {
	record, err := a.validateRecord(catName, amount)
	if err != nil {
		return err
	}
	record.Type = models.Expense
	record.Date = time.Now().Format(time.DateTime)
	return a.recordStorage.Save(record)
}

func (a *App) AddIncome(catName string, amount string) error {
	record, err := a.validateRecord(catName, amount)
	if err != nil {
		return err
	}
	record.Type = models.Income
	record.Date = time.Now().Format(time.DateTime)
	return a.recordStorage.Save(record)
}

func (a *App) AddCategory(name string) error {
	_, err := a.checkCategory(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cat := models.Category{Name: strings.ToLower(name)}
			if err = a.categoryStorage.Save(cat); err != nil {
				return fmt.Errorf("saving category failed: %w", err)
			}
		}
		return err
	}
	fmt.Println("Category already exist.")
	return nil
}

func (a *App) checkCategory(name string) (models.Category, error) {
	cat, err := a.categoryStorage.Get(strings.ToLower(name))
	if err != nil {
		return models.Category{}, err
	}
	return cat, nil
}

func (a *App) validateRecord(catName, amount string) (models.Record, error) {
	cat, err := a.checkCategory(catName)
	if err != nil {
		return models.Record{}, err
	}
	numAmount, err := parseAmount(amount)
	if err != nil {
		return models.Record{}, err
	}
	return models.Record{Category: cat, Amount: numAmount}, nil
}

func parseAmount(amount string) (int64, error) {
	t, err := regexp.MatchString("^-", amount)
	if err != nil {
		return 0, err
	}
	if t {
		return 0, errors.New("invalid amount format")
	}
	parts := regexp.MustCompile("[,.]").Split(amount, -1)
	unit, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, err
	}
	unit = unit * 100
	if len(parts) == 2 {
		if len(parts[1]) > 2 {
			return 0, errors.New("subunit len is too long")
		}
		subunit, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, err
		}
		if len(parts[1]) < 2 {
			return unit + subunit*10, nil
		}
		return unit + subunit, nil
	} else if len(parts) == 1 {
		return unit, nil
	} else {
		return 0, errors.New("invalid amount format")
	}
}
