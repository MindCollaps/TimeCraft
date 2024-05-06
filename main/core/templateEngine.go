package core

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"src/main/env"
	"src/main/web/webLogic"
	"strings"
	"text/template"
)

var tmpl = map[string]*template.Template{}

func LoadTemplates(r *gin.Engine) {
	tmpl = make(map[string]*template.Template)

	// Load templates files
	templateFiles := []string{}

	// Check if base template exists
	_, err := fs.ReadFile(env.Files, "main/web/templates/base.gohtml")
	if err != nil {
		panic("Base template not found, ensure it exists in main/web/templates/base.gohtml")
	}

	templateFiles = append(templateFiles, "main/web/templates/base.gohtml")

	log.Println("Loading templates...")
	// Walk through the "templates" folder and all its subdirectories
	nerr := WalkFS(env.Files, "main/web/templates", func(path string, name string, isDir bool) error {
		// Check if the file is an HTML templates
		if !isDir && strings.HasSuffix(name, ".gohtml") {
			if strings.HasSuffix(path, "base.gohtml") {
				return nil
			}
			// Replace backslashes with forward slashes (for Windows compatibility)
			templateName := strings.Replace(path, "\\", "/", -1)

			// Parse the file and add it to the "tmpl" map
			templateFiles = append(templateFiles, path)

			//console logChopper
			log.Println(templateName + " ")
		}
		return nil
	})

	if nerr != nil {
		panic(nerr)
	}

	log.Println("\n\nLoading sites...")

	// Walk through the "public" folder and all its subdirectories
	err = WalkFS(env.Files, "main/web/public", func(path string, name string, isDir bool) error {
		// Check if the file is an HTML template
		if !isDir && strings.HasSuffix(name, ".gohtml") {
			// Get the directory path (relative to the "public" folder)
			relPath, err := RelFS("main/web/public", path)
			if err != nil {
				return err
			}
			// Replace backslashes with forward slashes (for Windows compatibility)
			templateName := strings.Replace(relPath, "\\", "/", -1)

			//Cut first / from path
			if strings.HasPrefix(templateName, "/") {
				templateName = strings.Replace(templateName, "/", "", 1)
			}

			//cut index.gohtml to index
			if strings.HasSuffix(templateName, "index.gohtml") {
				templateName = strings.Replace(templateName, "index.gohtml", "", -1)
			}

			if strings.HasSuffix(path, "app.gohtml") {
				templateName += "app"
			}

			// Parse the file and add it to the "tmpl" map
			parsing := []string{}
			parsing = append(parsing, templateFiles...)
			parsing = append(parsing, path)

			if templateName == "" {
				templateName = "."
			}

			tmpl[templateName] = template.Must(template.ParseFS(env.Files, parsing...))

			log.Println("Serving " + relPath + " at /" + templateName)
			r.GET("/"+templateName, handler)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func handler(c *gin.Context) {
	c.Header("Content-Type", "text/html")

	path := c.Request.URL.Path
	path = strings.Trim(path, "/")

	// If the path is empty, default to "index"
	if path == "" {
		path = "."
	}

	// Look up the templates in the "tmpl" map
	t, ok := tmpl[path]

	if !ok {
		// If the templates doesn't exist, return a 404 error
		c.String(http.StatusNotFound, "Page Not Found")
		return
	}

	// Execute the templates with an empty data object
	logicData := webLogic.GetLogicData(c, path)
	err := t.Execute(c.Writer, logicData)
	if err != nil {
		// If there's an error executing the templates, return a 500 error
		c.String(http.StatusInternalServerError, "Internal Server Error")
		log.Println("Error executing template:", err)
		return
	}
}

func LoadServerAssets(r *gin.Engine) {
	log.Println("Loading assets...")
	// Walk through the "assets" folder and all its subdirectories
	err := WalkFS(env.Files, "main/web/assets", func(path string, name string, isDir bool) error {
		// Check if the file is not a directory
		if !isDir {
			// Get the directory path (relative to the "public" folder)
			relPath, err := RelFS("main/web/assets", path)
			if err != nil {
				return err
			}

			assetPath := strings.Replace(relPath, "\\", "/", -1)

			//Cut first / from path
			if strings.HasPrefix(assetPath, "/") {
				assetPath = strings.Replace(assetPath, "/", "", 1)
			}

			// Add the asset to a route
			log.Println("Serving " + path + " at /assets/" + assetPath)
			r.StaticFileFS("/assets/"+assetPath, path, http.FS(env.Files))

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	r.StaticFileFS("/favicon.ico", "main/web/static/favicon.ico", http.FS(env.Files))
}

// Funky self recursive function to walk through the embed.FS
func WalkFS(fs embed.FS, root string, fn func(path string, name string, isDir bool) error) error {
	entries, err := fs.ReadDir(root)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(root, entry.Name())
		fullPath = filepath.ToSlash(fullPath)
		if entry.IsDir() {
			if err := fn(fullPath, entry.Name(), true); err != nil {
				return err
			}
			err := WalkFS(fs, fullPath, fn)
			if err != nil {
				return err
			}
			continue
		}
		if err := fn(fullPath, entry.Name(), false); err != nil {
			return err
		}
	}
	return nil
}

// Funky function to rel the embed.FS
func RelFS(basepath, targpath string) (string, error) {
	basepath = filepath.ToSlash(basepath)
	targpath = filepath.ToSlash(targpath)
	if !strings.HasPrefix(targpath, basepath) {
		return "", filepath.ErrBadPattern
	}
	return targpath[len(basepath):], nil
}
