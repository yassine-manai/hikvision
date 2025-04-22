package hikvision

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// HandlerConfig configuration for the HTTP handler
type HandlerConfig struct {
	EndpointPath string
	SaveXML      bool
	SaveImages   bool
	LogLevel     string
}

// NewHandler creates a new Hikvision LPR handler
func NewHandler(config HandlerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		clientIP := c.ClientIP()

		// Check for multipart form
		form, err := c.MultipartForm()
		if err != nil {
			log.Warn().Str("ip", clientIP).Msg("Invalid multipart form data")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "multipart form data required",
				"code":  400,
			})
			return
		}

		// Process XML data
		xmlFile, err := getFileFromForm(form, "anpr.xml")
		if err != nil {
			log.Warn().Str("ip", clientIP).Msg("anpr.xml file missing")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "anpr.xml file required",
				"code":  400,
			})
			return
		}

		xmlData, err := readFileContent(xmlFile)
		if err != nil {
			log.Error().Str("ip", clientIP).Err(err).Msg("Failed to read XML file")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to process request",
				"code":  500,
			})
			return
		}

		// Parse XML
		alert, err := ParseXMLData(xmlData)
		if err != nil {
			log.Error().Str("ip", clientIP).Err(err).Msg("Failed to parse XML")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid XML format",
				"code":  400,
			})
			return
		}

		// Extract capture data
		capture := ExtractCaptureFromAlert(alert)
		capture.XMLData = xmlData // Store raw XML if needed

		// Process images
		images, err := ExtractImagesFromForm(form)
		if err != nil {
			log.Warn().Str("ip", clientIP).Err(err).Msg("Failed to extract images")
		} else {
			capture.Images = images
		}

		// Save data if configured
		if config.SaveXML {
			go saveXMLData(xmlData, capture.LicensePlate, clientIP)
		}

		if config.SaveImages && len(images) > 0 {
			go saveImages(images, clientIP)
		}

		// Process the capture (business logic)
		go ProcessCapture(capture)

		// Log and respond
		log.Info().
			Str("ip", clientIP).
			Str("licensePlate", capture.LicensePlate).
			Str("direction", capture.Direction).
			Int("confidence", capture.Confidence).
			Dur("duration", time.Since(startTime)).
			Msg("LPR event processed")

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"code":    200,
			"message": "LPR data processed successfully",
		})
	}
}

// Helper functions
func getFileFromForm(form *multipart.Form, filename string) (*multipart.FileHeader, error) {
	files, ok := form.File[filename]
	if !ok || len(files) == 0 {
		return nil, fmt.Errorf("file %s not found", filename)
	}
	return files[0], nil
}

func readFileContent(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return "", err
	}

	return buf.String(), nil
}
