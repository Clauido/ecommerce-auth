package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/auth-microservice/models"
	_ "github.com/lib/pq"
)

type PostgresRespository struct {
	db *sql.DB
}

func NewPostgresRespository(url string) (*PostgresRespository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRespository{db}, nil
}

// CRUD
func (repo *PostgresRespository) InsertUser(ctx context.Context, user *models.User) error{
	_, err := repo.db.ExecContext(ctx,
		 "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)",user.Id, user.Email, user.Password)
		 return err
}
func (repo *PostgresRespository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, 
		"SELECT id, email FROM users WHERE id = $1", id)
	if err != nil {
    return nil, err
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			log.Println("Error closing rows:", closeErr)
		}
	}()
	var user = models.User{}
	for rows.Next() {
    var user models.User
    if err = rows.Scan(&user.Id, &user.Email); err == nil {
        return &user, nil
    }
}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}
func (repo *PostgresRespository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, 
		"SELECT id, email, password FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			log.Println("Error closing rows:", closeErr)
		}
	}()
	if err != nil {
		return nil, err
	}
	var user = models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email, &user.Password); err == nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}
// UPDATE
// DELETE
func (repo *PostgresRespository) Close() error {
		return repo.db.Close()
}