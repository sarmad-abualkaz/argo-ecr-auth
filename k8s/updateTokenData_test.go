package callk8s

import (
	"context"
	"reflect"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestUpdateTokenData(t *testing.T) {
	
	// initialize empty client

	var client Client

	client = Client{
		Clientset: fake.NewSimpleClientset(),
	}

	type args struct {
		cm          string
		namespace   string
		expireyTime time.Time
		name        string
		exists      bool
	}
	tests := []struct {
		name            string
		args            args
		wantexpireyTime string
		wantname        string
		wantederr       string
	}{
		{
			name: "Should fail - no cm exists!",
			args: args{
				cm:          "my-cm",
				namespace:   "abualks",
				expireyTime: time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC),
				name:        "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				exists:      true,
			},
			wantexpireyTime: "2000-02-01T12:30:00Z",
			wantname:        "111111111111.dkr.ecr.us-east-1.amazonaws.com",
			wantederr:       "configmaps \"my-cm\" not found",
		},
		{
			name: "Should succeed - no cm exists, create it!",
			args: args{
				cm:          "my-cm",
				namespace:   "abualks",
				expireyTime: time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC),
				name:        "111111111111.dkr.ecr.us-east-1.amazonaws.com",
				exists:      false,
			},
			wantexpireyTime: "2000-02-01T12:30:00Z",
			wantname:        "111111111111.dkr.ecr.us-east-1.amazonaws.com",
			wantederr:       "",
		},
		{
			name: "Should succeed - cm exists, up it!",
			args: args{
				cm:          "my-cm",
				namespace:   "abualks",
				expireyTime: time.Date(2020, 3, 30, 12, 30, 0, 0, time.UTC),
				name:        "222222222222.dkr.ecr.us-east-1.amazonaws.com",
				exists:      true,
			},
			wantexpireyTime: "2020-03-30T12:30:00Z",
			wantname:        "222222222222.dkr.ecr.us-east-1.amazonaws.com",
			wantederr:       "",
		},
		{
			name: "Should fail - cm exists, cannot create!",
			args: args{
				cm:          "my-cm",
				namespace:   "abualks",
				expireyTime: time.Date(2020, 3, 30, 12, 30, 0, 0, time.UTC),
				name:        "222222222222.dkr.ecr.us-east-1.amazonaws.com",
				exists:      false,
			},
			wantexpireyTime: "2020-03-30T12:30:00Z",
			wantname:        "222222222222.dkr.ecr.us-east-1.amazonaws.com",
			wantederr:       "configmaps \"my-cm\" already exists",
		},
	} 
	for _, tt := range tests {
		
		t.Run(tt.name, func(t *testing.T) {
			
			updateErr := UpdateTokenData(client, tt.args.cm, tt.args.namespace, tt.args.expireyTime, tt.args.name, tt.args.exists)

			if tt.wantederr != "" {
				
				if updateErr.Error() != tt.wantederr {
					t.Errorf("UpdateTokenData() = %v, want %v", updateErr, tt.wantederr)
				}

			} else {

				cmContent, cmContentErr := client.Clientset.CoreV1().ConfigMaps(tt.args.namespace).Get(context.TODO(), tt.args.cm, metav1.GetOptions{})

				if cmContentErr != nil {
					t.Errorf("Recieved Error on UpdateTokenData() - %v", cmContentErr)
				}

				cmContentData := cmContent.Data
				
				if !reflect.DeepEqual(cmContentData["expireyTime"], tt.wantexpireyTime){
					t.Errorf("UpdateTokenData() = %v, want %v", cmContentData["expireyTime"], tt.wantexpireyTime)
				}

				if !reflect.DeepEqual(cmContentData["name"], tt.wantname){
					t.Errorf("UpdateTokenData() = %v, want %v", cmContentData["name"], tt.wantname)
				}
			}
		})
	}
}
