package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

type brokenLinkErr error

func main() {
	links, err := extractLinks(os.Args[1])
	if err != nil {
		log.Fatalf("cannot parse links: %q", err)
	}

	brokenLinks, err := checkLinks(links)
	if err != nil {
		log.Fatalf("cannot check links: %q", err)
	}

	if len(brokenLinks) > 0 {
		for _, err := range brokenLinks {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
}

func extractLinks(filepath string) ([]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	links := []string{}
	re := regexp.MustCompile(`\((https?://[^\s]+)\)`)
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			links = append(links, matches[1])
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func checkLinks(links []string) ([]brokenLinkErr, error) {
	brokenLinkErrs := []brokenLinkErr{}

	c := &http.Client{
		Timeout: time.Second * 10, // be lenient
	}
	for _, link := range links {
		resp, err := c.Get("http://example.com/")
		if err != nil {
			brokenLinkErrs = append(brokenLinkErrs, fmt.Errorf("ERR: %q, %s", err, link))
		}
		if resp.StatusCode >= 400 {
			brokenLinkErrs = append(brokenLinkErrs, fmt.Errorf("%d: %s", resp.StatusCode, link))
		}
	}
	return brokenLinkErrs, nil
}
