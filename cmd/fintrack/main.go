package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"finance-tracker/internal"
	"finance-tracker/internal/storage"

	_ "github.com/mattn/go-sqlite3"
)

const helpMessage = `
Finance tracker usage:

	ftrack <command> [argument]

Available commands:
	add (flags: [-e add expense | -i add income | -c add category]
		[-N category name | -amount]`

func main() {
	db, err := sql.Open("sqlite3", "./tracker.db?_fk=1")
	if err != nil {
		log.Fatalf("open database failed: %v", err)
	}
	defer db.Close()
	flag.Usage = func() {
		fmt.Println(helpMessage)
	}
	categoryStorage, err := storage.NewCategoryStorage(db)
	if err != nil {
		log.Fatalf("create category storage failed: %v", err)
	}
	recordStorage, err := storage.NewRecordStorage(db)
	if err != nil {
		log.Fatalf("create record storage failed: %v", err)
	}
	app := internal.NewApp(categoryStorage, recordStorage)
	if hasCommand := len(os.Args) > 1; !hasCommand {
		fmt.Println("Try -help to get more information")
		return
	}
	switch os.Args[1] {
	case "add":
		fs := flag.NewFlagSet("add", flag.ContinueOnError)
		expense := fs.Bool("e", false, "expense")
		income := fs.Bool("i", false, "income")
		category := fs.Bool("c", false, "category")
		name := fs.String("N", "", "category name")
		amount := fs.String("amount", "", "expense/income amount")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			log.Fatalf("failing to parse arguments: %v", err)
		}
		switch {
		case *expense:
			if err := app.AddExpense(*name, *amount); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					fmt.Printf("Category %s doesn't exist. Wouild you like to create category ?[y/n]: ", *name)
					var choose string
					_, err := fmt.Scan(&choose)
					if err != nil {
						log.Fatal(err)
					}
					if choose == "y" {
						if err := app.AddCategory(*name); err != nil {
							log.Fatal(err)
						}
						if err := app.AddExpense(*name, *amount); err != nil {
							log.Fatal(err)
						}
						fmt.Printf("Expense successfully added: Category %s, amount %s\n", *name, *amount)
						break
					} else {
						break
					}
				} else {
					log.Fatal(err)
				}
			}
			fmt.Printf("Expense successfully added: Category %s, amount %s\n", *name, *amount)
		case *income:
			if err := app.AddIncome(*name, *amount); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					fmt.Printf("Category %s doesn't exist. Would you like to create category ?[y/n]: ", *name)
					var choose string
					_, err = fmt.Scan(&choose)
					if err != nil {
						log.Fatal(err)
					}
					if choose == "y" {
						if err := app.AddCategory(*name); err != nil {
							log.Fatal(err)
						}
						if err := app.AddIncome(*name, *amount); err != nil {
							log.Fatal(err)
						}
						fmt.Printf("Income successfully added: Category %s, amount %s\n", *name, *amount)
						break
					} else {
						break
					}
				} else {
					log.Fatal(err)
				}
			}
			fmt.Printf("Income successfully added: Category %s, amount %s\n", *name, *amount)
		case *category:
			if err := app.AddCategory(*name); err != nil {
				log.Fatal(err)
			}
		}
	case "report":
		fs := flag.NewFlagSet("report", flag.ContinueOnError)
		from := fs.String("from", "", "Stat value from (date layout 2006-01-02 or 2006-01-02 15:04:05)")
		to := fs.String("to", "", "Stat value to (date layout 2006-01-02 or 2006-01-02 15:04:05)")
		format := fs.String("format", "raw", "report format (JSON, raw text)")
		fileName := fs.String("path", "./report.txt", "file path to import report")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
		if err := app.GetRecords(*from, *to, *fileName, *format); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("unknown command")
	}
}
