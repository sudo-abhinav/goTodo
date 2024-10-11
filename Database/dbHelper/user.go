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
	queryString := `INSERT INTO users (username, email, password) VALUES ($1, lower($2), $3)`
	_, err = Database.DBconn.Exec(queryString, UserName, email, hashPWD)

	if err != nil {
		return err
	}
	return nil

}
func CreateUserSession(UserId string) (string, error) {
	var sessionID string
	query := `INSERT INTO user_session(user_id) 
              VALUES ($1) RETURNING id`
	CusErr := Database.DBconn.QueryRow(query, UserId).Scan(&sessionID)

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
	err := Database.DBconn.QueryRowx(QueryString, email).Scan(&storedId, &storedEmail, &hashPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", nil // Return nil if no matching user is found
		}
		return "", "", err
	}

	var PasswordValidate = encryption.VerifyPassword(password, hashPassword)
	if PasswordValidate != nil {
		return "", "", err
	}

	return storedId, storedEmail, nil
}

func IsUserExists(email string) (bool, error) {
	SQLQuery := `SELECT email from users WHERE email = TRIM(LOWER($1))`
	var UserExisting bool
	err := Database.DBconn.Get(&UserExisting, SQLQuery, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	//if errors.Is(err, sql.ErrNoRows) {
	//	return false, nil
	//}
	return UserExisting, nil
}

func DeleteUser() {

}

func GetArchivedAt(sessionID string) (*time.Time, error) {
	var archivedAt *time.Time

	query := `SELECT archived_at FROM user_session WHERE id = $1`

	getErr := Database.DBconn.QueryRow(query, sessionID).Scan(&archivedAt)
	if getErr != nil {
		return nil, getErr // Return error if the query fails
	}

	return archivedAt, nil
}
