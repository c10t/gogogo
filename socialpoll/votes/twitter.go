package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/joeshaw/envdecode"
)

var conn net.Conn

func dial(netw, addr string) (net.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}

	netc, err := net.DialTimeout(netw, addr, 5*time.Second)
	if err != nil {
		return nil, err
	}

	conn = netc
	return netc, nil
}

var reader io.ReadCloser

func closeConn() {
	if conn != nil {
		conn.Close()
	}

	if reader != nil {
		reader.Close()
	}
}

var (
	authClient *oauth.Client
	creds      *oauth.Credentials
)

func setupTwitterAuth() {
	var ts struct {
		ConsumerKey    string `env:"SP_TWITTER_KEY,required"`
		ConsumerSecret string `env:"SP_TWITTER_SECRET,required"`
		AccessToken    string `env:"SP_TWITTER_KEY,required"`
		AccessSecret   string `env:"SP_TWITTER_ACCESSSECRET,required"`
	}

	if err := envdecode.Decode(&ts); err != nil {
		log.Fatalln(err)
	}

	creds = &oauth.Credentials{
		Token:  ts.AccessToken,
		Secret: ts.AccessSecret,
	}

	authClient = &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  ts.ConsumerKey,
			Secret: ts.ConsumerSecret,
		},
	}
}

var (
	authSetupOnce sync.Once
	httpClient    *http.Client
)

func makeRequest(query url.Values) (*http.Request, error) {
	authSetupOnce.Do(func() {
		setupTwitterAuth()
	})
	const endpoint = "https://stream.twitter.com/1.1/statuses/filter.json"

	req, err := http.NewRequest("POST", endpoint,
		strings.NewReader(query.Encode()))
	if err != nil {
		return nil, err
	}

	formEnc := query.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))

	ah := authClient.AuthorizationHeader(creds, "POST", req.URL, query)
	req.Header.Set("Authorization", ah)

	return req, nil
}

type tweet struct {
	Text string
}

func twitterStream(ctx context.Context, votes chan<- string) {
	defer close(votes)
	for {
		log.Println("start to read twitter stream...")
		readFromTwitterWithTimeout(ctx, 1*time.Minute, votes)
		log.Println("--- (waiting) ---")
		select {
		case <-ctx.Done():
			log.Println("stop to read twitter stream...")
			return
		case <-time.After(10 * time.Second):
		}
	}
}

func readFromTwitterWithTimeout(ctx context.Context,
	timeout time.Duration, votes chan<- string) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	readFromTwitter(ctx, votes)
}

func readFromTwitter(ctx context.Context, votes chan<- string) {
	options, err := loadOptions()
	if err != nil {
		log.Println("Failed to load options:", err)
		return
	}

	query := make(url.Values)
	query.Set("track", strings.Join(options, ","))
	req, err := makeRequest(query)
	if err != nil {
		log.Println("Failed to create filter request:", err)
		return
	}

	client := &http.Client{}
	if deadline, ok := ctx.Deadline(); ok {
		client.Timeout = deadline.Sub(time.Now())
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to request:", err)
		return
	}
	done := make(chan struct{})
	defer func() { <-done }()
	defer resp.Body.Close() // defer executes reversely

	go func() {
		defer close(done)
		log.Println("response:", resp.StatusCode)
		if resp.StatusCode != 200 {
			var buf bytes.Buffer
			io.Copy(&buf, resp.Body)
			log.Printf("response body: %s\n", buf.String())
			return
		}

		decoder := json.NewDecoder(resp.Body)
		for {
			var tweet tweet
			if err := decoder.Decode(&tweet); err != nil {
				break
			}

			log.Println("tweet:", tweet)

			for _, option := range options {
				twt := strings.ToLower(tweet.Text)
				opt := strings.ToLower(option)
				if strings.Contains(twt, opt) {
					log.Println("poll:", option)
					votes <- option
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}
}
