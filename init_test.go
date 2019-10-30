package fluentdpub

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	"reflect"
	"testing"
)

func Test_parseEnvVar(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name     string
		args     args
		want     *fluent.Config
		wantSubj string
		wantErr  bool
	}{
		{
			name: "simplest",
			args: args{
				env: "",
			},
			want: &fluent.Config{
				FluentPort:    24224,
				FluentHost:    "127.0.0.1",
				FluentNetwork: "tcp",
			},
			wantSubj: "",
			wantErr:  false,
		},
		{
			name: "hostname",
			args: args{
				env: "://localhost",
			},
			want: &fluent.Config{
				FluentPort:    24224,
				FluentHost:    "localhost",
				FluentNetwork: "tcp",
			},
			wantSubj: "",
			wantErr:  false,
		},
		{
			name: "hostname 2",
			args: args{
				env: "tcp://localhost:22222",
			},
			want: &fluent.Config{
				FluentPort:    22222,
				FluentHost:    "localhost",
				FluentNetwork: "tcp",
			},
			wantSubj: "",
			wantErr:  false,
		},
		{
			name: "udp",
			args: args{
				env: "udp://localhost:22222",
			},
			want: &fluent.Config{
				FluentPort:    22222,
				FluentHost:    "localhost",
				FluentNetwork: "udp",
			},
			wantSubj: "",
			wantErr:  false,
		},
		{
			name: "wrong scheme",
			args: args{
				env: "wss://localhost:22222",
			},
			wantErr: true,
		},
		{
			name: "tag prefix",
			args: args{
				env: "tcp://localhost/tag.prefix",
			},
			want: &fluent.Config{
				FluentPort:    24224,
				FluentHost:    "localhost",
				FluentNetwork: "tcp",
			},
			wantSubj: "tag.prefix",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotSubj, err := parseEnvVar(tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseEnvVar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEnvVar() got = %v, want %v", got, tt.want)
			}
			if gotSubj != tt.wantSubj {
				t.Errorf("parseEnvVar() got = %v, want %v", gotSubj, tt.wantSubj)
			}
		})
	}
}
