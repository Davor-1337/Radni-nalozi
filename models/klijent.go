package models

import (
	"Diplomski/database"
	"database/sql"
	
	
)

type Klijent struct {
	Klijent_ID   int64  
	Naziv        string
	KontaktOsoba string
	Email        string
	Tel          string 
	Adresa       string 
	User_ID int64
}


func GetAllClients() ([]Klijent, error) {
	query := "SELECT Klijent_ID, Naziv, KontaktOsoba, Email, Tel, Adresa, User_ID FROM klijenti"
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []Klijent

	for rows.Next(){
		var client Klijent 
		err := rows.Scan(&client.Klijent_ID, &client.Naziv, &client.KontaktOsoba, &client.Email, &client.Tel, &client.Adresa, &client.User_ID)
		if err != nil{
			return nil, err
		}

		clients = append(clients, client)
	}
	return clients, nil
  }
	 

func (k *Klijent) InsertClient() error {
    query := `INSERT INTO Klijenti(Klijent_ID, Naziv, KontaktOsoba, Email, Tel, Adresa, User_ID)
              VALUES (@Klijent_ID, @Naziv, @KontaktOsoba, @Email, @Tel, @Adresa, @User_ID)`

  
    stmt, err := database.DB.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()


    _, err = stmt.Exec(
				sql.Named("Klijent_ID", k.Klijent_ID),
        sql.Named("Naziv", k.Naziv),
        sql.Named("KontaktOsoba", k.KontaktOsoba),
        sql.Named("Email", k.Email),
        sql.Named("Tel", k.Tel),
        sql.Named("Adresa", k.Adresa),
				sql.Named("User_ID", k.User_ID))

    if err != nil {
        return err
    }

    return nil
}

func GetClientByID(id int64) (*Klijent, error) {
	query := "SELECT Klijent_ID, Naziv, KontaktOsoba, Email, Tel, Adresa, User_ID FROM Klijenti WHERE Klijent_id = @Klijent_id"
	
	// Koristimo `QueryRow` sa imenovanim parametrom
	row := database.DB.QueryRow(query, sql.Named("Klijent_id", id))

	var client Klijent

	// Mapiranje rezultata u `client` strukturu
	err := row.Scan(&client.Klijent_ID, &client.Naziv, &client.KontaktOsoba, &client.Email, &client.Tel, &client.Adresa, &client.User_ID)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (client Klijent) Update() error {
	query := `UPDATE Klijenti 
						SET Naziv = @Naziv, KontaktOsoba = @KontaktOsoba, Email = @Email, Tel = @Tel, Adresa = @Adresa, @User_ID = User_ID
						WHERE Klijent_id = @Klijent_id`

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Koristimo `sql.Named` za svaki parametar da bismo tačno mapirali vrednosti u upitu
	_, err = stmt.Exec(
		sql.Named("Naziv", client.Naziv),
		sql.Named("KontaktOsoba", client.KontaktOsoba),
		sql.Named("Email", client.Email),
		sql.Named("Tel", client.Tel),
		sql.Named("Adresa", client.Adresa),
		sql.Named("User_ID", client.User_ID),
		sql.Named("Klijent_id", client.Klijent_ID))

	return err
}

func (client Klijent)Delete() error {
	query := "DELETE FROM Klijenti WHERE Klijent_ID = @Klijent_ID"

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(sql.Named("Klijent_id", client.Klijent_ID))

	defer stmt.Close()

	return err
}
