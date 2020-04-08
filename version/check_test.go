package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestCheckNodeVersion(t *testing.T) {
	type args struct {
		resp    []*ReleaseVersion
		version string
	}
	tests := []struct {
		name        string
		mockHandler func() http.Handler
		args        args
		wantErr     bool
	}{
		{
			name: "version equal the latest",
			mockHandler: func() http.Handler {
				r := http.NewServeMux()
				r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					respB, _ := json.Marshal([]*ReleaseVersion{
						{
							HtmlUrl:    "https://github.com/Fantom-foundation/go-lachesis/releases/tag/v0.6.1-rc.2",
							Name:       "0.6.1-rc.2",
							CreatedAt:  time.Date(2020, 01, 4, 10, 20, 14, 0, time.UTC),
							ZipballUrl: "https://api.github.com/repos/Fantom-foundation/go-lachesis/zipball/v0.6.1-rc.2",
							TarballUrl: "https://api.github.com/repos/Fantom-foundation/go-lachesis/tarball/v0.6.1-rc.2",
						},
						{
							HtmlUrl:    "https://github.com/Fantom-foundation/go-lachesis/releases/tag/v0.5.0-rc.2",
							Name:       "0.5.0-rc.2",
							CreatedAt:  time.Date(2019, 12, 19, 10, 20, 30, 0, time.UTC),
							ZipballUrl: "https://api.github.com/repos/Fantom-foundation/go-lachesis/zipball/v0.5.0-rc.2",
							TarballUrl: "https://api.github.com/repos/Fantom-foundation/go-lachesis/tarball/v0.5.0-rc.2",
						},
					})
					fmt.Fprint(w, string(respB))
				})

				return r
			},
			args: args{
				version: "0.6.1-rc.2",
			},
			wantErr: false,
		},
		{
			name: "version not equal the latest",
			mockHandler: func() http.Handler {
				r := http.NewServeMux()
				r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					respB, _ := json.Marshal([]*ReleaseVersion{
						{
							HtmlUrl:    "https://github.com/Fantom-foundation/go-lachesis/releases/tag/v0.6.1-rc.2",
							Name:       "0.6.1-rc.2",
							CreatedAt:  time.Date(2020, 01, 4, 10, 20, 14, 0, time.UTC),
							ZipballUrl: "https://api.github.com/repos/Fantom-foundation/go-lachesis/zipball/v0.6.1-rc.2",
							TarballUrl: "https://api.github.com/repos/Fantom-foundation/go-lachesis/tarball/v0.6.1-rc.2",
						},
						{
							HtmlUrl:    "https://github.com/Fantom-foundation/go-lachesis/releases/tag/v0.5.0-rc.2",
							Name:       "0.5.0-rc.2",
							CreatedAt:  time.Date(2019, 12, 19, 10, 20, 30, 0, time.UTC),
							ZipballUrl: "https://api.github.com/repos/Fantom-foundation/go-lachesis/zipball/v0.5.0-rc.2",
							TarballUrl: "https://api.github.com/repos/Fantom-foundation/go-lachesis/tarball/v0.5.0-rc.2",
						},
					})
					fmt.Fprint(w, string(respB))
				})

				return r
			},
			args: args{
				version: "0.5.0-rc.2",
			},
			wantErr: true,
		},
		{
			name: "read releases error",
			mockHandler: func() http.Handler {
				r := http.NewServeMux()
				r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, "error")
				})
				return r
			},
			args: args{
				version: "0.6.1-rc.2",
			},
			wantErr: true,
		},
		{
			name: "versions empty",
			mockHandler: func() http.Handler {
				r := http.NewServeMux()
				r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					respB, _ := json.Marshal([]*ReleaseVersion{})
					fmt.Fprint(w, string(respB))
				})

				return r
			},
			args: args{
				version: "0.5.0-rc.2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.mockHandler())
			defer srv.Close()

			uri, _ := url.Parse(srv.URL)
			if err := CheckNodeVersion(uri, tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("CheckNodeVersion() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
