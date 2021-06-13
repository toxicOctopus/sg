package twitch

import "testing"

func TestMessageIsEmote(t *testing.T) {
	type args struct {
		message string
		list    EmoteList
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "is emote",
			args: args{
				message: "Kappa",
				list: EmoteList{
					{
						Name:      "Kappa",
					},
					{
						Name:      "Keepa",
					},
					{
						Name:      "Kapp",
					}},
			},
			want: true,
		},
		{
			name: "is not emote",
			args: args{
				message: "FeelsBadMan",
				list: EmoteList{
					{
						Name:      "FeelsBadManButActuallyGoodMan",
					},
					{
						Name:      "FeelsBadman",
					},
					{
						Name:      "feelsbadman",
					}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MessageIsEmote(tt.args.message, tt.args.list); got != tt.want {
				t.Errorf("MessageIsEmote() = %v, want %v", got, tt.want)
			}
		})
	}
}
