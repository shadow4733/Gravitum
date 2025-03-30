package repo

import (
	"Gravitum/internal/model"
	"database/sql"
	"log"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	query := `
        INSERT INTO users (id, first_name, last_name, email, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.DB.Exec(
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		log.Println("Error creating user:", err)
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByID(id string) (*model.User, error) {
	query := `SELECT id, first_name, last_name, email, password, created_at, updated_at FROM users WHERE id = $1`
	var user model.User
	err := r.DB.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println("Error getting user by ID:", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	query := `UPDATE users SET first_name = $1, last_name = $2, email = $3, password = $4, updated_at = $5 WHERE id = $6`
	_, err := r.DB.Exec(query, user.FirstName, user.LastName, user.Email, user.Password, user.UpdatedAt, user.ID)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}
	return nil
}
