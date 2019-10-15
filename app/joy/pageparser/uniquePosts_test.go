package pageparser

import (
	"github.com/stretchr/testify/assert"
	"joytotwi/app/joy"
	"testing"
)

func TestGetUniquePostsChan(t *testing.T) {

	done := make(chan struct{})
	allPostsArr := []*joy.Post{
		{Link: "link0", ImgURL: "url0"},
		{Link: "link1", ImgURL: "url1"},
		{Link: "link2", ImgURL: "url2"},
		{Link: "link0", ImgURL: "url0"},
		{Link: "link1", ImgURL: "url1"},
		{Link: "link3", ImgURL: "url3"},
	}
	defer close(done)

	uPosts := getUniquePostsChan(produceTestPosts(allPostsArr, done), done)

	assert.Equal(t, allPostsArr[0], <-uPosts)
	assert.Equal(t, allPostsArr[1], <-uPosts)
	assert.Equal(t, allPostsArr[2], <-uPosts)
	assert.Equal(t, allPostsArr[5], <-uPosts)
	assert.Nil(t, <-uPosts)
}

func produceTestPosts(postsArr []*joy.Post, done <-chan struct{}) <-chan *joy.Post {
	result := make(chan *joy.Post)
	go func() {
		defer close(result)
		for _, p := range postsArr {
			select {
			case <-done:
				return
			default:
				result <- p
			}
		}
	}()
	return result
}
