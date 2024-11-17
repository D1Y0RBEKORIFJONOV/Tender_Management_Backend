package postgres

import (
	"awesomeProject/internal/entity"
	"awesomeProject/logger"
	"errors"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(req *entity.CreateUsrRequest) (string, []interface{}, error) {
	id := uuid.New().String()
	query, args, err := squirrel.
		Insert("users").
		Columns("id", "username", "password", "role", "email").
		Values(id, req.Username, req.Password, req.Role, req.Email).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id,username,password,role,email").
		ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while Creating user %v", err))
		return "", nil, err
	}
	return query, args, err
}

func HaveUser(req string) (string, []interface{}, error) {
	query, args, err := squirrel.
		Select("EXISTS(SELECT *").From("users").Where(squirrel.Eq{"email": req}).
		PlaceholderFormat(squirrel.Dollar).
		Suffix(")").
		ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while Creating user %v", err))
		return "", nil, err
	}
	return query, args, err
}

func Getuser(req string) (string, []interface{}, error) {
	query, args, err := squirrel.Select("*").From("users").
		Where(squirrel.Eq{"email": req}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while getting user %v", err))
		return "", nil, err
	}
	return query, args, err
}

func CreateTender(req *entity.CreateTenderRequest) (string, []interface{}, error) {
	id := uuid.New().String()
	query, args, err := squirrel.Insert("tenders").
		Columns("id", "client_id", "title", "description", "deadline", "budget", "status", "fileattachment,created_at").
		Values(id, req.ClientID, req.Title, req.Description, req.Deadline, req.Budget, "open", req.FileAttachment, req.CreatedAt).
		PlaceholderFormat(squirrel.Dollar).Suffix("RETURNING *").
		ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while Creating tender %v", err))
		return "", nil, err
	}
	return query, args, err
}

func GetListTender(req *entity.GetListTender) (string, []interface{}, error) {
	queryBuilder := squirrel.
		Select("*").
		From("tenders").
		PlaceholderFormat(squirrel.Dollar)

	if req.Field != "" && req.Value != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{req.Field: req.Value})
	}

	if req.Limit > 0 && req.Page > 0 {
		offset := (req.Page - 1) * req.Limit
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("error generating tenders query: %v", err)
	}

	return query, args, nil
}

// func GetValidTenders() (string, []interface{}, error) {
// 	query, args, err := squirrel.
// 		Select("*").
// 		From("tenders").
// 		Where(squirrel.Gt{"deadline": "NOW()"}).
// 		Where(squirrel.Eq{"status": "open"}).
// 		PlaceholderFormat(squirrel.Dollar).
// 		ToSql()
// 	if err != nil {
// 		return "", nil, fmt.Errorf("failed to build query: %w", err)
// 	}

// 	return query, args, nil
// }

// func CloseExpiredTenders() (string, []interface{}, error) {
// 	query, args, err := squirrel.
// 		Update("tenders").
// 		Set("status", "closed").
// 		Where(squirrel.LtOrEq{"deadline": "NOW()"}).
// 		Where(squirrel.Eq{"status": "open"}).
// 		PlaceholderFormat(squirrel.Dollar).
// 		ToSql()

// 	if err != nil {
// 		return "", nil, fmt.Errorf("failed to build query: %w", err)
// 	}

// 	return query, args, nil
// }

func UpdateTender(req *entity.UpdateTenderStatusRequest) (string, []interface{}, error) {
	query, args, err := squirrel.
		Update("tenders").
		Set("status", req.NewStatus).
		Where(squirrel.Eq{"id": req.ID}).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while updating tenders %v", err))
		return "", nil, err
	}
	return query, args, err
}

func DeleteTender(req *entity.DeleteTenderRequest) (string, []interface{}, error) {
	query, args, err := squirrel.
		Delete("tenders").
		Where(squirrel.Eq{"id": req.ID, "client_id": req.ClientID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while deleting tenders %v", err))
		return "", nil, err
	}
	return query, args, err
}

func CreateBid(req *entity.CreateBidRequest) (string, []interface{}, error) {
	id := uuid.New().String()
	query, args, err := squirrel.Insert("bids").
		Columns("id", "tender_id", "contractor_id", "price", "delivery_time", "comments", "status").
		Values(id, req.TenderID, req.ContractorID, req.Price, req.DeliveryTime, req.Comments, "pending").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while creating bid %v", err))
		return "", nil, err
	}
	return query, args, err
}

func GetBids(req *entity.GetBidsRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Select("*").
		From("bids").
		Where(squirrel.Eq{"contractor_id": req.ContractorID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while getting bids %v", err))
		return "", nil, err
	}
	return query, args, err

}

func Update(req *entity.UpdateBidRequest) (string, []interface{}, error) {
	if req.TenderID == "" || req.ContractorID == "" {
		return "", nil, errors.New("tender_id and contractor_id are required")
	}

	updateFields := map[string]interface{}{}

	if req.Price > 0 {
		updateFields["price"] = req.Price
	}
	if req.Comments != "" {
		updateFields["comments"] = req.Comments
	}
	if req.Status != "" {
		updateFields["status"] = req.Status
	}

	if len(updateFields) == 0 {
		return "", nil, errors.New("no fields to update")
	}

	query, args, err := squirrel.
		Update("bids").
		SetMap(updateFields).
		Where(squirrel.Eq{"tender_id": req.TenderID, "contractor_id": req.ContractorID}).
		ToSql()

	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

func DeleteBid(req *entity.DeleteBidsRequest) (string, []interface{}, error) {
	query, args, err := squirrel.
		Delete("tenders").
		Where(squirrel.Eq{"contractor_id": req.ContractorID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while deleting bid %v", err))
		return "", nil, err
	}
	return query, args, err
}

func AnnounceWinner(req *entity.AnnounceWinnerRequest) (string, []interface{}, error) {
	if req.ContractorID == "" || req.BidID == "" {
		return "", nil, fmt.Errorf("missing required fields in AnnounceWinnerRequest")
	}

	updateMap := map[string]interface{}{
		"status": "awarded",
	}

	query, args, err := squirrel.
		Update("bids").
		SetMap(updateMap).
		Where(squirrel.Eq{
			"contractor_id": req.ContractorID,
			"id":            req.BidID,
		}).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		return "", nil, fmt.Errorf("error generating announce winner query: %v", err)
	}

	return query, args, nil
}

func RejectOtherBids(req *entity.AnnounceWinnerRequest) (string, []interface{}, error) {
	if req.BidID == "" {
		return "", nil, fmt.Errorf("missing required fields in AnnounceWinnerRequest")
	}
	updateMap := map[string]interface{}{
		"status": "rejected",
	}
	query, args, err := squirrel.
		Update("bids").
		SetMap(updateMap).
		Where(squirrel.Eq{"id": req.BidID}).
		Where(squirrel.NotEq{"contractor_id": req.ContractorID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return "", nil, fmt.Errorf("error generating reject other bids query: %v", err)
	}

	return query, args, nil
}

func ComparePassword(hashed, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func Hashing(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hashed)
}
