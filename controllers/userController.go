package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"myapp/models"
)

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var input RegisterInput

	// Cek apakah data dikirim dalam bentuk JSON atau form-data
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
	} else if contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data" {
		input.Name = r.FormValue("name")
		input.Email = r.FormValue("email")
		input.Password = r.FormValue("password")
		input.Phone = r.FormValue("phone")
		input.Address = r.FormValue("address")

		// Validasi input form-data
		if input.Name == "" || input.Email == "" || input.Password == "" || input.Phone == "" || input.Address == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Unsupported content type", http.StatusBadRequest)
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Phone:    input.Phone,
		Address:  input.Address,
	}

	if err := models.CreateUser(db, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User registered successfully")
}

// Login untuk menangani autentikasi user
func Login(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validasi input
	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByEmail(db, input.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Periksa password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode("Login successful")
}

// GetUsers untuk mengambil semua user
func GetUsers(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// UpdateUser untuk menghandle permintaan update data user
func UpdateUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari path parameters
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Cari user berdasarkan ID untuk memastikan user ada
	existingUser, err := models.GetUserByID(db, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Decode input JSON
	var input RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Gunakan field lama jika tidak ada input baru untuk field tersebut
	if input.Name != "" {
		existingUser.Name = input.Name
	}
	if input.Email != "" {
		existingUser.Email = input.Email
	}
	if input.Phone != "" {
		existingUser.Phone = input.Phone
	}
	if input.Address != "" {
		existingUser.Address = input.Address
	}
	if input.Password != "" {
		// Hash password baru jika ada
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		existingUser.Password = string(hashedPassword)
	}

	// Lakukan update user di database
	if err := models.UpdateUser(db, existingUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("User updated successfully")
}


// DeleteUser untuk menghandle permintaan penghapusan user
func DeleteUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := models.DeleteUser(db, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("User deleted successfully")
}
