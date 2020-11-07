package checker

type checker interface {
	LatestVersionUrl() (string, error)
}

func NewChecker(os string, includePrerelease bool) checker {
	if includePrerelease {
		return newBintrayChecker(os, getUrl)
	}

	return newGithubChecker(os, getUrl)
}
