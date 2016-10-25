package utils

import "strings"

func JoinURL(url1 string, url2 string, addTrailingSlash bool) (url string) {
	if strings.HasSuffix(url1, "/") {
		if strings.HasPrefix(url2, "/") {
			url = url1 + url2[1:]
		} else {
			url = url1 + url2
		}
	} else {
		if strings.HasPrefix(url2, "/") {
			url = url1 + url2
		} else {
			url = url1 + "/" + url2
		}
	}

	if addTrailingSlash && !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return
}
