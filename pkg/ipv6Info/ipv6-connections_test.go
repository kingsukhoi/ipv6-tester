package ipv6Info

import (
	"context"
	"testing"
)

func TestIpV6Google(t *testing.T) {
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name    string
		args    args
		want    *IpV6Result
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				ctx:  context.Background(),
				host: GoogleIpV6Addr,
			},
			want:    &IpV6Result{Success: true, ServerAddress: "2001:4860:4860::8888"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TestIpV6Connection(tt.args.ctx, tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestIpV6Google() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Success != tt.want.Success && got.ServerAddress != tt.want.ServerAddress {
				t.Errorf("TestIpV6Google() got = %v, want %v", got, tt.want)
			}
		})
	}
}
