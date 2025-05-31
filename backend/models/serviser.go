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

type RadniNalogServiser struct {
	Nalog_ID       int64     `json:"Nalog_ID"`
	Klijent_ID     int64     `json:"Klijent_ID"`
	OpisProblema   string    `json:"OpisProblema"`
	Prioritet      string    `json:"Prioritet"`
	DatumOtvaranja time.Time `json:"DatumOtvaranja"`
	Status         string    `json:"Status"`
	Lokacija       string    `json:"Lokacija"`
	BrojNaloga     string    `json:"BrojNaloga"`
	Serviser_ID    int64     `json:"Serviser_ID"`
}


type DetaljiRadnogNaloga struct {
    OpisProblema   string    `json:"OpisProblema"`
    NazivKlijenta  string    `json:"NazivKlijenta"`
    Lokacija       string    `json:"Lokacija"`
    DatumDodjele   time.Time `json:"DatumDodjele"`
    DatumZavrsetka time.Time `json:"DatumZavrsetka"`
    BrojRadnihSati float64   `json:"BrojRadnihSati"`
}

type RadniNalogSaKlijentom struct {
	Nalog_ID     int64  `json:"Nalog_ID"`
	BrojNaloga   string `json:"BrojNaloga"`
	Klijent      string `json:"Klijent"`
	OpisProblema string `json:"OpisProblema"`
	Serviser_ID int64 `json:"Serviser_ID"`
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

func GetWorkOrderByTehnicianID(serviserID int64) ([]RadniNalogSaKlijentom, error) {
	query :=  `SELECT 
            rn.Nalog_ID, 
            rn.BrojNaloga,
            k.Naziv AS Klijent, 
            rn.OpisProblema,
						dz.Serviser_ID
        FROM 
            RadniNalog rn
        INNER JOIN 
            DodjelaZadatka dz ON rn.Nalog_ID = dz.Nalog_ID
        INNER JOIN
            Klijenti k ON rn.Klijent_ID = k.Klijent_ID
        WHERE 
            dz.Serviser_ID = @Serviser_ID`
	
	rows, err := database.DB.Query(query, sql.Named("Serviser_ID", serviserID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workOrders []RadniNalogSaKlijentom

	for rows.Next() {
		var workOrder RadniNalogSaKlijentom
		err := rows.Scan(
			&workOrder.Nalog_ID,
			&workOrder.BrojNaloga,
			&workOrder.Klijent,
			&workOrder.OpisProblema,
			&workOrder.Serviser_ID,
		)
		if err != nil {
			return nil, err
		}
		workOrders = append(workOrders, workOrder)
	}

	return workOrders, nil
}

func GetWorkOrderByTehnician(serviserID int64) ([]RadniNalogServiser, error) {
	query := `
        SELECT 
    rn.Nalog_ID, 
    rn.Klijent_ID,
    rn.OpisProblema,
    rn.Prioritet,
    rn.DatumOtvaranja,
    rn.Status,
    rn.Lokacija,
    rn.BrojNaloga,
    dz.Serviser_ID
FROM 
    RadniNalog rn
INNER JOIN 
    DodjelaZadatka dz ON rn.Nalog_ID = dz.Nalog_ID
INNER JOIN
    Klijenti k ON rn.Klijent_ID = k.Klijent_ID
WHERE 
    dz.Serviser_ID = @Serviser_ID
ORDER BY 
    CASE 
        WHEN rn.Status = 'Otvoren' THEN 0
        WHEN rn.Status = 'Zavrsen' THEN 1
        ELSE 2
    END;

    `
	rows, err := database.DB.Query(query, sql.Named("Serviser_ID", serviserID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workOrders []RadniNalogServiser

	for rows.Next() {
		var workOrder RadniNalogServiser
		err := rows.Scan(
			&workOrder.Nalog_ID,
			&workOrder.Klijent_ID,
			&workOrder.OpisProblema,
			&workOrder.Prioritet,
			&workOrder.DatumOtvaranja,
			&workOrder.Status,
			&workOrder.Lokacija,
			&workOrder.BrojNaloga,
			&workOrder.Serviser_ID,
		)
		if err != nil {
			return nil, err
		}
		workOrders = append(workOrders, workOrder)
	}

	return workOrders, nil
}


func GetDetailedOrderForTehnician(serviserID int64) ([]DetaljiRadnogNaloga, error) {
	query := `SELECT 
                rn.OpisProblema, 
                k.Naziv, 
                rn.Lokacija, 
                dz.DatumDodjele, 
                dz.DatumZavrsetka, 
                es.BrojRadnihSati
            FROM 
                RadniNalog rn 
            JOIN 
                Klijenti k ON k.Klijent_ID = rn.Klijent_ID
            JOIN 
                DodjelaZadatka dz ON dz.Nalog_ID = rn.Nalog_ID 
            JOIN 
                EvidencijaSati es ON es.Nalog_ID = rn.Nalog_ID 
            JOIN 
                Serviser s ON s.Serviser_ID = es.Serviser_ID
            WHERE 
                s.Serviser_ID = @Serviser_ID`

	rows, err := database.DB.Query(query, sql.Named("Serviser_ID", serviserID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var detalji []DetaljiRadnogNaloga

	for rows.Next() {
		var (
			opisProblema   sql.NullString
			nazivKlijenta  sql.NullString
			lokacija       sql.NullString
			datumDodjele   sql.NullTime
			datumZavrsetka sql.NullTime
			brojRadnihSati sql.NullFloat64
		)

		err := rows.Scan(
			&opisProblema,
			&nazivKlijenta,
			&lokacija,
			&datumDodjele,
			&datumZavrsetka,
			&brojRadnihSati,
		)
		if err != nil {
			return nil, err
		}

		detalji = append(detalji, DetaljiRadnogNaloga{
			OpisProblema:   opisProblema.String,
			NazivKlijenta:  nazivKlijenta.String,
			Lokacija:       lokacija.String,
			DatumDodjele:   datumDodjele.Time,
			DatumZavrsetka: datumZavrsetka.Time,
			BrojRadnihSati: brojRadnihSati.Float64,
		})
	}

	return detalji, nil
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