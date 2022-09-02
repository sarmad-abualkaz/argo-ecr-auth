package cmd

import (
	"testing"
	"reflect"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
)

func TestValidateECRSecret(t *testing.T) {
	
	expEcrRegistry := "111111111111.dkr.ecr.us-east-1.amazonaws.com"

	type args struct {
		_type 		string
		ecrRegistry string
		enableOCI 	string
		label       string
		name 		string
		username 	string
	}
	tests := []struct{
		args 		   args
		expEcrRegistry string
		name 		   string
		unwantedstring string
		wantcomment    string
		wantresult 	   bool
		wantstring     string
	}{
		{
			name:            "Should return false - wrong ecr url",
			expEcrRegistry:   expEcrRegistry,
			args: args{
				_type: 		 "helm",
				ecrRegistry: "111111111122.dkr.ecr.us-east-1.amazonaws.com",
				enableOCI:   "true",
				label:       "repository",
				name:        "ecr",
				username:    "AWS",
			},
			wantresult:     false,
			wantcomment:    "secretContent.StringData['url']",
			unwantedstring: "111111111122.dkr.ecr.us-east-1.amazonaws.com",
			wantstring:     "111111111111.dkr.ecr.us-east-1.amazonaws.com",
		},
		{
			name:            "Should return false - wrong name",
			expEcrRegistry:   expEcrRegistry,
			args: args{
				_type: 		 "helm",
				ecrRegistry: "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				enableOCI:   "true",
				label:       "repository",
				name:        "notecr",
				username:    "AWS",
			},
			wantresult:     false,
			wantcomment:    "secretContent.StringData['name']",
			unwantedstring: "notecr",
			wantstring:     "ecr",
		},
		{
			name:            "Should return false - wrong type",
			expEcrRegistry:  expEcrRegistry,
			args: args{
				_type: 		 "nothelm",
				ecrRegistry: "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				enableOCI:   "true",
				label:       "repository",
				name:        "ecr",
				username:    "AWS",
			},
			wantresult:     false,
			wantcomment:    "secretContent.StringData['type']",
			unwantedstring: "nothelm",
			wantstring:     "helm",
		},
		{
			name:            "Should return false - wrong enableOCI",
			expEcrRegistry:  expEcrRegistry,
			args: args{
				_type: 		 "helm",
				ecrRegistry: "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				enableOCI:   "false",
				label:       "repository",
				name:        "ecr",
				username:    "AWS",
			},
			wantresult:     false,
			wantcomment:    "secretContent.StringData['enableOCI']",
			unwantedstring: "false",
			wantstring:     "true",
		},
		{
			name:            "Should return false - wrong username",
			expEcrRegistry:  expEcrRegistry,
			args: args{
				_type: 		 "helm",
				ecrRegistry: "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				enableOCI:   "true",
				label:       "repository",
				name:        "ecr",
				username:    "notAWS",
			},
			wantresult:     false,
			wantcomment:    "secretContent.StringData['username']",
			unwantedstring: "notAWS",
			wantstring:     "AWS",
		},
		{
			name:            "Should return false - wrong labels",
			expEcrRegistry:  expEcrRegistry,
			args: args{
				_type: 		 "helm",
				ecrRegistry: "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				enableOCI:   "true",
				label:       "notrepository",
				name:        "ecr",
				username:    "AWS",
			},
			wantresult:     false,
			wantcomment:    "secretContent.ObjectMeta.Labes['argocd.argoproj.io/secret-type']",
			unwantedstring: "notrepository",
			wantstring:     "repository",
		},
	}
	for _, tt := range tests {
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"argocd.argoproj.io/secret-type": tt.args.label,
				},
			},
			Data: map[string][]byte{
				"enableOCI": []byte(tt.args.enableOCI),
				"name":      []byte(tt.args.name),
				"type":      []byte(tt.args._type),
				"url":       []byte(tt.args.ecrRegistry),
				"username":  []byte(tt.args.username),
			},
		}

		gotresutl, gotcomment, gotunwantedstring, gotstring := ValidateECRSecret(tt.expEcrRegistry, secret)

		if !reflect.DeepEqual(gotresutl, tt.wantresult){
			t.Errorf("ValidateECRSecret() = gotresult %v, want %v", gotresutl, tt.wantresult)
		}

		if !reflect.DeepEqual(gotcomment, tt.wantcomment){
			t.Errorf("ValidateECRSecret() = gotcomment %v, want %v", gotcomment, tt.wantcomment)
		}

		if !reflect.DeepEqual(gotunwantedstring, tt.unwantedstring){
			t.Errorf("ValidateECRSecret() = gotunwantedstring %v, want %v", gotunwantedstring, tt.unwantedstring)
		}

		if !reflect.DeepEqual(gotstring, tt.wantstring){
			t.Errorf("ValidateECRSecret() = gotstring %v, want %v", gotunwantedstring, tt.unwantedstring)
		}
	}
}
