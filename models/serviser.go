package models

import (
	"Diplomski/database"
	"database/sql"
	"time"
)

type Zadatak struct{
	Zadatak_ID int64
	Nalog_ID int64
	Serviser_ID int64
	DatumDodjele time.Time
	DatumZavrsetka time.Time
}

type Serviser struct {
	Serviser_ID   int64  
	Ime        string 
	Prezime string 
	Specijalnost        string 
	Telefon          string
	User_ID int64
}

func GetAllTehnicians()([]Serviser, error){
query := `SELECT Serviser_ID, Ime, Prezime, Specijalnost, Telefon, User_ID  FROM Serviser` 
rows, err := database.DB.Query(query)
if err != nil {
		return nil, err
	}
defer rows.Close()

var tehnicians []Serviser

for rows.Next(){
	var tehnician Serviser
	err := rows.Scan(&tehnician.Serviser_ID, &tehnician.Ime, &tehnician.Prezime, &tehnician.Specijalnost, &tehnician.Telefon, &tehnician.User_ID)
		if err != nil {
		return nil, err
	}
	tehnicians = append(tehnicians, tehnician)
}
	return tehnicians, nil
}

func (s *Serviser) InsertTehnician()error{
	query := `INSERT INTO Serviser(Serviser_ID, Ime, Prezime, Specijalnost, Telefon, User_ID)
						VALUES(@Serviser_ID, @Ime, @Prezime, @Specijalnost, @Telefon, @User_ID)`
	
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
				sql.Named("Serviser_ID", s.Serviser_ID),
				sql.Named("Ime", s.Ime),
        sql.Named("Prezime", s.Prezime),
        sql.Named("Specijalnost", s.Specijalnost),
        sql.Named("Telefon", s.Telefon),
				sql.Named("User_ID", s.User_ID))

	if err != nil {
    return err
  }

    return nil
}

func GetTehnicianByID(id int64) (*Serviser, error) {
query := `SELECT Serviser_ID, Ime, Prezime, Specijalnost, Telefon, User_ID FROM Serviser WHERE Serviser_ID = @Serviser_ID`

row := database.DB.QueryRow(query, sql.Named("Serviser_ID", id))

var tehnician Serviser

err := row.Scan(&tehnician.Serviser_ID, &tehnician.Ime, &tehnician.Prezime, &tehnician.Specijalnost, &tehnician.Telefon, &tehnician.User_ID)
if err != nil {
	return nil, err
}

return &tehnician, nil
}


func (tehnician Serviser) Update() error {
	query := `UPDATE Serviser
						set Ime = @Ime, Prezime = @Prezime, Specijalnost = @Specijalnost, Telefon = @Telefon, User_ID = @User_ID
						WHERE Serviser_ID = @Serviser_ID`
	
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err 
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		sql.Named("Ime", tehnician.Ime),
		sql.Named("Prezime", tehnician.Prezime),
		sql.Named("Specijalnost", tehnician.Specijalnost),
		sql.Named("Telefon", tehnician.Telefon),
		sql.Named("User_ID", tehnician.User_ID),
		sql.Named("Serviser_ID", tehnician.Serviser_ID))

		return err
}

func GetWorkOrderByTehnicianID(serviserID int64) ([]RadniNalog, error) {
	query := `SELECT rn.Nalog_ID, rn.Klijent_ID, rn.OpisProblema, rn.Prioritet, rn.DatumOtvaranja, rn.Status, rn.Lokacija
        FROM RadniNalog rn
        INNER JOIN DodjelaZadatka dz ON rn.Nalog_ID = dz.Nalog_ID
        WHERE dz.Serviser_ID = @Serviser_ID`
	
	rows, err := database.DB.Query(query, sql.Named("Serviser_ID", serviserID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workOrders []RadniNalog

	for rows.Next(){
		var workOrder RadniNalog
		err := rows.Scan(&workOrder.Nalog_ID, &workOrder.Klijent_ID, &workOrder.OpisProblema, &workOrder.Prioritet, &workOrder.DatumOtvaranja, &workOrder.Status, &workOrder.Lokacija)

		if err != nil {
		return nil, err
	}
		workOrders = append(workOrders, workOrder)
	} 
	return workOrders, nil
}


func GetHoursForTehnician(serviserId int64) (float64, error) {
	query := `
        SELECT 
            COALESCE(SUM(BrojRadnihSati), 0) AS UkupniSati
        FROM 
            EvidencijaSati
        WHERE 
            Serviser_ID = @Serviser_ID`

	var totalHours float64

	row := database.DB.QueryRow(query, sql.Named("Serviser_ID", serviserId))
	err := row.Scan(&totalHours)

	if err != nil {
		return 0, err
	}

	return totalHours, nil
}