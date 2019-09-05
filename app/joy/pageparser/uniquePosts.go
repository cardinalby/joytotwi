package pageparser

import "joytotwi/app/joy"

func getUniquePostsChan(posts chan *joy.Post, done chan struct{}) chan *joy.Post {
	processed := make(map[string]bool)
	outPosts := make(chan *joy.Post)

	go func() {
		defer close(outPosts)

		for {
			select {
			case post := <-posts:
				if post == nil {
					return
				}
				_, exists := processed[post.Link]
				if !exists {
					processed[post.Link] = true
					outPosts <- post
				}
			case <-done:
				return
			}
		}
	}()

	return outPosts
}
