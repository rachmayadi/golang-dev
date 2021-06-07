package main

import (
	"os"

	"github.com/gocarina/gocsv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Entry defines both the CSV layout and database schema
type Entry struct {
	gorm.Model

	Field1 float64 `csv:"field1"`
	Field2 float64 `csv:"field2"`
}

func main() {
	// Open the CSV file for reading
	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Parse CSV into a slice of typed data `[]Entry` (just like json.Unmarshal() does)
	// The builtin package `encoding/csv` does not support unmarshaling into a struct
	// so you need to use an external library to avoid writing for-loops.
	var entries []Entry
	err = gocsv.Unmarshal(file, &entries)
	if err != nil {
		panic(err)
	}

	// Open a postgres database connection using GORM
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=dev dbname=foo port=5432 sslmode=disable TimeZone=Europe/Paris"))
	if err != nil {
		panic(err)
	}

	// Create `entries` table if not exists
	err = db.AutoMigrate(&Entry{})
	if err != nil {
		panic(err)
	}

	// Save all the records at once in the database
	result := db.Create(entries)
	if result.Error != nil {
		panic(result.Error)
	}
}
