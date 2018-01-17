package cashlessdevice

import (
	"testing"
)

func Test_validateCrc(t *testing.T) {
	type args struct {
		c []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test true",
			args: args{
				c: []byte{0x03, 0x00, 0x96, 0xF0, 0xF0, 0xF0, 0xF0, 0x00, 0x59},
			},
			want: true,
		},
		{
			name: "test false",
			args: args{
				c: []byte{0x01, 0x00, 0x96, 0xF0, 0xF0, 0xF0, 0xF0, 0x00, 0x59},
			},
			want: false,
		},
		{
			name: "test empty",
			args: args{
				c: []byte{},
			},
			want: false,
		},
		{
			name: "test 1 byte",
			args: args{
				c: []byte{0xFF},
			},
			want: false,
		},
		{
			name: "test 2 byte",
			args: args{
				c: []byte{0x06, 0x06},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateCrc(tt.args.c); got != tt.want {
				t.Errorf("validateCrc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateCrc(t *testing.T) {
	type args struct {
		c []byte
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{
			name: "test 1 byte",
			args: args{
				c: []byte{0x06},
			},
			want: 0x06,
		},
		{
			name: "test 1 byte",
			args: args{
				c: []byte{0x03, 0x00, 0x96, 0xF0, 0xF0, 0xF0, 0xF0, 0x00},
			},
			want: 0x59,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateCrc(tt.args.c); got != tt.want {
				t.Errorf("calculateCrc() = %v, want %v", got, tt.want)
			}
		})
	}
}
