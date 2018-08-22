package Vigor

import (
	"fmt"
	"net/url"
	"encoding/base64"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

var E_LOGIN_FAILED = errors.New("login failed")

func (this *Vigor) Login(username string, password string) (error) {
	v := url.Values{}

	token := make([]byte, 8)
	rand.Read(token)

	this.csrf = string([]byte(hex.EncodeToString(token))[:15])

	v.Set("aa", base64.StdEncoding.EncodeToString([]byte(username)))
	v.Add("ab", base64.StdEncoding.EncodeToString([]byte(password)))
	v.Add("sFormAuthStr", this.csrf)
	resp, err := this.client.PostForm(fmt.Sprintf("http://%s/cgi-bin/wlogin.cgi", this.ip), v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	cookies := resp.Header.Get("Set-Cookie")
	if cookies == "" {
		return E_LOGIN_FAILED
	}

	return nil
}
