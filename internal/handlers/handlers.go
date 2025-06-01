package handlers

import (
	"net/http"
	"strconv"
	"time"

	"educational-game-db/internal/database"
	"educational-game-db/internal/export"
	"educational-game-db/internal/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db            *database.Database
	exportService *export.ExportService
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{
		db:            db,
		exportService: export.NewExportService(db),
	}
}

func (h *Handler) CreateAccount(c *gin.Context) {
	var req models.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.db.CreateAccount(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *Handler) GetAccounts(c *gin.Context) {
	accounts, err := h.db.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h *Handler) GetAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := h.db.GetAccountByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *Handler) UpdateAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var req models.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.db.UpdateAccount(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *Handler) DeleteAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	err = h.db.DeleteAccount(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.db.GetAccountStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.db.VerifyPassword(req.Username, req.Password) {
		account, err := h.db.GetAccountByUsername(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "account": account})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

// ExportCSV exports accounts to CSV format
func (h *Handler) ExportCSV(c *gin.Context) {
	filename := "accounts_export_" + time.Now().Format("2006-01-02_15-04-05") + ".csv"

	err := h.exportService.ExportToCSV(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "text/csv")
	c.File(filename)
}

// ExportJSON exports accounts to JSON format
func (h *Handler) ExportJSON(c *gin.Context) {
	filename := "accounts_export_" + time.Now().Format("2006-01-02_15-04-05") + ".json"

	err := h.exportService.ExportToJSON(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/json")
	c.File(filename)
}

// ImportCSV imports accounts from uploaded CSV file
func (h *Handler) ImportCSV(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Save uploaded file temporarily
	tempFilename := "temp_import_" + time.Now().Format("20060102_150405") + ".csv"
	if err := c.SaveUploadedFile(header, tempFilename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	// Import from the temporary file
	err = h.exportService.ImportFromCSV(tempFilename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Accounts imported successfully"})
}
