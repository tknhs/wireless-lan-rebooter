package main

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mitchellh/osext"
)

const (
	BaseURL    string = "http://YOUR_LOCAL_IP_ADDRESS/index.cgi/reboot_main_set"
	SessionID  string = "YOUR_SESSION_ID"
	RouterID   string = "YOUR_ROUTER_ID"
	RouterPass string = "YOUR_ROUTER_PASSWORD"
)

func LogInit() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)

	binPath, _ := osext.ExecutableFolder()
	logPath := binPath + "wireless-lan-rebooter.log"

	f, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logger.Out = f

	return logger
}

func HTTPPost() (string, int, error) {
	values := url.Values{}
	values.Add("DUMMY", "")
	values.Add("DISABLED_CHECKBOX", "")
	values.Add("CHECK_ACTION_MODE", "0")
	values.Add("SESSION_ID", SessionID)

	req, err := http.NewRequest(
		"POST",
		BaseURL,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return "", -1, err
	}

	requestURL := BaseURL + "?" + values.Encode()

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(RouterID, RouterPass)

	client := &http.Client{Timeout: time.Duration(15 * time.Second)}
	resp, err := client.Do(req)
	if err != nil {
		return requestURL, -1, err
	}
	defer resp.Body.Close()

	return requestURL, resp.StatusCode, nil
}

func main() {
	logger := LogInit()
	if url, statusCode, err := HTTPPost(); err != nil {
		logger.WithFields(logrus.Fields{
			"url":         url,
			"status_code": statusCode,
		}).Error(err.Error())
	} else {
		logger.WithFields(logrus.Fields{
			"url":         url,
			"status_code": statusCode,
		}).Info("success")
	}
}
