package cmd

// AppOptions with all cli commands and flags
type AppOptions struct {
	Watch WatchCommand `command:"watch"`
	Fill  FillCommand  `command:"fill"`

	ConfigFile string `short:"c" long:"config-file" default:"" description:"Load app options from json file"`

	SourceType        string `short:"s" long:"source-type" default:"page" description:"Posts source: rss or page"`
	UserName          string `short:"u" long:"user-name" description:"JoyReactor user name"`
	AccessToken       string `long:"access-token" description:"Twitter API AccessToken"`
	AccessTokenSecret string `long:"access-token-secret" description:"Twitter API Access Token Secret"`
	ConsumerKey       string `long:"consumer-key" description:"Twitter API Consumer Key"`
	ConsumerSecret    string `long:"consumer-secret" description:"Twitter API Consumer Secret"`
}
