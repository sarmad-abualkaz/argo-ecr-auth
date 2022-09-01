package callecr

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"

	log "github.com/sirupsen/logrus"
)

type ECRClient struct {
	Client ecriface.ECRAPI
}

func GenerateECRTokent(sess *session.Session, region string, ecrRegistry string) (string, time.Time, error) {

	log.WithFields(log.Fields{
		"region": region,
		"ecr": ecrRegistry,
	}).Info("Call ECR for authorization token ...")

	ecrClient := ECRClient{
		Client: ecr.New(sess),
	}

	return GenerateECRTokentHelper(ecrClient)

}

func GenerateECRTokentHelper(ecrClient ECRClient) (string, time.Time, error) {

	input := &ecr.GetAuthorizationTokenInput{}

	result, err := ecrClient.Client.GetAuthorizationToken(input)
	
	if err != nil {

		return "", time.Time{}, err
	}

	decodedAuthorizationToken, _ := base64.StdEncoding.DecodeString(*result.AuthorizationData[0].AuthorizationToken)
	password := strings.Split(string(decodedAuthorizationToken), ":")[1]
	
	expiryTime := *result.AuthorizationData[0].ExpiresAt

	return password, expiryTime, nil
}
