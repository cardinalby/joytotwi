package pageparser

import (
	"errors"
	"joytotwi/app/joy"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func readPost(post *goquery.Selection) (*joy.Post, error) {
	img, imgExists := post.Find(".post_content .image img").First().Attr("src")
	link, linkExists := post.Find(".ufoot .link_wr a").First().Attr("href")
	if imgExists && linkExists {
		joyPost := joy.Post{
			Link:   createJoyPostLink(link),
			ImgURL: createJoyImgLink(img),
		}
		return &joyPost, nil
	}
	var messages []string
	if !linkExists {
		messages = append(messages, "link not found")
	}
	if !imgExists {
		messages = append(messages, "image not found")
	}

	return nil, errors.New(strings.Join(messages, ","))
}

func createJoyPostLink(relative string) string {
	if strings.HasPrefix(relative, "//") {
		return joyBaseURL[0:strings.Index(joyBaseURL, "//")] + relative
	}
	if strings.HasPrefix(relative, joyBaseURL) {
		return relative
	}
	return joyBaseURL + relative
}

func createJoyImgLink(relative string) string {
	if strings.HasPrefix(relative, "//") {
		return "http:" + relative
	}
	return relative
}
