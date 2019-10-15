package pageparser

import (
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type joyPagination struct {
	current  int
	nextLink string
	prevLink string
}

func readPagination(doc *goquery.Document) (*joyPagination, error) {
	paginationDiv := doc.Find("div.pagination")
	if paginationDiv.Length() == 0 {
		return &joyPagination{0, "", ""}, nil
	}

	var err error
	result := joyPagination{}
	result.current, err = getNumberFromElem(paginationDiv, ".pagination_expanded span.current")
	if err != nil {
		return nil, err
	}

	result.nextLink, _ = paginationDiv.Find("a.next").First().Attr("href")
	if result.nextLink != "" {
		result.nextLink = joyBaseURL + result.nextLink
	}
	result.prevLink, _ = paginationDiv.Find("a.prev").First().Attr("href")
	if result.prevLink != "" {
		result.prevLink = joyBaseURL + result.prevLink
	}

	return &result, nil
}

func getNumberFromElem(block *goquery.Selection, selector string) (int, error) {
	elem := block.Find(selector)
	if elem.Length() == 0 {
		return 0, fmt.Errorf("element '%s' not found in pagination block", selector)
	}
	if elem.Length() > 1 {
		return 0, fmt.Errorf("element '%s' found %d times in pagination block", selector, elem.Length())
	}
	res, err := strconv.Atoi(elem.First().Text())
	if err != nil {
		return 0, fmt.Errorf("invalid value in '%s': %s", selector, err.Error())
	}
	return res, nil
}
