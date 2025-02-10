package normalizer

import (
	"fmt"
	"github.com/PuerkitoBio/purell"
	"net/url"
)

func NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	if s := u.Scheme; s != "http" && s != "https" {
		return "", fmt.Errorf("unsupported scheme: %s", s)
	}

	flags := purell.FlagLowercaseScheme |
		purell.FlagLowercaseHost |
		purell.FlagRemoveDefaultPort |
		purell.FlagRemoveDotSegments |
		purell.FlagRemoveDuplicateSlashes |
		purell.FlagRemoveTrailingSlash |
		purell.FlagSortQuery

	normalizedURL, err := purell.NormalizeURLString(u.String(), flags)
	if err != nil {
		return "", err
	}

	return normalizedURL, nil
}
