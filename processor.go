package hikvision

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

// ProcessCapture handles the business logic for a license plate capture
func ProcessCapture(capture *Capture) {
	// Add your business logic here
	// This could include:
	// - Database operations
	// - External API calls
	// - Event publishing
	// - Alert generation

	log.Info().
		Str("licensePlate", capture.LicensePlate).
		Str("direction", capture.Direction).
		Str("vehicleType", capture.VehicleType).
		Msg("Processing license plate capture")

	// Example: Process in a goroutine with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Process the capture (replace with your actual logic)
	select {
	case <-processCaptureData(ctx, capture):
		log.Info().
			Str("licensePlate", capture.LicensePlate).
			Msg("Capture processed successfully")
	case <-ctx.Done():
		log.Error().
			Str("licensePlate", capture.LicensePlate).
			Msg("Timeout processing capture")
	}
}

func processCaptureData(ctx context.Context, capture *Capture) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		// Simulate processing time
		time.Sleep(1 * time.Second)

		// Your actual processing logic here
		// For example:
		// - Save to database
		// - Check against watch lists
		// - Trigger alerts

		// Example: Just log for now
		log.Debug().
			Str("licensePlate", capture.LicensePlate).
			Msg("Finished processing capture data")
	}()

	return done
}

// Convenience methods for accessing capture data
func (c *Capture) GetLicensePlate() string {
	return c.LicensePlate
}

func (c *Capture) GetDirection() string {
	return c.Direction
}

func (c *Capture) GetConfidence() int {
	return c.Confidence
}

func (c *Capture) GetVehicleType() string {
	return c.VehicleType
}

func (c *Capture) GetImages() []Image {
	return c.Images
}

func (c *Capture) GetLicensePlateImage() (Image, bool) {
	for _, img := range c.Images {
		if img.Type == "licensePlate" {
			return img, true
		}
	}
	return Image{}, false
}

func (c *Capture) GetDetectionImage() (Image, bool) {
	for _, img := range c.Images {
		if img.Type == "detection" {
			return img, true
		}
	}
	return Image{}, false
}
