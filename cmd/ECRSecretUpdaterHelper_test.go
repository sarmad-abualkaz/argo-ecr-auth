package cmd

import (
	"testing"
	"time"
	"reflect"

	callk8s "github.com/sarmad-abualkaz/argo-ecr-auth/k8s"

	"k8s.io/client-go/kubernetes/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
)

func TestECRSecretUpdaterHelper(t *testing.T) {

	type returndata struct{
		_type 		 string
		ecrRegistry  string
		enableOCI 	 string
		label        string
		name 		 string
		username 	 string
		expireyTime  time.Time
	}

	type args struct{
		ecrRegistry  string
		namespace    string
		returncm     bool
		returnsecret bool
		returndata   returndata
		secret       string
	}
	
	tests := []struct {
		name            string
		args            args
		wantupdate      bool
		wantsecretexits bool
		wantcmexits     bool
	}{
		{

			name: "Neither cm or secret exists, update both",
			
			args: args{
				
				ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				namespace:    "abualks",
				returncm:     false,
				returnsecret: false,

				returndata:  returndata{
					_type: 		  "helm",
					ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
					enableOCI: 	  "true",
					label:        "repository",
					name: 		  "ecr",
					username: 	  "AWS",
					expireyTime:  time.Now().UTC().Add(time.Hour * 10),
				},
				
				secret:      "my-secret",
			},
			
			wantupdate:      true,
			wantsecretexits: false,
			wantcmexits:     false,
		},
		{

			name: "Secret exists, but cm does not - update both",
			
			args: args{
				
				ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				namespace:    "abualks",
				returncm:     false,
				returnsecret: true,

				returndata:  returndata{
					_type: 		  "helm",
					ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
					enableOCI: 	  "true",
					label:        "repository",
					name: 		  "ecr",
					username: 	  "AWS",
					expireyTime:  time.Now().UTC().Add(time.Hour * 10),
				},
				
				secret:      "my-secret",
			},
			
			wantupdate:      true,
			wantsecretexits: true,
			wantcmexits:     false,
		},
		{

			name: "cm exists, but secret does not - update both",
			
			args: args{
				
				ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				namespace:    "abualks",
				returncm:     true,
				returnsecret: false,

				returndata:  returndata{
					_type: 		  "helm",
					ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
					enableOCI: 	  "true",
					label:        "repository",
					name: 		  "ecr",
					username: 	  "AWS",
					expireyTime:  time.Now().UTC().Add(time.Hour * 10),
				},
				
				secret:      "my-secret",
			},
			
			wantupdate:      true,
			wantsecretexits: false,
			wantcmexits:     true,
		},
		{

			name: "both exists - but token secret data is wrong - type",
			
			args: args{
				
				ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				namespace:    "abualks",
				returncm:     true,
				returnsecret: true,

				returndata:  returndata{
					_type: 		  "nothelm", //should fail due to this
					ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com", 
					enableOCI: 	  "true",
					label:        "repository",
					name: 		  "ecr",
					username: 	  "AWS",
					expireyTime:  time.Now().UTC().Add(time.Hour * 10), 
				},
				
				secret:      "my-secret",
			},
			
			wantupdate:      true,
			wantsecretexits: true,
			wantcmexits:     true,
		},
		{

			name: "both exists - but token secret data is wrong - ecr",
			
			args: args{
				
				ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				namespace:    "abualks",
				returncm:     true,
				returnsecret: true,

				returndata:  returndata{
					_type: 		  "helm",
					ecrRegistry:  "111111111112.dkr.ecr.us-east-1.amazonaws.com", //should fail due to this
					enableOCI: 	  "true",
					label:        "repository",
					name: 		  "ecr",
					username: 	  "AWS",
					expireyTime:  time.Now().UTC().Add(time.Hour * 10), 
				},
				
				secret:      "my-secret",
			},
			
			wantupdate:      true,
			wantsecretexits: true,
			wantcmexits:     true,
		},
		{

			name: "both exists - but token secret data is wrong - enableOCI",
			
			args: args{
				
				ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				namespace:    "abualks",
				returncm:     true,
				returnsecret: true,

				returndata:  returndata{
					_type: 		  "helm",
					ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
					enableOCI: 	  "false", // should fail due to this
					label:        "repository",
					name: 		  "ecr",
					username: 	  "AWS",
					expireyTime:  time.Now().UTC().Add(time.Hour * 10), 
				},
				
				secret:      "my-secret",
			},
			
			wantupdate:      true,
			wantsecretexits: true,
			wantcmexits:     true,
		},
		{

			name: "both exists - but token secret data is wrong - label",
			
			args: args{
				
				ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				namespace:    "abualks",
				returncm:     true,
				returnsecret: true,

				returndata:  returndata{
					_type: 		  "helm",
					ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
					enableOCI: 	  "true", 
					label:        "notrepository", // should fail due to this
					name: 		  "ecr",
					username: 	  "AWS",
					expireyTime:  time.Now().UTC().Add(time.Hour * 10), 
				},
				
				secret:      "my-secret",
			},
			
			wantupdate:      true,
			wantsecretexits: true,
			wantcmexits:     true,
		},
		{

			name: "both exists - but token secret data is wrong - name",
			
			args: args{
				
				ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				namespace:    "abualks",
				returncm:     true,
				returnsecret: true,

				returndata:  returndata{
					_type: 		  "helm",
					ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
					enableOCI: 	  "true", 
					label:        "repository", 
					name: 		  "notecr", // should fail due to this
					username: 	  "AWS",
					expireyTime:  time.Now().UTC().Add(time.Hour * 10), 
				},
				
				secret:      "my-secret",
			},
			
			wantupdate:      true,
			wantsecretexits: true,
			wantcmexits:     true,
		},
		{

			name: "both exists - but token secret data is wrong - username",
			
			args: args{
				
				ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				namespace:    "abualks",
				returncm:     true,
				returnsecret: true,

				returndata:  returndata{
					_type: 		  "helm",
					ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
					enableOCI: 	  "true", 
					label:        "repository", 
					name: 		  "tecr",
					username: 	  "notAWS", // should fail due to this
					expireyTime:  time.Now().UTC().Add(time.Hour * 10), 
				},
				
				secret:      "my-secret",
			},
			
			wantupdate:      true,
			wantsecretexits: true,
			wantcmexits:     true,
		},
		{

			name: "both exists - but token expired/expiring",
			
			args: args{
				
				ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				namespace:    "abualks",
				returncm:     true,
				returnsecret: true,

				returndata:  returndata{
					_type: 		  "helm",
					ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
					enableOCI: 	  "true",
					label:        "repository",
					name: 		  "ecr",
					username: 	  "AWS",
					expireyTime:  time.Now().UTC().Add(time.Hour * 1), //should fail due to this
				},
				
				secret:      "my-secret",
			},
			
			wantupdate:      true,
			wantsecretexits: true,
			wantcmexits:     true,
		},
		{

			name: "both exists - all data are clear - no update",
			
			args: args{
				
				ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				namespace:    "abualks",
				returncm:     true,
				returnsecret: true,

				returndata:  returndata{
					_type: 		  "helm",
					ecrRegistry:  "111111111111.dkr.ecr.us-east-1.amazonaws.com",
					enableOCI: 	  "true",
					label:        "repository",
					name: 		  "ecr",
					username: 	  "AWS",
					expireyTime:  time.Now().UTC().Add(time.Hour * 10),
				},
				
				secret:      "my-secret",
			},
			
			wantupdate:      false,
			wantsecretexits: true,
			wantcmexits:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var client callk8s.Client

			expireyTimeStr, _ := tt.args.returndata.expireyTime.MarshalText()

			cm := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: tt.args.namespace, 
					Name: tt.args.secret,
				},
				Data: map[string]string{
					"expireyTime": string(expireyTimeStr),
					"name":        tt.args.returndata.ecrRegistry,
				},
			}

			secret := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: tt.args.namespace, 
					Name: tt.args.secret,
					Labels: map[string]string{
						"argocd.argoproj.io/secret-type": tt.args.returndata.label,
					},
				},
				Data: map[string][]byte{
					"enableOCI": []byte(tt.args.returndata.enableOCI),
					"name":      []byte(tt.args.returndata.name),
					"type":      []byte(tt.args.returndata._type),
					"url":       []byte(tt.args.returndata.ecrRegistry),
					"username":  []byte(tt.args.returndata.username),
				},
			}
	
			if tt.args.returncm && tt.args.returnsecret {

				client = callk8s.Client{
					Clientset: fake.NewSimpleClientset(cm, secret),
				}
	
			} else if tt.args.returncm && !tt.args.returnsecret {

				client = callk8s.Client{
					Clientset: fake.NewSimpleClientset(cm),
				}

			} else if !tt.args.returncm && tt.args.returnsecret {

				client = callk8s.Client{
					Clientset: fake.NewSimpleClientset(secret),
				}

			} else {

				client = callk8s.Client{
					Clientset: fake.NewSimpleClientset(),
				}

			}

			gotUpdate, gotsecretexists, gotcmexist := ECRSecretUpdaterHelper(client, tt.args.namespace, tt.args.secret, tt.args.ecrRegistry)

			if !reflect.DeepEqual(gotUpdate, tt.wantupdate){
				t.Errorf("ECRSecretUpdaterHelper() = got for 'update' %v, want %v", gotUpdate, tt.wantupdate)
			}

			if !reflect.DeepEqual(gotsecretexists, tt.wantsecretexits){
				t.Errorf("ECRSecretUpdaterHelper() = got for 'secretexists' %v, want %v", gotsecretexists, tt.wantsecretexits)
			}

			if !reflect.DeepEqual(gotcmexist, tt.wantcmexits){
				t.Errorf("ECRSecretUpdaterHelper() = got for cmexists %v, want %v", gotcmexist, tt.wantcmexits)
			}
		})
	}
}