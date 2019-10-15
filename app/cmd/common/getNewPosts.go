package common

import (
	"github.com/sirupsen/logrus"
	"joytotwi/app/joy"
	"joytotwi/app/twisender"
)

// GetNewPosts return array of new posts (checked by client) sorted from old to new
func GetNewPosts(
	client *twisender.Client,
	postReader joy.PostsReader,
	userName string,
	done chan struct{},
) (posts []*joy.Post, err error) {
	doneReading := make(chan struct{})
	defer close(doneReading)

	logrus.Info("Checking for new posts..")
	readerPosts, readerErrors := postReader(userName, doneReading)
	var newPosts []*joy.Post

	offset := 0

	for {
		select {
		case post := <-readerPosts:
			if post == nil {
				return nil, nil
			}
			offset++
			if client.Exists(post.Link) {
				if len(newPosts) == 0 {
					return nil, nil
				}
				reversePosts(newPosts)
				return newPosts, nil
			}
			newPosts = append(newPosts, post)
		case err := <-readerErrors:
			if err != nil {
				offset++
				logrus.Warnf("Post (offset: %d) parse error: %s", offset, err.Error())
				return nil, err
			}
		case <-done:
			return nil, nil
		}
	}
}

func reversePosts(posts []*joy.Post) {
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}
}
