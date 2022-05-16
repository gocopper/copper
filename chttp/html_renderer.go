package chttp

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gocopper/copper/clogger"

	"github.com/gocopper/copper/cerrors"
)

type (
	// HTMLDir is a directory that can be embedded or found on the host system. It should contain sub-directories
	// and files to support the WriteHTML function in ReaderWriter.
	HTMLDir fs.FS

	// StaticDir represents a directory that holds static resources (JS, CSS, images, etc.)
	StaticDir fs.FS

	// HTMLRenderer provides functionality in rendering templatized HTML along with HTML components
	HTMLRenderer struct {
		htmlDir   HTMLDir
		staticDir StaticDir
	}

	// NewHTMLRendererParams holds the params needed to create HTMLRenderer
	NewHTMLRendererParams struct {
		HTMLDir   HTMLDir
		StaticDir StaticDir
		Config    Config
		Logger    clogger.Logger
	}
)

// NewHTMLRenderer creates a new HTMLRenderer with HTML templates stored in dir and registers the provided HTML
// components
func NewHTMLRenderer(p NewHTMLRendererParams) (*HTMLRenderer, error) {
	hr := HTMLRenderer{
		htmlDir:   p.HTMLDir,
		staticDir: p.StaticDir,
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

func (r *HTMLRenderer) funcMap(req *http.Request) template.FuncMap {
	return template.FuncMap{
		"partial": r.partial(req),
		"assets":  r.assets(req),
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

func (r *HTMLRenderer) assets(req *http.Request) func() (template.HTML, error) {
	return func() (template.HTML, error) {
		// todo: remove vite assets
		if false {
			return r.devAssets(req)
		}

		return r.prodAssets()
	}
}

func (r *HTMLRenderer) devAssets(req *http.Request) (template.HTML, error) {
	const reactRefreshURL = "http://localhost:3000/@react-refresh"
	var out strings.Builder

	reactReq, err := http.NewRequestWithContext(req.Context(), http.MethodGet, reactRefreshURL, nil)
	if err != nil {
		return "", cerrors.New(err, "failed to create request for @react-refresh", nil)
	}

	resp, err := http.DefaultClient.Do(reactReq)
	if err != nil {
		return "", cerrors.New(err, "failed to execute request for @react-refresh", nil)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusOK {
		out.WriteString(`
<script type="module">
  import RefreshRuntime from 'http://localhost:3000/@react-refresh'
  RefreshRuntime.injectIntoGlobalHook(window)
  window.$RefreshReg$ = () => {}
  window.$RefreshSig$ = () => (type) => type
  window.__vite_plugin_react_preamble_installed__ = true
</script>`)
	}

	// note: in dev mode only, the css is not part of the initial page load. since it is loaded async, there is
	// a brief time period where the page has no styles. to avoid this, the following snippet hides the
	// body until the css has been loaded.
	out.WriteString(`
<style type="text/css" id="copper-hide-body">
	body { visibility: hidden; }
</style>
<script type="text/javascript">
	(function() {
		let interval;
	
		function showBodyIfStylesPresent() {
			const styleEls = document.getElementsByTagName('style');
			const copperHideBodyStyleEl = document.getElementById('copper-hide-body');
			
			if (!copperHideBodyStyleEl || styleEls.length === 1) {
				return;
			}
			
			copperHideBodyStyleEl.remove();
			clearInterval(interval);
		}
	
		interval = setInterval(showBodyIfStylesPresent, 100);
	})();
</script>
`)

	out.WriteString(`
<script type="module" src="http://localhost:3000/@vite/client"></script>
<script type="module" src="http://localhost:3000/src/main.js"></script>`)

	// nolint:gosec
	return template.HTML(out.String()), nil
}

func (r *HTMLRenderer) prodAssets() (template.HTML, error) {
	var (
		out      strings.Builder
		manifest struct {
			MainJS struct {
				File string   `json:"file"`
				CSS  []string `json:"css"`
			} `json:"src/main.js"`
		}
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
