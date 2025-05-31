package models

import (
	"Diplomski/database"
	"database/sql"
	"errors"
	"fmt"
	"log"
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
	Arhiviran     bool
	
}

type PendingWorkOrderDTO struct {
	ID             int64       `json:"id"`
	BrojNaloga     string       `json:"brojNaloga"`
	Prioritet      string    `json:"prioritet"`
	Status         string    `json:"status"`
	Lokacija       string    `json:"lokacija"`
	DatumOtvaranja time.Time `json:"datumOtvaranja"`
	OpisProblema   string    `json:"opisProblema"`
	NazivKlijenta  string    `json:"nazivKlijenta"`
}

type WorkOrderActionRequest struct {
	ID     int    `json:"id"`
	Akcija string `json:"akcija"` 
}

func GetPendingOrders() ([]PendingWorkOrderDTO, error) {
    query := `SELECT rn.Nalog_ID, rn.BrojNaloga, rn.Prioritet, rn.Status, rn.Lokacija, 
                     rn.DatumOtvaranja, rn.OpisProblema, k.Naziv
              FROM RadniNalog rn
              JOIN Klijenti k on k.Klijent_ID = rn.Klijent_ID
              WHERE rn.Status = 'Na cekanju'`

    rows, err := database.DB.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to query pending orders: %w", err)
    }
    defer rows.Close()

    var workOrders []PendingWorkOrderDTO

    for rows.Next() {
        var workOrder PendingWorkOrderDTO
        if err := rows.Scan(
            &workOrder.ID,
            &workOrder.BrojNaloga,
            &workOrder.Prioritet,
            &workOrder.Status,
            &workOrder.Lokacija,
            &workOrder.DatumOtvaranja,
            &workOrder.OpisProblema,
            &workOrder.NazivKlijenta,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }
        workOrders = append(workOrders, workOrder)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error during rows iteration: %w", err)
    }

    return workOrders, nil
}

func GetAllWorkOrders() ([]RadniNalog, error) {
	
	query := "SELECT * FROM RadniNalog WHERE Arhiviran = 0"
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workOrders []RadniNalog

	for rows.Next(){
		var workOrder RadniNalog 
		err := rows.Scan(&workOrder.Nalog_ID, &workOrder.Klijent_ID, &workOrder.OpisProblema, &workOrder.Prioritet, &workOrder.DatumOtvaranja, &workOrder.Status, &workOrder.Lokacija, &workOrder.BrojNaloga, &workOrder.Arhiviran)
		if err != nil{
			return nil, err
		}

		workOrders = append(workOrders, workOrder)
	}
	return workOrders, nil
  }

	func GetTotalWorkOrderCount() (int, error) {
		query := `
				SELECT COUNT(*) FROM RadniNalog	
		`
	
		var total int
		err := database.DB.QueryRow(query).Scan(&total)
		if err != nil {
			return 0, err
		}
		return total, nil
	}
	

func GetWorkOrderStats() ([]struct {
    MonthNumber int
    Count       int
}, error) {
    query := `SELECT 
    MONTH(dz.DatumZavrsetka) AS MjesecBroj,
    COUNT(*) AS BrojNaloga
FROM 
    RadniNalog rn
JOIN 
    DodjelaZadatka dz ON dz.Nalog_ID = rn.Nalog_ID
WHERE 
    rn.Status = 'Zavrsen'
    AND YEAR(dz.DatumZavrsetka) = YEAR(GETDATE())
GROUP BY 
    MONTH(dz.DatumZavrsetka)
ORDER BY 
    MjesecBroj;`

    rows, err := database.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []struct {
        MonthNumber int
        Count       int
    }

    for rows.Next() {
        var monthNumber, count int
        if err := rows.Scan(&monthNumber, &count); err != nil {
            return nil, err
        }
        results = append(results, struct {
            MonthNumber int
            Count       int
        }{MonthNumber: monthNumber, Count: count})
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return results, nil
}

	func GetWorkOrderStatusCount() (int, int, error) {
		var completedCount int
		var inProgressCount int

		query := "SELECT COUNT(*) FROM RadniNalog WHERE Status = 'Otvoren'"
		err := database.DB.QueryRow(query).Scan(&inProgressCount)
		if err != nil {
			return 0, 0, err
		}

		query = `
		SELECT COUNT(*) FROM RadniNalog WHERE Status = 'Zavrsen'`
	
		err = database.DB.QueryRow(query).Scan(&completedCount)
		if err != nil {
			return 0, 0, err
		}

		return completedCount, inProgressCount, nil
		}

	
func Get4WorkOrders() ([]map[string]interface{}, error) {
    query := `SELECT TOP 4  rn.BrojNaloga, rn.Prioritet, rn.Status, rn.Lokacija, rn.DatumOtvaranja, rn.OpisProblema,
k.Naziv
FROM RadniNalog rn
JOIN Klijenti k on k.Klijent_ID = rn.Klijent_ID
ORDER BY Nalog_ID DESC`
    rows, err := database.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var workOrders []map[string]interface{}

    for rows.Next() {
        var brojNaloga, prioritet, status, lokacija, opisProblema, nazivKlijenta string
        var datumOtvaranja time.Time

        err := rows.Scan(&brojNaloga, &prioritet, &status, &lokacija, &datumOtvaranja, &opisProblema, &nazivKlijenta)
        if err != nil {
            return nil, err
        }

        workOrder := map[string]interface{}{
            "BrojNaloga":    brojNaloga,
            "Prioritet":     prioritet,
            "Status":        status,
            "Lokacija":      lokacija,
            "DatumOtvaranja": datumOtvaranja,
            "OpisProblema":  opisProblema,
            "NazivKlijenta": nazivKlijenta, 
        }

        workOrders = append(workOrders, workOrder)
    }
    return workOrders, nil
}

func GetActiveWorkOrders() ([]map[string]interface{}, error) {
    query := `SELECT rn.BrojNaloga, rn.Prioritet, rn.Status, rn.Lokacija, rn.DatumOtvaranja, rn.OpisProblema,
k.Naziv
FROM RadniNalog rn
JOIN Klijenti k on k.Klijent_ID = rn.Klijent_ID
WHERE Status = 'Otvoren'
ORDER BY DatumOtvaranja DESC`
    rows, err := database.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var workOrders []map[string]interface{}

    for rows.Next() {
        var brojNaloga, prioritet, status, lokacija, opisProblema, nazivKlijenta string
        var datumOtvaranja time.Time

        err := rows.Scan(&brojNaloga, &prioritet, &status, &lokacija, &datumOtvaranja, &opisProblema, &nazivKlijenta)
        if err != nil {
            return nil, err
        }

        workOrder := map[string]interface{}{
            "BrojNaloga":    brojNaloga,
            "Prioritet":     prioritet,
            "Status":        status,
            "Lokacija":      lokacija,
            "DatumOtvaranja": datumOtvaranja,
            "OpisProblema":  opisProblema,
            "NazivKlijenta": nazivKlijenta, 
        }

        workOrders = append(workOrders, workOrder)
    }
    return workOrders, nil
}
	 
	func GetWorkOrderByID(id int64) (*RadniNalog, error) {
			query := "SELECT Nalog_ID, Klijent_ID, OpisProblema, Prioritet, DatumOtvaranja, Status, Lokacija, BrojNaloga FROM RadniNalog WHERE Nalog_ID = @Nalog_ID"
	
	
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
              VALUES (@Klijent_ID, @OpisProblema, @Prioritet, GETDATE(), @Status, @Lokacija)`

   
    stmt, err := database.DB.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    
    err = stmt.QueryRow(
        sql.Named("Klijent_ID", wo.Klijent_ID),
        sql.Named("OpisProblema", wo.OpisProblema),
        sql.Named("Prioritet", wo.Prioritet),
        sql.Named("Status", wo.Status),
				sql.Named("Lokacija", wo.Lokacija),
    ).Scan(&wo.Nalog_ID)

    if err != nil {
        return err
    }

    return nil
}

func (r *RadniNalog) GetClientID() error {
	return database.DB.QueryRow(
			"SELECT Klijent_ID FROM RadniNalog WHERE Nalog_ID = @ID",
			sql.Named("ID", r.Nalog_ID),
	).Scan(&r.Klijent_ID)
}

func (workOrder *RadniNalog) UpdateStatus() error {
	if workOrder.Status == "" {
		return errors.New("status ne smije biti prazan")
	}

	query := `
		UPDATE RadniNalog 
		SET Status = @Status
		WHERE Nalog_ID = @Nalog_ID
	`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		sql.Named("Status", workOrder.Status),
		sql.Named("Nalog_ID", workOrder.Nalog_ID),
	)
	return err
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
	tx, err := database.DB.Begin()
	if err != nil {
			return err
	}
	defer func() {
			if err != nil {
					tx.Rollback()
			} else {
					tx.Commit()
			}
	}()

	
	updateDateQuery := `
			UPDATE DodjelaZadatka
			SET DatumZavrsetka = GETDATE()
			WHERE Nalog_ID = @Nalog_ID AND Serviser_ID = @Serviser_ID AND DatumZavrsetka IS NULL`
	if _, err = tx.Exec(updateDateQuery,
			sql.Named("Nalog_ID", radniNalogID),
			sql.Named("Serviser_ID", serviserID),
	); err != nil {
			return err
	}

	
	updateStatusQuery := `
			UPDATE RadniNalog
			SET Status = 'Zavrsen'
			WHERE Nalog_ID = @Nalog_ID`
	if _, err = tx.Exec(updateStatusQuery,
			sql.Named("Nalog_ID", radniNalogID),
	); err != nil {
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
	
	query := "SELECT Nalog_ID, Klijent_ID, OpisProblema, Prioritet, DatumOtvaranja, Status, Lokacija, BrojNaloga FROM RadniNalog WHERE Arhiviran = 1"
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

