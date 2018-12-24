# Copper - Web Framework
## Introduction
Copper is a dependency injection based web framework written in Go (GoLang). It provides a solid foundation that is easily extensible. You get out-of-the-box support for routing, database, validations, user management, and more!

``` go
package main

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/clogger"
)

func main() {
	copper.NewHTTPApp(
		clogger.StdFx,
	).Run()
}
```

The above code gives you an HTTP server, standardized logging, health endpoint, and graceful shutdown. You can continue to extend it by adding your own modules or choose from the provided modules.

## Status
Copper is currently in alpha stage. Feel free to try it out and open pull requests.

## Quick Start
The following snippet shows a ‘Hello World’ copper app. Even then, it does show some basic features of the framework.

``` go
// NewHelloWorldRoute defines a single http route. It configures the route with path, methods,
// middlewares, and the handler.
// It takes in `chttp.Responder` as a dependency that is injected via Fx. `chttp.Responder` provides various
// methods to respond with JSON easily.
// While the handler is defined inline here, in a real world app, you may want to define the handler
// separately.
func NewHelloWorldRoute(resp chttp.Responder) chttp.RouteResult {
	route := chttp.Route{
		Path:    "/hello",
		Methods: []string{http.MethodGet},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp.OK(w, map[string]string{
				"hello": "world",
			})
		}),
	}
	return chttp.RouteResult{Route: route}
}

// MyAppFx defines an Fx module for your app. All the route constructors are provided here.
var MyAppFx = fx.Provide(
	NewHelloWorldRoute,
)

func main() {
	// Start the copper app by providing your app's Fx module.
	copper.NewHTTPApp(
		clogger.StdFx,

		MyAppFx,
	).Run()
}
```

``` bash
> curl localhost:7450/hello
{"hello":"world"}

> curl localhost:7450/_health
OK
```

## Dependency Injection
Copper depends heavily on dependency injection to create a working application. To facilitate this, copper uses Uber’s DI framework: [FX](https://go.uber.org/fx). Head over to the Fx documentation to familiarize yourself with the basic concepts.

## Modules
Copper provides various modules out-of-the-box. You can add these modules to the copper app during initialization. For example,  the following snippet adds a `mailer` module that can be used to send emails using various implementations such as AWS.

``` go
func main() {
	copper.NewHTTPApp(
		...
		cmailer.AWSFx,
	).Run()
}
```

By adding `cmailer.AWSFx` to the app container, you’re providing an implementation for the `cmailer.Mailer` interface. As a result, in any of your modules, you can inject `cmailer.Mailer` and Fx will provide the AWS implementation for it. This paradigm makes it trivial to add new implementations for the `Mailer` interface as well as providing your own custom ones.

Included modules:
* **logger** (github.com/tusharsoni/copper/clogger) provides various logging methods that improves debugging experience. It is required to start a copper app.
* **http** (github.com/tusharsoni/copper/chttp) provides the primitives to start a http application. It supports routing, middleware, utilities to read and write JSON body, and more.
* **sql** (github.com/tusharsoni/copper/csql) is an Fx wrapper around a popular Go GORM, [GORM](http://github.com/jinzhu/gorm). 
* **mailer** (github.com/tusharsoni/copper/cmailer) provides various implementations for sending emails.
* **auth** (github.com/tusharsoni/copper/cauth) adds authentication to your copper app. It works with a Postgres database and provides authentication flows for signup, login, reset/change password, and verification emails.

## Configurations
In copper, configurations are handled as part of your Go code. They are defined as simple Go structs and work perfectly with the dependency injection paradigm.

Some modules may require config to be provided. By convention, each module that requires config, defines a `Config` object at the root package level. For example, http config is defined at `chttp.Config`. The http package also defines, by convention, a `chttp.GetDefaultConfig()` method that returns the config object with sane defaults.

It is recommended that as part of your app, you create a `Config` struct that holds all the configurations for the modules that you want to override. Then, provide it as part of the copper app initialization.

``` go
// config/config.go
type Config struct {
	fx.Out

	Auth cauth.Config
	SQL csql.Config
}

func NewConfig() Config {
	return Config{
		AWSMailer: cmailer.AWSConfig{
			Region:          “us-east-1”,
			AccessKeyId:     “ABCDEFJHIJ123456789”,
			SecretAccessKey: os.Getenv(“AWSSecretAccessKey”),
		},
		SQL: csql.Config{
			Host: “localhost”,
			Port: 5432,
			Name: “myapp”,
			User: “postgres”,
		},
	}
}

// config/fx.go
var Fx = fx.Provide(NewConfig)

// main.go
func main() {
	copper.NewHTTPApp(
		clogger.StdFx,
		csql.Fx,
		cmailer.AWSFx,
		
		config.Fx,
	).Run()
}

```

## License
Copper is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).
