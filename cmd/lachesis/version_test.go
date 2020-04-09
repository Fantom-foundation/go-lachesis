package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/urfave/cli.v1"

	vers "github.com/Fantom-foundation/go-lachesis/version"
)

func Test_checkNodeVersion(t *testing.T) {
	type args struct {
		uri     *url.URL
		ctx     *cli.Context
		cfg     *config
		version string
	}
	t.Run("node version is latest", func(t *testing.T) {
		srv := httptest.NewServer(func() http.Handler {
			r := http.NewServeMux()
			r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				respB, _ := json.Marshal([]*vers.ReleaseVersion{
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
		}())
		defer srv.Close()
		uri, _ := url.Parse(srv.URL)
		ctx := cli.NewContext(cli.NewApp(), flag.NewFlagSet("test", flag.ContinueOnError), nil)
		cfg := makeAllConfigs(ctx)

		tt := args{
			uri:     uri,
			ctx:     ctx,
			cfg:     &cfg,
			version: "0.6.1-rc.2",
		}
		checkNodeVersion(tt.uri, tt.ctx, tt.cfg, tt.version)
		assert.False(t, cfg.Lachesis.DisablePrivateAccountAPI)
	})

	t.Run("node version is nightly", func(t *testing.T) {
		srv := httptest.NewServer(func() http.Handler {
			r := http.NewServeMux()
			r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				respB, _ := json.Marshal([]*vers.ReleaseVersion{
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
		}())
		defer srv.Close()
		uri, _ := url.Parse(srv.URL)
		ctx := cli.NewContext(cli.NewApp(), flag.NewFlagSet("test", flag.ContinueOnError), nil)
		cfg := makeAllConfigs(ctx)

		tt := args{
			uri:     uri,
			ctx:     ctx,
			cfg:     &cfg,
			version: "0.6.1-rc.3",
		}
		checkNodeVersion(tt.uri, tt.ctx, tt.cfg, tt.version)
		assert.False(t, cfg.Lachesis.DisablePrivateAccountAPI)
	})

	t.Run("node version is less than latest", func(t *testing.T) {
		srv := httptest.NewServer(func() http.Handler {
			r := http.NewServeMux()
			r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				respB, _ := json.Marshal([]*vers.ReleaseVersion{
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
		}())
		defer srv.Close()
		uri, _ := url.Parse(srv.URL)
		ctx := cli.NewContext(cli.NewApp(), flag.NewFlagSet("test", flag.ContinueOnError), nil)
		cfg := makeAllConfigs(ctx)

		tt := args{
			uri:     uri,
			ctx:     ctx,
			cfg:     &cfg,
			version: "0.6.0",
		}
		checkNodeVersion(tt.uri, tt.ctx, tt.cfg, tt.version)
		assert.True(t, cfg.Lachesis.DisablePrivateAccountAPI)
	})
}
