<p align="center">
  <a href="https://gocopper.dev" target="_blank" rel="noopener noreferrer">
    <img width="180" src="https://gocopper.dev/static/logo.svg" alt="Copper logo">
  </a>
</p>

<p align="center">
    <a href="https://goreportcard.com/report/github.com/gocopper/copper" target="_blank" rel="noopener noreferrer"> 
        <img src="https://goreportcard.com/badge/github.com/gocopper/copper" alt="Go Report Card">
    </a>
    <a href="https://pkg.go.dev/github.com/gocopper/copper"  target="_blank" rel="noopener noreferrer">
        <img src="https://pkg.go.dev/badge/github.com/gocopper/copper?status.svg" alt="Go Doc">
    </a>
</p>

<br />

# Copper

<p>
Copper is a Go toolkit complete with everything you need to build web apps. It focuses on developer productivity and makes building web apps in Go more fun with less boilerplate and out-of-the-box support for common needs.
</p>

#### 🚀 Fullstack Toolkit
<p>Copper provides a toolkit complete with everything you need to build web apps quickly.</p>


#### 📦 One Binary
<p>Build frontend apps along with your backend and ship everything in a single binary.</p>


#### 📝 Server-side HTML
<p>Copper includes utilities that help build web apps with server rendered HTML pages.</p>

#### 💡 Auto Restarts
<p>Copper detects changes and automatically restarts server to save time.</p>

#### 🏗 Scaffolding
<p>Skip boilerplate and scaffold code for your packages, database queries and routes.</p>

#### 🔋 Batteries Included
<p>Includes CLI, lint, dev server, config management, and more!</p>

#### 🔩 First-party packages
<p>Includes packages for authentication, pub/sub, queues, emails, and websockets.</p>


<br />

## Current Status

💎️ Copper is currently in preview as new features are added. While the APIs are unlikely to change in any major ways, some details may change as they are refined. Feedback and contributions are welcome! 

<br />

## Getting Started

> Copper requires Go 1.16+

<br  />

1. Install the Copper CLI
```
❯ go install github.com/gocopper/cli/cmd/copper@latest
```

2. Install Wire CLI
```
❯ go install github.com/google/wire/cmd/wire@latest
```

3. Scaffold your project
```
❯ copper init
? What's the module name for your project? xnotes

# Create Project Files

 SUCCESS  Create xnotes/config/local.toml (Took 1ms)
 SUCCESS  Create xnotes/go.mod (Took 0s)
 SUCCESS  Create xnotes/pkg/app/handler.go (Took 0s)
 SUCCESS  Create xnotes/pkg/app/wire.go (Took 0s)
 SUCCESS  Create xnotes/pkg/web/public/favicon.svg (Took 1ms)
 SUCCESS  Create xnotes/pkg/web/public/logo.svg (Took 0s)
 SUCCESS  Create xnotes/pkg/web/src/pages/index.html (Took 0s)
 SUCCESS  Create xnotes/config/prod.toml (Took 0s)
 SUCCESS  Create xnotes/cmd/migrate/wire.go (Took 0s)
 SUCCESS  Create xnotes/pkg/app/migrations.go (Took 0s)
 SUCCESS  Create xnotes/pkg/web/wire.go (Took 0s)
 SUCCESS  Create xnotes/.golangci.yaml (Took 0s)
 SUCCESS  Create xnotes/cmd/app/main.go (Took 1ms)
 SUCCESS  Create xnotes/cmd/app/wire.go (Took 0s)
 SUCCESS  Create xnotes/cmd/migrate/main.go (Took 0s)
 SUCCESS  Create xnotes/config/dev.toml (Took 0s)
 SUCCESS  Create xnotes/pkg/web/src/layouts/main.html (Took 0s)
 SUCCESS  Create xnotes/pkg/web/src/main.js (Took 0s)
 SUCCESS  Create xnotes/pkg/web/vite.config.js (Took 0s)
 SUCCESS  Create xnotes/.gitignore (Took 0s)
 SUCCESS  Create xnotes/pkg/web/package.json (Took 0s)
 SUCCESS  Create xnotes/pkg/web/router.go (Took 0s)
 SUCCESS  Create xnotes/pkg/web/src/styles.css (Took 0s)
 SUCCESS  Create xnotes/config/base.toml (Took 0s)

# First Commands

┌──────────────┐
| cd xnotes    |
| copper build |
| copper watch |
└──────────────┘
```

4. Start Project
```cgo
❯ cd xnotes
❯ copper build
❯ copper watch
```

5. Open http://localhost:5901 in your browser

<br />

## License
MIT
