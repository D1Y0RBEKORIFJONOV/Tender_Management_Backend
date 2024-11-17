package userrepo

import (
	"awesomeProject/internal/entity"
	"awesomeProject/internal/infastructure/repository/databaseconnection"
	postgres "awesomeProject/internal/infastructure/repository/postgres/sqlbuilder"
	authusecase "awesomeProject/internal/usecase/auth"
	"awesomeProject/logger"
	"context"
)

type UserRepository struct {
	db *databaseconnection.Database
}

func NewUserRepository() authusecase.AuthDbUseCase {
	db, err := databaseconnection.Connect()
	if err != nil {
		logger.SetupLogger(err.Error())
	}
	return &UserRepository{db: db}
}

func (u *UserRepository) SaveUser(ctx context.Context, req *entity.CreateUsrRequest) (*entity.User, error) {
	query, args, err := postgres.CreateUser(req)
	if err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	var user entity.User

	if err := u.db.Db.QueryRow(query, args...).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.Email); err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) IsHaveUser(ctx context.Context, byEmail string) (bool, error) {
	query, args, err := postgres.HaveUser(byEmail)
	if err != nil {
		logger.SetupLogger(err.Error())
		return true, err
	}

	var check bool

	if err := u.db.Db.QueryRow(query, args...).Scan(&check); err != nil {
		logger.SetupLogger(err.Error())
		return true, err
	}
	return check, nil
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, byEmail string) (*entity.User, error) {
	query, args, err := postgres.Getuser(byEmail)
	if err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	var user entity.User

	if err := u.db.Db.QueryRow(query, args...).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.Email); err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	return &user, nil
}
