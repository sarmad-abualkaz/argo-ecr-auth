package cmd

import (

	"time"

	"github.com/sarmad-abualkaz/argo-ecr-auth/k8s"
	"k8s.io/apimachinery/pkg/api/errors"

	log "github.com/sirupsen/logrus"
)

func ECRSecretUpdater(awsProfile string, awsRegion string, ecrRegistry string, k8sConfig string, namespace string, secret string, sleep int) {
	
	// setup an err variable for all the secret/cm updates
	var err error

	// generate k8s session
	k8sclient, k8sclientErr := callk8s.SetupK8sClient(k8sConfig)
	
	if k8sclientErr != nil {

		log.WithFields(log.Fields{
			"location": k8sConfig,
			"Error": k8sclientErr.Error(),
		}).Fatal("Error setting up kubernetes client ...")
		
		// panic if client setup fails
		panic(k8sclientErr.Error())
	}

	for {

		// look for secret in argocd namespace
		secretContent, secretContentErr := callk8s.ReturnTokenSecret(k8sclient, secret , namespace)
		
		if secretContentErr != nil {

			// if not exist
			if errors.IsNotFound(secretContentErr) {
				
				log.WithFields(log.Fields{
					"namespace": namespace,
					"name": secret,
				}).Warn("Secret not found ...")

				// look for cm in argocd namespace
				_, cmContentErr := callk8s.ReturnTokenData(k8sclient, secret , namespace)
				
				if cmContentErr != nil {

					// if not exist - create secret and cm
					if errors.IsNotFound(cmContentErr) {

						log.WithFields(log.Fields{
							"namespace": namespace,
							"name": secret,
						}).Warn("Secret configmap data not found ...")

						err = UpdateECRSecret(awsProfile, awsRegion, false, ecrRegistry, k8sclient, namespace, false, secret)

					} else {

						log.WithFields(log.Fields{
							"namespace": namespace,
							"name": secret,
							"Error": secretContentErr.Error(),
						}).Error("Error getting secret configmap data ...")
					}
				
				// else - create secret and update cm
				} else {

					err = UpdateECRSecret(awsProfile, awsRegion, true, ecrRegistry, k8sclient, namespace, false, secret)
				}
			
			// if error retreiving secret is more than not found
			} else {
				
				log.WithFields(log.Fields{
					"namespace": namespace,
					"name": secret,
					"Error": secretContentErr.Error(),
				}).Error("Error getting secret ...")
			}
		
		// else - secret exists
		} else {

			// check cm
			cmContent, cmContentErr := callk8s.ReturnTokenData(k8sclient, secret , namespace)
			
			if cmContentErr != nil {

				// if cm not exist (assume token expired)
				if errors.IsNotFound(cmContentErr) {

					log.WithFields(log.Fields{
						"namespace": namespace,
						"name": secret,
					}).Warn("Secret configmap data not found ...")

					err = UpdateECRSecret(awsProfile, awsRegion, false, ecrRegistry, k8sclient, namespace, true, secret)

				} else {

					log.WithFields(log.Fields{
						"namespace": namespace,
						"name": secret,
						"Error": cmContentErr.Error(),
					}).Error("Error getting secret configmap data ...")
				}

			} else {
				
				// convert tokenExpirey to time.Time
				tokenExpirey, tokenExpireyErr := time.Parse(time.RFC3339, cmContent["expireyTime"])

				// if date is past 11hrs from now || could not convert token expirey
				if TokenExpired(tokenExpirey) || (tokenExpireyErr != nil) {

					if tokenExpireyErr != nil {
						log.WithFields(log.Fields{
							"Error": tokenExpireyErr.Error(),
						}).Error("Error converting token expirey to proper time.Time ...")
					}
					
					err = UpdateECRSecret(awsProfile, awsRegion, true, ecrRegistry, k8sclient, namespace, true, secret)
				
				// else 				
				} else {
					
					// validate secret data
					validData, contentArea, found, expect := ValidateECRSecret(ecrRegistry, secretContent)
					
					// if data is invalid
					if !validData {

						log.WithFields(log.Fields{
							"namespace": namespace,
							"name": secret,
						}).Warn("Secret data not valid ...")

						log.WithFields(log.Fields{
							"namespace": namespace,
							"name": secret,
						}).Warn("On ", contentArea, " found - ", found, " expects - ", expect)
						
						err = UpdateECRSecret(awsProfile, awsRegion, true, ecrRegistry, k8sclient, namespace, true, secret)
					
					// else	if data is valid
					} else {

						// no-op log
						log.WithFields(log.Fields{
							"namespace": namespace,
							"name": secret,
						}).Info("Secret found with correct valid token and data.")
					}
				}
			}
		}
		
		// check if an error was returned during update
		if err != nil {

			log.WithFields(log.Fields{
				"namespace": namespace,
				"name": secret,
			}).Error("Failed to updated data.", err.Error())
		}

		// sleep x seconds
		log.WithFields(log.Fields{
		}).Info("Sleeping for ", sleep, "seconds ...")
		
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}
