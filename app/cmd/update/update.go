package update

import (
	log "github.com/sirupsen/logrus"
	"joytotwi/app/cmd/common"
	"joytotwi/app/joy"
	"joytotwi/app/joy/selector"
	"joytotwi/app/twisender"
)

const postsReverse = false
const postsOffset = 0
const postsLimit = 0
const initByTweetsCount = twisender.TimelineMaxCount

// Command for checking for new joy posts once and post them to twitter
type Command struct {
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

	cmd.checkOnce(client, postReader, done)

	return nil
}

func (cmd *Command) checkOnce(
	client *twisender.Client,
	postReader joy.PostsReader,
	done <-chan struct{},
) {
	posts, err := common.GetNewPosts(client, postReader, cmd.UserName, done)
	if err != nil {
		log.Errorf("Can't get new posts: %s", err.Error())
		return
	}
	if posts == nil {
		log.Info("No new posts found")
		return
	}

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
