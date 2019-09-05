package pageparser

import (
	"joytotwi/app/joy"

	"github.com/PuerkitoBio/goquery"
)

type joyPage struct {
	number int
	posts  []*joy.Post
}

func countPagePosts(doc *goquery.Document) int {
	selection := doc.Find(".postContainer")
	return selection.Length()
}

func readPagePosts(doc *goquery.Document, reverse bool, done chan struct{}) (posts chan *joy.Post, errors chan error) {
	posts = make(chan *joy.Post)
	errors = make(chan error)

	go func() {
		defer close(posts)

		selection := doc.Find(".postContainer")
		iterateSelection(selection, func(i int, selection *goquery.Selection) bool {
			post, err := readPost(selection)
			if err != nil {
				select {
				case errors <- err:
					return true
				case <-done:
					return false
				}
			}

			select {
			case posts <- post:
				return true
			case <-done:
				return false
			}
		}, reverse)
	}()

	return posts, errors
}

func iterateSelection(s *goquery.Selection, f func(int, *goquery.Selection) bool, reverse bool) *goquery.Selection {
	if !reverse {
		return s.EachWithBreak(f)
	}
	selections := make([]*goquery.Selection, len(s.Nodes))
	s.Each(func(i int, s *goquery.Selection) {
		selections[i] = s
	})
	for i := len(selections) - 1; i >= 0; i-- {
		if !f(i, selections[i]) {
			break
		}
	}
	return s
}
