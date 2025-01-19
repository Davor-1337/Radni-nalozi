package models

import (
	"Diplomski/database"
	"database/sql"
	"time"
	// "database/sql"
)

type Faktura struct{
	Faktura_ID int64 
	Nalog_ID int64 `binding:"required"`
	Iznos float64 `binding:"required"`
	DatumFakture time.Time `binding:"required"`
}

func GetAllInvoices() ([]Faktura, error) {
	query := "SELECT * FROM Faktura"

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var invoices []Faktura

	for rows.Next(){
		var invoice Faktura
		err := rows.Scan(&invoice.Faktura_ID, &invoice.Nalog_ID, &invoice.Iznos, &invoice.DatumFakture)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func (f *Faktura) InsertInvoice() error {
	query := `INSERT INTO Faktura (Nalog_ID, Iznos, DatumFakture)
						OUTPUT INSERTED.Faktura_ID
						VALUES(@Nalog_ID, @Iznos, @DatumFakture)
						`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(
				sql.Named("Nalog_ID", f.Nalog_ID),
				sql.Named("Iznos", f.Iznos),
				sql.Named("DatumFakture", f.DatumFakture)).Scan(&f.Faktura_ID)
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
				SET Nalog_ID = @Nalog_ID, Iznos = @Iznos, DatumFakture = @DatumFakture
				WHERE Faktura_ID = @Faktura_ID`

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
					sql.Named("Nalog_id", invoice.Nalog_ID),
					sql.Named("Iznos", invoice.Iznos),
					sql.Named("DatumFakture", invoice.DatumFakture),
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