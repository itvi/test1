package mysqlite

import (
	"ams/pkg/models"
	"database/sql"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct{ DB *sql.DB }

func (m *UserModel) Insert(u *models.User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	q := `INSERT INTO app_user(sn,name,email,hashed_password)
		VALUES(?,?,?,?)`
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.SN, u.Name, u.Email, string(hashedPassword))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return models.ErrDuplicate
		}
	}
	return err
}

// Authenticate verify where a user exist with the user sn and password
// This will return the relevant user struct
func (m *UserModel) Authenticate(sn, password string) (*models.User, error) {
	user := &models.User{}
	row := m.DB.QueryRow(`SELECT id,hashed_password FROM app_user WHERE 
sn=?`, sn)
	err := row.Scan(&user.ID, &user.HashedPassword)
	if err == sql.ErrNoRows {
		return nil, models.ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, models.ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	q := "SELECT id,sn,name,email FROM app_user WHERE id=?"
	u := &models.User{}

	err := m.DB.QueryRow(q, id).Scan(&u.ID, &u.SN, &u.Name, &u.Email)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return u, nil
}
