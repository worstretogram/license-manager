// Package license содержит обработчики для управления лицензиями: генерация, проверка, получение, обновление, удаление и скачивание.
// @title License Service API
// @version 1.0
// @description API для генерации и управления лицензиями с использованием RSA-подписи.
// @host localhost:8080
// @BasePath /
// @schemes http
package license

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LicenseData представляет полную структуру лицензии
type LicenseData struct {
	Owner       string    `json:"owner"`
	LicenseID   string    `json:"license_id"`
	MaxUsers    int       `json:"max_users"`
	MaxMessages int       `json:"max_messages"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Signature   string    `json:"signature"`
}

// GenerateRequest представляет структуру входных данных для генерации лицензии
type GenerateRequest struct {
	Owner       string `json:"owner" example:"Acme Inc"`
	AccessLevel struct {
		MaxUsers    int `json:"max_users" example:"100"`
		MaxMessages int `json:"max_messages" example:"10000"`
	} `json:"access_level"`
	ExpirationDate time.Time `json:"expiration_date" example:"2025-12-31T23:59:59Z"`
}

// Generate godoc
// @Summary Генерация новой лицензии
// @Description Генерирует лицензию, подписывает её приватным RSA-ключом и сохраняет в БД. Лицензия возвращается в закодированном виде.
// @Tags license
// @Accept json
// @Produce json
// @Param license body GenerateRequest true "Данные для генерации лицензии"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/license/generate [post]
func Generate(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input GenerateRequest

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
			return
		}
		if input.Owner == "" || input.AccessLevel.MaxUsers <= 0 || input.AccessLevel.MaxMessages <= 0 || input.ExpirationDate.IsZero() || input.ExpirationDate.Before(time.Now().UTC()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input values"})
			return
		}

		var count int
		err := pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM licenses WHERE owner = $1`, input.Owner).Scan(&count)
		if err != nil || count > 0 {
			status := http.StatusInternalServerError
			msg := "DB error"
			if count > 0 {
				status = http.StatusConflict
				msg = "License already exists"
			}
			c.JSON(status, gin.H{"error": msg})
			return
		}

		privKey, err := loadPrivateKey(os.Getenv("RSA_PRIVATE_KEY_PATH"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Private key error"})
			return
		}

		license := LicenseData{
			Owner:       input.Owner,
			LicenseID:   uuid.NewString(),
			MaxUsers:    input.AccessLevel.MaxUsers,
			MaxMessages: input.AccessLevel.MaxMessages,
			IssuedAt:    time.Now().UTC(),
			ExpiresAt:   input.ExpirationDate,
		}

		licenseToSign := license
		licenseToSign.Signature = ""
		jsonData, _ := json.Marshal(licenseToSign)
		hash := sha256.Sum256(jsonData)
		signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hash[:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Sign error"})
			return
		}

		license.Signature = base64.StdEncoding.EncodeToString(signature)
		_, err = pool.Exec(context.Background(), `
			INSERT INTO licenses (id, owner, max_users, max_messages, issued_at, expires_at, signature)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, license.LicenseID, license.Owner, license.MaxUsers, license.MaxMessages, license.IssuedAt, license.ExpiresAt, license.Signature)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB insert error"})
			return
		}

		encoded, _ := json.Marshal(license)
		licenseKey := base64.StdEncoding.EncodeToString(encoded)
		c.JSON(http.StatusOK, gin.H{
			"license_key": licenseKey,
			"license_id":  license.LicenseID,
			"expires_at":  license.ExpiresAt,
		})
	}
}

// Verify godoc
// @Summary Проверка лицензии
// @Description Проверяет целостность лицензии с помощью публичного ключа и возвращает статус
// @Tags license
// @Accept json
// @Produce json
// @Param license body map[string]string true "Base64 лицензия: {\"license_key\": \"...\"}"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/license/verify [post]
func Verify(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			LicenseKey string `json:"license_key"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		pubKey, err := loadPublicKey(os.Getenv("RSA_PUBLIC_KEY_PATH"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Public key error"})
			return
		}
		raw, _ := base64.StdEncoding.DecodeString(input.LicenseKey)
		var license LicenseData
		if err := json.Unmarshal(raw, &license); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Decode error"})
			return
		}
		copy := license
		copy.Signature = ""
		jsonData, _ := json.Marshal(copy)
		hash := sha256.Sum256(jsonData)
		signature, _ := base64.StdEncoding.DecodeString(license.Signature)
		if rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], signature) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "invalid"})
			return
		}
		if license.ExpiresAt.Before(time.Now().UTC()) {
			c.JSON(http.StatusOK, gin.H{"status": "expired"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":     "valid",
			"license_id": license.LicenseID,
			"owner":      license.Owner,
			"limits": gin.H{
				"max_users":    license.MaxUsers,
				"max_messages": license.MaxMessages,
			},
			"expires_at": license.ExpiresAt,
		})
	}
}

// ListLicenses godoc
// @Summary Получить список всех лицензий
// @Tags license
// @Produce json
// @Success 200 {array} LicenseData
// @Failure 500 {object} map[string]string
// @Router /admin/licenses [get]
func ListLicenses(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := pool.Query(context.Background(), `SELECT id, owner, max_users, max_messages, issued_at, expires_at, signature FROM licenses`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}
		defer rows.Close()

		var list []LicenseData
		for rows.Next() {
			var l LicenseData
			if err := rows.Scan(&l.LicenseID, &l.Owner, &l.MaxUsers, &l.MaxMessages, &l.IssuedAt, &l.ExpiresAt, &l.Signature); err == nil {
				list = append(list, l)
			}
		}
		c.JSON(http.StatusOK, list)
	}
}

// GetLicense godoc
// @Summary Получить лицензию по ID
// @Tags license
// @Produce json
// @Param id path string true "UUID лицензии"
// @Success 200 {object} LicenseData
// @Failure 404 {object} map[string]string
// @Router /admin/licenses/{id} [get]
func GetLicense(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var l LicenseData
		err := pool.QueryRow(context.Background(), `SELECT id, owner, max_users, max_messages, issued_at, expires_at, signature FROM licenses WHERE id=$1`, id).Scan(&l.LicenseID, &l.Owner, &l.MaxUsers, &l.MaxMessages, &l.IssuedAt, &l.ExpiresAt, &l.Signature)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "License not found"})
			return
		}
		c.JSON(http.StatusOK, l)
	}
}

// UpdateLicense godoc
// @Summary Обновить лицензию
// @Tags license
// @Accept json
// @Produce json
// @Param id path string true "UUID лицензии"
// @Param license body LicenseData true "Обновленные данные"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/licenses/{id} [put]
func UpdateLicense(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input LicenseData
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		_, err := pool.Exec(context.Background(), `UPDATE licenses SET owner=$1, max_users=$2, max_messages=$3, expires_at=$4 WHERE id=$5`, input.Owner, input.MaxUsers, input.MaxMessages, input.ExpiresAt, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	}
}

// DeleteLicense godoc
// @Summary Удалить лицензию
// @Tags license
// @Produce json
// @Param id path string true "UUID лицензии"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/licenses/{id} [delete]
func DeleteLicense(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := pool.Exec(context.Background(), `DELETE FROM licenses WHERE id=$1`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}

// DownloadLicense godoc
// @Summary Скачать лицензию
// @Description Скачивает лицензию в формате .lic (base64 encoded JSON)
// @Tags license
// @Produce application/octet-stream
// @Param id path string true "UUID лицензии"
// @Success 200 {file} binary
// @Failure 404 {object} map[string]string
// @Router /admin/licenses/{id}/download [get]
func DownloadLicense(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var license LicenseData
		err := pool.QueryRow(context.Background(), `SELECT id, owner, max_users, max_messages, issued_at, expires_at, signature FROM licenses WHERE id = $1`, id).Scan(&license.LicenseID, &license.Owner, &license.MaxUsers, &license.MaxMessages, &license.IssuedAt, &license.ExpiresAt, &license.Signature)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "License not found"})
			return
		}
		raw, _ := json.Marshal(license)
		key := base64.StdEncoding.EncodeToString(raw)
		c.Header("Content-Disposition", "attachment; filename=license_"+id+".lic")
		c.Data(http.StatusOK, "application/octet-stream", []byte(key))
	}
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key.(*rsa.PrivateKey), nil
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}
