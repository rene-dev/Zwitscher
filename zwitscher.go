package main

import "gotter"

func main() {
	Connect()
	gotter.PostTweet(accounts.Credentials, "https://api.twitter.com/1/statuses/update.json", map[string]string{"status": "Eat this, @simonszu ! Hard codet Tweets!"})
}

