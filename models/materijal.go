package models

import (
	"Diplomski/database"
	"database/sql"
)

type Materijal struct {
	Materijal_ID       int64
	NazivMaterijala    string
	Cijena             float64
	KolicinaUSkladistu int64
}

func (m *Materijal) InsertMaterial() error {
	query := `INSERT INTO Materijal (NazivMaterijala, Cijena, KolicinaUSkladistu)
						OUTPUT INSERTED.Materijal_ID
						VALUES (@NazivMaterijala, @Cijena, @KolicinaUSkladistu)`

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(
			sql.Named("NazivMaterijala", m.NazivMaterijala),
			sql.Named("Cijena", m.Cijena),
			sql.Named("KolicinaUSkladistu", m.KolicinaUSkladistu)).Scan(&m.Materijal_ID)
	if err != nil {
		return err
	}
	
	return nil
}

func GetAllMaterials() ([]Materijal, error) {
	query := "SELECT * FROM Materijal"

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []Materijal

	for rows.Next(){
		var material Materijal
		err := rows.Scan(&material.Materijal_ID, &material.NazivMaterijala, &material.Cijena, &material.KolicinaUSkladistu)
		if err != nil {
			return nil, err
		}
		materials = append(materials, material)
	}

	return materials, nil
}

func GetMaterialByID(id int64) (*Materijal, error) {
	query := `SELECT * FROM Materijal WHERE Materijal_ID = @Materijal_ID`

	row := database.DB.QueryRow(query, sql.Named("Materijal_ID", id))

	var material Materijal 

	err := row.Scan(&material.Materijal_ID, &material.NazivMaterijala, &material.Cijena, &material.KolicinaUSkladistu)
	if err != nil {
		return nil, err
	}

	return &material, nil
}

func (material Materijal) Update() error {
	query := `UPDATE Materijal
						SET NazivMaterijala = @NazivMaterijala, Cijena = @Cijena, KolicinaUSkladistu = @KolicinaUSkladistu
						WHERE Materijal_ID = @Materijal_ID`
	
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		sql.Named("NazivMaterijala", material.NazivMaterijala),
		sql.Named("Cijena", material.Cijena),
		sql.Named("KolicinaUSkladistu", material.KolicinaUSkladistu),
		sql.Named("Materijal_ID", material.Materijal_ID))
	
		return err
}

func (material Materijal) Delete() error {
	query := "DELETE FROM Materijal WHERE Materijal_ID = @Materijal_ID"

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sql.Named("Materijal_ID", material.Materijal_ID))

	return err
}