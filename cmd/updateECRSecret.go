package cmd

import (
	"github.com/sarmad-abualkaz/argo-ecr-auth/aws"
	"github.com/sarmad-abualkaz/argo-ecr-auth/k8s"

)

func UpdateECRSecret(awsProfile string, awsRegion string, cmExist bool, ecrRegistry string, k8sClient callk8s.Client, namespace string, secretExist bool, secret string) error {

	// generate aws-session
	awsClient, awsClientErr := callecr.CreateSession(awsProfile, awsRegion)
	
	if awsClientErr != nil {

		return awsClientErr	
	}

	// initialize errors to update secert and configmap
	var secretErr error
	var cmErr     error

	// call ecr for token
	password, expireyTime, err := callecr.GenerateECRTokent(awsClient, awsRegion, ecrRegistry)

	if err != nil {
		return err
	}

	// create if !secretExist // update if secretExist
	secretErr = callk8s.UpdateTokenSecret(k8sClient, secret, namespace, password, ecrRegistry, secretExist)

	if secretErr != nil {
		return secretErr
	}

	// create if !cmExist // update if cmExist
	cmErr = callk8s.UpdateTokenData(k8sClient, secret, namespace, expireyTime, ecrRegistry, cmExist)

	if cmErr != nil {
		return cmErr
	}

	return nil
}
