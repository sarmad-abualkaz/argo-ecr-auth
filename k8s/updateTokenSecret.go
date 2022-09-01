package callk8s

import (
	"context"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


func UpdateTokenSecret(k8sClient Client, secret string, namespace string, password string, ecrRegistry string, exists bool) error {

	var err error

	secertsClient := k8sClient.Clientset.CoreV1().Secrets(namespace)

	secreManifest := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secret,
			Namespace: namespace,
			Labels: map[string]string{
				"argocd.argoproj.io/secret-type": "repository",
			},
		},
		Type: "Opaque",
		StringData: map[string]string{
			"url": ecrRegistry,
			"name": "ecr",
			"type": "helm",
			"enableOCI": "true",
			"username": "AWS",
			"password": password,
		},
	}

	if exists {
		
		log.WithFields(log.Fields{
			"namespace": namespace,
			"name": secret,
		}).Info("Updating secret ...")

		_, err = secertsClient.Update(context.TODO(), secreManifest, metav1.UpdateOptions{})
		
	
	} else {

		log.WithFields(log.Fields{
			"namespace": namespace,
			"name": secret,
		}).Info("Creating secret ...")

		_, err = secertsClient.Create(context.TODO(), secreManifest, metav1.CreateOptions{})

	}

	if err != nil {

		return err
	}
	
	return nil
}
