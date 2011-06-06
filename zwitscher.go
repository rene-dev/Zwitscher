// This File is based on mattn example of go-gtk programms.
// Actualy it is basicly the same, but not for long.
// https://github.com/mattn/go-gtk/blob/master/example/twitter/twitter.go

package main

import (
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"http"
	"json"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"path"
	"unsafe"
)

func url2pixbuf(url string) *gdkpixbuf.GdkPixbuf {
	r, err := http.Get(url)
	if err != nil {
		return nil
	}
	t := r.Header.Get("Content-Type")
	b := make([]byte, r.ContentLength)
	if _, err = io.ReadFull(r.Body, b); err != nil {
		return nil
	}
	var loader *gdkpixbuf.GdkPixbufLoader
	if strings.Index(t, "jpeg") >= 0 {
		loader, _ = gdkpixbuf.PixbufLoaderWithMimeType("image/jpeg")
	} else {
		loader, _ = gdkpixbuf.PixbufLoaderWithMimeType("image/png")
	}
	loader.SetSize(24, 24)
	loader.Write(b)
	loader.Close()
	return loader.GetPixbuf()
}

func sendTweet(text string) {
	println(text)
	}

func main() {
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

	buffer := textview.GetBuffer()

	tag := buffer.CreateTag("blue", map[string]string{
		"foreground": "#0000FF", "weight": "700"})
	button := gtk.ButtonWithLabel("Update Timeline")
	button.SetTooltipMarkup("update <b>public timeline</b>")
	button.Clicked(func() {
		go func() {
			gdk.ThreadsEnter()
			button.SetSensitive(false)
			gdk.ThreadsLeave()
			r, err := http.Get("http://twitter.com/statuses/public_timeline.json")
			if err == nil {
				var b []byte
				if r.ContentLength == -1 {
					b, err = ioutil.ReadAll(r.Body)
				} else {
					b = make([]byte, r.ContentLength)
					_, err = io.ReadFull(r.Body, b)
				}
				if err != nil {
					println(err.String())
					return
				}
				var j interface{}
				json.NewDecoder(bytes.NewBuffer(b)).Decode(&j)
				arr := j.([]interface{})
				for i := 0; i < len(arr); i++ {
					data := arr[i].(map[string]interface{})
					icon := data["user"].(map[string]interface{})["profile_image_url"].(string)
					var iter gtk.GtkTextIter
					gdk.ThreadsEnter()
					buffer.GetStartIter(&iter)
					buffer.InsertPixbuf(&iter, url2pixbuf(icon))
					gdk.ThreadsLeave()
					name := data["user"].(map[string]interface{})["screen_name"].(string)
					text := data["text"].(string)
					gdk.ThreadsEnter()
					buffer.Insert(&iter, " ")
					buffer.InsertWithTag(&iter, name, tag)
					buffer.Insert(&iter, ":"+text+"\n")
					gtk.MainIterationDo(false)
					gdk.ThreadsLeave()
				}
			}
			button.SetSensitive(true)
		}()
	})
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
	// GtkImage
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
		sendTweet(newTweetTextField.GetText())
		newTweetTextField.SetText("")
	})
	
	newTweetTextField.Connect("key-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		kev := *(**gdk.EventKey)(unsafe.Pointer(&arg))
		if(kev.Keyval == 65293 && newTweetTextField.GetText() != ""){//pressed enter, and text is not empty
			sendTweet(newTweetTextField.GetText())
			newTweetTextField.SetText("")
		}//Count remaining characters here
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

