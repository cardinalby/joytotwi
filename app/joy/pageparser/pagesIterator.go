package pageparser

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type joyPagesIterator struct {
	reverse    bool
	userName   string
	pagination *joyPagination
	offset     int
	Value      *goquery.Document
	Err        error
}

const pageNumberLast = -1

func createPagesIterator(userName string, reverse bool, offset int) *joyPagesIterator {
	return &joyPagesIterator{reverse: reverse, userName: userName, offset: offset}
}

func (it *joyPagesIterator) Next() bool {
	pageURL, exists, err := it.getNextPageURL()
	if err != nil {
		it.Err = err
		return false
	}
	if !exists {
		it.Err = nil
		return false
	}

	it.downloadPage(pageURL)
	return it.Err == nil
}

func (it *joyPagesIterator) getNextPageURL() (url string, exists bool, err error) {
	if it.pagination != nil {
		if it.reverse {
			return it.pagination.prevLink, it.pagination.prevLink != "", nil
		}
		return it.pagination.nextLink, it.pagination.nextLink != "", nil
	}

	if it.reverse {
		return getPageURL(it.userName, 1+it.offset), true, nil
	}
	if it.offset == 0 {
		return getPageURL(it.userName, pageNumberLast), true, nil
	}
	lastPageNumber, err := it.getLastPageNumber()
	if err != nil {
		return "", false, fmt.Errorf("can't get last page number: %s", err.Error())
	}
	if it.offset > lastPageNumber {
		return "", false, fmt.Errorf(
			"page offset %d is too big, last page number is: %d", it.offset, lastPageNumber,
		)
	}
	return getPageURL(it.userName, lastPageNumber-it.offset), true, nil
}

func (it *joyPagesIterator) GetCurrentPageNumber() int {
	if it.pagination == nil {
		return 0
	}
	return it.pagination.current
}

func (it *joyPagesIterator) downloadPage(url string) {
	html, err := goquery.NewDocument(url)
	if err != nil {
		it.Err = fmt.Errorf("page '%s' load error: %s", url, err.Error())
	} else {
		it.Err = nil
		it.Value = html
		it.pagination, it.Err = readPagination(html)
	}
}

func getPageURL(userName string, pageNumber int) string {
	userNameEscaped := url.QueryEscape(userName)
	if pageNumber == pageNumberLast {
		return fmt.Sprintf(joyBaseURL+joyUserPageURL, userNameEscaped)
	}
	return fmt.Sprintf(joyBaseURL+joyUserPageURL+joyUserPageNumber, userNameEscaped, pageNumber)
}

func (it *joyPagesIterator) getLastPageNumber() (int, error) {
	pageUrl := getPageURL(it.userName, pageNumberLast)
	html, err := goquery.NewDocument(pageUrl)
	if err != nil {
		return 0, fmt.Errorf("page '%s' load error: %s", pageUrl, err.Error())
	}
	pagination, err := readPagination(html)
	if err != nil {
		return 0, err
	}
	return pagination.current, nil
}

func getPostsPerPage(userName string) (int, error) {
	pageUrl := getPageURL(userName, 1)
	html, err := goquery.NewDocument(pageUrl)
	if err != nil {
		return 0, fmt.Errorf("page '%s' load error: %s", pageUrl, err.Error())
	}

	pagination, err := readPagination(html)
	if err != nil {
		return 0, err
	}
	if pagination.prevLink == "" {
		return 0, errors.New("only one page present, can't determine posts count per page")
	}
	return countPagePosts(html), nil
}
