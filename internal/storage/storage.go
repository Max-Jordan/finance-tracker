// Package storage describe app storage
package storage

import (
	"database/sql"
	"finance-tracker/internal/models"
	"fmt"
)

const initCategoryStorage = "create table if not exists categories (id integer not null primary key, name text not null unique)"
const initRecordStorage = `
create table if not exists records (
id integer not null primary key,
type text not null,
category integer not null,
amount integer not null,
date text not null,
foreign key (category) references categories(id)
)`

type CategoryRepository interface {
	Save(models.Category) error
	Get(string) (models.Category, error)
}

type RecordRepository interface {
	Save(models.Record) error
	GetByPeriod(start, end string) ([]models.Record, error)
}

type categoryStorage struct {
	db *sql.DB
}

func NewCategoryStorage(db *sql.DB) (*categoryStorage, error) {
	if _, err := db.Exec(initCategoryStorage); err != nil {
		return nil, fmt.Errorf("failed to init category storage: %w", err)
	}
	return &categoryStorage{db: db}, nil
}

func (c *categoryStorage) Save(cat models.Category) error {
	query := "insert into categories (name) values ($1)"
	if _, err := c.db.Exec(query, cat.Name); err != nil {
		return fmt.Errorf("failed to save category: %w", err)
	}
	return nil
}

func (c *categoryStorage) Get(name string) (models.Category, error) {
	query := "select * from categories where name = $1"
	var cat models.Category
	if err := c.db.QueryRow(query, name).Scan(&cat.ID, &cat.Name); err != nil {
		return models.Category{}, fmt.Errorf("failed to get a category: %w", err)
	}
	return cat, nil
}

type recordStorage struct {
	db *sql.DB
}

func NewRecordStorage(db *sql.DB) (*recordStorage, error) {
	if _, err := db.Exec(initRecordStorage); err != nil {
		return nil, fmt.Errorf("failed to init records storage: %w", err)
	}
	return &recordStorage{db: db}, nil
}

func (r *recordStorage) Save(rec models.Record) error {
	query := "insert into records (type, category, amount, date) values ($1, $2, $3, $4)"
	if _, err := r.db.Exec(query, rec.Type, rec.Category.ID, rec.Amount, rec.Date); err != nil {
		return fmt.Errorf("failed to save record: %w", err)
	}
	return nil
}

func (r *recordStorage) GetByPeriod(start, end string) ([]models.Record, error) {
	query := `select r.id, r.type, r.amount, r.date, c.id, c.name from records as r
	left join categories as c on r.category = c.id
	where r.date >= $1 and r.date < $2`
	rows, err := r.db.Query(query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get result by period: %w", err)
	}
	defer rows.Close()
	result := make([]models.Record, 0, 2)
	for rows.Next() {
		var record models.Record
		var category models.Category
		if err := rows.Scan(&record.ID, &record.Type, &record.Amount, &record.Date, &category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("failed to scan result: %w", err)
		}
		record.Category = category
		result = append(result, record)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return result, nil
}
