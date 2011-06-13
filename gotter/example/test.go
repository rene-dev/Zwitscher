package main

import (
	"io/ioutil"
	"json"
	"log"
	"os"
	"strings"
	"flag"
	"fmt"
	"gotter"
)

func main() {
	reply := flag.Bool("r", false, "show replies")
	list := flag.String("l", "", "show tweets")
	user := flag.String("u", "", "show user timeline")
	favorite := flag.String("f", "", "specify favorite ID")
	//search := flag.String("s", "", "search word")
	inreply := flag.String("i", "", "specify in-reply ID, if not specify text, it will be RT.")
	verbose := flag.Bool("v", false, "detail display")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage of twty:
  -f ID: specify favorite ID
  -i ID: specify in-reply ID, if not specify text, it will be RT.
  -l USER/LIST: show list's timeline (ex: mattn_jp/subtech)
  -u USER: show user's timeline
  -r: show replies
  -v: detail display
`)
	}
	flag.Parse()

	file, config := gotter.GetConfig()
	token, authorized, err := gotter.GetAccessToken(config)
	if err != nil {
		log.Fatal("faild to get access token:", err)
	}
	if authorized {
		b, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal("failed to store file:", err)
		}
		err = ioutil.WriteFile(file, b, 0700)
		if err != nil {
			log.Fatal("failed to store file:", err)
		}
	}


	if *reply {
		tweets, err := gotter.GetTweets(token, "https://api.twitter.com/1/statuses/mentions.json", map[string]string{})
		if err != nil {
			log.Fatal("failed to get tweets:", err)
		}
		gotter.ShowTweets(tweets, *verbose)
	} else if len(*list) > 0 {
		part := strings.Split(*list, "/", 2)
		tweets, err := gotter.GetTweets(token, "https://api.twitter.com/1/"+part[0]+"/lists/"+part[1]+"/statuses.json", map[string]string{})
		if err != nil {
			log.Fatal("failed to get tweets:", err)
		}
		gotter.ShowTweets(tweets, *verbose)
	} else if len(*user) > 0 {
		tweets, err := gotter.GetTweets(token, "https://api.twitter.com/1/statuses/user_timeline.json", map[string]string{"screen_name": *user})
		if err != nil {
			log.Fatal("failed to get tweets:", err)
		}
		gotter.ShowTweets(tweets, *verbose)
	} else if len(*favorite) > 0 {
		gotter.PostTweet(token, "https://api.twitter.com/1/favorites/create/"+*favorite+".json", map[string]string{})
	} else if flag.NArg() == 0 {
		if len(*inreply) > 0 {
			gotter.PostTweet(token, "https://api.twitter.com/1/statuses/retweet/"+*inreply+".json", map[string]string{})
		} else {
			tweets, err := gotter.GetTweets(token, "https://api.twitter.com/1/statuses/home_timeline.json", map[string]string{})
			if err != nil {
				log.Fatal("failed to get tweets:", err)
			}
			gotter.ShowTweets(tweets, *verbose)
		}
	} else {
		gotter.PostTweet(token, "https://api.twitter.com/1/statuses/update.json", map[string]string{"status": strings.Join(flag.Args(), " "), "in_reply_to_status_id": *inreply})
	}
}