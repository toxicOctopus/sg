package config

import "testing"

func TestEnv_String(t *testing.T) {
	tests := []struct {
		name string
		e    Env
		want string
	}{
		{
			name: "local",
			e:    0,
			want: "local",
		},
		{
			name: "production",
			e:    2,
			want: "production",
		},
		{
			name: "asdf",
			e:    55,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want Env
	}{
		{
			name: "local",
			args: args{
				s: "local",
			},
			want: Local,
		},
		{
			name: "production",
			args: args{
				s: "production",
			},
			want: Production,
		},
		{
			name: "asdf",
			args: args{
				s: "asdf",
			},
			want: Local,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnvFromString(tt.args.s); got != tt.want {
				t.Errorf("GetEnvFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}