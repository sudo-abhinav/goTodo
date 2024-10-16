package dbHelper

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sudo-abhinav/go-todo/Database"
	"github.com/sudo-abhinav/go-todo/utils/encryption"
	"time"
)

func CreateUserInDB(UserName, email, password string) error {
	//here i am creating hash password for security reason

	hashPWD, err := encryption.HashPassword(password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	queryString := `INSERT INTO users (username, email, password) VALUES ($1, trim(lower($2)), $3)`
	_, err = Database.DBconn.Exec(queryString, UserName, email, hashPWD)

	if err != nil {
		return err
	}
	return nil

}
func CreateUserSession(UserId string) (string, error) {
	var sessionID string
	query := `INSERT INTO user_session(user_id) 
              VALUES ($1)`
	CusErr := Database.DBconn.Get(&sessionID, query, UserId)

	if CusErr != nil {
		return "", CusErr
	}

	return sessionID, nil
}

func GetUser(email, password string) (string, string, error) {

	var storedId string
	var hashPassword string
	var storedEmail string

	QueryString := `SELECT u.id, u.email, u.password FROM users u
			  WHERE u.archived_at IS NULL
			    AND u.email = TRIM($1)`

	var results []struct {
		ID       string
		Email    string
		Password string
	}

	// TODO use get because you are not getting an array
	err := Database.DBconn.Select(&results, QueryString, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No matching user row found
			return "", "", nil
		}
		return "", "", err
	}
	// Check if we got a result then we check
	if len(results) == 0 {
		return "", "", nil
	}

	// Extract the user data from the results
	storedId = results[0].ID
	storedEmail = results[0].Email
	hashPassword = results[0].Password

	// Verify the password
	if err := encryption.VerifyPassword(password, hashPassword); err != nil {
		// Password verification failed
		return "", "", err
	}

	var PasswordValidate = encryption.VerifyPassword(password, hashPassword)
	if PasswordValidate != nil {
		return "", "", err
	}

	return storedId, storedEmail, nil
}

func IsUserExists(email string) (bool, error) {
	SQL := `SELECT count(id) > 0 as is_exist
			  FROM users
			  WHERE email = TRIM($1)
			    AND archived_at IS NULL`

	var Usercheck bool
	Err := Database.DBconn.Get(&Usercheck, SQL, email)
	return Usercheck, Err
}

/*
func IsUserExists(email string) (bool, error) {
	// language=SQL
	SQL := `SELECT id FROM users WHERE email = TRIM(LOWER($1)) AND archived_at IS NULL`
	var id string
	err := database.Todo.Get(&id, SQL, email)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}
*/

func DeleteUser(userID string) error {
	query := `UPDATE users SET
                    archived_at=now() 
                		WHERE  id= $1
                		  AND archived_at IS NULL`
	//res, err := Database.DBconn.Exec("DELETE FROM usertodo WHERE id=$1", param.Id)
	_, Err := Database.DBconn.Exec(query, userID)
	if Err != nil {
		return Err
	}
	return nil

}
func UserSessionDelete(SessionId string) error {
	query := `UPDATE user_session SET
                    archived_at=now() 
                		WHERE  id= $1
                		  AND archived_at IS NULL`
	_, Err := Database.DBconn.Exec(query, SessionId)
	if Err != nil {
		return Err
	}
	return nil

}

func GetArchivedAt(sessionID string) (*time.Time, error) {
	var archivedAt *time.Time

	query := `SELECT archived_at FROM user_session WHERE id = $1`

	//TODO use get not queryRowx
	getErr := Database.DBconn.QueryRow(query, sessionID).Scan(&archivedAt)
	if getErr != nil {
		return nil, getErr // Return error if the query fails
	}

	return archivedAt, nil
}
