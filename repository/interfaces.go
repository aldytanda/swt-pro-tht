// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error)

	// user
	CreateUser(ctx context.Context, v User) (User, error)
	UpdateUser(ctx context.Context, v User) (User, error)
	GetUser(ctx context.Context, vid int64) (User, error)

	// login
	Login(ctx context.Context, v Login) (string, error)
}
