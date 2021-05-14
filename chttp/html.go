package chttp

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/cerrors"
)

// HTMLDir is a directory that can be embedded or found on the host system. It should contain sub-directories
// and files to support the WriteHTML function in ReaderWriter.
type HTMLDir fs.FS

// HTMLRenderer provides functionality in rendering templatized HTML along with HTML components
type HTMLRenderer struct {
	htmlDir   HTMLDir
	staticDir StaticDir
	env       cconfig.Env
}

// NewHTMLRendererParams holds the params needed to create HTMLRenderer
type NewHTMLRendererParams struct {
	HTMLDir   HTMLDir
	StaticDir StaticDir
	AppConfig cconfig.Config
}

// NewHTMLRenderer creates a new HTMLRenderer with HTML templates stored in dir and registers the provided HTML
// components
func NewHTMLRenderer(p NewHTMLRendererParams) (*HTMLRenderer, error) {
	var config config

	hr := HTMLRenderer{
		htmlDir:   p.HTMLDir,
		staticDir: p.StaticDir,
		env:       p.AppConfig.Env(),
	}

	err := p.AppConfig.Load("chttp", &config)
	if err != nil {
		return nil, cerrors.New(err, "failed to load chttp config", nil)
	}

	if config.WebDir != "" {
		hr.htmlDir = os.DirFS(config.WebDir)
	}

	return &hr, nil
}

func (r *HTMLRenderer) funcMap(req *http.Request) template.FuncMap {
	return template.FuncMap{
		"partial": r.partial(req),
		"assets":  r.assets,
	}
}

func (r *HTMLRenderer) render(req *http.Request, layout, page string, data interface{}) (template.HTML, error) {
	var dest strings.Builder

	tmpl, err := template.New(layout).
		Funcs(r.funcMap(req)).
		ParseFS(r.htmlDir,
			path.Join("src", "layouts", layout),
			path.Join("src", "pages", page),
		)
	if err != nil {
		return "", cerrors.New(err, "failed to parse templates in html dir", map[string]interface{}{
			"layout": layout,
			"page":   page,
		})
	}

	err = tmpl.Execute(&dest, data)
	if err != nil {
		return "", cerrors.New(err, "failed to execute template", nil)
	}

	// nolint:gosec
	return template.HTML(dest.String()), nil
}

func (r *HTMLRenderer) partial(req *http.Request) func(name string, data interface{}) (template.HTML, error) {
	return func(name string, data interface{}) (template.HTML, error) {
		var dest strings.Builder

		tmpl, err := template.New(name+".html").
			Funcs(r.funcMap(req)).
			ParseFS(r.htmlDir,
				path.Join("src", "partials", "*.html"),
			)
		if err != nil {
			return "", cerrors.New(err, "failed to parse partial template", map[string]interface{}{
				"name": name,
			})
		}

		err = tmpl.Execute(&dest, data)
		if err != nil {
			return "", cerrors.New(err, "failed to execute partial template", map[string]interface{}{
				"name": name,
			})
		}

		// nolint:gosec
		return template.HTML(dest.String()), nil
	}
}

func (r *HTMLRenderer) assets() (template.HTML, error) {
	if r.env == "dev" {
		return `<script type="module" src="http://localhost:3000/@vite/client"></script>
    <script type="module" src="http://localhost:3000/src/main.js"></script>`, nil
	}

	var (
		manifest struct {
			MainJS struct {
				File string   `json:"file"`
				CSS  []string `json:"css"`
			} `json:"src/main.js"`
		}
		out strings.Builder
	)

	manifestFile, err := r.staticDir.Open("static/manifest.json")
	if err != nil {
		return "", cerrors.New(err, "failed to open manifest.json", nil)
	}

	err = json.NewDecoder(manifestFile).Decode(&manifest)
	if err != nil {
		return "", cerrors.New(err, "failed to decode manifest.json", nil)
	}

	if len(manifest.MainJS.CSS) == 1 {
		out.WriteString(fmt.Sprintf("<link rel=\"stylesheet\" href=\"/static/%s\" />\n", manifest.MainJS.CSS[0]))
	}

	out.WriteString(fmt.Sprintf("<script type=\"module\" src=\"/static/%s\"></script>", manifest.MainJS.File))

	//nolint:gosec
	return template.HTML(out.String()), nil
}
