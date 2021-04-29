package chttp

import (
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/cerrors"
)

// HTMLDir is a directory that can be embedded or found on the host system. It should contain sub-directories
// and files to support the WriteHTML function in ReaderWriter.
type HTMLDir fs.FS

// HTMLComponent represents a HTML partial with a corresponding Go struct that can render the partial with additional
// logic
type HTMLComponent interface {
	Name() string
}

// HTMLRenderer provides functionality in rendering templatized HTML along with HTML components
type HTMLRenderer struct {
	dir        HTMLDir
	components []HTMLComponent
}

// NewHTMLRenderer creates a new HTMLRenderer with HTML templates stored in dir and registers the provided HTML
// components
func NewHTMLRenderer(dir HTMLDir, components []HTMLComponent, appConfig cconfig.Config) (*HTMLRenderer, error) {
	var config config

	hr := HTMLRenderer{
		dir:        dir,
		components: components,
	}

	err := appConfig.Load("chttp", &config)
	if err != nil {
		return nil, cerrors.New(err, "failed to load chttp config", nil)
	}

	if config.WebDir != "" {
		hr.dir = os.DirFS(config.WebDir)
	}

	return &hr, nil
}

func (r *HTMLRenderer) render(req *http.Request, layout, page string, data interface{}) (template.HTML, error) {
	var dest strings.Builder

	tmpl, err := template.New(layout).
		Funcs(template.FuncMap{
			"component": r.component(req),
		}).
		ParseFS(r.dir,
			path.Join("html", "layouts", layout),
			path.Join("html", "pages", page),
			path.Join("html", "partials", "*.html"),
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

// nolint: funlen,lll
func (r *HTMLRenderer) component(req *http.Request) func(componentName string, props ...interface{}) (template.HTML, error) {
	return func(componentName string, props ...interface{}) (template.HTML, error) {
		var (
			component HTMLComponent
			dest      strings.Builder
		)

		for i := range r.components {
			if r.components[i].Name() == componentName {
				component = r.components[i]
				break
			}
		}

		if component == nil {
			return "", cerrors.New(nil, "invalid component name", map[string]interface{}{
				"name": componentName,
			})
		}

		renderMethod := reflect.ValueOf(component).MethodByName("Render")
		if !renderMethod.IsValid() {
			return "nil", cerrors.New(errors.New("no Render method"), "invalid component", map[string]interface{}{
				"component": componentName,
			})
		}

		renderMethodArgs := []reflect.Value{
			reflect.ValueOf(req),
		}

		for i := range props {
			renderMethodArgs = append(renderMethodArgs, reflect.ValueOf(props[i]))
		}

		renderResult := renderMethod.Call(renderMethodArgs)

		if !renderResult[1].IsNil() {
			return "", cerrors.New(renderResult[1].Interface().(error), "failed to render component", map[string]interface{}{
				"component": componentName,
			})
		}

		tmpl, err := template.New(component.Name()+".html").
			Funcs(template.FuncMap{
				"component": r.component(req),
			}).
			ParseFS(r.dir,
				path.Join("html", "components", "*.html"),
			)
		if err != nil {
			return "", cerrors.New(err, "failed to parse component template", map[string]interface{}{
				"name": component.Name(),
			})
		}

		err = tmpl.Execute(&dest, renderResult[0].Interface())
		if err != nil {
			return "", cerrors.New(err, "failed to execute component template", map[string]interface{}{
				"name": component.Name(),
			})
		}

		// nolint:gosec
		return template.HTML(dest.String()), nil
	}
}
