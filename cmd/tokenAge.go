package cmd

import (

	"time"

	log "github.com/sirupsen/logrus"
)

func TokenExpired(tokenGenTime time.Time) bool {
	
	currentTime := time.Now().UTC()

	log.WithFields(log.Fields{
		"tokenGenTime": tokenGenTime,
		"currentTime": currentTime,
	}).Info("Validating token age ...")
	
	diff := currentTime.Sub(tokenGenTime)

	if diff.Hours() > 11.0 {

		log.WithFields(log.Fields{
			"tokenGenTime": tokenGenTime,
			"currentTime": currentTime,
		}).Warn("Token expired ...")
		
		return true
	}

	return false

}
