package joy

// Post represents individual post on joyReactor
type Post struct {
	Link   string
	ImgURL string
}

// PostsReader describes interface for exported readers exported to outside
type PostsReader func(
	userName string,
	done chan struct{},
) (chan *Post, chan error)
