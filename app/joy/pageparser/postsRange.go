package pageparser

import "joytotwi/app/joy"

func limitPostsRange(
	offset int,
	limit int,
	posts <-chan *joy.Post,
	postErrors <-chan error,
	done <-chan struct{},
) (<-chan *joy.Post, <-chan error) {
	outPosts := make(chan *joy.Post)
	outErrors := make(chan error)

	go func() {
		defer close(outPosts)
		defer close(outErrors)

		for postNumber := 1; ; postNumber++ {
			isBeforeRange := postNumber < offset
			isAfterRange := limit != 0 && postNumber >= offset+limit

			select {
			case post := <-posts:
				if isBeforeRange {
					continue
				}
				if isAfterRange {
					return
				}
				outPosts <- post
			case err := <-postErrors:
				if isBeforeRange {
					continue
				}
				if isAfterRange {
					return
				}
				outErrors <- err
			case <-done:
				return
			}
		}
	}()
	return outPosts, outErrors
}
