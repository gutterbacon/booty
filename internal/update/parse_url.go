package update

import (
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

type repoInfo struct {
	repoUrl  RepositoryURL
	repoUser string
	repoName string
	curVer   Version
}

func parseGithubUrl(repoUrl RepositoryURL, v Version) (*repoInfo, error) {
	u, err := url.Parse(string(repoUrl))
	if err != nil {
		return nil, err
	}
	path := u.Path
	return &repoInfo{
		repoUser: strings.Trim(filepath.Dir(path), "/"),
		repoName: filepath.Base(path),
		curVer:   v,
		repoUrl:  repoUrl,
	}, nil
}

type semver []int64

func parseReleaseTag(ver string) semver {
	ver = strings.TrimLeft(ver, "v")
	ver = strings.ReplaceAll(ver, "-", ".")
	parts := strings.Split(ver, ".")
	ret := make(semver, len(parts))
	for n, p := range parts {
		i, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			i = 0
		}
		ret[n] = i
	}
	return ret
}

func (v Version) toSemver() []int64 {
	return parseReleaseTag(string(v))
}

func isTagNewer(current semver, latest semver) bool {
	maximum := len(current)
	if len(latest) != maximum {
		return false
	}
	for i := 0; i < maximum; i++ {
		if current[i] == latest[i] {
			continue
		}
		return latest[i] > current[i]
	}
	return false
}
