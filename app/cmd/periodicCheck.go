package cmd

import (
	"joytotwi/app/joy"
	"time"

	log "github.com/sirupsen/logrus"
)

// returns chan with posts. After publishing each posts waits for acknowledge in postAck chan.
// If ack is false stops emitting posts and waits for next timer
func watchForPosts(
	postReader joy.PostsReader,
	userName string,
	period time.Duration,
	done chan struct{},
) (posts chan *joy.Post, postAck chan bool) {
	// start first attempt immediately
	timer := time.NewTimer(time.Millisecond * 0)
	posts = make(chan *joy.Post)
	postAck = make(chan bool)

	go func() {
		for {
			select {
			case <-timer.C:
				sendPost := func(post *joy.Post) bool {
					posts <- post
					return <-postAck
				}
				checkForPosts(postReader, userName, done, sendPost)
				timer.Reset(period)
			case <-done:
				return
			}
		}
	}()

	return
}

func checkForPosts(
	postReader joy.PostsReader,
	userName string,
	done chan struct{},
	// returns true: expects next, false: stop producing
	consumePost func(*joy.Post) bool,
) {
	doneReading := make(chan struct{})
	defer close(doneReading)

	log.Info("Checking for new posts..")
	readerPosts, readerErrors := postReader(userName, doneReading)

	offset := 0
	for {
		select {
		case post := <-readerPosts:
			if post == nil {
				return
			}
			offset++
			if !consumePost(post) {
				return
			}
		case err := <-readerErrors:
			if err != nil {
				offset++
				log.Warnf("Post (offset: %d) parse error: %s", offset, err.Error())
			}
		case <-done:
			return
		}
	}
}
