package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/Masterminds/semver"
	"github.com/ethereum/go-ethereum/log"
)

var ErrVersionCheckFailed = errors.New("failed to check the version, try again later")

const (
	maxIdleConns    = 10
	idleConnTimeout = 15 * time.Second

	nightlyBuildMsg = `You use nightly build - %s. 
Nightly builds are intended for development testing and may include bugs and other issues. 
You might want to use the stable releases instead.
The latest stable(recommended) version of lachesis is published on the page: %s.
`

	outdatedBuildMsg = `The latest lachesis version is %s, but you are currently running %s, 
The latest stable(recommended) version of lachesis is published on the page: %s.
`

	githubApiUrl     = "api.github.com"
	releasesListPath = "repos/Fantom-foundation/go-lachesis/releases"
)

type ReleaseVersion struct {
	HtmlUrl         string    `json:"html_url"` // url with page to release
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"` // release name version
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"` // publish date
	ZipballUrl      string    `json:"zipball_url"`  // zip archive download url
	TarballUrl      string    `json:"tarball_url"`  // tar archive download url
}

type BuildStatus uint8

const (
	Actual BuildStatus = iota
	Outdated
	Nightly
)

// CheckRelease - checks the version of the build for the latest release.
func CheckRelease(uri *url.URL, version string) (
	status BuildStatus,
	message string,
	err error,
) {
	if uri == nil {
		uri = &url.URL{Scheme: "https", Host: githubApiUrl, Path: releasesListPath}
	}
	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:       maxIdleConns,
			IdleConnTimeout:    idleConnTimeout,
			DisableCompression: true,
		},
	}
	latestReleaseVersion, err := getLatestReleaseVersion(client, *uri)
	if err != nil {
		err = ErrVersionCheckFailed
		return
	}
	compare, err := compareVersions(version, latestReleaseVersion.Name)
	if err != nil {
		return
	}
	switch compare {
	case -1: // less than the latest
		status = Outdated
		message = fmt.Sprintf(outdatedBuildMsg,
			latestReleaseVersion.Name,
			version,
			latestReleaseVersion.HtmlUrl,
		)
		return
	case 1: // night build
		status = Nightly
		message = fmt.Sprintf(nightlyBuildMsg,
			version,
			latestReleaseVersion.HtmlUrl,
		)
		return
	default:
		status = Actual
		return
	}
}

// getLatestReleaseVersion - gets the list of releases from github API
// and finds among them the latest release by creation date
func getLatestReleaseVersion(client http.Client, url url.URL) (*ReleaseVersion, error) {
	var releases = make([]*ReleaseVersion, 0)
	resp, err := client.Get(url.String())
	if err != nil {
		log.Error("http client Get releases failed", "ur+l", url.String())
		return nil, ErrVersionCheckFailed
	}
	defer resp.Body.Close()

	respB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("read response body failed", "url", url.String())
		return nil, ErrVersionCheckFailed
	}

	err = json.Unmarshal(respB, &releases)
	if err != nil {
		log.Error("unmarshal response body failed", "error", err, "body", string(respB))
		return nil, ErrVersionCheckFailed
	}
	if len(releases) == 0 {
		return nil, ErrVersionCheckFailed
	}

	latestRelease := findLatestRelease(releases)
	if latestRelease == nil {
		return nil, ErrVersionCheckFailed
	}

	return latestRelease, nil
}

// findLatestRelease finds the latest release by field "created_at"
func findLatestRelease(releases []*ReleaseVersion) *ReleaseVersion {
	var (
		latestRelease   *ReleaseVersion
		latestCreatedAt time.Time
	)
	for _, release := range releases {
		if release.CreatedAt.Before(latestCreatedAt) {
			continue
		}
		latestCreatedAt = release.CreatedAt
		latestRelease = release
	}

	return latestRelease
}

// compareVersions returns:
// 	-1 - if version < latest;
// 	0 - if version is the latest;
// 	1 - if version > latest;
func compareVersions(version, latestVersion string) (int, error) {
	currentV, err := semver.NewVersion(version)
	if err != nil {
		return 0, err
	}
	latestV, err := semver.NewVersion(latestVersion)
	if err != nil {
		return 0, err
	}

	return currentV.Compare(latestV), nil
}
