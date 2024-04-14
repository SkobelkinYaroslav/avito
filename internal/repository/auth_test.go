package repository_test

import (
	"avito/internal/domain"
	"avito/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewAuthRepository(db)

	expectedUser := domain.AuthStruct{
		Email:    "test@test.com",
		Password: "password",
		IsAdmin:  false,
	}

	mock.ExpectExec("^INSERT INTO users").
		WithArgs(expectedUser.Email, expectedUser.Password, expectedUser.IsAdmin).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.RegisterRepo(expectedUser)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewAuthRepository(db)

	expectedUser := domain.AuthStruct{
		ID:       1,
		Email:    "test@test.com",
		Password: "password",
		IsAdmin:  false,
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "is_admin"}).
		AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Password, expectedUser.IsAdmin)
	mock.ExpectQuery("^SELECT \\* FROM users WHERE email = \\$1").
		WithArgs(expectedUser.Email).
		WillReturnRows(rows)

	resultUser, err := repo.GetUserRepo(expectedUser)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, resultUser)

	assert.NoError(t, mock.ExpectationsWereMet())
}
