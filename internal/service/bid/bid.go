package bid

import (
	"awesomeProject/internal/entity"
	bidusecase "awesomeProject/internal/usecase/bid"
	notificationusecase "awesomeProject/internal/usecase/notification"
	tenderusecase "awesomeProject/internal/usecase/tender"
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
)

type Bid struct {
	log          *slog.Logger
	bid          *bidusecase.BidUseCaseIml
	tender       *tenderusecase.TenderRepositoryImpl
	notification *notificationusecase.NotificationUseCase
}

func NewBid(log *slog.Logger,
	bid *bidusecase.BidUseCaseIml,
	tender *tenderusecase.TenderRepositoryImpl,
	notification *notificationusecase.NotificationUseCase) *Bid {
	return &Bid{log, bid, tender, notification}
}

func (b *Bid) CreateBid(ctx context.Context, req *entity.CreateBidRequest) (*entity.Bid, error) {
	const operation = "bid.CreateBid"
	log := b.log.With(slog.String("operation", operation))

	log.Info("Star")
	defer log.Info("End")
	bid, err := b.bid.CreateBid(ctx, req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	tender, err := b.getTender(ctx, bid.TenderID)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	message, err := json.Marshal(bid)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	err = b.notification.AddNotification(ctx, &entity.AddNotificationReq{
		UserId: tender.ClientID,
		CreateMessage: &entity.CreateMessageReq{
			Status:     string(message),
			SenderName: "bid",
		},
	})
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return bid, nil
}

func (b *Bid) GetBids(ctx context.Context, req *entity.GetBidsRequest) ([]entity.Bid, error) {
	const operation = "bid.GetBids"
	log := b.log.With(slog.String("operation", operation))
	log.Info("Star")
	defer log.Info("End")
	bid, err := b.bid.GetBids(ctx, req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return bid, nil
}

func (b *Bid) UpdateBid(ctx context.Context, req *entity.UpdateBidRequest) (*entity.Bid, error) {
	const operation = "bid.UpdateBid"
	log := b.log.With(slog.String("operation", operation))
	log.Info("Star")
	defer log.Info("End")
	bid, err := b.bid.UpdateBid(ctx, req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return bid, nil
}

func (b *Bid) DeleteBid(ctx context.Context, req *entity.DeleteBidsRequest) (message string, err error) {
	const operation = "bid.DeleteBid"
	log := b.log.With(slog.String("operation", operation))
	log.Info("Star")
	defer log.Info("End")
	_, err = b.bid.DeleteBid(ctx, req)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	return "successfully deleted", err
}

func (b *Bid) AnnounceWinner(ctx context.Context, req *entity.AnnounceWinnerRequest) (*entity.Bid, error) {
	const operation = "bid.AnnounceWinner"
	log := b.log.With(slog.String("operation", operation))
	log.Info("Star")
	defer log.Info("End")
	bid, err := b.bid.AnnounceWinner(ctx, req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	tender, err := b.getTender(ctx, bid.TenderID)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	_, err = b.tender.UpdateTenderStatus(ctx, &entity.UpdateTenderStatusRequest{
		ID:        bid.TenderID,
		NewStatus: "closed",
		ClientID:  tender.ClientID,
	})
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	message, err := json.Marshal(tender)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	err = b.notification.AddNotification(ctx, &entity.AddNotificationReq{
		UserId: bid.ContractorID,
		CreateMessage: &entity.CreateMessageReq{
			Status:     string(message),
			SenderName: "tender",
		},
	})
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return bid, nil
}

func (b *Bid) getTender(ctx context.Context, id string) (*entity.Tender, error) {
	tender, err := b.tender.GetTenders(ctx, &entity.GetListTender{
		Field: "id",
		Value: id,
	})
	if err != nil {
		log.Panicln(err.Error())
		return nil, err
	}
	if len(tender) == 0 {
		log.Panicln(err.Error())
		return nil, errors.New("error:tender not found")
	}
	return &tender[0], nil
}
