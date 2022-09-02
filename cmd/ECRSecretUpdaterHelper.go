package cmd

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"github.com/sarmad-abualkaz/argo-ecr-auth/k8s"
	
	log "github.com/sirupsen/logrus"
)


func ECRSecretUpdaterHelper(k8sclient callk8s.Client, namespace string, secret string, ecrRegistry string) (bool, bool, bool){

	var update      bool
	var secretExist bool
	var cmExist    bool

	// look for secret in argocd namespace
	secretContent, secretContentErr := callk8s.ReturnTokenSecret(k8sclient, secret , namespace)
	
	if secretContentErr != nil {

		// if not exist
		if errors.IsNotFound(secretContentErr) {
			
			log.WithFields(log.Fields{
				"name":      secret,
				"namespace": namespace,
			}).Warn("Secret not found ...")

			secretExist = false

		} else {

			log.WithFields(log.Fields{
				"error":     secretContentErr.Error(),
				"name":      secret,
				"namespace": namespace,
			}).Error("Error getting secret ...")

		} 
	} else {

		secretExist = true
	}

	// look for cm in argocd namespace
	cmContent, cmContentErr := callk8s.ReturnTokenData(k8sclient, secret , namespace)

	if cmContentErr != nil {

		if cmContentErr != nil {

			// if not exist - create secret and cm
			if errors.IsNotFound(cmContentErr) {

				log.WithFields(log.Fields{
					"name":      secret,
					"namespace": namespace,
				}).Warn("Secret configmap data not found ...")

				cmExist = false

			} else {

				log.WithFields(log.Fields{
					"error":     secretContentErr.Error(),
					"name":      secret,
					"namespace": namespace,
				}).Error("Error getting secret configmap data ...")
			}
		} 

	} else {
	
		cmExist = true
	}

	// if secret or cm cannot be found
	if !secretExist || !cmExist {

		update = true

	// if both exist validate token then data
	} else {

		// action on token exipred
		if TokenExpired(cmContent) {

			update = true

		} else {

			// validate secret data
			validData, contentArea, found, expect := ValidateECRSecret(ecrRegistry, secretContent)

			// if data is invalid
			if !validData {

				log.WithFields(log.Fields{
					"name":      secret,
					"namespace": namespace,
				}).Warn("Secret data not valid ...")

				log.WithFields(log.Fields{
					"name":      secret,
					"namespace": namespace,
				}).Warn("On ", contentArea, " found - ", found, " expects - ", expect)
				
				update = true

			// else	if data is valid
			} else {
				
				update = false

			}
		}
	}

	return update, secretExist, cmExist

}
