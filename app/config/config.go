package config

// Config root of config
type Config struct {
	UserName     string             `json:"userName" env:"JOY_USER_NAME"`
	TwitterCreds TwitterCredsConfig `json:"twitterCreds"`
}

// TwitterCredsConfig creds for Twitter API access
type TwitterCredsConfig struct {
	AccessToken       string `json:"accessToken" env:"TW_ACCESS_TOKEN"`
	AccessTokenSecret string `json:"accessTokenSecret" env:"TW_ACCESS_TOKEN_SECRET"`
	ConsumerKey       string `json:"consumerKey" env:"TW_CONSUMER_KEY"`
	ConsumerSecret    string `json:"consumerSecret" env:"TW_CONSUMER_SECRET"`
}
