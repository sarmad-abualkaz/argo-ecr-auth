package callk8s

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestReturnTokenData(t *testing.T) {
	type args struct {
		cm          string
		namespace   string
		expireyTime string
		name        string
	}
	tests := []struct {
		name            string
		args            args
		wantexpireyTime string
		wantname        string
		wantederr       string
	}{
		{
			name: "Should succeed",
			args: args{
				cm:          "my-cm",
				namespace:   "abualks",
				expireyTime: "2000-02-01T10:25:30.456Z",
				name:        "111111111111.dkr.ecr.us-east-1.amazonaws.com",
			},
			wantexpireyTime: "2000-02-01T10:25:30.456Z",
			wantname:        "111111111111.dkr.ecr.us-east-1.amazonaws.com",
			wantederr:       "",
		},
		{
			name: "Should fail and log not found",
			args: args{
				cm:          "",
				namespace:   "",
				expireyTime: "",
				name:        "",
			},
			wantexpireyTime: "",
			wantname:        "",
			wantederr:       "configmaps \"\" not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var client Client

			if tt.args.cm != "" {

				cm := &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: tt.args.namespace, 
						Name: tt.args.cm,
					},
					Data: map[string]string{
						"expireyTime": tt.args.expireyTime,
						"name":        tt.args.name,
					},
				}
				
				client = Client{
					Clientset: fake.NewSimpleClientset(cm),
				}

			} else {

				client = Client{
					Clientset: fake.NewSimpleClientset(),
				}

			}
			
			cmContentData, cmErr := ReturnTokenData(client, tt.args.cm, tt.args.namespace)

			if !reflect.DeepEqual(cmContentData["expireyTime"], tt.wantexpireyTime){
				t.Errorf("ReturnTokenData() = %v, want %v", cmContentData["expireyTime"], tt.wantexpireyTime)
			}

			if !reflect.DeepEqual(cmContentData["name"], tt.wantname){
				t.Errorf("ReturnTokenData() = %v, want %v", cmContentData["name"], tt.wantname)
			}

			if cmErr != nil {
				
				if cmErr.Error() != tt.wantederr {
					t.Errorf("ReturnTokenData() = %v, want %v", cmErr, tt.wantederr)
				}

			}
		})
	}
}
