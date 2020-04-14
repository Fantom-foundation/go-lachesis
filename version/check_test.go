package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Masterminds/semver"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestCheckNodeVersion(t *testing.T) {
	require := require.New(t)

	tests := []struct {
		name    string
		server  func() http.Handler
		version string
		result  BuildStatus
		err     error
	}{
		{
			name: "actual build",
			server: func() http.Handler {
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
			version: "0.6.1-rc.2",
			result:  Actual,
		},
		{
			name: "outdated build",
			server: func() http.Handler {
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
			version: "0.5.0-rc.2",
			result:  Outdated,
		},
		{
			name: "nightly build",
			server: func() http.Handler {
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
			version: "0.6.1-rc.3",
			result:  Nightly,
		},
		{
			name: "current version semantic error",
			server: func() http.Handler {
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
			version: "sdfsds0fsd.6.ssf1fsd-sdfrsdc.sdf2",
			err:     semver.ErrInvalidSemVer,
		},
		{
			name: "latest version semantic error",
			server: func() http.Handler {
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
			version: "0.6.1-rc.2",
			err:     semver.ErrInvalidSemVer,
		},
		{
			name: "read releases error",
			server: func() http.Handler {
				r := http.NewServeMux()
				r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, "error")
				})
				return r
			},
			version: "0.6.1-rc.2",
			err:     ErrVersionCheckFailed,
		},
		{
			name: "versions empty",
			server: func() http.Handler {
				r := http.NewServeMux()
				r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					respB, _ := json.Marshal([]*ReleaseVersion{})
					fmt.Fprint(w, string(respB))
				})

				return r
			},
			version: "0.5.0-rc.2",
			err:     ErrVersionCheckFailed,
		},
	}
	for _, exp := range tests {
		t.Run(exp.name, func(t *testing.T) {
			srv := httptest.NewServer(exp.server())
			defer srv.Close()

			uri, _ := url.Parse(srv.URL)
			got, _, err := CheckRelease(uri, exp.version)
			require.Equal(exp.err, err)
			require.Equal(exp.result, got)
		})
	}

	t.Run("current", func(t *testing.T) {
		cur := params.VersionWithCommit("anyCommit", "anyDate")
		_, _, err := CheckRelease(nil, cur)
		require.NotEqual(semver.ErrInvalidSemVer, err)
		require.NotEqual(semver.ErrInvalidSemVer, err)
	})
}
