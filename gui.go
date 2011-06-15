package main

import (
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/gdk"
	"os"
	"path/filepath"
	"strconv"
	"utf8"
	"gotter"
	"time"
)

func Gui() {
	//--------------------------------------------------------
	// Setting up the GTK-Foo
	//--------------------------------------------------------
	gdk.ThreadsInit()
	gtk.Init(&os.Args)
	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
	window.SetTitle("Zwitscher!")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	vbox := gtk.VBox(false, 1)

	notebook := gtk.Notebook()
	//--------------------------------------------------------
	// Home View
	//--------------------------------------------------------
	vboxHome := gtk.VBox(false, 1)
	scrolledwinHome := gtk.ScrolledWindow(nil, nil)
	scrolledwinHome.SetPolicy(gtk.GTK_POLICY_NEVER, gtk.GTK_POLICY_ALWAYS) //Disable hscrollbar, enable vscrollbar
	vboxHome.Add(scrolledwinHome)
	vboxscrolledwinHome := gtk.VBox(false, 1)
	scrolledwinHome.AddWithViewPort(vboxscrolledwinHome)

	buttonUT := gtk.ButtonWithLabel("Update Timeline")
	buttonUT.Clicked(func() {
		var tweet gotter.Tweet
		tweets, err := gotter.GetTweets(accounts.Credentials, "https://api.twitter.com/1/statuses/home_timeline.json", map[string]string{})
		if err != nil {
			println("failed to get tweets:", err)
		}
		for i := len(tweets) - 1; i >= 0; i-- {
			tweet = tweets[i]
			id, _ := strconv.Atoi64(tweet.Identifier)
			if accounts.Maxreadid < id {
				tweetwidget := TweetWidget(tweet)
				vboxscrolledwinHome.PackEnd(tweetwidget, false, false, 0)
				tweetwidget.ShowAll()
				accounts.Maxreadid = id
			}
		}
	})
	vboxHome.PackEnd(buttonUT, false, false, 0)

	notebook.AppendPage(vboxHome, gtk.Label("Home"))
	//--------------------------------------------------------
	// gtk.Notebook
	//--------------------------------------------------------
	scrolledwin := gtk.ScrolledWindow(nil, nil)
	scrolledwin = gtk.ScrolledWindow(nil, nil)
	notebook.AppendPage(scrolledwin, gtk.Label("Mentions"))
	scrolledwin = gtk.ScrolledWindow(nil, nil)
	notebook.AppendPage(scrolledwin, gtk.Label("Messages"))
	//scrolledwin = gtk.ScrolledWindow(nil, nil)
	//notebook.AppendPage(scrolledwin, gtk.Label("Faviores"))
	//scrolledwin = gtk.ScrolledWindow(nil, nil)
	//notebook.AppendPage(scrolledwin, gtk.Label("Retweets"))
	//scrolledwin = gtk.ScrolledWindow(nil, nil)
	//notebook.AppendPage(scrolledwin, gtk.Label("Search"))

	//--------------------------------------------------------
	// Public Timeline View
	//--------------------------------------------------------
	vboxPT := gtk.VBox(false, 1)
	scrolledwinPT := gtk.ScrolledWindow(nil, nil)
	textview := gtk.TextView()
	textview.SetEditable(false)
	textview.SetCursorVisible(false)
	scrolledwinPT.Add(textview)

	buffer := textview.GetBuffer()
	tag := buffer.CreateTag("blue", map[string]string{
		"foreground": "#0000FF", "weight": "700"})

	button := gtk.ButtonWithLabel("Update Timeline")
	button.SetTooltipMarkup("update <b>public timeline</b>")
	button.Clicked(func() {
		UpdatePublicTimeline(func(tweet *Tweet) {
			var iter gtk.GtkTextIter
			gdk.ThreadsEnter()
			buffer.GetStartIter(&iter)
			buffer.InsertPixbuf(&iter, tweet.User.ProfileImagePixbuf)
			gdk.ThreadsLeave()
			gdk.ThreadsEnter()
			buffer.Insert(&iter, " ")
			buffer.InsertWithTag(&iter, tweet.User.Name, tag)
			buffer.Insert(&iter, ":"+tweet.Text+"\n")
			gtk.MainIterationDo(false)
			gdk.ThreadsLeave()
		})
	})

	//	button.Clicked()
	vboxPT.Add(scrolledwinPT)
	vboxPT.PackEnd(button, false, false, 0)
	//--------------------------------------------------------
	// End Public Timeline View
	//--------------------------------------------------------

	notebook.AppendPage(vboxPT, gtk.Label("Public Timeline"))
	vbox.Add(notebook)

	//--------------------------------------------------------
	// Fild for Tweets
	//--------------------------------------------------------
	hbox := gtk.HBox(false, 1)

	dir, _ := filepath.Split(os.Args[0])
	imagefile := filepath.Join(dir, "Awesome Smiley Original.jpg")
	image := gtk.ImageFromFile(imagefile)
	hbox.PackStart(image, false, false, 0)

	buttonZwitscher := gtk.ButtonWithLabel("Zwitscher!")
	newTweetTextField := gtk.Entry()
	charCounterLabel := gtk.Label("140")

	buttonZwitscher.SetTooltipMarkup("Tweet")

	buttonZwitscher.Clicked(func() {
		charCounterLabel.SetLabel("140")
		SendTweet(newTweetTextField.GetText())
		newTweetTextField.SetText("")
	})

	newTweetTextField.Connect("key-release-event", func() {
		length := utf8.RuneCountInString(newTweetTextField.GetText())
		charCounterLabel.SetLabel((string)(strconv.Itoa(140 - length)))
	})

	newTweetTextField.Connect("activate", func() {
		if newTweetTextField.GetText() != "" { //pressed enter, and text is not empty
			charCounterLabel.SetLabel("140")
			SendTweet(newTweetTextField.GetText())
			newTweetTextField.SetText("")
		}
	})
	hbox.PackStartDefaults(newTweetTextField)
	hbox.PackStart(charCounterLabel, false, false, 0)
	hbox.PackEnd(buttonZwitscher, false, false, 0)

	vbox.PackEnd(hbox, false, false, 0)

	//--------------------------------------------------------
	// Event
	//--------------------------------------------------------
	window.Add(vbox)
	window.SetSizeRequest(400, 550)
	window.ShowAll()

	gdk.ThreadsEnter()
	gtk.Main()
	gdk.ThreadsLeave()
}

func TweetWidget(tweet gotter.Tweet) *gtk.GtkFrame {
	frame := gtk.Frame(tweet.User.ScreenName)
	hbox := gtk.HBox(false, 1)
	dir, _ := filepath.Split(os.Args[0])
	//gtk_image_new_from_pixbuf is unsupported in go-gtk )=
	imagefile := filepath.Join(dir, "Awesome Smiley Original.jpg")
	image := gtk.ImageFromFile(imagefile)
	vbox := gtk.VBox(false, 1)
	tweettext := gtk.TextView()
	tweettext.SetWrapMode(gtk.GTK_WRAP_WORD)
	tweettext.SetEditable(false)
	tweetbuffer := tweettext.GetBuffer()

	tweetbuffer.SetText(tweet.Text)

	tweettime, _ := time.Parse(time.RubyDate, tweet.CreatedAt)
	hour := (string)(strconv.Itoa(tweettime.Hour))
	var minute string
	if tweettime.Minute < 10 {
		minute = "0" + (string)(strconv.Itoa(tweettime.Minute))
	} else {
		minute = (string)(strconv.Itoa(tweettime.Minute))
	}
	whenfromtext := gtk.Label(hour + ":" + minute)
	//	wherefromtext := gtk.Label(tweet.Source)

	hbox.PackStart(image, false, false, 0)
	hbox.PackEndDefaults(vbox)
	vbox.PackStart(tweettext, false, false, 0)
	vbox.PackEnd(whenfromtext, false, false, 0)
	//	vbox.PackEnd(wherefromtext, false, false, 0)

	frame.Add(hbox)

	return frame
}

