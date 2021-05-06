package centrifugo

import (
	"reflect"
	"testing"
)

func TestFormMessage(t *testing.T) {
	type args struct {
		message string
	}
	var tests = []struct {
		name string
		args args
		want []byte
	}{
		{name: "simple", args: args{"test"}, want: []byte("{\"message\":\"test\"}")},
		{name: "slashes", args: args{"////\\"}, want: []byte("{\"message\":\"////\\\"}")},
		{name: "numbers", args: args{"123"}, want: []byte("{\"message\":\"123\"}")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormMessage(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FormMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConnToken(t *testing.T) {
	type args struct {
		user  string
		token string
		exp   int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "exp > 0",
			args: args{"112", "safasdf", 123},
			want: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEyMywic3ViIjoiMTEyIn0._UI8MLF8TvsaL9PHKeS88taquQC0xpgQLUdCsVlZD_I",
		},
		{
			name: "exp = 0",
			args: args{"112", "safasdf", 0},
			want: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMTIifQ.qs3LXhIHVvgx39EQ1J_8xf8MDJdsf4ZSgPooTh02JEo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConnToken(tt.args.user, tt.args.token, tt.args.exp); got != tt.want {
				t.Errorf("GetConnToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
