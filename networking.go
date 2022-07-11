package main

import (
	"io/ioutil"
	"fmt"
	"net/http"
    "crypto/tls"
)

func get_json(subreddit string) []byte {
	client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{},
        },
    }

    url := fmt.Sprintf("https://www.reddit.com/r/%s/top.json", subreddit)
    
    req, err := http.NewRequest("GET", url, nil)
    checkErr(err)
    
    req.Header.Set("User-Agent", "linux:go-postgrabber:v0.1")
    
    resp, err := client.Do(req)
    checkErr(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
    checkErr(err)

    return body
}
