package tweet

import (
	"time"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/go-twitter/twitter"
)

type TweetQueue struct {
	queued []string
	quit chan interface{}

	client *twitter.Client
}

func NewTweetQueue(consumerKey, consumerSecret, accessToken, accessSecret string) (TweetQueue, chan interface{}) {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)
	quit := make(chan interface{})

	return TweetQueue{
		queued: make([]string, 0),
		quit: quit,

		client: client,
	}, quit
}

func (t *TweetQueue) QueueTweet(tweet string) {
	t.queued = append(t.queued, tweet)
}

func (t *TweetQueue) HandleTweetQueue() {
	for {
		select {
		case <-t.quit:
			break
		default:
			remaining := make([]string, 0)
			for _, tweet := range t.queued {
				_, _, err := t.client.Statuses.Update(tweet, nil)
				if err != nil {
					remaining = append(remaining, tweet)
				}
				t.queued = remaining
			}
		}
		time.Sleep(5 * time.Second)
	}
}
