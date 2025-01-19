package models

import (
	"Diplomski/database"
	"database/sql"
	"log"

	// "database/sql"
	"time"
)

type EvidencijaSati struct {
	Evidencija_ID int64
	Nalog_ID        int64
	Serviser_ID     int64
	BrojRadnihSati float64
	DatumUnosa     time.Time
}

type UtroseniMaterijal struct{
	UtroseniMaterijal_ID int64
	Nalog_ID int64
	Materijal_ID int64
	KolicinaUtrosena int64
	DatumUnosa time.Time
}

type RadniNalog struct {
	Nalog_ID       int64
	Klijent_ID     int64
	OpisProblema   string
	Prioritet      string
	DatumOtvaranja time.Time
	Status string
	Lokacija string
	BrojNaloga string
}

func GetAllWorkOrders() ([]RadniNalog, error) {
	
	query := "SELECT * FROM RadniNalog"
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workOrders []RadniNalog

	for rows.Next(){
		var workOrder RadniNalog 
		err := rows.Scan(&workOrder.Nalog_ID, &workOrder.Klijent_ID, &workOrder.OpisProblema, &workOrder.Prioritet, &workOrder.DatumOtvaranja, &workOrder.Status, &workOrder.Lokacija, &workOrder.BrojNaloga)
		if err != nil{
			return nil, err
		}

		workOrders = append(workOrders, workOrder)
	}
	return workOrders, nil
  }
	 
	func GetWorkOrderByID(id int64) (*RadniNalog, error) {
			query := "SELECT * FROM RadniNalog WHERE Nalog_ID = @Nalog_ID"
	
	
	row := database.DB.QueryRow(query, sql.Named("Nalog_ID", id))

	var workOrder RadniNalog

	
	err := row.Scan(&workOrder.Nalog_ID, &workOrder.Klijent_ID, &workOrder.OpisProblema, &workOrder.Prioritet, &workOrder.DatumOtvaranja, &workOrder.Status, &workOrder.Lokacija, &workOrder.BrojNaloga )
	if err != nil {
		return nil, err
	}
 
	return &workOrder, nil
	}

func (wo *RadniNalog) InsertWorkOrder() error {
    query := `INSERT INTO RadniNalog(Klijent_ID, OpisProblema, Prioritet, DatumOtvaranja, Status, Lokacija)
              OUTPUT INSERTED.Nalog_id
              VALUES (@Klijent_ID, @OpisProblema, @Prioritet, @DatumOtvaranja, @Status, @Lokacija)`

   
    stmt, err := database.DB.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    
    err = stmt.QueryRow(
        sql.Named("Klijent_ID", wo.Klijent_ID),
        sql.Named("OpisProblema", wo.OpisProblema),
        sql.Named("Prioritet", wo.Prioritet),
        sql.Named("DatumOtvaranja", wo.DatumOtvaranja),
        sql.Named("Status", wo.Status),
				sql.Named("Lokacija", wo.Lokacija),
    ).Scan(&wo.Nalog_ID)

    if err != nil {
        return err
    }

    return nil
}


func (workOrder RadniNalog) Update() error {
	query := `UPDATE RadniNalog 
						SET Klijent_ID = @Klijent_ID,
						OpisProblema = @OpisProblema, 
						Prioritet = @Prioritet, 
						DatumOtvaranja = @DatumOtvaranja, 
						Status = @Status, 
						Lokacija = @Lokacija
						WHERE Nalog_ID = @Nalog_ID`

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	
	_, err = stmt.Exec(
        sql.Named("Klijent_ID", workOrder.Klijent_ID),
        sql.Named("OpisProblema", workOrder.OpisProblema),
        sql.Named("Prioritet", workOrder.Prioritet),
        sql.Named("DatumOtvaranja", workOrder.DatumOtvaranja),
        sql.Named("Status", workOrder.Status),
				sql.Named("Lokacija", workOrder.Lokacija),
				sql.Named("Nalog_ID", workOrder.Nalog_ID),)

	return err
}

func (workOrder RadniNalog) Delete() error {
	query := "DELETE FROM RadniNalog WHERE Nalog_ID = @Nalog_ID"
	

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		log.Println("Error preparing query:", err) 
		return err
	}
	defer func() {
	
		stmt.Close()
	}()

	_, err = stmt.Exec(sql.Named("Nalog_ID", workOrder.Nalog_ID))
	if err != nil {
		return err
	}

	
	return nil
}


func Finish(radniNalogID int64, serviserID int64) error {
	query := `UPDATE DodjelaZadatka
						SET DatumZavrsetka = GETDATE()
						WHERE Nalog_ID = @Nalog_ID AND Serviser_ID = @Serviser_ID AND DatumZavrsetka IS NULL`
	
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
					sql.Named("Nalog_ID", radniNalogID),
					sql.Named("Serviser_ID", serviserID))		
	if err != nil {
		return err
	}	

	return nil
}

func (um UtroseniMaterijal) InputMaterial(nalogID int64) error {
	query := `INSERT INTO UtroseniMaterijal (Nalog_ID, Materijal_ID, KolicinaUtrosena, DatumUnosa)
						VALUES (@Nalog_ID, @Materijal_ID, @KolicinaUtrosena, GETDATE())`
	
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
				sql.Named("Nalog_ID", nalogID),
				sql.Named("Materijal_ID", um.Materijal_ID),
				sql.Named("KolicinaUtrosena", um.KolicinaUtrosena))	
	if err != nil {
		return err
	}
	return nil
}

func (es EvidencijaSati) AssignTask() error {
	query := `INSERT INTO DodjelaZadatka (Nalog_ID, Serviser_ID, DatumDodjele)
						VALUES (@Nalog_ID, @Serviser_ID, GETDATE())`
	
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return  err
	}
	defer stmt.Close()

	 _, err = stmt.Exec(sql.Named("Nalog_ID", es.Nalog_ID),
									sql.Named("Serviser_ID", es.Serviser_ID))	
	if err != nil {
		return err
	}		

	return nil
}

func CheckTaskAssignment(nalogId, serviserId int64) (bool, error) {
	query := `SELECT COUNT(*)
						FROM DodjelaZadatka
						WHERE Nalog_ID = @Nalog_ID AND Serviser_ID = @Serviser_ID`
	
	var count int 
	err := database.DB.QueryRow(query, sql.Named("Nalog_ID", nalogId), sql.Named("Serviser_ID", serviserId)).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (e *EvidencijaSati) InsertHours(nalogId int64) error {
	query := `INSERT INTO EvidencijaSati (Nalog_ID, Serviser_ID, BrojRadnihSati, DatumUnosa)
						VALUES (@Nalog_ID, @Serviser_ID, @BrojRadnihSati, GETDATE())`

	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return nil
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		sql.Named("Nalog_ID", nalogId),
		sql.Named("Serviser_ID", e.Serviser_ID),
		sql.Named("BrojRadnihSati", e.BrojRadnihSati))
	
	if err != nil {
		return err
	}

	return nil
}

func GetHours(id int64) (*EvidencijaSati, error) {
query := `SELECT * FROM EvidencijaSati WHERE Nalog_ID = @Nalog_ID`

row := database.DB.QueryRow(query, sql.Named("Nalog_ID", id))

var hours EvidencijaSati

err := row.Scan(&hours.Evidencija_ID, &hours.Nalog_ID, &hours.Serviser_ID, &hours.BrojRadnihSati, &hours.DatumUnosa)
if err != nil {
	return nil, err
}

return &hours, nil
}

func IsTehnicianAssignedToWorkOrder(tehnicianID, workOrderID int64) (bool, error) {
    query := `SELECT COUNT(*) FROM DodjelaZadatka WHERE Nalog_ID = @Nalog_ID AND Serviser_ID = @Serviser_ID`
    var count int
    err := database.DB.QueryRow(query, sql.Named("Nalog_ID", workOrderID), sql.Named("Serviser_ID", tehnicianID)).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func IsWorkOrderOwnedByClient(userID, workOrderID int64) (bool, error) {
    query := `SELECT COUNT(*) FROM RadniNalog WHERE Nalog_ID = @Nalog_ID AND Klijent_ID = @Klijent_ID`
    var count int
    err := database.DB.QueryRow(query, sql.Named("Nalog_ID", workOrderID), sql.Named("Klijent_ID", userID)).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func GetAllArchivedWorkOrders() ([]RadniNalog, error) {
	
	query := "SELECT Nalog_ID, Klijent_ID, OpisProblema, Prioritet, DatumOtvaranja, Status, Lokacija, BrojNaloga FROM ArhiviraniRadniNalozi"
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workOrders []RadniNalog

	for rows.Next(){
		var workOrder RadniNalog 
		err := rows.Scan(&workOrder.Nalog_ID, &workOrder.Klijent_ID, &workOrder.OpisProblema, &workOrder.Prioritet, &workOrder.DatumOtvaranja, &workOrder.Status, &workOrder.Lokacija, &workOrder.BrojNaloga)
		if err != nil{
			return nil, err
		}

		workOrders = append(workOrders, workOrder)
	}
	return workOrders, nil
  }

	func Archive(nalogId int64) error {
		query := "EXEC Arhiviranje @NalogID = @ID"
		
		_, err := database.DB.Exec(query, sql.Named("ID", nalogId))
		if err != nil {
			return err 
		}
		return nil
	}

