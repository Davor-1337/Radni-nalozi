package models

import (
	"Diplomski/database"
	"database/sql"
	"time"
)

type Notification struct {
	ID           int64     `json:"obavestenje_id"`
	UserID       int64     `json:"korisnik_id"`
	Type         string    `json:"tip"`
	Message      string    `json:"poruka"`
	CreatedAt    time.Time `json:"datum_kreiranja"`
	IsRead       bool      `json:"procitano"`
}

func (n *Notification) Create() error {
	const query = `
			INSERT INTO Obavjestenja (Korisnik_ID, Tip, Poruka)
			VALUES (@UserID, @Type, @Message)`
	_, err := database.DB.Exec(query,
			sql.Named("UserID", n.UserID),
			sql.Named("Type", n.Type),
			sql.Named("Message", n.Message),
	)
	return err
}

func DoesNotificationBelongToUser(userID, notificationID int64) (bool, error) {
	var exists bool
	

	

	query := `
	SELECT CASE 
			WHEN EXISTS (
					SELECT 1 FROM Obavjestenja 
					WHERE Obavjestenje_ID = @NotificationID AND Korisnik_ID = @UserID
			) THEN 1
			ELSE 0
	END`
	

	
	err := database.DB.QueryRow(query,
			sql.Named("NotificationID", notificationID),
			sql.Named("UserID", userID),
	).Scan(&exists)
	

	if err != nil {
			
			return false, err
	}
	
	return exists, nil
}

func (n *Notification) Delete() error {
	const query = `
		DELETE FROM Obavjestenja 
		WHERE Obavjestenje_ID = @NotificationID`
	
	_, err := database.DB.Exec(query,
		sql.Named("NotificationID", n.ID),
	)
	return err
}

func GetNotificationsForUser(userID int64) ([]Notification, error) {
	const query = `
			SELECT Obavjestenje_ID, Korisnik_ID, Tip, Poruka, DatumKreiranja
				FROM Obavjestenja
			 WHERE Korisnik_ID = @UserID
			 ORDER BY DatumKreiranja DESC`
	rows, err := database.DB.Query(query, sql.Named("UserID", userID))
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	var list []Notification
	for rows.Next() {
			var n Notification
			if err := rows.Scan(
					&n.ID, &n.UserID, &n.Type, &n.Message,
					&n.CreatedAt,
			); err != nil {
					return nil, err
			}
			list = append(list, n)
	}
	return list, nil
}

