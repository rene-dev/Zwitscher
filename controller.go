package zwitscher

import (
//	"github.com/mattn/go-gtk/gtk"
//"github.com/mattn/go-gtk/gdk"
//	"github.com/mattn/go-gtk/gdkpixbuf"
//"http"
//	"json"
//	"bytes"
//	"io"
//	"io/ioutil"
//	"os"
//	"strings"
//	"path"
//	"unsafe"
)

/*
func UpdatePublicTimeline() {
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
}
*/

/*
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
*/

func SendTweet(text string) {
	println(text)
}

