package cmd

import (

	"time"

	log "github.com/sirupsen/logrus"
)

func TokenExpired(cmContent map[string]string) bool {

	// convert tokenExpirey to time.Time
	tokenExpirey, tokenExpireyErr := time.Parse(time.RFC3339, cmContent["expireyTime"])

	// if could not convert token expirey
	if tokenExpireyErr != nil {

		log.WithFields(log.Fields{
			"error": tokenExpireyErr.Error(),
		}).Error("Error converting token expirey to proper time.Time ...")
	}

	currentTime := time.Now().UTC()

	log.WithFields(log.Fields{
		"tokenExpireyTime": tokenExpirey,
		"currentTime":      currentTime,
	}).Info("Validating token age ...")
	
	diff := tokenExpirey.Sub(currentTime)

	if diff.Hours() < 2.0 {

		log.WithFields(log.Fields{
			"tokenExpireyTime": tokenExpirey,
			"currentTime":      currentTime,
		}).Warn("Token expired ...")
		
		return true
	}

	return false

}
