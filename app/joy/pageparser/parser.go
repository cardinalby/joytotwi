package pageparser

import (
	"joytotwi/app/joy"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

// GetPosts iterates over posts
func GetPosts(
	userName string,
	reverse bool,
	offset int,
	limit int,
	done <-chan struct{},
) (<-chan *joy.Post, <-chan error) {
	outPosts := make(chan *joy.Post)
	outErrors := make(chan error)

	pageOffset, postsOffset := getOffsets(userName, offset)
	pagesIt := createPagesIterator(userName, reverse, pageOffset)

	go func() {
		defer close(outPosts)
		defer close(outErrors)

		processAllPagesPosts(pagesIt, outPosts, outErrors, done)
	}()

	uniquePosts := getUniquePostsChan(outPosts, done)
	limitedPosts, limitedErrors := limitPostsRange(postsOffset, limit, uniquePosts, outErrors, done)
	return limitedPosts, limitedErrors
}

func getOffsets(userName string, postOffset int) (pageOffset, newPostsOffset int) {
	pageOffset = 0
	newPostsOffset = postOffset
	if postOffset > 0 {
		postsPerPage, err := getPostsPerPage(userName)
		if err == nil {
			pageOffset = postOffset / postsPerPage
			newPostsOffset = postOffset % postsPerPage
		} else {
			log.Warnf("Can't determine needed page offset: %s, will iterate from the beginning", err.Error())
		}
	}
	return
}

func processAllPagesPosts(
	pagesIt *joyPagesIterator,
	outPosts chan<- *joy.Post,
	outErrors chan<- error,
	done <-chan struct{},
) {
	for pagesIt.Next() {
		if pagesIt.Err != nil {
			outErrors <- pagesIt.Err
			return
		}
		select {
		case <-done:
			return
		default:
			log.Infof("Page %d loaded", pagesIt.GetCurrentPageNumber())
			processPagePosts(pagesIt.Value, pagesIt.reverse, outPosts, outErrors, done)
		}
	}
}

func processPagePosts(
	page *goquery.Document,
	reverse bool,
	outPosts chan<- *joy.Post,
	outErrors chan<- error,
	done <-chan struct{},
) {
	pagePosts, pagePostErrors := readPagePosts(page, reverse, done)

	for {
		select {
		case err := <-pagePostErrors:
			outErrors <- err
		case post := <-pagePosts:
			if post == nil {
				return
			}
			outPosts <- post
		case <-done:
			return
		}
	}
}
