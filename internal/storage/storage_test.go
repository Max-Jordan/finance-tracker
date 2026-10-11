package storage

import (
	"database/sql"
	"finance-tracker/internal/models"
	"path/filepath"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCategoryStorage(t *testing.T) {
	db := createTempDB(t)

	categoryStorage, err := NewCategoryStorage(db)
	if err != nil {
		t.Fatalf("Creating category storage: want error false, got %v", err)
	}

	category := models.Category{Name: "Groceries"}

	saveTests := []struct {
		name    string
		input   models.Category
		wantErr bool
	}{
		{name: "Save category", input: category, wantErr: false},
		{name: "Duplicate category", input: category, wantErr: true},
	}

	for _, test := range saveTests {
		t.Run(test.name, func(t *testing.T) {
			err := categoryStorage.Save(test.input)
			getErr := err != nil
			if getErr != test.wantErr {
				t.Errorf("Test: %s, want error %t, get error %t", test.name, test.wantErr, getErr)
			}
		})
	}

	getTests := []struct {
		name       string
		input      string
		wantResult models.Category
		wantErr    bool
	}{
		{name: "Get category", input: category.Name, wantResult: models.Category{ID: 1, Name: category.Name}, wantErr: false},
		{name: "Get doesn't exist categoty", input: "hello", wantResult: models.Category{}, wantErr: true},
	}
	for _, test := range getTests {
		t.Run(test.name, func(t *testing.T) {
			res, err := categoryStorage.Get(test.input)
			getErr := err != nil
			if getErr != test.wantErr {
				t.Fatalf("Test: %s, want error %t, get error %t", test.name, test.wantErr, getErr)
			}
			if res != test.wantResult {
				t.Errorf("Test: %s, want %v, got %v", test.name, test.wantResult, res)
			}
		})
	}
}

func TestRecordStorage(t *testing.T) {
	db := createTempDB(t)

	recordStorage, err := NewRecordStorage(db)
	if err != nil {
		t.Fatalf("Creating storage: want error false, got %v", err)
	}

	categoryStorage, err := NewCategoryStorage(db)
	if err != nil {
		t.Fatalf("Creating category storage: want error false, got %v", err)
	}

	expenseCat := models.Category{Name: "Groceries"}
	err = categoryStorage.Save(expenseCat)
	if err != nil {
		t.Fatalf("Saving category: want error false, got error %v", err)
	}
	incomeCat := models.Category{Name: "Salary"}
	err = categoryStorage.Save(incomeCat)
	if err != nil {
		t.Fatalf("Saving category: want error false, got error %v", err)
	}

	expenseCat, err = categoryStorage.Get(expenseCat.Name)
	if err != nil {
		t.Fatalf("Getting expense category: want err false, got error %v", err)
	}
	incomeCat, err = categoryStorage.Get(incomeCat.Name)
	if err != nil {
		t.Fatalf("Getting income category: want error false, got error %v", err)
	}

	expenseRecord := models.Record{Type: models.Expense, Category: expenseCat, Amount: 12_064, Date: "1999-01-01 10:00:00"}
	incomeRecord := models.Record{Type: models.Income, Category: incomeCat, Amount: 20_000, Date: "1999-01-01 12:00:00"}
	saveRecordTests := []struct {
		name    string
		input   models.Record
		wantErr bool
	}{
		{name: "Save expense record", input: expenseRecord, wantErr: false},
		{name: "Save income record", input: incomeRecord, wantErr: false},
	}

	for _, test := range saveRecordTests {
		t.Run(test.name, func(t *testing.T) {
			err := recordStorage.Save(test.input)
			gotErr := err != nil
			if gotErr != test.wantErr {
				t.Errorf("Test: %s, want error %t, got error %t", test.name, test.wantErr, gotErr)
			}
		})
	}

	expenseRecord.ID = 1
	incomeRecord.ID = 2
	recordResults := []models.Record{expenseRecord, incomeRecord}
	getRecordTests := []struct {
		name       string
		inputStart string
		inputEnd   string
		wantResult []models.Record
		wantErr    bool
	}{
		{
			name: "Get records: correct period", inputStart: "1999-01-01 09:00:00",
			inputEnd: "1999-01-01 13:00:00", wantResult: recordResults, wantErr: false,
		}, {
			name: "Get records: included end date", inputStart: "1999-01-01 10:00:00",
			inputEnd: "1999-01-01 12:00:00", wantResult: []models.Record{expenseRecord}, wantErr: false,
		}, {
			name: "Get records: period doesn't exist", inputStart: "2000-01-01 10:00:00",
			inputEnd: "2000-01-01 12:00:00", wantResult: []models.Record{}, wantErr: false,
		}, {
			name: "Get records: missing period", wantResult: []models.Record{}, wantErr: false,
		},
	}

	for _, test := range getRecordTests {
		t.Run(test.name, func(t *testing.T) {
			result, err := recordStorage.GetByPeriod(test.inputStart, test.inputEnd)
			gotErr := err != nil
			if gotErr != test.wantErr {
				t.Fatalf("Test: %s, want error %t, got error %t", test.name, test.wantErr, gotErr)
			}
			if !reflect.DeepEqual(test.wantResult, result) {
				t.Errorf("Test: %s, want %v, got %v", test.name, test.wantResult, result)
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
