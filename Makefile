
.PHONY: all serve tailwind browser-sync build

# Default target executed when no arguments are given to make.
all: serve tailwind

# Target for starting the Go application with live reload using air
serve:
	@echo Starting Go server with AIR...
	@air -d

# Target for compiling Tailwind CSS and watching for changes
tailwind:
	@echo Compiling Tailwind CSS...
	@bunx tailwindcss -i './public/tailwind.css' -o './public/styles.css' --watch

# Target for starting browser-sync for live reloading HTML and CSS
browser-sync:
	@bunx browser-sync start --files 'views/**/*.html, public/**/*.css' --port 3001 --proxy 'localhost:3000' --middleware 'function(req, res, next) { res.setHeader("Cache-Control", "no-cache, no-store, must-revalidate"); return next(); }' 

# Target for running tests
test:
	@go	test --cover ./...

# Target for building production assets
build:
	@echo Building production assets...
	@echo Compiling minified Tailwind CSS...
	@bunx tailwindcss -i './public/tailwind.css' -o './public/styles.css' --minify
	@echo Minifying JavaScript...
	@bunx terser ./public/index.js -o ./public/index.js
	@bunx terser ./public/htmx/*.js -o ./public/htmx/*.js
	@echo Generating embedded files...
	@go generate ./...
	@echo Building Go binary...
	@go build -o ./bin/app -race
	@echo Copying static files...
	@cp -r ./public ./bin/
	@rm ./bin/public/tailwind.css
	@echo Done! 🚀

