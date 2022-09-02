# argo-ecr-auth
A tool to automate authentication to ECR helm (OCI) for ArgoCD

## Purpose
Argo CD does not have a native mechanism to stay authetnicated to ECR (given that ECR requires a short-lived token, with maximum 12 hours expirey time). 

Argo CD gives the option to provied a secret with specific labels to provide authentication details (e.g. username/password etc.) to authetnicate to private registries.

This tool creates and continously updates a secret in the namespace where Argo CD is running. The tool authenticates against the ECR endpoint - `<account-id>.dkr.ecr.us-east-1.amazonaws.com>` and updates the secret with the token.

## How to install
This project can be deployed to a Kubenretes cluster via Helm using the following:

```
helm repo add sarmad-helm-charts https://sarmad-abualkaz.github.io/my-helm-charts/

helm install <Release-Name> sarmad-helm-charts/argo-ecr-auth --set args='{"--ecr-registry=<ECR-Endpoint>","--aws-region=<AWS-Region>", "--namespace=<ArgoCD-Namspace>}"'
```


Flags to note:

| flag | purpose | default |
| --- | --- | --- | 
|`--aws-profile` | aws profile to use in `~/<user>/.aws` folder. If set to empty it will perform proper aws-cred cascade. Set to empty to make use of AssumeWebIdentity through a service account. |``|
|`--aws-region` | aws region to target. |`us-east-1` |
| `--ecr-registry` | Name of ECR registry to authenticate to. | `` |
| `--kube-config` | where the process is running, i.e. how kubeconfig will be setup. `"in-cluster"` and `local` are the only other acceptable options. | `in-cluster` |
| `--namespace` | namespace where the secret is stored. | `"argocd"` |
| `--secret-name` | kubernetes secret name to sync from/to. | `ecr-auth` |
| `--sleep-between-checks` | sleep time between syncs in seconds. | 120 |

## How it works
The tool calls AWS ECR to retrieve the token (i.e. auth password) and expirey time. It will then updated the ecr-auth secret (the name of the secret can be specified at start time). 

The secret will consist of data equivalent to the below structure:

```
apiVersion: v1
kind: Secret
metadata:
  name: <secret-name>
  namespace: <argocd-namespace>
  labels:
    argocd.argoproj.io/secret-type: repository
type: Opaque
StringData:
  enableOCI: true
  name: ecr
  password: <passpword>
  type: helm
  url: <ecr url>
  username: AWS
```

The tool will also create and continously update a ConfigMap on the same namepsace with data below:

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: <secret-name>
  namespace: <argocd-namespace>
  labels:
    argo-ecr-auth: managed-resource
data:
  expireyTime: <string of time.Time formate>
  name: <ecr url>
```

## Permission Reqiured
The tool requires a cluster-admin permission on the namespace (not across the cluster).

For authetnication to an ECR in the same AWS account as where the pod is running (same account as the EKS cluster), relying on the pod to use the same IAM role as the ec2 instance profile is sufficient.

While cross-account authetnication is not tested, however using an IAM role for service account might be the correct option here.
