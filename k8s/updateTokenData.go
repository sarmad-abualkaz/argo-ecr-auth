package callk8s

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// creates and updates configmap where token expirey is stored
func UpdateTokenData(k8sClient Client, cm string, namespace string, expireyTime time.Time, ecrRegistry string, exists bool) error {

	var err error

	cmClient := k8sClient.Clientset.CoreV1().ConfigMaps(namespace)
	
	expireyTimeStr, _ := expireyTime.MarshalText()

	cmManifest := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cm,
			Namespace: namespace,
			Labels:    map[string]string{
				         "argo-ecr-auth": "managed-resource",
			           },
		},
		Data: map[string]string{
			"expireyTime": string(expireyTimeStr),
			"name":        ecrRegistry,
		},
	}

	if exists {
		
		log.WithFields(log.Fields{
			"name":      cm,
			"namespace": namespace,
		}).Info("Updating configmap with token expirey timestamp ...")

		_, err = cmClient.Update(context.TODO(), cmManifest, metav1.UpdateOptions{})
		
	
	} else {

		log.WithFields(log.Fields{
			"name":      cm,
			"namespace": namespace,
		}).Info("Creating configmap with token expirey timestamp ...")

		_, err = cmClient.Create(context.TODO(), cmManifest, metav1.CreateOptions{})

	}

	if err != nil {

		return err
	}
	
	return nil
}
