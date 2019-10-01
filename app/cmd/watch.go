package cmd

import (
	"joytotwi/app/joy"
	"joytotwi/app/joy/selector"
	"joytotwi/app/twisender"
	"time"

	log "github.com/sirupsen/logrus"
)

const postsReverse = false
const postsOffset = 0
const postsLimit = 0
const initByTweetsCount = 1

// WatchCommand for checking for new posts periodically and post them to twitter
type WatchCommand struct {
	Period int `short:"p" long:"period" default:"43200" description:"Period of checking for new posts in seconds"`
	CommonOptions
}

// SetCommonOptions sets common options in command
func (cmd *WatchCommand) SetCommonOptions(opts *CommonOptions) {
	cmd.CommonOptions = *opts
}

// Execute command method for flags.Commander
func (cmd *WatchCommand) Execute(args []string) error {
	postReader, err := selector.GetPostReader(cmd.SourceType, postsReverse, postsOffset, postsLimit)
	if err != nil {
		return err
	}

	client := twisender.CreateNewClient(getTwiCredsFromOpts(cmd.CommonOptions))
	err = client.Init(initByTweetsCount)
	if err != nil {
		return err
	}

	done := make(chan struct{})
	defer close(done)

	cmd.startWatch(client, postReader, done)

	return nil
}

func (cmd *WatchCommand) startWatch(
	client *twisender.Client,
	postReader joy.PostsReader,
	done chan struct{},
) {
	posts, postAck := watchForPosts(
		postReader,
		cmd.UserName,
		time.Duration(cmd.Period)*time.Second,
		done,
	)
	for {
		select {
		case post := <-posts:
			tweetID, exists, postErr := client.PostNew(post.Link, post.ImgURL)

			if postErr != nil {
				log.Errorf("Post '%s': error publishing tweet: %s", post.Link, postErr.Error())
				postAck <- true
			} else if exists {
				log.Infof("Post '%s' exists, waiting for next check...", post.Link)
				postAck <- false
			} else {
				log.Infof("Post '%s' added, tweet id: %d", post.Link, tweetID)
				postAck <- true
			}

		case <-done:
			return
		}
	}
}
