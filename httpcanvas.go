package httpcanvas

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
  "strconv"
  "io/ioutil"
)

type mouseMovement struct {
	command string
	x       float64
	y       float64
}

type CanvasHandler func(*Context)

type Canvas struct {
	handler CanvasHandler
	Width   float64
	Height  float64
	Unique  string
	started bool
	command chan string
	mouse   chan mouseMovement
}

func newCanvas(handler CanvasHandler) *Canvas {
	return &Canvas{handler, 640, 480, "", false,
		make(chan string, 10000),
		make(chan mouseMovement, 10000)}
}

func (c *Canvas) updateUnique() {
	c.Unique = fmt.Sprintf("%f", rand.Float64())
}

func (c *Canvas) renderHtml(w http.ResponseWriter, r *http.Request) error {
  data, err := ioutil.ReadFile("resource/container.html")
  container := string( data )
  if err != nil {
		return err
	}
	template,  err := template.New("basic").Parse(container)
	if err != nil {
		return err
	}
	err = template.Execute(w, c)
	return err
}

func (c *Canvas) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		http.NotFound(w, r)
		log.Println(err)
		return
	}
	command := u.Path

	if command == "/jquery.js" {
		// TODO: set mime type
		w.Write(jQuery)
		return
	}

	if command == "/" && r.Method == "GET" {
		c.updateUnique()
		err := c.renderHtml(w, r)
		if err != nil {
			return
		}
		if !c.started {
			c.started = true
			go func() {
				c.handler(newContext(c.Width, c.Height, c.command, c.mouse))
				c.command <- "END"
				close(c.command)
				close(c.mouse)
			}()
		}
		return
	}

	q := u.Query()
	unique := ""
	if _, ok := q["id"]; !ok {
		unique = r.PostFormValue("id")
		if unique == "" {
			http.NotFound(w, r)
			log.Println("missing id", r)
			return
		}
	} else {
		unique = q["id"][0]
	}

	if unique != c.Unique {
		http.NotFound(w, r)
		return
	}

	if command == "/command" && r.Method == "GET" {
		commandGroup := ""
		for command := range c.command {
			if len(commandGroup) > 0 {
				commandGroup += "~"
			}
			commandGroup += command
			if command == "SHOWFRAME" {
				break
			}
		}
		w.Write([]byte(commandGroup))
		return
	}

	if command == "/command" && r.Method == "POST" {
		cmd := strings.Fields(r.PostFormValue("cmd"))
		if len(cmd) == 0 {
			log.Println("missing command")
			http.NotFound(w, r)
			return
		}
		if cmd[0] == "MOUSEMOVE" || cmd[0] == "MOUSECLICK" {
			if len(cmd) == 3 {
				x, err := strconv.Atoi(cmd[1])
				if err != nil {
					http.NotFound(w, r)
					log.Println("invalid x")
					return
				}
				y, err := strconv.Atoi(cmd[2])
				if err != nil {
					http.NotFound(w, r)
					log.Println("invalid y")
					return
				}
				c.mouse <- mouseMovement{cmd[0], float64(x), float64(y)}
				return
			}
		}
	}
	http.NotFound(w, r)
}

func ListenAndServe(addr string, handler CanvasHandler) (err error) {
	return http.ListenAndServe(addr, newCanvas(handler))
}

func stringPartition(s, sep string) (string, string, string) {
	sepPos := strings.Index(s, sep)
	if sepPos == -1 { // no seperator found
		return s, "", ""
	}
	split := strings.SplitN(s, sep, 2)
	return split[0], sep, split[1]
}
