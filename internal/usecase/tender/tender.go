package tenderusecase

import (
	"awesomeProject/internal/entity"
	"context"
)

type tenderUseCase interface {
	CreateTender(ctx context.Context, tender entity.CreateTenderRequest) (entity.Tender, error)
	GetTenders(ctx context.Context, req entity.GetListTender) ([]entity.Tender, error)
	UpdateTenderStatus(ctx context.Context, req entity.UpdateTenderStatusRequest) (entity.Tender, error)
	DeleteTender(ctx context.Context, req entity.DeleteTenderRequest) (message string, err error)
}

type TenderUseCaseIml struct {
	tender tenderUseCase
}

func NewTenderUseCase(tender tenderUseCase) *TenderUseCaseIml {
	return &TenderUseCaseIml{
		tender: tender,
	}
}

func (t *TenderUseCaseIml) CreateTender(ctx context.Context, tender entity.CreateTenderRequest) (entity.Tender, error) {
	return t.tender.CreateTender(ctx, tender)
}

func (t *TenderUseCaseIml) GetTenders(ctx context.Context, req entity.GetListTender) ([]entity.Tender, error) {
	return t.tender.GetTenders(ctx, req)
}

func (t *TenderUseCaseIml) UpdateTenderStatus(ctx context.Context, req entity.UpdateTenderStatusRequest) (entity.Tender, error) {
	return t.tender.UpdateTenderStatus(ctx, req)
}

func (t *TenderUseCaseIml) DeleteTender(ctx context.Context, req entity.DeleteTenderRequest) (message string, err error) {
	return t.tender.DeleteTender(ctx, req)
}
