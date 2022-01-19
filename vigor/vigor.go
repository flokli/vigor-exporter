package Vigor

import (
	"net/http"
	"net/http/cookiejar"
)

type Vigor struct {
	jar    *cookiejar.Jar
	client *http.Client
	host   string
	csrf   string
}

func New(host string) (*Vigor, error) {
	var err error

	v := Vigor{host: host}
	v.jar, err = cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	v.client = &http.Client{
		Jar: v.jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &v, nil
}
