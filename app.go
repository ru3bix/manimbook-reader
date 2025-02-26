package main

import (
	"context"
	"manimbook-reader/book"

	"github.com/spf13/afero"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx         context.Context
	currentBook book.Book // just some global data we would like to store
	fs          afero.Fs
	port        int
}

func NewApp() *App {
	// unzipped manimbook goes here
	fs := afero.NewMemMapFs()

	return &App{
		fs: fs,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if a.GetBook().Title != "" {
		rt.LogDebugf(a.ctx, "initialized book %v", a.GetBook())
		rt.EventsEmit(a.ctx, "bookOpen")
	}
}

func (a *App) setBook(b book.Book) {
	a.currentBook = b
}

func (a *App) GetBook() book.Book {
	return a.currentBook
}

func (a *App) setPort(port int) {
	a.port = port
}

func (a *App) GetPort() int {
	return a.port
}
