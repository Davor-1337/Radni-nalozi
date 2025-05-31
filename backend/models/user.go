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

func GetUserById(id int64) (*User, error) {
			query := "SELECT User_ID, Username, Email, Password, Role FROM Korisnici WHERE User_ID = @User_ID"
	
	
	row := database.DB.QueryRow(query, sql.Named("User_ID", id))

	var users User

	
	err := row.Scan(&users.User_ID, &users.Username, &users.Email, &users.Password, &users.Role )
	if err != nil {
		return nil, err
	}
 
	return &users, nil
	}

func GetUserByIdForPassword(id int64) (*User, error) {
		query := "SELECT User_ID, Username, Password, Role FROM Korisnici WHERE User_ID = @User_ID"


row := database.DB.QueryRow(query, sql.Named("User_ID", id))

var users User


err := row.Scan(&users.User_ID, &users.Username, &users.Password, &users.Role )
if err != nil {
	return nil, err
}

return &users, nil
}

func (u User)ChangePassword(id int64, password string) (error) {
		query := `UPDATE Korisnici SET Password = @Password WHERE User_ID = @User_ID`
		
		stmt, err := database.DB.Prepare(query)
	if err != nil {
		return  err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		sql.Named("Password", password),
		sql.Named("User_ID", id))
		
	
		return err
 
	}

