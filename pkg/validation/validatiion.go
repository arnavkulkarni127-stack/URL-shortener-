package validation

import "net/url"

func IsValidURL(URL string) bool {

	u, err := url.Parse(URL)
	return err == nil && u.Scheme != "" && u.Host != ""
}
