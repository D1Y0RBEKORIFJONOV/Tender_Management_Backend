package bids

import (
	"awesomeProject/internal/entity"
	"awesomeProject/internal/infastructure/repository/databaseconnection"
	postgres "awesomeProject/internal/infastructure/repository/postgres/sqlbuilder"
	bidusecase "awesomeProject/internal/usecase/bid"
	"awesomeProject/logger"
	"context"
)

type BidRepository struct {
	db *databaseconnection.Database
}

func NewBidRepository() bidusecase.Bid {
	db, err := databaseconnection.Connect()
	if err != nil {
		logger.SetupLogger(err.Error())
	}
	return &BidRepository{db: db}
}

func (u *BidRepository) CreateBid(ctx context.Context, req *entity.CreateBidRequest) (*entity.Bid, error) {
	query, args, err := postgres.CreateBid(req)
	if err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}

	var bid entity.Bid

	if err := u.db.Db.QueryRow(query, args...).Scan(&bid.ID, &bid.TenderID, &bid.ContractorID, &bid.Price, &bid.DeliveryTime, &bid.Comments, &bid.Status); err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	return &bid, nil
}

func (u *BidRepository) GetBids(ctx context.Context, req *entity.GetBidsRequest) ([]entity.Bid, error) {
	query, args, err := postgres.GetBids(req)
	if err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}

	var bids []entity.Bid

	dbids, err := u.db.Db.Query(query, args...)
	if err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}

	for dbids.Next() {
		var bid entity.Bid

		if err := dbids.Scan(&bid.ID, &bid.TenderID, &bid.ContractorID, &bid.Price, &bid.DeliveryTime, &bid.Comments, &bid.Status); err != nil {
			logger.SetupLogger(err.Error())
			return nil, err
		}
		bids = append(bids, bid)
	}
	return bids, nil

}

func (u *BidRepository) UpdateBid(ctx context.Context, req *entity.UpdateBidRequest) (*entity.Bid, error) {
	query, args, err := postgres.Update(req)
	if err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}

	var bid entity.Bid

	if err := u.db.Db.QueryRow(query, args...).Scan(&bid.ID, &bid.TenderID, &bid.ContractorID, &bid.Price, &bid.DeliveryTime, &bid.Comments, &bid.Status); err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	return &bid, nil
}

func (u *BidRepository) DeleteBid(ctx context.Context, req *entity.DeleteBidsRequest) (message string, err error) {
	query, args, err := postgres.DeleteBid(req)
	if err != nil {
		logger.SetupLogger(err.Error())
		return "", err
	}

	_, err = u.db.Db.Exec(query, args...)
	if err != nil {
		logger.SetupLogger(err.Error())
		return "", err
	}
	return "Bid has been deleted", nil
}

func (u *BidRepository) AnnounceWinner(ctx context.Context, req *entity.AnnounceWinnerRequest) (*entity.Bid, error) {
	query, args, err := postgres.AnnounceWinner(req)
	if err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	var bid entity.Bid
	if err := u.db.Db.QueryRow(query, args...).Scan(&bid.ID, &bid.TenderID, &bid.ContractorID, &bid.Price, &bid.DeliveryTime, &bid.Comments, &bid.Status); err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}

	q, a, e := postgres.RejectOtherBids(req)
	if e != nil {
		logger.SetupLogger(e.Error())
		return nil, e
	}

	_, err = u.db.Db.Exec(q, a...)
	if err != nil {
		logger.SetupLogger(err.Error())
		return nil, err
	}
	return &bid, nil
}
