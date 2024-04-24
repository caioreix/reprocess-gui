package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCaller(t *testing.T) {
	type args struct {
		skip int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get caller skip 0",
			args: args{skip: 0},
			want: "runtime.Callers",
		},
		{
			name: "Get caller skip 1",
			args: args{skip: 1},
			want: "utils.GetCaller",
		},
		{
			name: "Get caller skip 2",
			args: args{skip: 2},
			want: "utils.TestGetCaller.func1",
		},
		{
			name: "Get caller skip 3",
			args: args{skip: 3},
			want: "testing.tRunner",
		},
		{
			name: "Get caller skip out of range",
			args: args{skip: 10},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCaller(tt.args.skip)
			assert.Equal(t, tt.want, got)
		})
	}
}
