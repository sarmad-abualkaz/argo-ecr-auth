package cmd 

import (
	"testing"
	"time"
	"reflect"
)

func TestTokenExpired(t *testing.T) {
	type args struct {
		tokenExpireyTime time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			// over 2 hrs mark - should return false 
			name: "should return false - 10 hours remaining",
			args: args{
				tokenExpireyTime: time.Now().UTC().Add(time.Hour * 10),
			},
			want: false,
		},
		{
			// over 2 hrs mark - should return false 
			name: "should return false - 3 hours remaining",
			args: args{
				tokenExpireyTime: time.Now().UTC().Add(time.Hour * 3),
			},
			want: false,
		},
		{
			// exactly 1 hrs mark - should return true
			name: "should return true - 1 hours remaining only",
			args: args{
				tokenExpireyTime: time.Now().UTC().Add(time.Hour * 1),
			},
			want: true,
		},
		{
			// exactly now - should return true
			name: "should return true - 1 hours remaining only",
			args: args{
				tokenExpireyTime: time.Now().UTC(),
			},
			want: true,
		},
		{
			// a day in the past - should return true
			name: "should return true - 1 hours remaining only",
			args: args{
				tokenExpireyTime: time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC),
			},
			want: true,
		},
	}
	for _, tt := range tests {

		expireyTimeStr, _ := tt.args.tokenExpireyTime.MarshalText()
		cmData := map[string]string{
			"expireyTime": string(expireyTimeStr),
		}

		got := TokenExpired(cmData)
		if !reflect.DeepEqual(got, tt.want){
			t.Errorf("TokenExpired() = %v, want %v", got, tt.want)
		}
	} 
}
