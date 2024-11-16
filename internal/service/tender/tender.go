package tender

import (
	"awesomeProject/internal/entity"
	tenderusecase "awesomeProject/internal/usecase/tender"
	"context"
	"log/slog"
	"time"
)

type Tender struct {
	log    *slog.Logger
	tender tenderusecase.TenderRepositoryImpl
}

func NewTender(
	log *slog.Logger,
	tender tenderusecase.TenderRepositoryImpl) *Tender {
	return &Tender{
		log:    log,
		tender: tender,
	}
}

func (t *Tender) CreateTender(ctx context.Context, req entity.CreateTenderRequest) (entity.Tender, error) {
	const op = "Tender.CreateTender"
	log := t.log.With(slog.String("method", op))
	log.Info("Start")
	defer log.Info("End")

	tender, err := t.tender.SaveTender(ctx, &req)
	if err != nil {
		return entity.Tender{}, err
	}

	return *tender, nil
}

func (t *Tender) GetTenders(ctx context.Context, req entity.GetListTender) ([]entity.Tender, error) {
	const op = "Tender.GetTenders"
	log := t.log.With(slog.String("method", op))
	log.Info("Start")
	defer log.Info("End")
	var resultTenners []entity.Tender
	tenders, err := t.tender.GetTenders(ctx, &req)
	if err != nil {
		log.Error(err.Error())
		return resultTenners, err
	}
	for _, tender := range tenders {
		if tender.Deadline.Before(time.Now()) {
			_, err = t.tender.UpdateTenderStatus(ctx, &entity.UpdateTenderStatusRequest{
				ID:        tender.ID,
				NewStatus: "expired",
				ClientID:  tender.ClientID,
			})
			if err != nil {
				log.Error(err.Error())
				return resultTenners, err
			}
			continue
		}
		resultTenners = append(resultTenners, tender)
	}
	return resultTenners, nil
}

func (t *Tender) UpdateTenderStatus(ctx context.Context, req entity.UpdateTenderStatusRequest) (entity.Tender, error) {
	const op = "Tender.UpdateTenderStatus"
	log := t.log.With(slog.String("method", op))
	log.Info("Start")
	defer log.Info("End")
	tender, err := t.tender.UpdateTenderStatus(ctx, &req)
	if err != nil {
		return entity.Tender{}, err
	}
	return *tender, nil
}

func (t *Tender) DeleteTender(ctx context.Context, req entity.DeleteTenderRequest) (message string, err error) {
	const op = "Tender.DeleteTender"
	log := t.log.With(slog.String("method", op))
	log.Info("Start")
	defer log.Info("End")
	err = t.tender.DeleteTender(ctx, &req)
	if err != nil {
		log.Error("err", err.Error())
		return "error delete", err
	}
	return "success", nil
}
