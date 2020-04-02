package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

func main() {
	// Load config
	c := ReadConfig()

	repoList := ListImages(c)
	var wg sync.WaitGroup

	// Loop over all image repositories
	for _, registry := range repoList.Repositories {
		ichan := make(chan ImageDigest)
		t := ListTags(c, string(registry))
		var img Images
		img.Name = string(registry)
		imageMap := make(map[string]string)
		var mySlice []int
		fmt.Printf("%s tags list:\n", registry)

		// Loop over all tags
		for _, tg := range t.Tags {
			tag := string(tg)
			wg.Add(1)
			go manifestWorker(ichan, &wg, c, string(registry), tag)
			im := <-ichan

			for date := range im.Digest {
				intUnixTime, _ := strconv.Atoi(date)
				mySlice = append(mySlice, intUnixTime)
			}

			for key, val := range im.Digest {
				imageMap[key] = val
			}
		}

		img.Digests = ImageDigest{imageMap}
		sort.Ints(mySlice)
		var x int

		if (len(mySlice) - c.DrImgCount) > 0 {
			x = len(mySlice) - c.DrImgCount
		} else {
			x = len(mySlice)
		}

		for _, date := range mySlice[:x] {
			dateString := strconv.Itoa(date)
			if len(mySlice) > c.DrImgCount {
				DelManifest(c, string(registry), img.Digests.Digest[dateString])
				fmt.Printf("    image was deleted %s \n", unixToTime(dateString))
			}
		}
	}

	wg.Wait()
}

// Get digests and creation dates for all image tags
func manifestWorker(ichan chan ImageDigest, wg *sync.WaitGroup, c *Config, registry, tag string) {
	var dateCreated string
	digest, resp := GetManifest(c, registry, tag)
	if resp == "200 OK" {
		dateCreated = GetCreationDate(c, registry, tag)
	} else {
		dateCreated = "empty"
	}

	m := make(map[string]string)
	m[dateCreated] = digest
	ichan <- ImageDigest{m}
	fmt.Printf("    %s -> %v\n", tag, unixToTime(dateCreated))
	wg.Done()
}
