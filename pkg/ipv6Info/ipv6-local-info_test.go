package ipv6Info

import (
	"context"
	"fmt"
	"testing"
)

func TestGetIpv6Addresses(t *testing.T) {
	type args struct {
		ctx           context.Context
		interfaceName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				ctx:           context.Background(),
				interfaceName: "en0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLocalIpv6Addresses(tt.args.ctx, tt.args.interfaceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLocalIpv6Addresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.Addresses) == 0 {
				t.Errorf("GetLocalIpv6Addresses() got empty addresses")
				return
			}
			fmt.Printf("Addresses %s\n", got.Addresses)
		})
	}
}
