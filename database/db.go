package database

import (
    "database/sql"
    "fmt"
    _ "github.com/denisenkom/go-mssqldb"
)

var DB *sql.DB

func ConnectDB() error {
    connString := "sqlserver://davor:davor@localhost:1433?database=RadniNaloziDB"

    // Povezivanje sa bazom podataka
    var err error
    DB, err = sql.Open("sqlserver", connString) // Ovdje postavljamo globalnu varijablu
    if err != nil {
        return fmt.Errorf("greska pri otvaranju konekcije: %v", err)
    }

    // Testiranje konekcije
    err = DB.Ping()
    if err != nil {
        return fmt.Errorf("greska pri pingovanju baze: %v", err)
    }

    fmt.Println("Uspesno povezano sa bazom podataka!")
    return nil
}


     

