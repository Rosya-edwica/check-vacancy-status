package utils

import (
	"io"
	"net/http"
	"time"
)



func GetJson(url string, headers map[string]string) (json string, err error) {
	client := http.Client{
		Timeout: 120 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	CheckErr(err)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	response, err := client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	
	return string(data), nil
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}