package chttp

import (
	_ "embed"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/clogger"
)

type (
	// HTMLDir is a directory that can be embedded or found on the host system. It should contain sub-directories
	// and files to support the WriteHTML function in ReaderWriter.
	HTMLDir fs.FS

	// StaticDir represents a directory that holds static resources (JS, CSS, images, etc.)
	StaticDir fs.FS

	// HTMLRenderer provides functionality in rendering templatized HTML along with HTML components
	HTMLRenderer struct {
		htmlDir     HTMLDir
		staticDir   StaticDir
		renderFuncs []HTMLRenderFunc
	}

	// HTMLRenderFunc can be used to register new template functions
	HTMLRenderFunc struct {
		// Name for the function that can be invoked in a template
		Name string

		// Func should return a function that takes in any number of params and returns either a single return value,
		// or two return values of which the second has type error.
		Func func(r *http.Request, html *HTMLRenderer) interface{}
	}

	// NewHTMLRendererParams holds the params needed to create HTMLRenderer
	NewHTMLRendererParams struct {
		HTMLDir     HTMLDir
		StaticDir   StaticDir
		RenderFuncs []HTMLRenderFunc
		Config      Config
		Logger      clogger.Logger
	}
)

// NewHTMLRenderer creates a new HTMLRenderer with HTML templates stored in dir and registers the provided HTML
// components
func NewHTMLRenderer(p NewHTMLRendererParams) (*HTMLRenderer, error) {
	hr := HTMLRenderer{
		htmlDir:     p.HTMLDir,
		staticDir:   p.StaticDir,
		renderFuncs: p.RenderFuncs,
	}

	if p.Config.UseLocalHTML {
		wd, err := os.Getwd()
		if err != nil {
			return nil, cerrors.New(err, "failed to get current working directory", nil)
		}

		hr.htmlDir = os.DirFS(filepath.Join(wd, "web"))
	}

	return &hr, nil
}

func (r *HTMLRenderer) Render(req *http.Request, path string, data map[string]interface{}) (template.HTML, error) {
	return r.renderPartialFromDirWithFuncs(path, r.funcMap(req), data)
}

func (r *HTMLRenderer) funcMap(req *http.Request) template.FuncMap {
	var funcMap = template.FuncMap{
		"partial": func(req *http.Request, html *HTMLRenderer) interface{} {
			return func(name string, data interface{}) (template.HTML, error) {
				return r.renderPartialFromDirWithFuncs(path.Join("partials", name), r.funcMap(req), data)
			}
		},
		"dict": dict,
	}

	for i := range r.renderFuncs {
		funcMap[r.renderFuncs[i].Name] = r.renderFuncs[i].Func(req, r)
	}

	return funcMap
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

func (r *HTMLRenderer) renderPartialFromDirWithFuncs(tmplPath string, fnMap template.FuncMap, data interface{}) (template.HTML, error) {
	var (
		name = path.Base(tmplPath) + ".html"
		dir  = path.Dir(tmplPath)

		dest strings.Builder
	)

	tmpl, err := template.New(name).
		Funcs(fnMap).
		ParseFS(r.htmlDir,
			path.Join("src", dir, "*.html"),
		)
	if err != nil {
		return "", cerrors.New(err, "failed to parse partial template", map[string]interface{}{
			"dir":  dir,
			"name": name,
		})
	}

	err = tmpl.Execute(&dest, data)
	if err != nil {
		return "", cerrors.New(err, "failed to execute partial template", map[string]interface{}{
			"dir":  dir,
			"name": name,
		})
	}

	// nolint:gosec
	return template.HTML(dest.String()), nil
}
