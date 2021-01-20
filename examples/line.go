package main

import (
	"github.com/jwatson-CO-edu/httpcanvas"
	"time"
)

func app(context *httpcanvas.Context) {
	context.BeginPath()
	context.MoveTo(50, 50)
	context.LineTo(200, 200)
	context.Stroke()

	time.Sleep(5 * time.Second)
	context.BeginPath()
	context.MoveTo(200, 50)
	context.LineTo(50, 200)
	context.Stroke()
	
	context.ShowFrame()
}

func main() {
	httpcanvas.ListenAndServe(":8080", app)
}
