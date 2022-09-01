package callecr

import (
	"encoding/base64"
	"testing"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
)

type mockedGeneratedToken struct {
	ecriface.ECRAPI
	Resp ecr.GetAuthorizationTokenOutput
}

func (m mockedGeneratedToken) GetAuthorizationToken(in *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
	// Only need to return mocked response output
	return &m.Resp, nil
}


func TestGenerateECRTokent(t *testing.T){

	type args struct {
		username string
		password string
		expirationTime time.Time
	}

	tests := []struct {
		name string
		args args
		wantpwd string
		wanttime time.Time
	}{
		{
			name: "should return a response",
			args: args{
				username: 		"AWS",
				password: 		"testpw1",
				expirationTime: time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC),
			},
			wantpwd: "testpw1",
			wanttime: time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC),
		},
	} 
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
		
			planeAuthorizationToken := tt.args.username + ":" + tt.args.password

			encodedAuthorizationToken := base64.StdEncoding.EncodeToString([]byte(planeAuthorizationToken))

			response := ecr.GetAuthorizationTokenOutput{
				AuthorizationData: []*ecr.AuthorizationData{
					{
						AuthorizationToken: &encodedAuthorizationToken,
						ExpiresAt: &tt.args.expirationTime,
					},	
				},
			}	

			ecrMockedClient := ECRClient{
				Client: mockedGeneratedToken{Resp: response},
			}
			
			pwd, genTime, _ := GenerateECRTokentHelper(ecrMockedClient)

			if !reflect.DeepEqual(pwd, tt.wantpwd){
				t.Errorf("GenerateECRTokentHelper() = %v, want %v", pwd, tt.wantpwd)
			}

			if !reflect.DeepEqual(genTime, tt.wanttime){
				t.Errorf("GenerateECRTokentHelper() = %v, want %v", genTime, tt.wanttime)
			}

		})
	}

}
