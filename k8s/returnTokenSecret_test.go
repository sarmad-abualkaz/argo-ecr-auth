package callk8s

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestReturnTokenSecret(t *testing.T) {
	type args struct {
		secret      string
		namespace string
		fakesecret string
	}
	tests := []struct {
		name string
		args args
		wantedSecret *v1.Secret
		wantederr string
	}{
		{
			name: "Should succeed",
			args: args{
				secret:      "my-secret",
				namespace: "abualks",
				fakesecret: "ZmFrZXNlY3JldA==",
			},
			wantedSecret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "abualks", 
					Name: "my-secret",
				},
				Data: map[string][]byte{
					"fakesecret": []byte("ZmFrZXNlY3JldA=="),
				},
			},
			wantederr: "",
		},
		{
			name: "Should fail and log not found",
			args: args{
				secret:    "not-found",
				namespace: "abualks",
			},
			wantedSecret: &v1.Secret{},
			wantederr: "secrets \"not-found\" not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var client Client

			if tt.args.secret != "not-found" {

				secret := &v1.Secret{
					ObjectMeta: metav1.ObjectMeta{Namespace: tt.args.namespace, Name: tt.args.secret},
					Data: map[string][]byte{
						"fakesecret": []byte(tt.args.fakesecret),
					},
				}
				
				client = Client{
					Clientset: fake.NewSimpleClientset(secret),
				}

			} else {

				client = Client{
					Clientset: fake.NewSimpleClientset(),
				}

			}
			
			secretcontent, secretcontentErr := ReturnTokenSecret(client, tt.args.secret, tt.args.namespace)

			if secretcontentErr != nil {
				
				if secretcontentErr.Error() != tt.wantederr {
					t.Errorf("ReturnTokenSecret() = %v, want %v", secretcontentErr, tt.wantederr)
				}

			} else if !reflect.DeepEqual(secretcontent, tt.wantedSecret){
				t.Errorf("ReturnTokenSecret() = %v, want %v", secretcontent, tt.wantedSecret)
			}
		})
	}
}
