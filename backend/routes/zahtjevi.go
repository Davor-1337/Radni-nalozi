package routes

import (
	"Diplomski/database"
	"Diplomski/models"
	"database/sql"
	"fmt"
	
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func createZahtjev(context *gin.Context) {
	var zahtjev models.Zahtjev

	if err := context.ShouldBindJSON(&zahtjev); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := zahtjev.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Greška prilikom čuvanja zahtjeva."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Zahtjev uspješno poslan."})
}


func GetAllZahtjevi(context *gin.Context) {
	rows, err := database.DB.Query(`SELECT Zahtjev_ID, username, role, status, VrijemeKreiranja FROM Zahtjevi WHERE status = 'na cekanju'`)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse requested data."})
		return
	}
	defer rows.Close()

	var zahtjevi []models.Zahtjev
	for rows.Next() {
		var z models.Zahtjev
		err := rows.Scan(&z.Zahtjev_ID, &z.Username, &z.Role, &z.Status, &z.VrijemeKreiranja)
		if err != nil {
			continue
		}
		zahtjevi = append(zahtjevi, z)
	}

	context.JSON(http.StatusOK, zahtjevi)
}


func handleRequest(c *gin.Context) {
	var payload struct {
		ZahtjevID int    `json:"zahtjev_id"`
		Akcija    string `json:"akcija"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload."})
		return
	}

	if payload.Akcija != "prihvati" && payload.Akcija != "odbij" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid action."})
		return
	}

	
	var zahtjev models.Zahtjev
	query := `
	SELECT zahtjev_id, username, email, password, role,
	       ime, prezime, specijalnost, telefon,
	       naziv, kontaktOsoba, adresa
	FROM Zahtjevi
	WHERE Zahtjev_ID = @zahtjev_id
`
err := database.DB.QueryRow(
	query,
	sql.Named("zahtjev_id", payload.ZahtjevID),
).Scan(
	&zahtjev.Zahtjev_ID,
	&zahtjev.Username,
	&zahtjev.Email,
	&zahtjev.Password,
	&zahtjev.Role,
	&zahtjev.Ime,
	&zahtjev.Prezime,
	&zahtjev.Specijalnost,
	&zahtjev.Tel,
	&zahtjev.Naziv,
	&zahtjev.KontaktOsoba,
	&zahtjev.Adresa,
)
if err != nil {
	c.JSON(http.StatusNotFound, gin.H{"message": "Request not found."})
	return
}
fmt.Println(zahtjev.Role)
	
	if payload.Akcija == "prihvati" {
		if err := createUser(zahtjev); err != nil {
			fmt.Println("Greška kod kreiranja korisnika:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Greška kod kreiranja korisnika."})
			return
		}
	}

	
	newStatus := "odbijen"
	if payload.Akcija == "prihvati" {
		newStatus = "prihvacen"
	}
	_, err = database.DB.Exec(
		`UPDATE Zahtjevi SET status = @status WHERE Zahtjev_ID = @zahtjev_id`,
		sql.Named("status", newStatus),
		sql.Named("zahtjev_id", payload.ZahtjevID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Greška kod ažuriranja statusa."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Zahtjev obrađen uspješno."})
}



func getNextUserID(role string) (int, error) {
    var lastID int
    var query string

    if role == "serviser" {
        query = `SELECT MAX(user_id) FROM Korisnici WHERE user_id >= 100 AND user_id < 200`
    } else if role == "klijent" {
        query = `SELECT MAX(user_id) FROM Korisnici WHERE user_id >= 200 AND user_id < 300`
    } else {
        return 0, fmt.Errorf("invalid role")
    }

    err := database.DB.QueryRow(query).Scan(&lastID)
    if err != nil {
        return 0, err
    }

 
    if lastID == 0 {
        if role == "serviser" {
            return 105, nil  
        } else if role == "klijent" {
            return 205, nil  
        }
    }

    
    return lastID + 1, nil
}

func createUser(zahtjev models.Zahtjev) error {
	fmt.Println("ddd",zahtjev.Role)
	zahtjev.Role = strings.ToLower(zahtjev.Role)
	nextID, err := getNextUserID(zahtjev.Role)
	if err != nil {
			return err
	}

	id64 := int64(nextID)
	

	
	_, err = database.DB.Exec(
			`INSERT INTO Korisnici (user_id, username, password, role, email)
			 VALUES (@user_id, @username, @password, @role, @email)`,
			sql.Named("user_id",  nextID),
			sql.Named("username",  zahtjev.Username),
			sql.Named("password",  zahtjev.Password),
			sql.Named("role",      zahtjev.Role),
			sql.Named("email",     zahtjev.Email),
	)
	if err != nil {
			return fmt.Errorf("greška pri unosu korisnika: %w", err)
	}

	
	if zahtjev.Role == "serviser" {
			s := models.Serviser{
					Serviser_ID:  id64,
					Ime:          zahtjev.Ime,
					Prezime:      zahtjev.Prezime,
					Specijalnost: zahtjev.Specijalnost,
					Telefon:      zahtjev.Tel,
					User_ID:      id64,
			}
			if err := s.InsertTehnician(); err != nil {
					return fmt.Errorf("greška pri unosu servisera: %w", err)
			}
	}

	
	if zahtjev.Role == "klijent" {
			c := models.Klijent{
					Klijent_ID:     id64,
					Naziv:          zahtjev.Naziv,
					KontaktOsoba:   zahtjev.KontaktOsoba,
					Email:          zahtjev.Email,
					Tel:             zahtjev.Tel,
					Adresa:         zahtjev.Adresa,
					User_ID:        id64,
			}
			if err := c.InsertClient(); err != nil {
					return fmt.Errorf("greška pri unosu klijenta: %w", err)
			}
	}

	return nil
}
