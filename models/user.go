package models

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

func CreateUser(db *sql.DB, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	_, err = db.Exec("INSERT INTO users (name, email, password, phone, address) VALUES (?, ?, ?, ?, ?)",
		user.Name, user.Email, user.Password, user.Phone, user.Address)
	return err
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	row := db.QueryRow("SELECT id, name, email, password, phone, address FROM users WHERE email = ?", email)
	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Address)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, name, email, phone, address FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Address)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func GetUserByID(db *sql.DB, id int) (*User, error) {
	row := db.QueryRow("SELECT id, name, email, password, phone, address FROM users WHERE id = ?", id)
	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Address)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser memperbarui informasi user berdasarkan ID
func UpdateUser(db *sql.DB, user *User) error {
	_, err := db.Exec("UPDATE users SET name = ?, email = ?, password = ?, phone = ?, address = ? WHERE id = ?",
		user.Name, user.Email, user.Password, user.Phone, user.Address, user.ID)
	return err
}

// DeleteUser menghapus user berdasarkan ID
func DeleteUser(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
