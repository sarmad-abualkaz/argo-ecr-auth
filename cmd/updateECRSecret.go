package cmd

import (
	"github.com/sarmad-abualkaz/argo-ecr-auth/aws"
	"github.com/sarmad-abualkaz/argo-ecr-auth/k8s"
	
	"k8s.io/client-go/kubernetes"
)

func UpdateECRSecret(awsProfile string, awsRegion string, cmExist bool, ecrRegistry string, k8sClient *kubernetes.Clientset, namespace string, secretExist bool, secret string) error {

	// generate aws-session
	awsClient, awsClientErr := callecr.CreateSession(awsProfile, awsRegion)
	
	if awsClientErr != nil {

		return awsClientErr	
	}

	// initialize errors to update secert and configmap
	var secretErr error
	var cmErr error

	// call ecr for token
	password, expireyTime, err := callecr.GenerateECRTokent(awsClient, awsRegion, ecrRegistry)

	if err != nil {
		return err
	}

	// if secret not exist
	if !secretExist {
	
		// create secret with password
		secretErr = callk8s.UpdateTokenSecret(k8sClient, secret, namespace, password, ecrRegistry, false)

		if secretErr != nil {
			return secretErr
		}
		
	// else
	} else {

		// update secret with password
		secretErr = callk8s.UpdateTokenSecret(k8sClient, secret, namespace, password, ecrRegistry, true)

		if secretErr != nil {
			return secretErr
		}
	}

	// if cm not exist
	if !cmExist {

		// create cm
		cmErr = callk8s.UpdateTokenData(k8sClient, secret, namespace, expireyTime, ecrRegistry, false)

		if cmErr != nil {
			return cmErr
		}
	
	// else
	} else {

		// update cm
		cmErr = callk8s.UpdateTokenData(k8sClient, secret, namespace, expireyTime, ecrRegistry, true)

		if cmErr != nil {
			return cmErr
		}
	}

	return nil
}
