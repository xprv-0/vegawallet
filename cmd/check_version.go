package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/blang/semver/v4"
)

const (
	vegaWalletReleasesPage = "https://github.com/vegaprotocol/go-wallet/releases"
	vegaWalletReleaseAPI   = "https://api.github.com/repos/vegaprotocol/go-wallet/releases"
)

// returns a newer version, or an error or nil for both
// if no error happened, an no updates are needed
func checkVersion(currentVersion string) (*semver.Version, error) {
	req, err := http.NewRequest("GET", vegaWalletReleaseAPI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	releases := []struct {
		Name string `json:"name"`
	}{}

	err = json.Unmarshal(body, &releases)
	if err != nil {
		return nil, err
	}

	last, _ := semver.Make(strings.TrimPrefix(currentVersion, "v"))
	cur := last

	for _, v := range releases {
		newV, err := semver.Make(strings.TrimPrefix(v.Name, "v"))
		if err != nil {
			// unsupported version
			continue
		}
		if newV.GT(last) {
			last = newV
		}
	}

	if cur.EQ(last) {
		// no updates
		return nil, nil
	}

	return &last, nil
}
