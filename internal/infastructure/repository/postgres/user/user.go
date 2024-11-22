package userrepo

import (
	"awesomeProject/internal/entity"
	"awesomeProject/internal/infastructure/repository/databaseconnection"
	postgres "awesomeProject/internal/infastructure/repository/postgres/sqlbuilder"
	authusecase "awesomeProject/internal/usecase/auth"
	"awesomeProject/logger"
	"context"
	"fmt"
)

// UserRepository - структура репозитория пользователей
type UserRepository struct {
	db *databaseconnection.Database
}

// NewUserRepository - конструктор для UserRepository
func NewUserRepository() (authusecase.AuthDbUseCase, error) {
	db, err := databaseconnection.Connect()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Database connection error: %v", err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &UserRepository{db: db}, nil
}

// SaveUser - сохранение нового пользователя
func (u *UserRepository) SaveUser(ctx context.Context, req *entity.CreateUsrRequest) (*entity.User, error) {
	// Проверяем подключение к базе данных
	if u.db == nil || u.db.Db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Генерируем SQL-запрос
	query, args, err := postgres.CreateUser(req)
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Error generating query in SaveUser: %v", err))
		return nil, err
	}

	var user entity.User

	// Выполняем запрос
	if err := u.db.Db.QueryRow(query, args...).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.Email); err != nil {
		logger.SetupLogger(fmt.Sprintf("Error executing query in SaveUser: %v, query: %s, args: %+v", err, query, args))
		return nil, err
	}

	return &user, nil
}

// IsHaveUser - проверка, существует ли пользователь с указанным email
func (u *UserRepository) IsHaveUser(ctx context.Context, byEmail string) (bool, error) {
	// Проверяем подключение к базе данных
	if u.db == nil || u.db.Db == nil {
		return false, fmt.Errorf("database connection is nil")
	}

	// Генерируем SQL-запрос
	query, args, err := postgres.HaveUser(byEmail)
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Error generating query in IsHaveUser: %v", err))
		return false, err
	}

	var check bool

	// Выполняем запрос
	if err := u.db.Db.QueryRow(query, args...).Scan(&check); err != nil {
		logger.SetupLogger(fmt.Sprintf("Error executing query in IsHaveUser: %v, query: %s, args: %+v", err, query, args))
		return false, err
	}

	return check, nil
}

// GetUserByEmail - получение пользователя по email
func (u *UserRepository) GetUserByEmail(ctx context.Context, byEmail string) (*entity.User, error) {
	// Проверяем подключение к базе данных
	if u.db == nil || u.db.Db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Генерируем SQL-запрос
	query, args, err := postgres.Getuser(byEmail)
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Error generating query in GetUserByEmail: %v", err))
		return nil, err
	}

	var user entity.User

	// Выполняем запрос
	if err := u.db.Db.QueryRow(query, args...).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.Email); err != nil {
		logger.SetupLogger(fmt.Sprintf("Error executing query in GetUserByEmail: %v, query: %s, args: %+v", err, query, args))
		return nil, err
	}

	return &user, nil
}
