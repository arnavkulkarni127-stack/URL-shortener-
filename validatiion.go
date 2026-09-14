package main

import "net/url"

func isValidURL(URL string) bool {

	u, err := url.Parse(URL)
	return err == nil && u.Scheme != "" && u.Host != ""
}
