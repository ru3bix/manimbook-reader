package main

import (
	"embed"
	"log"
	"manimbook-reader/book"
	"manimbook-reader/helpers"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var frontendFiles embed.FS // copying this to a different, possibly afero fs is not a good idea

func main() {
	app := NewApp()

	// this is all because of THAT ONE WEBKIT BUG NEVER GETTING FIXED SKDJFSKDJFBKXCSODFSJKDFHKSDFHBBKCK
	// reference:
	// 1. https://github.com/wailsapp/wails/issues/1568
	// 2. https://github.com/tauri-apps/tauri/issues/3725#issuecomment-1160842638
	// 3. https://bugs.webkit.org/show_bug.cgi?id=146351
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()
	addr := listener.Addr().(*net.TCPAddr)
	port := addr.Port
	app.setPort(port)
	go func(app *App) {
		http.Handle("/", helpers.FileServer(app.fs))
		log.Printf("File server running on http://localhost:%d\n", port)
		http.Serve(listener, nil)
	}(app)

	AppMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.AppMenu()) // On macOS platform, this must be done right after `NewMenu()`
	}
	FileMenu := AppMenu.AddSubmenu("File")
	FileMenu.AddText("Open", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
		result, err := helpers.GetManimBookFile(app.ctx)
		if err != nil {
			rt.LogError(app.ctx, "could not open file")
			helpers.DisplayErrorMsg(app.ctx, err)
			return
		}
		rt.LogDebugf(app.ctx, "selected file %s", result)
		book, err := book.InitializeBook(result, app.fs)
		if err != nil {
			rt.LogError(app.ctx, err.Error())
			helpers.DisplayErrorMsg(app.ctx, err)
			return
		}
		app.setBook(*book)
		rt.EventsEmit(app.ctx, "bookOpen")
		rt.LogDebugf(app.ctx, "initialized book %v", app.GetBook())
		return
	})
	FileMenu.AddSeparator()
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		rt.Quit(app.ctx)
	})

	ViewMenu := AppMenu.AddSubmenu("View")
	ViewMenu.AddText("Zoom In", keys.CmdOrCtrl("+"), func(_ *menu.CallbackData) {
		rt.EventsEmit(app.ctx, "incZoom", 0.05)
	})
	ViewMenu.AddSeparator()
	ViewMenu.AddText("Zoom Out", keys.CmdOrCtrl("-"), func(cd *menu.CallbackData) {
		rt.EventsEmit(app.ctx, "decZoom", 0.05)
	})
	ViewMenu.AddSeparator()
	ViewMenu.AddText("Zoom Reset", keys.CmdOrCtrl("="), func(cd *menu.CallbackData) {
		rt.EventsEmit(app.ctx, "setZoom", 1)
	})
	ViewMenu.AddSeparator()
	ShowViewMenu := ViewMenu.AddSubmenu("Show")
	ShowViewMenu.AddCheckbox("Navbar", true, nil, func(cd *menu.CallbackData) {
		rt.LogDebugf(app.ctx, "is checked %t", cd.MenuItem.Checked)
		if !cd.MenuItem.Checked {
			rt.EventsEmit(app.ctx, "hideNavbar")
			cd.MenuItem.SetChecked(false)
		} else {
			rt.EventsEmit(app.ctx, "showNavbar")
			cd.MenuItem.SetChecked(true)
		}
	})
	ShowViewMenu.AddCheckbox("Contents", false, nil, func(cd *menu.CallbackData) {
		if !cd.MenuItem.Checked {
			rt.EventsEmit(app.ctx, "hideContents")
			cd.MenuItem.SetChecked(false)
		} else {
			rt.EventsEmit(app.ctx, "showContents")
			cd.MenuItem.SetChecked(true)
		}
	})
	ViewMenu.AddSeparator()
	ViewMenu.AddText("Fullscreen", keys.Key("f11"), func(cd *menu.CallbackData) {
		if !rt.WindowIsFullscreen(app.ctx) {
			rt.WindowFullscreen(app.ctx)
			rt.EventsEmit(app.ctx, "hideNavbar")
			rt.EventsEmit(app.ctx, "hideContents")
		} else {
			rt.WindowUnfullscreen(app.ctx)
			rt.EventsEmit(app.ctx, "showNavbar")
		}
	})

	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.EditMenu()) // On macOS platform, EditMenu should be appended to enable Cmd+C, Cmd+V, Cmd+Z... shortcuts
	}

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 1 {
		log.Fatalln("more than 1 argument passed")
	}
	if len(argsWithoutProg) == 1 {
		file := os.Args[1]
		if _, err := os.Stat(file); err != nil {
			log.Fatalln("cannot open file")
		}
		book, err := book.InitializeBook(file, app.fs)
		if err != nil {
			log.Fatalln("cannot initialized given book")
		}
		app.setBook(*book)
	}

	err = wails.Run(&options.App{
		Title:  "manimbook reader",
		Width:  1024,
		Height: 768,
		Menu:   AppMenu,
		AssetServer: &assetserver.Options{
			Assets: frontendFiles,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
