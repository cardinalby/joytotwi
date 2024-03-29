package fill

import (
	"errors"
	"fmt"
	"joytotwi/app/cmd/common"
	"joytotwi/app/joy"
	"joytotwi/app/joy/selector"
	"joytotwi/app/twisender"

	log "github.com/sirupsen/logrus"
)

// Command is intended to initially fill account by all posts
type Command struct {
	Offset         int  `short:"o" long:"offset" default:"1" description:"Number of first post to start"`
	Limit          int  `short:"l" long:"limit" default:"0" description:"How many posts to process. 0 to process all"`
	StopOnExistent bool `long:"stop-on-existent" description:"Stop parse posts after first existent tweet"`
	StopOnError    bool `long:"stop-on-error" description:"Stop parse posts after first error"`
	common.Options
}

// SetCommonOptions sets common options in command
func (cmd *Command) SetCommonOptions(opts *common.Options) {
	cmd.Options = *opts
}

// Execute command method for flags.Commander
func (cmd *Command) Execute(args []string) error {
	postReader, err := selector.GetPostReader(cmd.SourceType, true, cmd.Offset, cmd.Limit)
	if err != nil {
		return err
	}

	client := twisender.CreateNewClient(cmd.Options.GetTwiCreds())
	err = client.Init(twisender.TimelineMaxCount)
	if err != nil {
		return err
	}

	done := make(chan struct{})
	defer close(done)

	stats, fillErr := cmd.performFill(client, postReader, cmd.Offset, done)
	if fillErr != nil {
		return fillErr
	}
	stats.print()
	return nil
}

func (cmd *Command) performFill(
	client *twisender.Client,
	postReader joy.PostsReader,
	offset int,
	done <-chan struct{},
) (fillStats, error) {
	posts, postErrors := postReader(cmd.UserName, done)

	stats := fillStats{}
	for {
		var itemError error

		select {
		case post := <-posts:
			if post == nil {
				log.Info("End of posts reached")
				return stats, nil
			}
			offset++
			postLogName := getPostLogName(post, offset)
			tweetID, exists, postErr := client.PostNew(post.Link, post.ImgURL)

			if postErr != nil {
				stats.errorsCount++
				itemError = fmt.Errorf("%s: error publishing tweet: %s", postLogName, postErr.Error())
			} else if exists {
				stats.existentCount++
				existsMessage := fmt.Sprintf("%s: tweet already exists", postLogName)
				if cmd.StopOnExistent {
					return stats, errors.New(existsMessage)
				}
				log.Warn(existsMessage)
			} else {
				stats.postedCount++
				log.Infof("%s: added, tweet id: %d", postLogName, tweetID)
			}
		case err := <-postErrors:
			if err != nil {
				offset++
				stats.errorsCount++
				itemError = fmt.Errorf("post (offset: %d) parse error: %s", offset, err.Error())
			}
		case <-done:
			return stats, nil
		}

		if itemError != nil {
			if cmd.StopOnError {
				return stats, itemError
			}
			log.Warn(itemError.Error())
		}
	}
}

func getPostLogName(post *joy.Post, offset int) string {
	return fmt.Sprintf("post '%s' (offset: %d)", post.Link, offset)
}
