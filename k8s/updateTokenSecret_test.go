package callk8s

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestUpdateTokenSecret(t *testing.T) {
	
	// initialize empty client

	var client Client

	client = Client{
		Clientset: fake.NewSimpleClientset(),
	}

	wantuser      := "AWS"
	wantname      := "ecr"
	wanttype      := "helm"
	wantenableOCI := "true"
	wantLabelVal  := "repository"

	type args struct {
		secret    string
		namespace string
		password  string
		ecr       string
		exists    bool
	}
	tests := []struct {
		name       string
		args       args
		wantpasswd string
		wantecr    string
		wantederr  string
	}{
		{
			name: "Should fail - no secret exists!",
			args: args{
				secret:    "my-secret",
				namespace: "abualks",
				password:  "my-pwd",
				ecr:       "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				exists:    true,
			},
			wantpasswd: "my-pwd",
			wantecr:    "111111111111.dkr.ecr.us-east-1.amazonaws.com",
			wantederr:  "secrets \"my-secret\" not found",
		},
		{
			name: "Should succeed - no secret exists, create it!",
			args: args{
				secret:    "my-secret",
				namespace: "abualks",
				password:  "my-pwd",
				ecr:       "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				exists:    false,
			},
			wantpasswd: "my-pwd",
			wantecr:    "111111111111.dkr.ecr.us-east-1.amazonaws.com",
			wantederr:  "",
		},
		{
			name: "Should succeed - secret exists, update it!",
			args: args{
				secret:    "my-secret",
				namespace: "abualks",
				password:  "my-pwd",
				ecr:       "222222222222.dkr.ecr.us-east-1.amazonaws.com",
				exists:    true,
			},
			wantpasswd: "my-pwd",
			wantecr:    "222222222222.dkr.ecr.us-east-1.amazonaws.com",
			wantederr:  "",
		},
		{
			name: "Should fail - secret exists, cannot create!",
			args: args{
				secret:    "my-secret",
				namespace: "abualks",
				password:  "my-pwd",
				ecr:       "222222222222.dkr.ecr.us-east-1.amazonaws.com",
				exists:    false,
			},
			wantpasswd: "my-pwd",
			wantecr:    "222222222222.dkr.ecr.us-east-1.amazonaws.com",
			wantederr:  "secrets \"my-secret\" already exists",
		},
	} 
	for _, tt := range tests {
		
		t.Run(tt.name, func(t *testing.T) {
			
			updateErr := UpdateTokenSecret(client, tt.args.secret, tt.args.namespace, tt.args.password, tt.args.ecr, tt.args.exists)

			if tt.wantederr != "" {
				
				if updateErr.Error() != tt.wantederr {
					t.Errorf("UpdateTokenSecret() = %v, want %v", updateErr, tt.wantederr)
				}

			} else {

				secretContent, secretContentErr := client.Clientset.CoreV1().Secrets(tt.args.namespace).Get(context.TODO(), tt.args.secret, metav1.GetOptions{})

				if secretContentErr != nil {
					t.Errorf("Recieved Error on UpdateTokenData() - %v", secretContentErr)
				}

				fmt.Println(secretContent)

				secretContentLabels := secretContent.ObjectMeta.Labels

				if !reflect.DeepEqual(secretContentLabels["argocd.argoproj.io/secret-type"], wantLabelVal){
					t.Errorf("UpdateTokenSecret() = %v, want %v", secretContentLabels["argocd.argoproj.io/secret-type"], wantLabelVal)
				}

				secretContentData := secretContent.StringData
				
				if !reflect.DeepEqual(secretContentData["url"], tt.wantecr){
					t.Errorf("UpdateTokenSecret() = got url %v, want %v", secretContentData["url"], tt.wantecr)
				}

				if !reflect.DeepEqual(secretContentData["name"], wantname){
					t.Errorf("UpdateTokenSecret() = got name %v, want %v", secretContentData["name"], wantname)
				}

				if !reflect.DeepEqual(secretContentData["type"], wanttype){
					t.Errorf("UpdateTokenSecret() = got type %v, want %v", secretContentData["type"], wanttype)
				}

				if !reflect.DeepEqual(secretContentData["enableOCI"], wantenableOCI){
					t.Errorf("UpdateTokenSecret() = got enableOCI %v, want %v", secretContentData["enableOCI"], wantenableOCI)
				}

				if !reflect.DeepEqual(secretContentData["username"], wantuser){
					t.Errorf("UpdateTokenSecret() = got username %v, want %v", secretContentData["username"], wantuser)
				}

				if !reflect.DeepEqual(secretContentData["password"], tt.wantpasswd){
					t.Errorf("UpdateTokenSecret() = got password %v, want %v", secretContentData["password"], tt.wantpasswd)
				}
			}
		})
	}
}
