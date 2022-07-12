package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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
