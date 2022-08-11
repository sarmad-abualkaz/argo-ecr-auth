package cmd

import (
	v1 "k8s.io/api/core/v1"
)

func ValidateECRSecret(ecrRegistry string, secretContent *v1.Secret) (bool, string, string, string) {

	_ecrRegistry := string(secretContent.Data["url"])
	_name := string(secretContent.Data["name"])
	_type := string(secretContent.Data["type"])
	_enableOCI := string(secretContent.Data["enableOCI"])
	_username := string(secretContent.Data["username"])
	_labels := secretContent.ObjectMeta.Labels["argocd.argoproj.io/secret-type"]

	if _ecrRegistry != ecrRegistry {
		return false, "secretContent.StringData['url']", _ecrRegistry, ecrRegistry
	}

	if _name != "ecr" {
		return false, "secretContent.StringData['name']", _name, "ecr"
	}

	if _type != "helm" {
		return false, "secretContent.StringData['type']", _type, "helm"
	}

	if _enableOCI != "true" {
		return false, "secretContent.StringData['enableOCI']", _enableOCI, "true"
	}

	if _username != "AWS" {
		return false, "secretContent.StringData['username']", _username, "AWS"
	}

	if _labels != "repository" {
		return false, "secretContent.ObjectMeta.Labes['argocd.argoproj.io/secret-type']", _labels, "repository"
	}

	return true, "", "", ""
}
