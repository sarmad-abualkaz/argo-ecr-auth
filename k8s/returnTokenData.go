package callk8s

import (
	"context"

	"k8s.io/client-go/kubernetes"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ReturnTokenData(k8sClient *kubernetes.Clientset, cm string, namespace string) (map[string]string, error) {

	log.WithFields(log.Fields{
		"namespace": namespace,
		"name": cm,
	}).Info("Looking for configmap data for secret ...")

	cmcontent, err := k8sClient.CoreV1().ConfigMaps(namespace).Get(context.TODO(), cm, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"namespace": namespace,
		"name": cm,
	}).Info("configmap for secret found ...")
	
	return cmcontent.Data, nil
}
