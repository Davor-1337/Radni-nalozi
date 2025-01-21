package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDB() error {

	err := godotenv.Load("databaseConn.env")
	if err != nil {
	fmt.Printf("Error loading .env file: %v\n", err)
	return err}
 

	// Ucitavanje koneckionog stringa iz env. fajla
	connString := os.Getenv("DB_CONNECTION")
	if connString == "" {
		return fmt.Errorf("DB_CONNECTION is not defined in .env file")
	}
   

	// Povezivanje sa bazom podataka
	DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		return fmt.Errorf("Error opening database connection: %v", err)
	}

	
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("Error pinging database: %v", err)
	}

	return nil
}
