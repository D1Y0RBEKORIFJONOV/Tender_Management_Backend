package tenderusecase

import (
	"awesomeProject/internal/entity"
	"context"
)

type tenderRepository interface {
	SaveTender(ctx context.Context, req *entity.CreateTenderRequest) (*entity.Tender, error)
	GetTenders(ctx context.Context, req *entity.GetListTender) ([]entity.Tender, error)
	UpdateTenderStatus(ctx context.Context, req *entity.UpdateTenderStatusRequest) (*entity.Tender, error)
	DeleteTender(ctx context.Context, req *entity.DeleteTenderRequest) error
}

type TenderRepositoryImpl struct {
	tenderRepository tenderRepository
}

func NewTenderRepository(tenderRepository tenderRepository) *TenderRepositoryImpl {
	return &TenderRepositoryImpl{tenderRepository: tenderRepository}
}

func (t *TenderRepositoryImpl) SaveTender(ctx context.Context, req *entity.CreateTenderRequest) (*entity.Tender, error) {
	return t.tenderRepository.SaveTender(ctx, req)
}

func (t *TenderRepositoryImpl) GetTenders(ctx context.Context, req *entity.GetListTender) ([]entity.Tender, error) {
	return t.tenderRepository.GetTenders(ctx, req)
}

func (t *TenderRepositoryImpl) UpdateTenderStatus(ctx context.Context, req *entity.UpdateTenderStatusRequest) (*entity.Tender, error) {
	return t.tenderRepository.UpdateTenderStatus(ctx, req)
}

func (t *TenderRepositoryImpl) DeleteTender(ctx context.Context, req *entity.DeleteTenderRequest) error {
	return t.tenderRepository.DeleteTender(ctx, req)
}
	