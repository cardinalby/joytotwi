package rssparser

import (
	"errors"
	"fmt"
	"joytotwi/app/joy"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

const feedURL = "http://joyreactor.cc/rss/user/username/%s"

// GetPosts from RSS (has only last 10 posts)
func GetPosts(
	userName string,
	reverse bool,
	offset int,
	limit int,
	done chan struct{},
) (chan *joy.Post, chan error) {
	if reverse || offset < 0 || limit < 0 {
		panic("Invalid parser args")
	}
	posts := make(chan *joy.Post)
	errs := make(chan error)

	go func() {
		defer close(posts)
		defer close(errs)

		parser := gofeed.NewParser()
		feed, err := parser.ParseURL(fmt.Sprintf(feedURL, url.QueryEscape(userName)))
		if err != nil {
			errs <- err
			return
		}

		rangeEnd := len(feed.Items) - 1
		if limit != 0 {
			limitedRangeEnd := offset + limit - 1
			if limitedRangeEnd < rangeEnd {
				rangeEnd = limitedRangeEnd
			}
		}

		for i := offset; i <= rangeEnd; i++ {
			post, err := feedItemToPost(feed.Items[i])
			if err != nil {
				select {
				case errs <- err:
					continue
				case <-done:
					return
				}
			}

			select {
			case posts <- post:
				continue
			case <-done:
				return
			}
		}
	}()

	return posts, errs
}

func feedItemToPost(item *gofeed.Item) (*joy.Post, error) {
	post := &joy.Post{
		Link: item.Link,
	}
	if item.Image != nil {
		post.ImgURL = item.Image.URL
	} else {
		imgURL, err := getImageURLFromItemDescr(item.Description)
		if err != nil {
			return nil, err
		}
		post.ImgURL = imgURL
	}

	return post, nil
}

func getImageURLFromItemDescr(descr string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(descr))
	if err != nil {
		return "", fmt.Errorf("can't parse description as HTML: %s", err.Error())
	}
	src, exists := doc.Find("img").First().Attr("src")
	if !exists {
		return "", errors.New("can't find <img> tag in description")
	}
	return src, nil
}
