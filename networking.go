package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

func get_json(subreddit string) []byte {
    // Setup http client with empty tls config for reddit url to work
	client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{},
        },
    }

    url := fmt.Sprintf("https://www.reddit.com/r/%s/top.json?limit=100", subreddit)
    
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

func get_images(body []byte) []string {
    var dat map[string]interface{}
    var images_list []string
    match, _ := regexp.Compile("^https?://.*/.*.(png|gif|webp|jpeg|jpg)$")

    err := json.Unmarshal(body, &dat)
    checkErr(err)

    json_data := dat["data"].(map[string]interface{})
    post_list := json_data["children"].([]interface{})
    for i := 0; i < len(post_list); i++ {
        post := post_list[i].(map[string]interface{})
        post_data := post["data"].(map[string]interface{})
        checkErr(err)
        url := post_data["url_overridden_by_dest"]
        if url != nil {
             if match.MatchString(url.(string)) {
                 images_list = append(images_list, url.(string))
             }
        }
    }

    return images_list
}

func downloadFile(URL string) ([]byte, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
        defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
        return nil, errors.New(response.Status)
	}
	var data bytes.Buffer
	_, err = io.Copy(&data, response.Body)
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}

func downloadMultipleFiles(urls []string) {
    var waiter sync.WaitGroup
    waiter.Add(len(urls))
	for _, URL := range urls {
        go func(URL string, waiter *sync.WaitGroup) {
            b, err := downloadFile(URL)
            url_comp := strings.Split(URL, "/")
            file_name := url_comp[len(url_comp) - 1]
            fmt.Println(file_name)
            if err != nil {
                return
            }
            ioutil.WriteFile(file_name, b, 0644)
            waiter.Done()
        } (URL, &waiter)
	}
    waiter.Wait()
}
