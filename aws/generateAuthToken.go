package callecr

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"

	log "github.com/sirupsen/logrus"
)

func GenerateECRTokent(sess *session.Session, region string, ecrRegistry string) (string, time.Time, error) {

	log.WithFields(log.Fields{
		"region": region,
		"ecr": ecrRegistry,
	}).Info("Call ECR for authorization token ...")

	ecrClient := ecr.New(sess)

	input := &ecr.GetAuthorizationTokenInput{}

	result, err := ecrClient.GetAuthorizationToken(input)
	
	if err != nil {

		return "", time.Time{}, err
	}

	decodedAuthorizationToken, _ := base64.StdEncoding.DecodeString(*result.AuthorizationData[0].AuthorizationToken)
	password := strings.Split(string(decodedAuthorizationToken), ":")[1]
	
	expiryTime := *result.AuthorizationData[0].ExpiresAt

	return password, expiryTime, nil

}
