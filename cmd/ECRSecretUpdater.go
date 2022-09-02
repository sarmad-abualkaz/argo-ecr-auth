package cmd

import (

	"time"

	"github.com/sarmad-abualkaz/argo-ecr-auth/k8s"

	log "github.com/sirupsen/logrus"
)

func ECRSecretUpdater(awsProfile string, awsRegion string, ecrRegistry string, k8sConfig string, namespace string, secret string, sleep int) {
	
	// setup an err variable for all the secret/cm updates
	var err error

	// generate k8s session
	k8sclient, k8sclientErr := callk8s.SetupK8sClient(k8sConfig)
	
	if k8sclientErr != nil {

		log.WithFields(log.Fields{
			"error":    k8sclientErr.Error(),
			"location": k8sConfig,
		}).Fatal("Error setting up kubernetes client ...")
		
		// panic if client setup fails
		panic(k8sclientErr.Error())
	}

	for {

		// look for secret and cm in argocd namespace
		update, secretExist, cmExist := ECRSecretUpdaterHelper(k8sclient, namespace, secret, ecrRegistry)

		// if an update is required - e.g. 
		// 1. secret or cm do not exist
		// 2. token expired
		// 3. secret data are invalid
		
		if update {
			
			// run the update command
			err = UpdateECRSecret(awsProfile, awsRegion, cmExist, ecrRegistry, k8sclient, namespace, secretExist, secret)
		
		// if update is not required
		} else {

			// no-op log
			log.WithFields(log.Fields{
				"name":      secret,
				"namespace": namespace,
			}).Info("Secret found with correct valid token and data.")
		}	
		
		// check if an error was returned during update
		if err != nil {

			log.WithFields(log.Fields{
				"name":      secret,
				"namespace": namespace,
			}).Error("Failed to updated data.", err.Error())
		}

		// sleep x seconds
		log.WithFields(log.Fields{
		}).Info("Sleeping for ", sleep, "seconds ...")
		
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}
