package main

import (
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/gdk"
	"os"
)


func main() {
	//--------------------------------------------------------
	// Setting up the GTK-Foo
	//--------------------------------------------------------
	gdk.ThreadsInit()
	gtk.Init(&os.Args)
	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
	window.SetTitle("Zwitscher")
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
	vbox.Add(notebook)
	
	gtk.
	
	entry := gtk.Entry()
	

	button := gtk.ButtonWithLabel("Update Timeline")
	button.SetTooltipMarkup("update <b>public timeline</b>")
	button.Clicked(func() {
	})

	vbox.PackEnd(button, false, false, 0)

	//--------------------------------------------------------
	// Event
	//--------------------------------------------------------
	window.Add(vbox)
	//	window.SetSizeRequest(800, 500)
	window.ShowAll()
	gdk.ThreadsEnter()
	gtk.Main()
	gdk.ThreadsLeave()
}

