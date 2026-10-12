# CLi Finance tracker app

Simple command-line finance tracker written in GO.

The app allows you to add income and expense records, manage categories, filter records by date period and export reports.

## Features

- Add expense records
- Add income records
- Add categories
- Store data in SQLite
- Filter records by date period
- Export report in raw text format
- Export reports in JSON format
- Foreign key support between records and categories
- Tests for application logic and storage layer

## Tech Stack

- Go
- SQLite3
- Git

## Project Structure

```text
cmd/
    fintrack/
        main.go

internal/
    app.go
    app_test.go

    models/
        models.go

    sotrage/
        storage.go
        storage_test.go
```

## Build

```bash
    git clone https://github.com/Max-Jordan/finance-tracker
    cd finance-tracker

    go build -o fintrack ./cmd/fintrack
```

## Usage

### Add a category

```bash
./fintrack add -c -N Groreries
```

### Add an expense

```bash
./fintrack add -e -N Groceries -amount 10.64(10,64)
```

### Add an income

```bash
./fintrack add -i -N Salary -amount 1000
```

Money is stored internally as integer subinit:

```raw
120,64 -> 12064
```

This approach avoids floating-ping arithemtic for monetary values.

### Reports

Get records for a period:

```bash
./fintrack report -from '1999-01-01' -to '1999-01-02'
```

DateTime also supported:

```bash
./fintrack report -from '1999-01-01 10:00:00' -to '1999-01-02 12:00:00'
```

JSON export:

```bash
./fintrack report -from 'date' -to 'date' -json -path 'file_path.json'
```

Text expost(default):

```bash
./fintrack report -from 'date' -to 'date' -path ''file_path.txt''
```

Default path to export is ./report.txt

## Database

The project uses two tables.

```sql
    create table if not exists categories(
        id integer not null primary key,
        name text not null unique
    );

    create table if not exists records(
        id integer not null primary key,
        type text not null,
        category integer not null,
        amount integer not null,
        date text not null,
        foreign key(category) references categories(id)
    );
```
