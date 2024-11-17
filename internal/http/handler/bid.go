package handler

import (
	"awesomeProject/internal/entity"
	bidusecase "awesomeProject/internal/usecase/bid"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Bid struct {
	bid bidusecase.BidUseCaseIml
}

func NewBid(bid bidusecase.BidUseCaseIml) *Bid {
	return &Bid{bid: bid}
}

// CreateBid(ctx context.Context, req *entity.CreateBidRequest) (*entity.Bid, error)
// GetBids(ctx context.Context, req *entity.GetBidsRequest) ([]*entity.Bid, error)
// UpdateBid(ctx context.Context, req *entity.UpdateBidRequest) (*entity.Bid, error)
// DeleteBid(ctx context.Context, req *entity.DeleteBidsRequest) (message string, err error)
// AnnounceWinner(ctx context.Context, req *entity.AnnounceWinnerRequest) (*entity.Bid, error)

// CreateBid godoc
// @Summary Create a new bid
// @Description Create a new bid
// @Tags bid
// @Accept json
// @Produce json
// @Param bid body entity.CreateBidRequest true "Bid request body"
// @Success 200 {object} entity.Bid
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /tenders/{id}/bids [post]
func (b *Bid) CreateBid(c *gin.Context) {
	var req entity.CreateBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := b.bid.CreateBid(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetBids godoc
// @Summary Get all bids for a tender
// @Description Get a list of bids for a specific tender
// @Tags bid
// @Accept json
// @Produce json
// @Param request body entity.GetBidsRequest true "Get Bids request body"
// @Success 200 {array} entity.Bid
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /tenders/{id}/bids [get]
func (b *Bid) GetBids(c *gin.Context) {
	var req entity.GetBidsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := b.bid.GetBids(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// UpdateBid godoc
// @Summary Update a bid
// @Description Update a bid with new information
// @Tags bid
// @Accept json
// @Produce json
// @Param bid body entity.UpdateBidRequest true "Update Bid request body"
// @Success 200 {object} entity.Bid
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /tenders/{id}/bids [put]
func (b *Bid) UpdateBid(c *gin.Context) {
	var req entity.UpdateBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	  return
	}
	id, ok := c.Get("user_id")
	if !ok {
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "user id not found"})
	  return
	}
  
	req.ContractorID = id.(string)
  
	res, err := b.bid.UpdateBid(c.Request.Context(), &req)
	if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	  return
	}
  
	c.JSON(http.StatusOK, res)
  }

// DeleteBid godoc
// @Summary Delete a bid by ID
// @Description Delete a specific bid by its ID
// @Tags bid
// @Accept json
// @Produce json
// @Param id path string true "Contractor ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /tender/{id} [delete]
func (b *Bid) DeleteBid(c *gin.Context) {
	var req entity.DeleteBidsRequest
	id := c.Param("id")
	req.ContractorID = id
	res, err := b.bid.DeleteBid(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// AnnounceWinner godoc
// @Summary Announce the winner for a tender
// @Description Announce the winner for a specific tender from the list of bids
// @Tags bid
// @Accept json
// @Produce json
// @Param request body entity.AnnounceWinnerRequest true "Announce Winner request body"
// @Success 200 {object} entity.Bid
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /tenders/{id}/bids/winner [post]
func (b *Bid) AnnounceWinner(c *gin.Context) {
	var req entity.AnnounceWinnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := b.bid.AnnounceWinner(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
