package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Max-Jordan/finance-tracker/internal/models"
	"github.com/Max-Jordan/finance-tracker/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

func TestAddExpense(t *testing.T) {
	app := createTempApp(t)

	if err := app.categoryStorage.Save(models.Category{Name: "groceries"}); err != nil {
		t.Fatal("Saving category failed")
	}

	tests := []struct {
		name          string
		inputCategory string
		inputAmount   string
		wantErr       bool
	}{
		{name: "existing category", inputCategory: "Groceries", inputAmount: "123.12", wantErr: false},
		{name: "incorrect amount", inputCategory: "Groceries", inputAmount: "123,123", wantErr: true},
		{name: "category isn't exist", inputCategory: "Other", inputAmount: "123", wantErr: true},
		{name: "empty amount", inputCategory: "Groceries", inputAmount: "", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := app.AddExpense(test.inputCategory, test.inputAmount)
			gotErr := err != nil
			if gotErr != test.wantErr {
				t.Errorf("Test: %s, want error %t, got error %v", test.name, test.wantErr, err)
			}
		})
	}
}

func TestAddIncome(t *testing.T) {
	app := createTempApp(t)

	if err := app.categoryStorage.Save(models.Category{Name: "salary"}); err != nil {
		t.Fatal("Saving category failed")
	}

	tests := []struct {
		name          string
		inputCategory string
		inputAmount   string
		wantErr       bool
	}{
		{name: "existring category", inputCategory: "salary", inputAmount: "1000,42", wantErr: false},
		{name: "incorrect amount", inputCategory: "salary", inputAmount: "1000,42,50", wantErr: true},
		{name: "category isn't exist", inputCategory: "other", inputAmount: "1000,50", wantErr: true},
		{name: "empty amount", inputCategory: "salary", wantErr: true},
		{name: "empty category", inputAmount: "1000,42", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := app.AddIncome(test.inputCategory, test.inputAmount)
			gotErr := err != nil
			if gotErr != test.wantErr {
				t.Errorf("Test: %s, want error %t, got %v", test.name, test.wantErr, err)
			}
		})
	}
}

func TestAddCategory(t *testing.T) {
	app := createTempApp(t)

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "correct category", input: "groceries", wantErr: false},
		{name: "duplicate category", input: "groceries", wantErr: false},
		{name: "empty name", wantErr: true},
	}

	for _, test := range tests {
		err := app.AddCategory(test.input)
		gotErr := err != nil
		if gotErr != test.wantErr {
			t.Errorf("Test: %s, want error %t, got error %v", test.name, test.wantErr, err)
		}
	}
}

func TestAddCategorySaveCategory(t *testing.T) {
	app := createTempApp(t)

	if err := app.AddCategory("Groceries"); err != nil {
		t.Fatalf("Adding category erro: %v", err)
	}

	result, err := app.categoryStorage.Get("groceries")
	if err != nil {
		t.Errorf("Getting category: want error false, got error %v", err)
	}

	if result.Name != "groceries" {
		t.Errorf("Want category name groceries, got %v", result)
	}
}

func TestCreateReport(t *testing.T) {
	app := createTempApp(t)

	if err := app.AddCategory("groceries"); err != nil {
		t.Fatalf("Adding category error: %v", err)
	}

	if err := app.AddExpense("groceries", "123,45"); err != nil {
		t.Fatalf("Adding record error: %v", err)
	}
	dir := t.TempDir()
	testJSON := filepath.Join(dir, "test_json.json")
	testRaw := filepath.Join(dir, "test_raw.txt")

	startPeriod := time.Now().Add(-time.Second).Format(time.DateTime)
	endPeriod := time.Now().Add(time.Second).Format(time.DateTime)

	result, err := app.recordStorage.GetByPeriod(startPeriod, endPeriod)
	if err != nil {
		t.Fatal("Getting record failed")
	}
	if len(result) != 1 {
		t.Fatalf("Failed: want len result 1, got %d", len(result))
	}
	model := models.NewRecordDTO(result[0])

	tests := []struct {
		name        string
		fromInput   string
		toInput     string
		pathInput   string
		formatInput string
		wantResult  models.RecordDTO
		wantErr     bool
	}{
		{name: "JSON report", fromInput: startPeriod,
			toInput: endPeriod, pathInput: testJSON,
			formatInput: "JSON", wantResult: model, wantErr: false},
		{name: "raw text report", fromInput: startPeriod,
			toInput: endPeriod, pathInput: testRaw,
			formatInput: "raw", wantResult: model, wantErr: false},
		{name: "empty format report", fromInput: startPeriod,
			toInput: endPeriod, pathInput: testRaw, wantErr: true},
		{name: "no records period", fromInput: "1000-01-01 10:10:10", toInput: "1000-01-01 10:10:10", wantErr: true},
		{name: "wrong from date time format", fromInput: "01, Jan 26 10:10:10", wantErr: true},
		{name: "wrong to date time format", fromInput: startPeriod, toInput: "abc", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := app.CreateReport(test.fromInput, test.toInput, test.pathInput, test.formatInput)
			gotErr := err != nil
			if gotErr != test.wantErr {
				t.Fatalf("Test: %s, want error %t, got error %v", test.name, test.wantErr, err)
			} else if test.wantErr{
				return
			}
			data, err := os.ReadFile(test.pathInput)
			if err != nil {
				t.Fatalf("Reading from file failed: %v", err)
			}

			if test.formatInput == "JSON" {
				var record models.RecordDTO
				if err := json.Unmarshal(data, &record); err != nil {
					t.Fatalf("Unmarshalling record failed: %v", err)
				}
				if record != test.wantResult {
					t.Errorf("Test: %s, want result %v, got result %v", test.name, test.wantResult, record)
				}
			} else {
				if string(data) != fmt.Sprintln(test.wantResult) {
					t.Errorf("Test: %s, want result %v, got result %v", test.name, test.wantResult, string(data))
				}
			}
		})
	}
}

func TestAddExpenseSaveRecords(t *testing.T) {
	app := createTempApp(t)

	if err := app.categoryStorage.Save(models.Category{Name: "groceries"}); err != nil {
		t.Fatal("Saving category failed")
	}

	if err := app.AddExpense("Groceries", "123,00"); err != nil {
		t.Fatal("Adding expense failed")
	}

	startTime := time.Now().Add(-time.Second).Format(time.DateTime)
	endTime := time.Now().Add(time.Second).Format(time.DateTime)

	result, err := app.recordStorage.GetByPeriod(startTime, endTime)
	if err != nil {
		t.Fatal("Getting record failed")
	}

	if result[0].Category.Name != "groceries" {
		t.Errorf("Record category name: want groceries, got %s", result[0].Category.Name)
	}
	if result[0].Amount != 12300 {
		t.Errorf("Record amount: want %d, got %d", 12300, result[0].Amount)
	}
}

func TestParseAmount(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNum int64
		wantErr bool
	}{
		{name: "Num without subunit", input: "123", wantNum: 12300, wantErr: false},
		{name: "One symbol after dot", input: "100.4", wantNum: 10040, wantErr: false},
		{name: "One symbol after comma", input: "100,4", wantNum: 10040, wantErr: false},
		{name: "Two symbols after dot with zero", input: "100.40", wantNum: 10040, wantErr: false},
		{name: "Two symbols after dot without zero", input: "100.44", wantNum: 10044, wantErr: false},
		{name: "string argument", input: "abc", wantNum: 0, wantErr: true},
		{name: "unpositive argument", input: "-123,00", wantNum: 0, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := parseStringAmount(test.input)
			if gotErr := err != nil; gotErr != test.wantErr {
				t.Fatalf("Test: %s, want error %t, got error %t", test.name, test.wantErr, gotErr)
			}
			if res != test.wantNum {
				t.Errorf("Test: %s, want %d, got %d", test.name, test.wantNum, res)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	date := time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)
	dateTime := time.Date(1999, time.January, 1, 10, 10, 10, 0, time.UTC)

	tests := []struct {
		name       string
		input      string
		wantResult time.Time
		wantLayout string
		wantErr    bool
	}{
		{name: "Parse date", input: "1999-01-01", wantResult: date, wantLayout: time.DateOnly, wantErr: false},
		{name: "Parse date time", input: "1999-01-01 10:10:10", wantResult: dateTime, wantLayout: time.DateTime, wantErr: false},
		{name: "Incorrect date format", input: "01, Jan 1999 10:10:10", wantErr: true},
		{name: "Invalid calendar date", input: "1999-01-100", wantErr: true},
		{name: "Invalid time", input: "1999-01-01 25:00:00", wantErr: true},
		{name: "Missing seconds", input: "1999-01-01 10:10", wantErr: true},
		{name: "Empty input", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, layout, err := parseDate(test.input)
			gotErr := err != nil
			if gotErr != test.wantErr {
				t.Fatalf("Test: %s, want error %t, got error %t", test.name, test.wantErr, gotErr)
			}
			if !res.Equal(test.wantResult) {
				t.Errorf("Test: %s, want result %v, got gesult %v", test.name, test.wantResult, res)
			}
			if layout != test.wantLayout {
				t.Errorf("Test: %s, want layout %s, got layout %s", test.name, test.wantLayout, layout)
			}
		})
	}
}

func createTempDB(t *testing.T) *sql.DB {
	t.Helper()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := sql.Open("sqlite3", dbPath+"?_fk=1")
	if err != nil {
		t.Fatalf("failed to init temp CategoryStorage: %v", err)
	}

	defer t.Cleanup(func() {
		db.Close()
	})

	return db
}

func createTempApp(t *testing.T) *App {
	t.Helper()

	db := createTempDB(t)

	recordStorage, err := storage.NewRecordStorage(db)
	if err != nil {
		t.Fatalf("Creating records storage failed: %v", err)
	}
	categoryStorage, err := storage.NewCategoryStorage(db)
	if err != nil {
		t.Fatalf("Creating category storage filed: %v", err)
	}
	return NewApp(categoryStorage, recordStorage)
}

func createTempRecordStorage(t *testing.T, db *sql.DB) storage.RecordRepository {
	t.Helper()

	recordStorage, err := storage.NewRecordStorage(db)
	if err != nil {
		t.Fatalf("Creating record storage failed: %v", err)
	}
	return recordStorage
}

func createTempCategoryStorage(t *testing.T, db *sql.DB) storage.CategoryRepository {
	t.Helper()

	categoryStorage, err := storage.NewCategoryStorage(db)
	if err != nil {
		t.Fatalf("Creating category storage failed: %v", err)
	}
	return categoryStorage
}
