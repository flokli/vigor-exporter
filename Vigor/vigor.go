package Vigor

import (
	"net/http"
	"net/http/cookiejar"
)

type Vigor struct {
	jar    *cookiejar.Jar
	client *http.Client
	ip     string
	csrf   string
}

func New(ip string) (*Vigor, error) {
	var err error

	v := Vigor{ip: ip}
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
