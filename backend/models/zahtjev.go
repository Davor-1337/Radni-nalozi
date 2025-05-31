package models

import (
	"Diplomski/database"
	"Diplomski/utils"
	"database/sql"
	"fmt"
	"time"
)

type Zahtjev struct {
    Zahtjev_ID      int       `json:"zahtjev_id"`
    Username        string    `json:"username"`
    Email           string    `json:"email"`
    Password        string    `json:"password"`
    Role            string    `json:"role"`
    Status          string    `json:"status"`
    VrijemeKreiranja time.Time `json:"vrijemeKreiranja"`

    // Polja za servisera
    Ime           string `json:"name,omitempty"`
    Prezime       string `json:"surname,omitempty"`
    Specijalnost  string `json:"specialty,omitempty"`
    Tel           string `json:"tel,omitempty"`

    // Polja za klijenta
    Naziv         string `json:"naziv,omitempty"`
    KontaktOsoba  string `json:"contactPerson,omitempty"`
    Adresa        string `json:"address,omitempty"`
}



func FetchAllZahtjevi() ([]Zahtjev, error) {
	rows, err := database.DB.Query(`SELECT Zahtjev_ID, username, role, status, VrijemeKreiranja FROM Zahtjevi WHERE Status = 'na cekanju'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zahtjevi []Zahtjev
	for rows.Next() {
		var z Zahtjev
		if err := rows.Scan(&z.Zahtjev_ID, &z.Username, &z.Role, &z.Status, &z.VrijemeKreiranja); err == nil {
			zahtjevi = append(zahtjevi, z)
		}
	}

	return zahtjevi, nil
}


func (z *Zahtjev) Save() error {
    
    hashedPassword, err := utils.HashPassword(z.Password)
    if err != nil {
        return fmt.Errorf("greška pri hashiranju lozinke: %v", err)
    }

    
    query := `
    INSERT INTO Zahtjevi 
    (username, email, password, role, ime, prezime, specijalnost, telefon, naziv, kontaktOsoba, adresa)
    VALUES
    (@Username, @Email, @Password, @Role, @Ime, @Prezime, @Specijalnost, @Tel, @Naziv, @KontaktOsoba, @Adresa)
`

_, err = database.DB.Exec(
    query,
    sql.Named("Username", z.Username),
    sql.Named("Email", z.Email),
    sql.Named("Password", hashedPassword),
    sql.Named("Role", z.Role),
    sql.Named("Ime", z.Ime),
    sql.Named("Prezime", z.Prezime),
    sql.Named("Specijalnost", z.Specijalnost),
    sql.Named("Tel", z.Tel),
    sql.Named("Naziv", z.Naziv),
    sql.Named("KontaktOsoba", z.KontaktOsoba),
    sql.Named("Adresa", z.Adresa),
)
    return err
}


