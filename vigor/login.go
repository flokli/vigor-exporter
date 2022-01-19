package Vigor

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
)

var ErrLoginFailed = errors.New("login failed")

func (v *Vigor) Login(username string, password string) error {
	urlValues := url.Values{}

	token := make([]byte, 8)
	rand.Read(token)

	v.csrf = string([]byte(hex.EncodeToString(token))[:15])

	urlValues.Set("aa", base64.StdEncoding.EncodeToString([]byte(username)))
	urlValues.Add("ab", base64.StdEncoding.EncodeToString([]byte(password)))
	urlValues.Add("sFormAuthStr", v.csrf)
	resp, err := v.client.PostForm(fmt.Sprintf("http://%s/cgi-bin/wlogin.cgi", v.host), urlValues)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	cookies := resp.Header.Get("Set-Cookie")
	if cookies == "" {
		return ErrLoginFailed
	}

	return nil
}
