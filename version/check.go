package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

const (
	maxIdleConns    = 10
	idleConnTimeout = 15 * time.Second

	FailedGetNodeVersionMsg = "failed to check the latest version of the node, try again later"
	versionsNotEqualMsgPtrn = `The current version - %s of the node is not the latest - %s, 
please update the node from here - %s to continue.
If it doesn’t work, download the zip archive from here - %s, 
or tar archive from here - %s.`

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

// CheckNodeVersion - checks the version of the node for the latest
func CheckNodeVersion(uri *url.URL, version string) error {
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
		return err
	}
	if compare := strings.Compare(version, latestReleaseVersion.Name); compare < 0 {
		return fmt.Errorf(
			versionsNotEqualMsgPtrn,
			version,
			latestReleaseVersion.Name,
			latestReleaseVersion.HtmlUrl,
			latestReleaseVersion.ZipballUrl,
			latestReleaseVersion.TarballUrl,
		)
	}
	return nil
}

// getLatestReleaseVersion - gets the list of releases from "api.github.com/repos/Fantom-foundation/go-lachesis/releases"
// and finds among them the latest release by creation date
func getLatestReleaseVersion(client http.Client, url url.URL) (*ReleaseVersion, error) {
	var releases = []*ReleaseVersion{}
	resp, err := client.Get(url.String())
	if err != nil {
		log.Error("http client Get releases failed", "ur+l", url.String())
		return nil, errors.New(FailedGetNodeVersionMsg)
	}
	defer resp.Body.Close()
	respB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("read response body failed", "url", url.String())
		return nil, errors.New(FailedGetNodeVersionMsg)
	}
	json.Unmarshal(respB, &releases)
	if len(releases) == 0 {
		return nil, errors.New(FailedGetNodeVersionMsg)
	}
	latestRelease := findLatestRelease(releases)
	if latestRelease == nil {
		return nil, errors.New(FailedGetNodeVersionMsg)
	}

	return latestRelease, nil
}

// findLatestRelease - finds the latest release by field - "created_at"
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
