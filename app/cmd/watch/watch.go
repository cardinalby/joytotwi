package watch

import (
	"joytotwi/app/cmd/common"
	"joytotwi/app/joy"
	"joytotwi/app/joy/selector"
	"joytotwi/app/twisender"
	"time"

	log "github.com/sirupsen/logrus"
)

const postsReverse = false
const postsOffset = 0
const postsLimit = 0
const initByTweetsCount = twisender.TimelineMaxCount

// Command for checking for new posts periodically and post them to twitter
type Command struct {
	Period int `short:"p" long:"period" default:"43200" description:"Period of checking for new posts in seconds"`
	common.Options
}

// SetCommonOptions sets common options in command
func (cmd *Command) SetCommonOptions(opts *common.Options) {
	cmd.Options = *opts
}

// Execute command method for flags.Commander
func (cmd *Command) Execute(args []string) error {
	//noinspection GoBoolExpressions
	postReader, err := selector.GetPostReader(cmd.SourceType, postsReverse, postsOffset, postsLimit)
	if err != nil {
		return err
	}

	client := twisender.CreateNewClient(cmd.Options.GetTwiCreds())
	err = client.Init(initByTweetsCount)
	if err != nil {
		return err
	}

	done := make(chan struct{})
	defer close(done)

	cmd.startWatch(client, postReader, done)

	return nil
}

func (cmd *Command) startWatch(
	client *twisender.Client,
	postReader joy.PostsReader,
	done <-chan struct{},
) {
	consumePosts := func(posts []*joy.Post) {
		log.Infof("%d new posts found", len(posts))
		for _, post := range posts {
			select {
			case <-done:
				return
			default:
				tweetID, _, postErr := client.PostNew(post.Link, post.ImgURL)
				if postErr != nil {
					log.Errorf("Post '%s': error publishing tweet: %s", post.Link, postErr.Error())
					return
				}
				log.Infof("Post '%s' added, tweet id: %d", post.Link, tweetID)
			}
		}
	}

	watchForPosts(
		client,
		postReader,
		cmd.UserName,
		time.Duration(cmd.Period)*time.Second,
		consumePosts,
		done,
	)

	<-done
}
