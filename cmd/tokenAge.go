package cmd

import (

	"time"

	log "github.com/sirupsen/logrus"
)

func TokenExpired(tokenExpireyTime time.Time) bool {
	
	currentTime := time.Now().UTC()

	log.WithFields(log.Fields{
		"tokenExpireyTime": tokenExpireyTime,
		"currentTime": currentTime,
	}).Info("Validating token age ...")
	
	diff := tokenExpireyTime.Sub(currentTime)

	if diff.Hours() < 2.0 {

		log.WithFields(log.Fields{
			"tokenExpireyTime": tokenExpireyTime,
			"currentTime": currentTime,
		}).Warn("Token expired ...")
		
		return true
	}

	return false

}
