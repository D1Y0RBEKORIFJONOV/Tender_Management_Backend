package tenders

import (
	"awesomeProject/internal/entity"
	"awesomeProject/internal/infastructure/repository/databaseconnection"
	postgres "awesomeProject/internal/infastructure/repository/postgres/sqlbuilder"
	tenderusecase "awesomeProject/internal/usecase/tender"
	"awesomeProject/logger"
	"context"
)

type TenderRepository struct {
	db *databaseconnection.Database
}

func NewTenderRepository() tenderusecase.TenderRepository {
	db, err := databaseconnection.Connect()
	if err != nil {
		logger.SetupLogger(err.Error())
	}
	return &TenderRepository{db: db}
}

func (u *TenderRepository) SaveTender(ctx context.Context, req *entity.CreateTenderRequest) (*entity.Tender, error) {
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

func (u *TenderRepository) GetTenders(ctx context.Context, req *entity.GetListTender) ([]entity.Tender, error) {
	// q, a, e := postgres.CloseExpiredTenders()
	// if e != nil {
	// 	logger.SetupLogger(e.Error())
	// 	return nil, e
	// }
	// _, err := u.db.Db.Exec(q, a...)
	// if err != nil {
	// 	logger.SetupLogger(err.Error())
	// 	return nil, err
	// }
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

func (u *TenderRepository) UpdateTenderStatus(ctx context.Context, req *entity.UpdateTenderStatusRequest) (*entity.Tender, error) {
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

func (u *TenderRepository) DeleteTender(ctx context.Context, req *entity.DeleteTenderRequest) error {
	query, args, err := postgres.DeleteTender(req)
	if err != nil {
		logger.SetupLogger(err.Error())
		return err
	}

	_, err = u.db.Db.Exec(query, args...)
	if err != nil {
		if err != nil {
			logger.SetupLogger(err.Error())
			return err
		}
		logger.SetupLogger(err.Error())
		return err
	}
	return nil
}
