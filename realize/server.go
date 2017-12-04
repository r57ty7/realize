package realize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
)

// Dafault host and port
const (
	Host = "localhost"
	Port = 5002
)

// Server settings
type Server struct {
	Parent *Realize `yaml:"-" json:"-"`
	Status bool     `yaml:"status" json:"status"`
	Open   bool     `yaml:"open" json:"open"`
	Port   int      `yaml:"port" json:"port"`
	Host   string   `yaml:"host" json:"host"`
}

// Websocket projects
func (s *Server) projects(c echo.Context) (err error) {
	websocket.Handler(func(ws *websocket.Conn) {
		msg, _ := json.Marshal(s.Parent)
		err = websocket.Message.Send(ws, string(msg))
		go func() {
			for {
				select {
				case <-s.Parent.Sync:
					msg, _ := json.Marshal(s.Parent)
					err = websocket.Message.Send(ws, string(msg))
					if err != nil {
						break
					}
				}
			}
		}()
		for {
			// Read
			text := ""
			err = websocket.Message.Receive(ws, &text)
			if err != nil {
				break
			} else {
				err := json.Unmarshal([]byte(text), &s.Parent)
				if err == nil {
					s.Parent.Settings.Write(s.Parent)
					break
				}
			}
		}
		ws.Close()
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

// Render return a web pages defined in bindata
func (s *Server) render(c echo.Context, path string, mime int) error {
	data, err := Asset(path)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	rs := c.Response()
	// check content type by extensions
	switch mime {
	case 1:
		rs.Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
		break
	case 2:
		rs.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJavaScriptCharsetUTF8)
		break
	case 3:
		rs.Header().Set(echo.HeaderContentType, "text/css")
		break
	case 4:
		rs.Header().Set(echo.HeaderContentType, "image/svg+xml")
		break
	case 5:
		rs.Header().Set(echo.HeaderContentType, "image/png")
		break
	}
	rs.WriteHeader(http.StatusOK)
	rs.Write(data)
	return nil
}

// Start the web server
func (s *Server) Start() (err error) {
	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 2,
	}))
	e.Use(middleware.Recover())

	// web panel
	e.GET("/", func(c echo.Context) error {
		return s.render(c, "assets/index.html", 1)
	})
	e.GET("/assets/js/all.min.js", func(c echo.Context) error {
		return s.render(c, "assets/assets/js/all.min.js", 2)
	})
	e.GET("/assets/css/app.css", func(c echo.Context) error {
		return s.render(c, "assets/assets/css/app.css", 3)
	})
	e.GET("/app/components/settings/index.html", func(c echo.Context) error {
		return s.render(c, "assets/app/components/settings/index.html", 1)
	})
	e.GET("/app/components/project/index.html", func(c echo.Context) error {
		return s.render(c, "assets/app/components/project/index.html", 1)
	})
	e.GET("/app/components/index.html", func(c echo.Context) error {
		return s.render(c, "assets/app/components/index.html", 1)
	})
	e.GET("/assets/img/svg/settings.svg", func(c echo.Context) error {
		return s.render(c, "assets/assets/img/svg/settings.svg", 4)
	})
	e.GET("/assets/img/svg/fullscreen.svg", func(c echo.Context) error {
		return s.render(c, "assets/assets/img/svg/fullscreen.svg", 4)
	})
	e.GET("/assets/img/svg/add.svg", func(c echo.Context) error {
		return s.render(c, "assets/assets/img/svg/add.svg", 4)
	})
	e.GET("/assets/img/svg/backspace.svg", func(c echo.Context) error {
		return s.render(c, "assets/assets/img/svg/backspace.svg", 4)
	})
	e.GET("/assets/img/svg/error.svg", func(c echo.Context) error {
		return s.render(c, "assets/assets/img/svg/error.svg", 4)
	})
	e.GET("/assets/img/svg/remove.svg", func(c echo.Context) error {
		return s.render(c, "assets/assets/img/remove.svg", 4)
	})
	e.GET("/assets/img/svg/logo.svg", func(c echo.Context) error {
		return s.render(c, "assets/assets/img/svg/logo.svg", 4)
	})
	e.GET("/assets/img/fav.png", func(c echo.Context) error {
		return s.render(c, "assets/assets/img/fav.png", 5)
	})
	e.GET("/assets/img/svg/circle.svg", func(c echo.Context) error {
		return s.render(c, "assets/assets/img/svg/circle.svg", 4)
	})

	//websocket
	e.GET("/ws", s.projects)
	e.HideBanner = true
	e.Debug = false
	go func() {
		log.Println(s.Parent.Prefix("Started on " + string(s.Host) + ":" + strconv.Itoa(s.Port)))
		e.Start(string(s.Host) + ":" + strconv.Itoa(s.Port))
	}()
	return nil
}

// OpenURL in a new tab of default browser
func (s *Server) OpenURL() (io.Writer, error) {
	url := "http://" + string(s.Parent.Server.Host) + ":" + strconv.Itoa(s.Parent.Server.Port)
	stderr := bytes.Buffer{}
	cmd := map[string]string{
		"windows": "start",
		"darwin":  "open",
		"linux":   "xdg-open",
	}
	if s.Open {
		open, err := cmd[runtime.GOOS]
		if !err {
			return nil, fmt.Errorf("operating system %q is not supported", runtime.GOOS)
		}
		cmd := exec.Command(open, url)
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return cmd.Stderr, err
		}
	}
	return nil, nil
}
