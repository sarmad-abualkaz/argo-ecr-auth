package callk8s

import (
	"context"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
)

func ReturnTokenSecret(k8sClient Client, secret string, namespace string) (*v1.Secret, error) {

	log.WithFields(log.Fields{
		"name":      secret,
		"namespace": namespace,
	}).Info("Looking for secret ...")

	secretcontent, err := k8sClient.Clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secret, metav1.GetOptions{})
	

	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"name":      secret,
		"namespace": namespace,
	}).Info("secret found ...")
	
	return secretcontent, nil
}
