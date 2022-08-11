package main

import (
	"flag"
	"github.com/sarmad-abualkaz/argo-ecr-auth/cmd"

	log "github.com/sirupsen/logrus"
)

func main() {

	awsProfile := flag.String("aws-profile", "", "name of aws profile")
	awsRegion := flag.String("aws-region", "us-east-1", "aws region")
	ecrRegistry :=  flag.String("ecr-registry", "", "name of object for cert data")
	k8sConfig := flag.String("kube-config", "in-cluster", "kubeconfig setup")
	namespace := flag.String("namespace", "argocd", "kubernetes namespace where secret exists")
	secretName := flag.String("secret-name", "ecr-auth", "kubernetes secret name to sync from/to")
	sleep := flag.Int("sleep-between-checks", 20, "sleep time between syncs")
	
	flag.Parse()

	// log program starting
	log.WithFields(log.Fields{
		"aws-profile": *awsProfile,
		"aws-region": *awsRegion,
		"ecrRegistry": *ecrRegistry,
		"kube-config": *k8sConfig,
		"namespace": *namespace,
		"secret-name": *secretName,
		"sleeping": *sleep,
	  }).Info("program started ...")


	// call cmd.updatECRSecret with params
	cmd.ECRSecretUpdater(*awsProfile, *awsRegion, *ecrRegistry, *k8sConfig, *namespace, *secretName, *sleep)

}
