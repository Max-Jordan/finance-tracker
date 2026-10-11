package internal

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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
			res, err := parseAmount(test.input)
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
