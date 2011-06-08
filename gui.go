package zwitscher

import (
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/gdk"
	//	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	//	"http"
	//	"json"
	//	"bytes"
	//	"io"
	//	"io/ioutil"
	"os"
	//	"strings"
	"path"
	"unsafe"
	"strconv"
	"utf8"
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

	//--------------------------------------------------------
	// gtk.Notebook
	//--------------------------------------------------------
	notebook := gtk.Notebook()

	scrolledwin := gtk.ScrolledWindow(nil, nil)
	notebook.AppendPage(scrolledwin, gtk.Label("Home"))
	scrolledwin = gtk.ScrolledWindow(nil, nil)
	notebook.AppendPage(scrolledwin, gtk.Label("Mentions"))
	scrolledwin = gtk.ScrolledWindow(nil, nil)
	notebook.AppendPage(scrolledwin, gtk.Label("Messages"))
	scrolledwin = gtk.ScrolledWindow(nil, nil)
	notebook.AppendPage(scrolledwin, gtk.Label("Faviores"))
	scrolledwin = gtk.ScrolledWindow(nil, nil)
	notebook.AppendPage(scrolledwin, gtk.Label("Retweets"))
	scrolledwin = gtk.ScrolledWindow(nil, nil)
	notebook.AppendPage(scrolledwin, gtk.Label("Search"))

	//--------------------------------------------------------
	// Public Timeline View
	//--------------------------------------------------------
	vboxPT := gtk.VBox(false, 1)
	scrolledwinPT := gtk.ScrolledWindow(nil, nil)
	textview := gtk.TextView()
	textview.SetEditable(false)
	textview.SetCursorVisible(false)
	scrolledwinPT.Add(textview)

	//	buffer := textview.GetBuffer()

	//	tag := buffer.CreateTag("blue", map[string]string{
	//		"foreground": "#0000FF", "weight": "700"})
	button := gtk.ButtonWithLabel("Update Timeline")
	button.SetTooltipMarkup("update <b>public timeline</b>")
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

	//--------------------------------------------------------
	// Tweetbar
	//--------------------------------------------------------
	dir, _ := path.Split(os.Args[0])
	imagefile := path.Join(dir, "/Awesome Smiley Original.jpg")
	image := gtk.ImageFromFile(imagefile)
	hbox.Add(image)

	buttonZwitscher := gtk.ButtonWithLabel("Zwitscher!")
	newTweetTextField := gtk.Entry()
	charCounterLabel := gtk.Label("140")

	buttonZwitscher.SetTooltipMarkup("Tweet")

	buttonZwitscher.Clicked(func() {
		SendTweet(newTweetTextField.GetText())
		newTweetTextField.SetText("")
	})

	newTweetTextField.Connect("key-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		kev := *(**gdk.EventKey)(unsafe.Pointer(&arg))
		if kev.Keyval == 65293 && newTweetTextField.GetText() != "" { //pressed enter, and text is not empty
			SendTweet(newTweetTextField.GetText())
			newTweetTextField.SetText("")
		} else {
			length := utf8.RuneCountInString(newTweetTextField.GetText())
			charCounterLabel.SetLabel((string)(strconv.Itoa(140 - length))) //update does not work correct :(
		}
	})

	hbox.Add(newTweetTextField)
	hbox.Add(buttonZwitscher)
	hbox.Add(charCounterLabel)

	vbox.PackEnd(hbox, false, false, 0)

	//--------------------------------------------------------
	// Event
	//--------------------------------------------------------
	window.Add(vbox)
	window.SetSizeRequest(500, 600)
	window.ShowAll()
	gdk.ThreadsEnter()
	gtk.Main()
	gdk.ThreadsLeave()
}

