package handler

import (
	"awesomeProject/internal/entity"
	tenderusecase "awesomeProject/internal/usecase/tender"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

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
		// Tekshirilgan xatoliklarni e'tiborga olmaslik
		if errResp.Code != "BucketAlreadyOwnedByYou" && errResp.Code != "BucketAlreadyExists" {
			panic(fmt.Sprintf("Failed to create bucket: %v", err))
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
// @Accept multipart/form-data
// @Param pdf formData file false "Upload PDF"
// @Param data body entity.CreateTenderRequest true "Tender data"
// @Success 201 {object} string
// @Failure 400 {object} gin.H
// @Router /tenders [post]
func (t *Tender) CreateTender(c *gin.Context) {
	var req entity.CreateTenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, _ := c.FormFile("pdf")
	if file != nil {
		ext := filepath.Ext(file.Filename)
		if ext != ".pdf" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type, must be .pdf"})
			return
		}

		newFileName := uuid.NewString() + ext
		err := t.uploadPDF(c, file, newFileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		req.PDFURL = newFileName
	}

	message, err := t.tender.CreateTender(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}

// GetTenders godoc
// @Summary Get tenders
// @Description Retrieve tenders and their associated PDFs
// @Tags tenders
// @Accept json
// @Param filter body entity.GetListTender true "Tender filter"
// @Success 200 {array} entity.TenderResponse
// @Failure 400 {object} gin.H
// @Router /tenders [get]
func (t *Tender) GetTenders(c *gin.Context) {
	var req entity.GetListTender
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := t.tender.GetTenders(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, tender := range res {
		if tender.FileAttachment != "" {
			presignedURL, err := t.miniClient.PresignedGetObject(context.Background(), t.bucketName, tender.FileAttachment, 24*time.Hour, nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			res[i].FileAttachment = presignedURL.String()
		}
	}

	c.JSON(http.StatusOK, res)
}

// UpdateTender godoc
// @Summary Update Tender
// @Description Update information of a specific Tender by its ID
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "Tender ID"
// @Param tender body entity.UpdateTenderStatusRequest true "Tender update request body"
// @Success 200 {object} entity.Tender
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /tenders/{id} [put]
func (t *Tender) UpdateTenderStatus(c *gin.Context) {
	var req entity.UpdateTenderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := t.tender.UpdateTenderStatus(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// DeleteTender godoc
// @Summary Delete a tender by ID
// @Description Delete a specific tender by its ID
// @Tags tender
// @Accept json
// @Produce json
// @Param id path string true "Tender ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security Bearer
// @Router /tender/{id} [delete]
func (t *Tender) DeleteTender(c *gin.Context) {
	var req entity.DeleteTenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := t.tender.DeleteTender(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (t *Tender) uploadPDF(c *gin.Context, file *multipart.FileHeader, filename string) error {

	path := "./temp/" + filename
	if err := c.SaveUploadedFile(file, path); err != nil {
		return err
	}

	_, err := t.miniClient.FPutObject(context.Background(), t.bucketName, filename, path, minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	return err
}