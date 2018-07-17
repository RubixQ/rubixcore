package api

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

func sendSMS(msg, to string) error {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	form := url.Values{}
	form.Add("username", "AWA")
	form.Add("password", "ttvpass101")
	form.Add("numbers", to)
	form.Add("message", msg)
	form.Add("from", "RUBIXCORE")

	url := "https://infoline.nandiclient.com/AWA/campaigns/sendmsg"
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
