package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	drCatalogUrl = "/v2/_catalog?n=1000"
	drTagsList = "/tags/list?n=1000"
)

type Tag string
type Repo string
type Tags struct {
	Name string
	Tags []Tag `json:"tags"`
}
type Repos struct {
	Repositories []Repo `json:"repositories"`
}
type ImageDigest struct {
	Digest map[string]string
}
type Images struct {
	Name string
	Digests ImageDigest
}
type V1Manifest struct {
	Name string `json:"name"`
	Tag string `json:"tag"`
	History []History `json:"history"`
}
type History struct {
	V1Compatibility string `json:"v1Compatibility"`
}


func makeV2Request(url, m string, c *Config) *http.Response {
	req, err := http.NewRequest(m, url, nil)
	if err != nil {
		log.Println(err)
	}
	req.SetBasicAuth(c.DrUser, c.DrPass)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	return resp
}

func makeV1Request(url, m string, c *Config) *http.Response {
	req, err := http.NewRequest(m, url, nil)
	if err != nil {
		log.Println(err)
	}
	req.SetBasicAuth(c.DrUser, c.DrPass)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v1+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	return resp
}

func DelManifest(c *Config, registry, manifest string) string {
	resp := makeV2Request(c.DrHost+":"+c.DrPort+"/v2/"+registry+"/manifests/"+manifest, "DELETE", c)
	_, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return resp.Status
}

func timeToUnix(ts string) string  {
	t, _ := time.Parse(time.RFC3339, ts)
	return strconv.FormatInt(t.Unix(),10)
}

func GetCreationDate(c *Config, registry, tagName string) string {

	resp := makeV1Request(c.DrHost+":"+c.DrPort+"/v2/"+registry+"/manifests/"+tagName, "GET", c)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var v1manifest V1Manifest
	parseErr := json.Unmarshal(body, &v1manifest)
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	in := []byte(v1manifest.History[0].V1Compatibility)
	var raw map[string]interface{}
	json.Unmarshal(in, &raw)
	ct := timeToUnix(raw["created"].(string))
	return ct
}

func GetManifest(c *Config, registry, tagName string) (string, string) {
	resp := makeV2Request(c.DrHost+":"+c.DrPort+"/v2/"+registry+"/manifests/"+tagName, "GET", c)
	_, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return resp.Header.Get("Docker-Content-Digest"), resp.Status
}

func ListTags(c *Config, registry string) *Tags {
	resp := makeV2Request(c.DrHost+":"+c.DrPort+"/v2/"+registry+drTagsList, "GET", c)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var tags Tags
	parseErr := json.Unmarshal(body, &tags)
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	return &tags
}

func ListImages(c *Config) *Repos {
	resp := makeV2Request(c.DrHost+":"+c.DrPort+drCatalogUrl, "GET", c)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var repos Repos
	parseErr := json.Unmarshal(body, &repos)
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	return &repos
}
