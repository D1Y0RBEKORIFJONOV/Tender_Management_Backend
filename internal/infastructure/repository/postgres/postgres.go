package postgres

import (
	"awesomeProject/internal/entity"
	"awesomeProject/internal/infastructure/repository/databaseconnection"
	postgres "awesomeProject/internal/infastructure/repository/postgres/sqlbuilder"
	"awesomeProject/logger"
)

type Repository struct {
	db *databaseconnection.Database
}

func NewRepository() *Repository {
	db, err := databaseconnection.Connect()
	if err != nil {
		logger.SetupLogger(err.Error())
	}
	return &Repository{db: db}
}

func (u *Repository) SaveUser(req *entity.CreateUsrRequest) (*entity.User, error) {
	query, args, err := postgres.CreateUser(req)
	if err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	var user entity.User

	if err := u.db.Db.QueryRow(query, args...).Scan(&user.ID, &user.Username, &user.Role, &user.Email); err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	return &user, nil
}

func (u *Repository) SaveTender(req *entity.CreateTenderRequest) (*entity.Tender, error) {
	query, args, err := postgres.CreateTender(req)
	if err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}

	var tender entity.Tender

	if err := u.db.Db.QueryRow(query, args...).Scan(&tender.ID, &tender.ClientID, &tender.Title, &tender.Description, &tender.Deadline, &tender.Budget, &tender.Status, &tender.FileAttachment); err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	return &tender, nil
}

func (u *Repository) GetTenders(req *entity.GetListTender) ([]entity.Tender, error) {
	query, args, err := postgres.GetListTender(req)
	if err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}

	var tenders []entity.Tender

	dtenders, err := u.db.Db.Query(query, args...)
	if err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}

	for dtenders.Next() {
		var tender entity.Tender

		if err := dtenders.Scan(&tender.ID, &tender.ClientID, &tender.Title, &tender.Description, &tender.Deadline, &tender.Budget, &tender.Status, &tender.FileAttachment); err != nil {
			logger.SetupLogger(err.Error())
			return nil, err
		}

		tenders = append(tenders, tender)
	}
	return tenders, nil
}

func (u *Repository) UpdateTenderStatus(req *entity.UpdateTenderStatusRequest) (*entity.Tender, error) {
	query, args, err := postgres.UpdateTender(req)
	if err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	var tender entity.Tender

	if err := u.db.Db.QueryRow(query, args...).Scan(&tender.ID, &tender.ClientID, &tender.Title, &tender.Description, &tender.Deadline, &tender.Budget, &tender.Status, &tender.FileAttachment); err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	return &tender, nil
}

func (u *Repository) DeleteTender(req *entity.DeleteTenderRequest) error {
	query, args, err := postgres.DeleteTender(req)
	if err != nil {
		logger.SetupLogger(err.Error())
		return err
	}

	_, err = u.db.Db.Exec(query, args...)
	if err != nil {if err != nil {
		logger.SetupLogger(err.Error())
		return err
	}
		logger.SetupLogger(err.Error())
		return err
	}
	return nil
}
