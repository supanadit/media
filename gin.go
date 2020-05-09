package media

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var ginMedia *GinMedia

type GinMedia struct {
	Engine               *gin.Engine
	SingleUploadRoute    string
	BrowseDirectoryRoute string
	CreateDirectoryRoute string
	Destination          string
}

type UploadResponse struct {
	Message string `json:"message"`
	File    string `json:"file"`
}

func Gin(e *gin.Engine) *GinMedia {
	ginMedia = &GinMedia{
		Engine:               e,
		SingleUploadRoute:    "upload",
		BrowseDirectoryRoute: "directory/browse",
		CreateDirectoryRoute: "directory/create",
	}
	return ginMedia
}

func (m GinMedia) SetDestination(d string) GinMedia {
	m.Destination = d
	return m
}

func (m GinMedia) GetDestination() string {
	return m.Destination
}

func (m GinMedia) SetSingleUploadRoute(p string) GinMedia {
	m.SingleUploadRoute = p
	return m
}

func (m GinMedia) SetBrowseDirectoryRoute(p string) GinMedia {
	m.BrowseDirectoryRoute = p
	return m
}

func (m GinMedia) SetCreateDirectoryRoute(p string) GinMedia {
	m.CreateDirectoryRoute = p
	return m
}

func (m GinMedia) Browse() (md []Media) {
	root := "./"
	if m.Destination != "" {
		root = m.Destination
	}
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		t := "file"
		if f.IsDir() {
			t = "directory"
		}
		nm := Media{
			Name: f.Name(),
			Type: t,
		}
		md = append(md, nm)
	}
	return md
}

func (m GinMedia) Create() GinMedia {
	m.Engine.POST(fmt.Sprintf("/%s", m.SingleUploadRoute), func(c *gin.Context) {
		file, err := c.FormFile("file")
		path := c.PostForm("path")

		if err != nil {
			ur := UploadResponse{
				Message: err.Error(),
			}
			c.JSON(http.StatusInternalServerError, ur)
			return
		}

		f := filepath.Base(file.Filename)
		lf := f
		if m.GetDestination() != "" {
			if path != "" {
				lf = filepath.Join(m.GetDestination(), path, f)
			} else {
				lf = filepath.Join(m.GetDestination(), f)
			}
		}
		if err = c.SaveUploadedFile(file, lf); err != nil {
			ur := UploadResponse{
				Message: err.Error(),
			}
			c.JSON(http.StatusInternalServerError, ur)
			return
		}

		ur := UploadResponse{
			Message: "File uploaded successfully",
			File:    fmt.Sprintf("%s", file.Filename),
		}
		c.JSON(http.StatusOK, ur)
	})
	m.Engine.GET(fmt.Sprintf("/%s", m.BrowseDirectoryRoute), func(c *gin.Context) {
		c.JSON(http.StatusOK, m.Browse())
	})
	m.Engine.POST(fmt.Sprintf("/%s", m.CreateDirectoryRoute), func(c *gin.Context) {
		n := c.PostForm("name")
		p := c.PostForm("path")

		valid := true
		if p != "" {
			valid = !strings.ContainsAny(p, "../")
		}

		fj := filepath.Join(m.Destination, n)
		if p != "" {
			fj = filepath.Join(m.Destination, p, n)
		}

		_, err := os.Stat(fj)
		if os.IsNotExist(err) {
			ck := strings.ContainsAny(n, "/") || strings.ContainsAny(n, "\\")
			if !ck {
				if valid {
					fmt.Println(p)
					err := os.MkdirAll(fj, 0755)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"message": fmt.Sprintf("%s", err.Error()),
						})
						return
					} else {
						d := strings.ReplaceAll(n, "/", " > ")
						c.JSON(http.StatusOK, gin.H{
							"message": fmt.Sprintf("Success create directory %s", d),
						})
						return
					}
				} else {
					c.JSON(http.StatusOK, gin.H{
						"message": "Cannot create directory outside root directory",
					})
				}
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": "Directory name cannot contain '/' or '\\' symbol",
				})
			}
		} else {
			d := strings.ReplaceAll(n, "/", " > ")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Directory %s already exist", d),
			})
			return
		}
	})
	return m
}
