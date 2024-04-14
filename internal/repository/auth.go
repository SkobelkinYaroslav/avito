package repository

import (
	"avito/internal/domain"
	"database/sql"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (a AuthRepository) RegisterRepo(user domain.AuthStruct) error {
	_, err := a.db.Exec("INSERT INTO users(email, password, is_admin) VALUES ($1, $2, $3)", user.Email, user.Password, user.IsAdmin)
	if err != nil {
		return err
	}
	return nil
}

func (a AuthRepository) GetUserRepo(user domain.AuthStruct) (domain.AuthStruct, error) {
	err := a.db.QueryRow("SELECT * FROM users WHERE email = $1", user.Email).Scan(&user.ID, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		return domain.AuthStruct{}, err
	}

	return user, nil
}
