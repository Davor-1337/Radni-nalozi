package models

import (
	"Diplomski/database"
	"Diplomski/utils"
	"database/sql"
	"errors"
)


type User struct {
	User_ID int64
	Username string 
	Email string 
	Password string 
	Role string
}

func (u User) Save() error {
	query := `INSERT INTO Korisnici(User_ID, Username, Email, Password, Role)
						VALUES(@User_ID, @Username, @Email, @Password, @Role)`
	
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		
	}

	result, err := stmt.Exec(
		sql.Named("User_ID", u.User_ID),
		sql.Named("Username", u.Username),
		sql.Named("Email", u.Email),
		sql.Named("Password", hashedPassword),
		sql.Named("Role", u.Role))


	userId, err := result.LastInsertId()
	u.User_ID = userId

	return nil
}

func (u *User) ValidateCredentials() error {
	query := "SELECT User_ID, Role, Password FROM Korisnici WHERE Username = @Username"


	row := database.DB.QueryRow(query, sql.Named("Username", u.Username))

	var retrievedPassword string
	err := row.Scan(&u.User_ID, &u.Role, &retrievedPassword)

	if err != nil {
		return errors.New("Credentials Invalid")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)

	if !passwordIsValid {
		return errors.New("Credentials Invalid")
	}

	return nil
}

