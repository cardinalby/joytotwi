package watch

import (
	"joytotwi/app/cmd/common"
	"joytotwi/app/joy"
	"joytotwi/app/twisender"
	"time"

	log "github.com/sirupsen/logrus"
)

func watchForPosts(
	client *twisender.Client,
	postReader joy.PostsReader,
	userName string,
	period time.Duration,
	consumePosts func([]*joy.Post),
	done <-chan struct{},
) {
	// start first attempt immediately
	timer := time.NewTimer(time.Millisecond * 0)

	go func() {
		for {
			select {
			case <-timer.C:
				posts, err := common.GetNewPosts(client, postReader, userName, done)
				if err == nil {
					if posts != nil {
						consumePosts(posts)
					} else {
						log.Info("No new posts found")
					}
				}
				timer.Reset(period)
			case <-done:
				return
			}
		}
	}()

	return
}
