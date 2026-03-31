package repository

import (
	"errors"

	"github.com/Ayan25844/netflix/dto"
	"github.com/Ayan25844/netflix/model"
	"github.com/Ayan25844/netflix/properties"
	"github.com/lib/pq"
)

var roles pq.StringArray

// Login

func GetUserByEmail(name string) (model.User, error) {
	sqlStatement := `SELECT * FROM users WHERE email = $1`
	var user model.User
	err := properties.Db.QueryRow(sqlStatement, name).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &roles)
	user.Role = []string(roles)
	if err != nil {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}

// Insert 1 record

func InsertOneUser(user model.User) (model.User, error) {
	sqlStatement := `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id`
	var id string
	err := properties.Db.QueryRow(sqlStatement, user.Name, user.Email, user.Password, pq.Array(user.Role)).Scan(&id)
	if err != nil {
		return model.User{}, err
	}
	user.ID = id
	return user, nil
}

// Update 1 record

func UpdateOneUser(id string, payload dto.Payload) (model.User, error) {
	sqlStatement := `
        UPDATE users 
        SET 
            name = COALESCE(NULLIF($2, ''), name), 
            email = COALESCE(NULLIF($3, ''), email),
            password = COALESCE(NULLIF($4, ''), password)
        WHERE id = $1 
        RETURNING id, name, email, password, role`
	var user model.User
	err := properties.Db.QueryRow(sqlStatement, id, payload.Name, payload.Email, payload.Password).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &roles)
	user.Role = []string(roles)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// Delete 1 record

func DeleteOneUser(userId string) (string, error) {
	sqlStatement := `DELETE FROM users WHERE id = $1 RETURNING id`
	var id string
	err := properties.Db.QueryRow(sqlStatement, userId).Scan(&id)
	if err != nil {
		return "", err
	}
	return "User with id: " + id + " deleted", nil
}

// Delete all records

func DeleteAllRecords() (string, error) {
	sqlStatement := `TRUNCATE TABLE users`
	_, err := properties.Db.Exec(sqlStatement)
	if err != nil {
		return "", err
	}
	return "All records deleted from users table", nil
}

// Get all users from database

func GetAll() ([]model.User, error) {
	sqlStatement := `SELECT * FROM users`
	rows, err := properties.Db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var m model.User
		err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Password, &roles)
		m.Role = []string(roles)
		if err != nil {
			return nil, err
		}
		users = append(users, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Get user by id

func GetById(id string) (model.User, error) {
	sqlStatement := `SELECT * from users where id = $1`
	var user model.User
	err := properties.Db.QueryRow(sqlStatement, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &roles)
	user.Role = []string(roles)
	if err != nil {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}
