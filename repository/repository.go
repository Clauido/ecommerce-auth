package repository

import (
	"context"

	"github.com/auth-microservice/models"
)

type Repository interface {

	InsertUser(ctx context.Context, user *models.User) (int32,error)

	GetUserById(ctx context.Context, id int32) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	// UpdateUser(ctx context.Context, user *models.User) error

	// DeleteUserById(ctx context.Context,id string) error
	// DeleteUserByEmail(ctx context.Context, email string) error

	Close() error
}

var implementation Repository

func SetRepository(repositoty Repository){
	implementation=repositoty
}

func InsertUser(ctx context.Context, user *models.User) (int32, error) {
	return implementation.InsertUser(ctx,user)
} 

func GetUserById(ctx context.Context, id int32) (*models.User, error) {
	return implementation.GetUserById(ctx,id)
}
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx,email)
}
// func UpdateUser(ctx context.Context, user *models.User) error {
// 	return implementation.UpdateUser(ctx, user)
// }

// func DeleteUserById(ctx context.Context, id string) error {
// 	return implementation.DeleteUserById(ctx, id)
// }
// func DeleteUserByEmail(ctx context.Context, email string) error {
// 	return implementation.DeleteUserByEmail(ctx, email)
// }
func Close() error {
	return implementation.Close()
}