package multihashmapper

import (
	"net/url"
	"testing"

	"github.com/multiformats/go-multihash"
)

func TestMakeMap(t *testing.T) {
	type fields struct {
		multihashCode uint64
	}
	type args struct {
		link string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "default sha2-256",
			fields:  fields{multihashCode: multihash.Names["sha2-256"]},
			args:    args{link: "https://localhost:8080/foo?bar"},
			want:    "Qmf44mf3WDz3a1hpAg9zP3ZjgirGc7wBh7PsGyGxJUQsSM",
			wantErr: false,
		},
		{
			name:    "README.md sha2-256",
			fields:  fields{multihashCode: multihash.Names["sha2-256"]},
			args:    args{link: "https://example.com/foo?bar"},
			want:    "Qma3YMYZUNAY7Dp7UhtZfqKAfsLkHyF9jf1yFXjZbYjWqt",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mapper{
				multihashCode: tt.fields.multihashCode,
			}
			link, err := url.Parse(tt.args.link)
			if err != nil {
				t.Errorf("error parsing url %v, %v", tt.args.link, err)
			}

			got, err := m.Map(link)
			if (err != nil) != tt.wantErr {
				t.Errorf("mapper.Map() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("mapper.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
