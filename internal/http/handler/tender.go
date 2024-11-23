package handler

import (
	"awesomeProject/internal/entity"
	tenderusecase "awesomeProject/internal/usecase/tender"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const maxFileSize = 10 * 1024 * 1024

type Tender struct {
	tender     tenderusecase.TenderUseCaseIml
	miniClient *minio.Client
	bucketName string
}

func NewTender(tender tenderusecase.TenderUseCaseIml) *Tender {
	client, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("mini", "minio123", ""),
		Secure: false,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize MinIO client: %v", err))
	}

	bucketName := "tender-pdfs"
	err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "BucketAlreadyOwnedByYou" || errResp.Code == "BucketAlreadyExists" {
			log.Printf("Bucket %s already exists or is owned by you: %v", bucketName, errResp)
		} else {
			log.Printf("Failed to create bucket: %v", err)
		}
	}

	return &Tender{
		tender:     tender,
		miniClient: client,
		bucketName: bucketName,
	}
}

// CreateTender godoc
// @Summary Create a new tender
// @Description Create a new tender and optionally upload a PDF
// @Tags tenders
// @Accept json
// @Param data body entity.CreateTenderRequest true "Tender data"
// @Success 201 {object} string
// @Failure 400 {object} string
// @Security Bearer
// @Router /api/client/tenders [post]
func (t *Tender) CreateTender(c *gin.Context) {
	var req entity.CreateTenderRequest

	id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user_id not found"})
		return
	}
	req.ClientID = id.(string)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
		})
		return
	}
	if req.Deadline.IsZero() || req.Budget <= 0 {
		c.JSON(400, gin.H{"message": "Invalid tender data"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
			"details": err.Error(),
		})
		return
	}

	message, err := t.tender.CreateTender(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}

var validate = validator.New()

// GetTenders godoc
// @Summary Get tenders
// @Description Retrieve tenders for the authenticated client
// @Tags tenders
// @Accept json
// @Success 200 {object} []entity.Tender
// @Failure 400 {object} string
// @Security Bearer
// @Router /api/client/tenders [get]
func (t *Tender) GetTenders(c *gin.Context) {
	var req entity.GetListTender
	id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
		return
	}
	if req.Value == "" {
		req.Field = "client_id"
		req.Value = id.(string)
	}
	res, err := t.tender.GetTenders(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// UpdateTenderStatus godoc
// @Summary Update Tender
// @Description Update information of a specific Tender by its ID
// @Tags tenders
// @Accept json
// @Produce json
// @Param tenderId path string true "Tender ID"
// @Param tender body entity.UpdateTenderStatusRequest true "Tender update request body"
// @Success 200 {object} entity.Tender
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /api/client/tenders/{tenderId} [put]
func (t *Tender) UpdateTenderStatus(c *gin.Context) {
	var req entity.UpdateTenderStatusRequest
	id, ok := c.Get("user_id")
	if !ok {
		log.Println("11111111111111111111111111111111111")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
		return
	}
	if req.NewStatus == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tender status"})
		return
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("222222222222222222222222222", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = c.Param("tenderId")
	req.ClientID = id.(string)
	log.Println("REEEEEEEEEEEEEEEEEEQ", req.NewStatus, req.ID, req.ClientID)
	_, err := t.tender.UpdateTenderStatus(c.Request.Context(), req)
	if err != nil {
		log.Println("3333333333333333333333333333333333333333", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tender status updated"})
}

// DeleteTender godoc
// @Summary Delete a tender by ID
// @Description Delete a specific tender by its ID
// @Tags tenders
// @Accept json
// @Produce json
// @Param id path string true "Tender ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /api/client/tenders/{tenderId} [delete]
func (t *Tender) DeleteTender(c *gin.Context) {
	var req entity.DeleteTenderRequest
	id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found"})
		return
	}
	req.ClientID = id.(string)

	_, err := t.tender.DeleteTender(c.Request.Context(), req)
	if err != nil {
		c.JSON(404, gin.H{"message": "Tender not found or access denied"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tender deleted successfully"})
}

func (t *Tender) uploadPDF(c *gin.Context, file *multipart.FileHeader, filename string) error {
	path := "./temp/" + filename
	if err := c.SaveUploadedFile(file, path); err != nil {
		return err
	}
	defer os.Remove(path)

	_, err := t.miniClient.FPutObject(context.Background(), t.bucketName, filename, path, minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	return err
}

//// GetALlTenders godoc
//// @Summary Get all a tender by ID
//// @Description get all a specific tender by its ID
//// @Tags tenders
//// @Produce json
//// @Success 200 {object} string
//// @Failure 400 {object} string
//// @Failure 500 {object} string
//// @Security Bearer
//// @Router /api/client/tenders [get]
//func (t *Tender) GetALlTenders(c *gin.Context) {
//	res, err := t.tender.GetTenders(c.Request.Context(), entity.GetListTender{})
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, res)
//}
