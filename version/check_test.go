package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckNodeVersion(t *testing.T) {
	assertar := assert.New(t)
	type args struct {
		resp    []*ReleaseVersion
		version string
	}
	tests := []struct {
		name        string
		mockHandler func() http.Handler
		args        args
		want        struct {
			IsNightlyBuild bool
			Message        string
		}
		wantErr error
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
			want: struct {
				IsNightlyBuild bool
				Message        string
			}{},
			wantErr: nil,
		},
		{
			name: "version less than the latest",
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
			want: struct {
				IsNightlyBuild bool
				Message        string
			}{
				Message: `The latest lachesis version is 0.6.1-rc.2, but you are currently running 0.5.0-rc.2, 
The latest stable(recommended) version of lachesis is published on the page: https://github.com/Fantom-foundation/go-lachesis/releases/tag/v0.6.1-rc.2.
Zip archive latest lachesis version: https://api.github.com/repos/Fantom-foundation/go-lachesis/zipball/v0.6.1-rc.2, 
Tar archive latest lachesis version: https://api.github.com/repos/Fantom-foundation/go-lachesis/tarball/v0.6.1-rc.2.`,
			},
			wantErr: nil,
		},
		{
			name: "version is nightly build",
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
				version: "0.6.1-rc.3",
			},
			want: struct {
				IsNightlyBuild bool
				Message        string
			}{
				IsNightlyBuild: true,
				Message: `You use nightly build - 0.6.1-rc.3. 
Nightly builds are intended for development testing and may include bugs and other issues. 
You might want to use the stable releases instead.
The latest stable(recommended) version of lachesis is published on the page: https://github.com/Fantom-foundation/go-lachesis/releases/tag/v0.6.1-rc.2.
Zip archive latest lachesis version: https://api.github.com/repos/Fantom-foundation/go-lachesis/zipball/v0.6.1-rc.2, 
Tar archive latest lachesis version: https://api.github.com/repos/Fantom-foundation/go-lachesis/tarball/v0.6.1-rc.2.`,
			},
			wantErr: nil,
		},
		{
			name: "current version semantic error",
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
				version: "sdfsds0fsd.6.ssf1fsd-sdfrsdc.sdf2",
			},
			want: struct {
				IsNightlyBuild bool
				Message        string
			}{},
			wantErr: semver.ErrInvalidSemVer,
		},
		{
			name: "latest version semantic error",
			mockHandler: func() http.Handler {
				r := http.NewServeMux()
				r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					respB, _ := json.Marshal([]*ReleaseVersion{
						{
							HtmlUrl:    "https://github.com/Fantom-foundation/go-lachesis/releases/tag/v0.6.1-rc.2",
							Name:       "fsad0.adf6.1-fdsrc.2",
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
			want: struct {
				IsNightlyBuild bool
				Message        string
			}{},
			wantErr: semver.ErrInvalidSemVer,
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
			wantErr: errors.New(FailedGetNodeVersionMsg),
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
			wantErr: errors.New(FailedGetNodeVersionMsg),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.mockHandler())
			defer srv.Close()

			uri, _ := url.Parse(srv.URL)
			got, gotErr := CheckNodeVersion(uri, tt.args.version)
			assertError(t, tt.wantErr, gotErr)
			assertar.Equal(tt.want, got)
		})
	}
}

func assertError(t *testing.T, want, got error) {
	switch {
	case want != nil && got != nil:
		require.EqualError(t, got, want.Error())
	case want == nil && got != nil:
		require.Failf(t, "errors not equal", "want: nil, got: %v", got)
	case want != nil && got == nil:
		require.Failf(t, "errors not equal", "want: %v, got: nil", want)
	case want == nil && got == nil:
		return
	}
}
