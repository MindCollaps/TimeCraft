package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"src/main/middleware"
	"src/main/web/webLogic"
	"strings"
	"text/template"
)

var tmpl = map[string]*template.Template{}

func LoadTemplates(r *gin.Engine) {
	tmpl = make(map[string]*template.Template)

	// Load templates files
	templateFiles := []string{}

	fmt.Println("Loading templates...")
	// Walk through the "templates" folder and all its subdirectories
	nerr := filepath.Walk("main/web/templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is an HTML templates
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".gohtml") {
			// Replace backslashes with forward slashes (for Windows compatibility)
			templateName := strings.Replace(path, "\\", "/", -1)

			// Parse the file and add it to the "tmpl" map
			templateFiles = append(templateFiles, path)

			//console log
			fmt.Print(templateName + " ")
		}
		return nil
	})

	if nerr != nil {
		panic(nerr)
	}

	fmt.Println("\n\nLoading sites...")
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.LoginToken())
	adminGroup.Use(middleware.VerifyAdmin())
	devGroup := r.Group("/dev")
	devGroup.Use(middleware.LoginToken())

	// Walk through the "public" folder and all its subdirectories
	err := filepath.Walk("main/web/public", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is an HTML template
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".gohtml") {
			// Get the directory path (relative to the "public" folder)
			relPath, err := filepath.Rel("main/web/public", filepath.Dir(path))
			if err != nil {
				return err
			}
			// Replace backslashes with forward slashes (for Windows compatibility)
			templateName := strings.Replace(relPath, "\\", "/", -1)

			if strings.HasSuffix(path, "app.gohtml") {
				templateName += "app"
			}

			// Parse the file and add it to the "tmpl" map
			parsing := []string{}
			parsing = append(parsing, templateFiles...)
			parsing = append(parsing, path)

			tmpl[templateName] = template.Must(template.ParseFiles(parsing...))

			if strings.HasPrefix(templateName, "admin") {
				fmt.Println("Serving " + relPath + " as admin at /" + templateName)
				themPath := strings.Replace(templateName, "admin", "", 1)
				adminGroup.GET("/"+themPath, handler)
			} else if strings.HasPrefix(templateName, "dev") {
				fmt.Println("Serving " + relPath + " as dev at /" + templateName)
				themPath := strings.Replace(templateName, "dev", "", 1)
				devGroup.GET("/"+themPath, handler)
			} else {
				fmt.Println("Serving " + relPath + " at /" + templateName)
				r.GET("/"+templateName, handler)
			}
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
		fmt.Println("Error executing template:", err)
		return
	}
}

func LoadServerAssets(r *gin.Engine) {
	fmt.Println("Loading assets...")
	// Walk through the "assets" folder and all its subdirectories
	err := filepath.Walk("main/web/assets", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is not a directory
		if !info.IsDir() {
			// Get the directory path (relative to the "public" folder)
			relPath, err := filepath.Rel("main/web/assets", path)
			if err != nil {
				return err
			}

			assetPath := strings.Replace(relPath, "\\", "/", -1)
			// Add the asset to a route
			fmt.Println("Serving " + path + " at /assets/" + assetPath)
			r.StaticFile("/assets/"+assetPath, path)

			if err != nil {
				return err
			}
		}

		return nil
	})

	r.StaticFile("/favicon.ico", "main/web/static/favicon.ico")

	if err != nil {
		panic(err)
	}
}
