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

It is recommended that each app has its own Fx module. An Fx module can be defined in a package by using `fx.Provide` , and providing all the constructors.
```go
// fx.go
var MyAppFx = fx.Provide(
	NewRouteA,
	NewRouteB,

	NewSvcA,
	NewSvcB,

  ...
)
```

Your app module can be provided during copper app initialization:
``` go
// main.go
func main() {
	copper.NewHTTPApp(
		// required modules
		clogger.StdFx,

		// other modules

		// my app's module
		MyAppFx,
	).Run()
}
```

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

Some modules may require config to be provided. By convention, each module that requires config, defines a `Config` struct at the root package level. For example, http config is defined at `chttp.Config`. The http package also defines, by convention, a `chttp.GetDefaultConfig()` method that returns the config object with sane defaults.

It is recommended that as part of your app, you create a `Config` struct that holds all the configurations for the modules that you want to override. Then, provide it as part of the copper app initialization.

``` go
// config/config.go
type Config struct {
	fx.Out

	Auth cauth.Config
	SQL  csql.Config
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

## Routing
Copper uses [Gorilla Mux](https://github.com/gorilla/mux) under the hood for routing HTTP requests. Combined with Fx, routes can be added to the application container from any module. For example, the `cauth` package registers routes for login, signup, forgot password, and more.

A copper route is defined using the `chttp.Route` struct. 
``` go
route := chttp.Route{
	Path:    "/hello",
	Methods: []string{http.MethodGet},
	MiddlewareFuncs: []chttp.MiddlewareFunc{},
	Handler: http.HandlerFunc(HandleHello),
}
```

The route can be provided to the app by creating a route constructor like so:
``` go
func NewHelloRoute() chttp.RouteResult {
	route := ...
	return chttp.RouteResult{Route: route}
}
```

Paths can also have variables. For example, to add an `id` variable to the path, you can modify the path to `/hello/{id}`. This variable can be accessed in the handler func by calling `mux.Vars(r)`.

A route constructor can be provided in your app’s Fx module. When the app starts, you can see the route getting registered in the logs. For example, the following log shows the health route getting registered by the `chttp` module.
```
2018/12/24 17:51:48 [INFO] Registering route.. where path=/_health,methods=GET
```

## Middleware
Copper allows you to register one or more middleware functions on your route using the `MiddlewareFuncs` field on the route object.

``` go
func NewUserProfileRoute(authMw chttp.AuthMiddleware) chttp.RouteResult {
	route := chttp.Route{
		MiddlewareFuncs: []chttp.MiddlewareFunc{authMw.AllowVerified},
		Path:            "/profile",
		Methods:         []string{http.MethodGet},
		Handler:         http.HandlerFunc(HandleUserProfile),
	}
	return chttp.RouteResult{Route: route}
}
```

The route above uses the auth middleware provided by the `cauth`package. The middleware verifies that the `Authorization` header is valid and that the user’s email has been verified. If any of the checks fail, it sends back `401`. Otherwise, it calls the handler for the route.

Every middleware is of the following type:
``` go 
type MiddlewareFunc func(http.Handler) http.Handler
```

A logger middleware may have the following implementation:
``` go
type LoggerMiddleware struct {
	logger clogger.Logger
}

// NewLoggerMiddleware can be used with fx.Provide to add this middleware to the application container.
func NewLoggerMiddleware(logger cloger.Logger) *LoggerMiddleware {...}

// LogAndNext is a middleware that simply logs each request's path and method, and calls the route handler.
func (m *LoggerMiddleware) LogAndNext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.info("Handling request..", map[string]string{
			"path":   r.URL.Path,
			"method": r.Method,
		})
		next.ServeHTTP(w, r)
	})
}
```

The following modifies our example user profile route to add the logger middleware:
``` go
func NewUserProfileRoute(loggerMw LoggerMiddleware, authMw chttp.AuthMiddleware) chttp.RouteResult {
	route := chttp.Route{
		MiddlewareFuncs: []chttp.MiddlewareFunc{loggerMw.LogAndNext, authMw.AllowVerified},
	...
```

## Reading JSON Body
Copper’s `chttp` package provides `chttp.BodyReader` that can be used to read and validate JSON body. Copper uses [govalidator](https://github.com/asaskevich/govalidator) under-the-hood for validations.

To use `chttp.BodyReader`, inject it into your route handler.
``` go
type RouteHandlers struct {
	req chttp.BodyReader
}

// NewRouterHandlers can be used with fx.Provide to add these route handlers to the application container.
func NewRouterHandlers(req chttp.BodyReader) *RouteHandlers {...}

// NewRouteA creates and configures the route. This route constructor should also be added to the application
// container.
func NewRouteA(handlers *RouteHandlers) chttp.RouteResult {...}

func (ro *RouteHandlers) HandleRouteA(w http.ResponseWriter, r *http.Request) {
	// Define the struct that will hold the values from the request. Use govalidator's struct tags to
	// define validation rules on each field.
	var body struct {
		Email string `valid:"email"`
	}

	// Call the Read method on chttp.BodyReader to read the request's body into the body struct. If the
	// JSON is malformed or if the validation fails, BodyReader responds with a `Bad Request` response
	// and returns false.
	if !ro.req.Read(w, r, &body) {
		return
	}

	// your code goes here..
}
```

## Writing JSON Responses
Copper’s `chttp` package provides `chttp.Responder` that can be used to write JSON responses easily.

To use `chttp.Responder`, inject it into your route handler (check above for example). `Responder` provides various methods such as `OK`, `Created`, `InternalErr` that will set the proper status code, headers, and marshal the provided object as JSON.

For example, the following responds with a `200 OK` and marshals the user object as JSON:
``` go
func (ro *RouteHandlers) HandleRouteA(w http.ResponseWriter, r *http.Request) {
	user := ro.users.GetCurrentUser()
	
	// chttp.Responder has been injected into RouteHandlers
	ro.resp.OK(w, user)
}
```

To override the marshaling of your structs, override the `MarshalJSON() ([]byte, error)` method on the pointer receiver.

## Errors
Copper provides a `cerror` package that has a richer type for errors. `cerrer.Error` can hold a message, tags, and the cause of an error.

To create a new error, use the `cerror.New` method. The `clogger` package recognizes this error type and is helpful when reading logs.
``` go
err := testOp()
if err != nil {
	logger.Error(“Test op has failed miserably”, err)
}
```

The above generates the following log that provides stack trace, contextual tags, and more.
```
2018/12/24 19:31:45 [ERROR] Test op has failed miserably because
> failed to send message to Mars where address=Mars Road, Mars,method=usps because
> failed to connect to satellite
```

## License
Copper is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).
