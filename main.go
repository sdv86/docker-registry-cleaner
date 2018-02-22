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
	var wg sync.WaitGroup
	ichan := make(chan ImageDigest)

	// Loop over all image repositories
	for _, registry := range c.DrRegs{

		t := ListTags(c, registry)
		var img Images
		img.Name = registry
		imageMap := make(map[string]string)
		var mySlice []int
		fmt.Printf("Tags were found for %s:\n", registry)

		// Loop over all tags
		for _, tg := range t.Tags {
			tag := string(tg)
			wg.Add(1)
			go manifestWorker(ichan, &wg, c, registry, tag)
			im := <- ichan

			for date := range im.Digest {
				intUnixTime, _ := strconv.Atoi(date)
				mySlice = append(mySlice, intUnixTime)
			}

			fmt.Printf("  %s -> %v\n", tag, im.Digest)

			for key, val := range im.Digest {
				imageMap[key] = val
			}
		}

		img.Digests = ImageDigest{imageMap}
		sort.Ints(mySlice)
		var x int

		if (len(mySlice)-c.DrImgCount) > 0 {
			x = len(mySlice)-c.DrImgCount
		} else {
			x = len(mySlice)
		}

		for _, date := range mySlice[:x] {
			dateString := strconv.Itoa(date)
			if len(mySlice) > c.DrImgCount {
				DelManifest(c, registry, img.Digests.Digest[dateString])
				fmt.Printf("Tag was deleted %s \n", img.Digests.Digest[dateString])
			}
		}
	}

	wg.Wait()
}

// Get digests and creation dates for all image tags
func manifestWorker(ichan chan ImageDigest, wg *sync.WaitGroup, c *Config, registry, tag string) {
	var dateCreated string
	defer wg.Done()
	digest, resp := GetManifest(c, registry, tag)

	if resp == "200 OK" {
		dateCreated = GetCreationDate(c, registry, tag)
	} else {
		dateCreated = "empty"
	}

	m := make(map[string]string)
	m[dateCreated] = digest
	ichan <- ImageDigest{m}
}
