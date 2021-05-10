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

// HTMLComponent represents a HTML partial with a corresponding Go struct that can render the partial with additional
// logic
type HTMLComponent interface {
	Name() string
}

// HTMLComponentSubscriber can be implemented by HTML components to subscribe to events broadcasted by other components
// in the ComponentTree. When a subscribed event is published, the component's render method is called.
type HTMLComponentSubscriber interface {
	Subscribe() []string
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

func (r *HTMLRenderer) funcMap(req *http.Request) template.FuncMap {
	return template.FuncMap{
		"component": r.component(req, nil),
		"partial":   r.partial(req),
	}
}

func (r *HTMLRenderer) render(req *http.Request, layout, page string, data interface{}) (template.HTML, error) {
	var dest strings.Builder

	tmpl, err := template.New(layout).
		Funcs(r.funcMap(req)).
		ParseFS(r.dir,
			path.Join("html", "layouts", layout),
			path.Join("html", "pages", page),
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

// nolint: lll
func (r *HTMLRenderer) callComponentMethod(req *http.Request, componentID, componentName, methodName string, propValues, argValues []json.RawMessage) (template.HTML, error) {
	var component HTMLComponent

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

	reflector, err := newHTMLComponentReflector(component)
	if err != nil {
		return "", cerrors.New(err, "failed to create html component reflector", map[string]interface{}{
			"component": component.Name(),
		})
	}

	switch methodName {
	case "$_Refresh":
		props, err := reflector.createPropsInterfaceValuesFromJSONValues(propValues)
		if err != nil {
			return "", cerrors.New(err, "failed to create props struct", map[string]interface{}{
				"component": componentName,
			})
		}

		return r.component(req, &componentID)(componentName, props...)
	default:
		props, err := reflector.createPropsStructFromJSONValues(propValues)
		if err != nil {
			return "", cerrors.New(err, "failed to create props struct", map[string]interface{}{
				"component": componentName,
			})
		}

		args, err := reflector.createActionMethodArgs(methodName, argValues)
		if err != nil {
			return "", cerrors.New(err, "failed to create action args", map[string]interface{}{
				"component": componentName,
				"action":    methodName,
			})
		}

		newProps, err := reflector.callActionMethod(req, methodName, props, args)
		if err != nil {
			return "", cerrors.New(err, "failed to call action method", map[string]interface{}{
				"component": componentName,
				"action":    methodName,
			})
		}

		return r.component(req, &componentID)(componentName, newProps...)
	}
}

func (r *HTMLRenderer) partial(req *http.Request) func(name string, data interface{}) (template.HTML, error) {
	return func(name string, data interface{}) (template.HTML, error) {
		var dest strings.Builder

		tmpl, err := template.New(name+".html").
			Funcs(r.funcMap(req)).
			ParseFS(r.dir,
				path.Join("html", "partials", "*.html"),
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

// nolint: funlen,lll
func (r *HTMLRenderer) component(req *http.Request, id *string) func(componentName string, propValues ...interface{}) (template.HTML, error) {
	return func(componentName string, propValues ...interface{}) (template.HTML, error) {
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

		reflector, err := newHTMLComponentReflector(component)
		if err != nil {
			return "", cerrors.New(err, "failed to create html component reflector", map[string]interface{}{
				"component": component.Name(),
			})
		}

		props := reflector.createPropsStructFromInterfaceValues(propValues)

		renderResult, err := reflector.callRenderMethod(req, props)
		if err != nil {
			return "", cerrors.New(err, "failed to render component", map[string]interface{}{
				"component": componentName,
			})
		}

		tmpl, err := template.New(component.Name()+".html").
			Funcs(r.funcMap(req)).
			ParseFS(r.dir,
				path.Join("html", "components", "*.html"),
			)
		if err != nil {
			return "", cerrors.New(err, "failed to parse component template", map[string]interface{}{
				"name": component.Name(),
			})
		}

		err = tmpl.Execute(&dest, map[string]interface{}{
			"Props": props,
			"View":  renderResult,
		})
		if err != nil {
			return "", cerrors.New(err, "failed to execute component template", map[string]interface{}{
				"name": component.Name(),
			})
		}

		propsJSON, err := json.Marshal(propValues)
		if err != nil {
			return "", cerrors.New(err, "failed to marshal component props as json", map[string]interface{}{
				"name": component.Name(),
			})
		}

		events := make([]string, 0)

		eventSubscriber, ok := component.(HTMLComponentSubscriber)
		if ok {
			events = eventSubscriber.Subscribe()
		}

		eventsJSON, err := json.Marshal(events)
		if err != nil {
			return "", cerrors.New(err, "failed to marshal component events as json", map[string]interface{}{
				"name":   component.Name(),
				"events": events,
			})
		}

		wrappedComponent := fmt.Sprintf(`<copper-component name="%s" props='%s' events='%s'>%s</copper-component>`,
			componentName,
			string(propsJSON),
			string(eventsJSON),
			dest.String(),
		)

		// nolint:gosec
		return template.HTML(wrappedComponent), nil
	}
}
