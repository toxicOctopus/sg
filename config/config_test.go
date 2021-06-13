package config

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"successful",
			args{
				from: GetDefaultValuesPath(),
				to:   filepath.Join(os.TempDir(), "testResultStruct.go"),
			},
			false,
		},
	}
	makeFunctionRunOnRootFolder()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := os.Open(tt.args.to)
			if err == nil {
				err = os.Remove(tt.args.to)
				if err != nil {
					t.Fatal(errors.Wrap(err, "failed to delete temporary file"))
				}
			}
			if err := Generate(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
			content, err := ioutil.ReadFile(tt.args.to)
			if err != nil {
				t.Fatal(err)
			}
			if len(string(content)) == 0 {
				t.Error("Generated config is empty")
			}
			err = os.Remove(tt.args.to)
			if err != nil {
				t.Fatal(errors.Wrap(err, "failed to delete temporary file"))
			}
		})
	}
}

func makeFunctionRunOnRootFolder() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestGetDefaultValuesPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "only 1 case",
			want: filepath.Join("config", "env", "values.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDefaultValuesPath(); got != tt.want {
				t.Errorf("GetDefaultValuesPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRead(t *testing.T) {
	type args struct {
		env        Env
		valuesPath string
	}
	tests := []struct {
		name         string
		args         args
		want         Config
		wantErr      bool
		wantMismatch bool
	}{
		{
			name: "success",
			args: args{
				env:        Local,
				valuesPath: "test_fixtures/values.json",
			},
			want: Config{
				Centrifugo: struct {
					BackendUserID     string `json:"backendUserID"`
					JwtToken          string `json:"jwtToken"`
					TwitchBossChannel string `json:"twitchBossChannel"`
					URL               string `json:"url"`
				}{
					BackendUserID:     "123",
					JwtToken:          "zzz",
					TwitchBossChannel: "bbb",
					URL:               "ws://localhost:8000/connection/websocket",
				},
				ConfigReadInterval: "1m",
				LogLevel:           "debug",
				Twitch: struct {
					Nick string `json:"nick"`
					Pass string `json:"pass"`
				}{
					Nick: "billy",
					Pass: "master",
				},
				Web: struct {
					Host string `json:"host"`
					Port int64  `json:"port"`
				}{
					Host: "localhost",
					Port: 8182,
				},
			},
			wantErr:      false,
			wantMismatch: false,
		},
		{
			name: "fail",
			args: args{
				env:        Local,
				valuesPath: "test_fixtures/values.json",
			},
			want: Config{
				Centrifugo: struct {
					BackendUserID     string `json:"backendUserID"`
					JwtToken          string `json:"jwtToken"`
					TwitchBossChannel string `json:"twitchBossChannel"`
					URL               string `json:"url"`
				}{
					BackendUserID:     "123",
					JwtToken:          "zz123z",
					TwitchBossChannel: "bbb",
					URL:               "ws://localhost:8000/connection/websocket",
				},
				ConfigReadInterval: "1m",
				LogLevel:           "debug",
				Twitch: struct {
					Nick string `json:"nick"`
					Pass string `json:"pass"`
				}{
					Nick: "billy",
					Pass: "masssssster",
				},
				Web: struct {
					Host string `json:"host"`
					Port int64  `json:"port"`
				}{
					Host: "localhost",
					Port: 8182,
				},
			},
			wantErr:      false,
			wantMismatch: true,
		},
	}
	makeFunctionRunOnRootFolder()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(tt.args.env, tt.args.valuesPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && !tt.wantMismatch {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getConfigInterval(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "1m",
			args: args{
				s: "1m",
			},
			want: time.Minute,
		},
		{
			name: "hjgkhjkghj",
			args: args{
				s: "hjgkhjkghj",
			},
			want: defaultUpdateInterval,
		},
		{
			name: "1ms",
			args: args{
				s: "1ms",
			},
			want: minimalUpdateInterval,
		},
		{
			name: "-2h",
			args: args{
				s: "-2h",
			},
			want: minimalUpdateInterval,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConfigInterval(tt.args.s); got != tt.want {
				t.Errorf("getConfigInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}
