package twitch

import (
	"reflect"
	"testing"
)

func TestRegisteredChannels_GetChannel(t *testing.T) {
	type fields struct {
		Channels []Channel
	}
	type args struct {
		name string
	}

	rc := RegisteredChannels{
		Channels: []Channel{
			{
				Name: "weebdog",
			},
			{
				Name: "reckfuru",
			},
			{
				Name: "podasoppin",
			},
			{
				Name: "drrespect",
			},
			{
				Name: "horsen",
			},
		},
	}
	tests := []struct {
		name    string
		args    args
		want    Channel
		wantErr bool
	}{
		{
			name: "channel not found",
			args: args{
				name: "drrrespect",
			},
			want:    Channel{},
			wantErr: true,
		},
		{
			name: "channel found",
			args: args{
				name: "HORSEN",
			},
			want: Channel{
				Name: "horsen",
			},
			wantErr: false,
		},
		{
			name: "channel not found horsen",
			args: args{
				name: "horseny",
			},
			want:    Channel{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rc.GetChannel(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChannel() got = %v, want %v", got, tt.want)
			}
		})
	}
}
