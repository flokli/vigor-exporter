package Vigor

import (
	"net/http/cookiejar"
	"net/http"
)

type Vigor struct {
	jar    *cookiejar.Jar
	client *http.Client
	ip     string
	csrf   string
}

func New(ip string) (*Vigor, error) {
	var err error

	this := Vigor{ip: ip}
	this.jar, err = cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	this.client = &http.Client{
		Jar: this.jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &this, nil
}
