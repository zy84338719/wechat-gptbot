package ability

import "testing"

func Test_moyu_textFunc(t *testing.T) {
	type args struct {
		strin string
		g2    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//mo := moyu{}
			//if got := mo.TextFunc(tt.args.strin, tt.args.g2); got != tt.want {
			//	t.Errorf("TextFunc() = %v, want %v", got, tt.want)
			//}
		})
	}
}
