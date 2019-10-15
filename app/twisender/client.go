package twisender

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
)

const TimelineMaxCount = 200

// Client to send Posts
type Client struct {
	client       *anaconda.TwitterApi
	existedLinks map[string]bool
}

// ClientCreds is creds needed for client
type ClientCreds struct {
	AccessToken       string
	AccessTokenSecret string
	ConsumerKey       string
	ConsumerSecret    string
}

// CreateNewClient creates new Client
func CreateNewClient(creds ClientCreds) *Client {
	return &Client{
		client: anaconda.NewTwitterApiWithCredentials(
			creds.AccessToken,
			creds.AccessTokenSecret,
			creds.ConsumerKey,
			creds.ConsumerSecret,
		),
		existedLinks: make(map[string]bool),
	}
}

// Init loads user tweets to track duplicates
func (client *Client) Init(loadCount int) error {
	timeline, err := client.client.GetHomeTimeline(url.Values{"count": []string{strconv.Itoa(loadCount)}})
	if err != nil {
		return err
	}
	for _, tweet := range timeline {
		lastTweetLink, err := getOriginPostLinkFromTweet(&tweet)
		if err == nil {
			client.existedLinks[lastTweetLink] = true
		}
	}
	return nil
}

// PostNew post new tweet if it didn't exist at the moment of Init() call
func (client *Client) PostNew(text string, imageURL string) (tweetID int64, exists bool, err error) {
	if _, exists := client.existedLinks[text]; exists {
		return 0, true, nil
	}
	media, err := client.uploadExternalImg(imageURL)
	if err != nil {
		return 0, false, err
	}
	mediaIds := []string{strconv.FormatInt(media.MediaID, 10)}
	tweet, err := client.client.PostTweet(text, url.Values{"media_ids": mediaIds})
	if err != nil {
		return 0, false, err
	}
	client.existedLinks[text] = true
	return tweet.Id, false, nil
}

func getOriginPostLinkFromTweet(tweet *anaconda.Tweet) (string, error) {
	urls := tweet.Entities.Urls
	if len(urls) != 1 {
		return "", fmt.Errorf("tweet has %d links", len(urls))
	}
	return urls[0].Expanded_url, nil
}

func (client *Client) uploadExternalImg(imgURL string) (*anaconda.Media, error) {
	base64img, err := downloadImageBase64(imgURL)
	if err != nil {
		return nil, err
	}

	media, err := client.client.UploadMedia(base64img)
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func downloadImageBase64(imgURL string) (string, error) {
	resp, err := http.Get(imgURL)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("wrong response code: %d", resp.StatusCode)
	}

	imageData, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return "", readErr
	}
	return base64.StdEncoding.EncodeToString(imageData), nil
}
