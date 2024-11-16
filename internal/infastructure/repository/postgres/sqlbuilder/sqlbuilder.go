package postgres

import (
	"awesomeProject/internal/entity"
	"awesomeProject/logger"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(req *entity.CreateUsrRequest) (string, []interface{}, error) {
	id := uuid.New().String()
	password := Hashing(req.Password)
	query, args, err := squirrel.
		Insert("users").
		Columns("id", "username", "password", "role", "email").
		Values(id, req.Username, password, req.Role, req.Email).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id,username,role,email").
		ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while Creating user %v", err))
		return "", nil, err
	}
	return query, args, err
}

func CreateTender(req *entity.CreateTenderRequest) (string, []interface{}, error) {
	id := uuid.New().String()
	query, args, err := squirrel.Insert("tenders").
		Columns("id", "client_id", "title", "description", "deadline", "budget", "status", "fileattachment").
		Values(id, req.ClientID, req.Title, req.Description, req.Deadline, req.Budget, "open", req.FileAttachment).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while Creating tender %v", err))
		return "", nil, err
	}
	return query, args, err
}

func GetListTender(req *entity.GetListTender)(string, []interface{}, error){
	query, args, err:=squirrel.
	Select("*").From("tenders").Limit(uint64(req.Limit)).
	Offset(uint64(req.Page)).
	PlaceholderFormat(squirrel.Dollar).
	ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while getting tenders %v", err))
		return "", nil, err
	}
	return query, args, err
}

func UpdateTender(req *entity.UpdateTenderStatusRequest)(string, []interface{}, error){
	query, args, err:=squirrel.
	Update("tenders").
	Set("status",req.NewStatus).
	Where(squirrel.Eq{"id":req.ID,"client_id":req.ClientID}).
	PlaceholderFormat(squirrel.Dollar).
	Suffix("RETURNING *").
	ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while updating tenders %v", err))
		return "", nil, err
	}
	return query, args, err
}

func DeleteTender(req *entity.DeleteTenderRequest)(string, []interface{}, error){
	query, args, err:=squirrel.
	Delete("tenders").
	Where(squirrel.Eq{"id":req.ID,"client_id":req.ClientID}).
	PlaceholderFormat(squirrel.Dollar).
	ToSql()
	if err != nil {
		logger.SetupLogger(fmt.Sprintf("Something went wrong while deleting tenders %v", err))
		return "", nil, err
	}
	return query, args, err
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
