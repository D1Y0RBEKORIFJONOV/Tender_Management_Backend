package bidusecase

import (
	"awesomeProject/internal/entity"
	"context"
)

type bid interface {
	CreateBid(ctx context.Context, req *entity.CreateBidRequest) (*entity.Bid, error)
	GetBids(ctx context.Context, req *entity.GetBidsRequest) ([]*entity.Bid, error)
	UpdateBid(ctx context.Context, req *entity.UpdateBidRequest) (*entity.Bid, error)
	DeleteBid(ctx context.Context, req *entity.DeleteBidsRequest) (message string, err error)
	AnnounceWinner(ctx context.Context, req *entity.AnnounceWinnerRequest) (*entity.Bid, error)
}

type BidUseCaseIml struct {
	bid bid
}

func NewBidUseCaseIml(bid bid) *BidUseCaseIml {
	return &BidUseCaseIml{
		bid: bid,
	}
}

func (b *BidUseCaseIml) CreateBid(ctx context.Context, req *entity.CreateBidRequest) (*entity.Bid, error) {
	return b.bid.CreateBid(ctx, req)
}
func (b *BidUseCaseIml) GetBids(ctx context.Context, req *entity.GetBidsRequest) ([]*entity.Bid, error) {
	return b.bid.GetBids(ctx, req)
}
func (b *BidUseCaseIml) UpdateBid(ctx context.Context, req *entity.UpdateBidRequest) (*entity.Bid, error) {
	return b.bid.UpdateBid(ctx, req)
}
func (b *BidUseCaseIml) DeleteBid(ctx context.Context, req *entity.DeleteBidsRequest) (message string, err error) {
	return b.bid.DeleteBid(ctx, req)
}
func (b *BidUseCaseIml) AnnounceWinner(ctx context.Context, req *entity.AnnounceWinnerRequest) (*entity.Bid, error) {
	return b.bid.AnnounceWinner(ctx, req)
}
