package models

import (
	"Diplomski/database"
	"database/sql"
	"fmt"
	"time"
)

type Faktura struct{
	Faktura_ID int64 
	Nalog_ID int64 
	Iznos float64 
	DatumFakture time.Time 
}

func GetAllInvoices() ([]struct {
    BrojNaloga    string    `json:"broj_naloga"`
    DatumFakture  time.Time `json:"datum_fakture"`
    Iznos         float64   `json:"iznos"`
		Faktura_ID    int       `json:"faktura_id"`
		
}, error) {
    query := `
    SELECT rn.BrojNaloga, f.DatumFakture, f.Iznos, f.Faktura_ID
    FROM Faktura f 
    JOIN RadniNalog rn ON rn.Nalog_ID = f.Nalog_ID
`

    rows, err := database.DB.Query(query)
    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }
    defer rows.Close()

    var results []struct {
        BrojNaloga    string    `json:"broj_naloga"`
        DatumFakture  time.Time `json:"datum_fakture"`
        Iznos         float64   `json:"iznos"`
				Faktura_ID    int       `json:"faktura_id"`
				
    }

    for rows.Next() {
        var result struct {
            BrojNaloga    string    `json:"broj_naloga"`
            DatumFakture  time.Time `json:"datum_fakture"`
            Iznos         float64   `json:"iznos"`
						Faktura_ID    int       `json:"faktura_id"`
						
        }
        err := rows.Scan( &result.BrojNaloga, &result.DatumFakture, &result.Iznos,  &result.Faktura_ID)
        if err != nil {
            return nil, fmt.Errorf("row scan failed: %w", err)
        }
        results = append(results, result)
    }
		

    return results, nil
}

func GetAllInvoicesForClient(clientID int64) ([]struct {
	BrojNaloga   string    `json:"broj_naloga"`
	DatumFakture time.Time `json:"datum_fakture"`
	Iznos        float64   `json:"iznos"`
	Faktura_ID   int       `json:"faktura_id"`
}, error) {
	const query = `
			SELECT rn.BrojNaloga,
						 f.DatumFakture,
						 f.Iznos,
						 f.Faktura_ID
				FROM Faktura f
				JOIN RadniNalog rn ON rn.Nalog_ID = f.Nalog_ID
			 WHERE rn.Klijent_ID = @ClientID
			 ORDER BY f.DatumFakture DESC
	`

	rows, err := database.DB.Query(
			query,
			sql.Named("ClientID", clientID),
	)
	if err != nil {
			return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []struct {
			BrojNaloga   string    `json:"broj_naloga"`
			DatumFakture time.Time `json:"datum_fakture"`
			Iznos        float64   `json:"iznos"`
			Faktura_ID   int       `json:"faktura_id"`
	}

	for rows.Next() {
			var inv struct {
					BrojNaloga   string `json:"broj_naloga"`
					DatumFakture time.Time `json:"datum_fakture"`
					Iznos        float64 `json:"iznos"`
					Faktura_ID   int `json:"faktura_id"`
			}
			if err := rows.Scan(
					&inv.BrojNaloga,
					&inv.DatumFakture,
					&inv.Iznos,
					&inv.Faktura_ID,
			); err != nil {
					return nil, fmt.Errorf("row scan failed: %w", err)
			}
			results = append(results, inv)
	}
	if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return results, nil
}

func (f *Faktura) InsertInvoice() error {
	query := `INSERT INTO Faktura (Nalog_ID, Iznos, DatumFakture)
						OUTPUT INSERTED.Faktura_ID
						VALUES(@Nalog_ID, @Iznos, GETDATE())
						`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(
				sql.Named("Nalog_ID", f.Nalog_ID),
				sql.Named("Iznos", f.Iznos)).Scan(&f.Faktura_ID)
	if err != nil {
		return err
	}

	


	return nil
}

func GetInvoiceByID(id int64) (*Faktura, error) {
	query := "SELECT * FROM Faktura WHERE Faktura_ID = @Faktura_ID"

	row := database.DB.QueryRow(query, sql.Named("Faktura_ID", id))

	var invoice Faktura

	err := row.Scan(&invoice.Faktura_ID, &invoice.Nalog_ID, &invoice.Iznos, &invoice.DatumFakture)
	if err != nil {
		return nil, err 
	}

	return &invoice, nil

}

func (invoice Faktura) Update() error {
	
	query := `UPDATE Faktura
				SET  Iznos = @Iznos
				WHERE Faktura_ID = @Faktura_ID`

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	fmt.Println("UPDATE podaci:", invoice)

	_, err = stmt.Exec(
					sql.Named("Iznos", invoice.Iznos),
					sql.Named("Faktura_ID", invoice.Faktura_ID))	
	
	return err
}

func (invoice Faktura) Delete() error {
	query := "DELETE FROM Faktura WHERE Faktura_ID = @Faktura_ID"

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sql.Named("Faktura_ID", invoice.Faktura_ID))
	if err != nil {
		return err
	}

	return nil 
}