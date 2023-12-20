package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/auth-microservice/models"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLRepository struct {
	db *sql.DB
}

// Crear nueva instancia del repositorio MySQL
func NewMySQLRepository(url string) (*MySQLRepository, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	return &MySQLRepository{db}, nil
}

// CRUD

func (repo *MySQLRepository) InsertUser(ctx context.Context, user *models.User) (int32, error) {
	result, err := repo.db.ExecContext(ctx, "INSERT INTO users (name, middle_name, rut, phone_number, email, password) VALUES (?, ?, ?, ?, ?, ?)",
		user.Name, user.MiddleName, user.Rut, user.PhoneNumber, user.Email, user.Password)

	if err != nil {
		log.Println("Error executing SQL:", err)
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Println("Error getting LastInsertId:", err)
		return 0, err
	}

	user.Id = int32(id)

	return user.Id, nil
}


func (repo *MySQLRepository) GetUserById(ctx context.Context, id int32) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name, middle_name, rut, phone_number, email FROM users WHERE id = ?", id)
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
		if err = rows.Scan(&user.Id, &user.Name,&user.MiddleName,&user.Rut,&user.PhoneNumber,&user.Email); err == nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *MySQLRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name, middle_name, rut, phone_number, email, password FROM users WHERE email = ?", email)
	if err != nil {
		log.Fatal(err)
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
		err = rows.Scan(&user.Id, &user.Name, &user.MiddleName, &user.Rut, &user.PhoneNumber, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}

		// If a user is found, return it immediately
		return &user, nil
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nil, sql.ErrNoRows
}


// UPDATE

// DELETE

func (repo *MySQLRepository) Close() error {
	return repo.db.Close()
}
