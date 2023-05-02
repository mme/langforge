.DEFAULT_GOAL := all

# Project Variables
VERSION := $(file < VERSION)

# JupyterLab Variables
JL_WHEEL_FILE = python/files/langforge-0.1.0-py3-none-any.whl
JL_IPYNB_FILES = $(wildcard jupyterlab_extension/notebooks/*.ipynb)
JL_TS_FILES = $(wildcard jupyterlab_extension/src/*.ts)
JL_CSS_JS_FILES = $(wildcard jupyterlab_extension/style/*.css) $(wildcard jupyterlab_extension/style/*.js)
JL_JSON_FILES = jupyterlab_extension/package.json jupyterlab_extension/tsconfig.json
JL_WEBPACK_CONFIG = jupyterlab_extension/webpack.config.js

# Go Variables
GO_BUILD_FILE=build/golang/.done
GO_SOURCES = $(shell find cmd environment python system tui -type f)

# Python Variables
PY_SOURCES = $(wildcard pypi/*.py) $(wildcard pypi/langforge/*.py)

# JS Variables
JS_SOURCES = $(wildcard npm/*.js)
JS_PACKAGE_JSON = npm/package.json

all: jupyterlab go compress copy pypi npm archive

$(JL_WHEEL_FILE): $(JL_IPYNB_FILES) $(JL_TS_FILES) $(JL_CSS_JS_FILES) $(JL_JSON_FILES) $(JL_WEBPACK_CONFIG)
	cd jupyterlab_extension && \
	jlpm install && \
	jlpm run build && \
	python -m build && \
	mv dist/langforge-0.1.0-py3-none-any.whl ../python/files/langforge-0.1.0-py3-none-any.whl && \
	rm -rf dist

jupyterlab: $(JL_WHEEL_FILE)

$(GO_BUILD_FILE): $(GO_SOURCES)
	mkdir -p build/golang && \
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/golang/langforge-windows-amd64.exe && \
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/golang/langforge-macos-amd64 && \
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o build/golang/langforge-macos-arm64 && \
	CGO=0 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/golang/langforge-linux-amd64 && \
	CGO=0 CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o build/golang/langforge-linux-arm64 && \
	touch $(GO_BUILD_FILE)

go: $(GO_BUILD_FILE)

compress: $(GO_BUILD_FILE)
	upx --brute build/golang/langforge-windows-amd64.exe && \
	upx --brute build/golang/langforge-macos-amd64 && \
	upx --brute build/golang/langforge-linux-amd64 && \
	upx --brute build/golang/langforge-linux-arm64 && \
	touch $(GO_BUILD_FILE)

copy: $(GO_BUILD_FILE)
	mkdir -p pypi/langforge/bin && \
	mkdir -p npm/bin && \
	cp build/golang/* pypi/langforge/bin && \
	cp build/golang/* npm/bin && \
	touch $(GO_BUILD_FILE)

pypi: $(GO_BUILD_FILE) $(PY_SOURCES)
	cd pypi && \
  python setup.py bdist_wheel && \
	cd .. && touch $(GO_BUILD_FILE)

npm: $(GO_BUILD_FILE) $(JS_SOURCES) $(JS_PACKAGE_JSON)
	cd npm && \
	npm pack && \
	cd .. && touch $(GO_BUILD_FILE)

archive: $(GO_BUILD_FILE)
	rm -rf build/pypi && \
	rm -rf build/npm && \
	cp -r pypi build/pypi && \
	cp -r npm build/npm && \
	touch $(GO_BUILD_FILE)

clean:
	rm -rf build
	rm -rf pypi/build
	rm -rf pypi/dist
	rm -rf pypi/langforge_ai.egg-info
	rm -rf pypi/.venv
	rm -rf python/files/langforge-0.1.0-py3-none-any.whl
	rm -rf jupyterlab_extension/lib
	rm -rf jupyterlab_extension/node_modules
	rm -rf jupyterlab_extension/tsconfig.tsbuildinfo
	rm -rf pypi/langforge/bin
	rm -rf npm/bin
	rm -rf npm/langforge*.tgz

bump:
	$(eval VERSION := $(shell cat VERSION))
	@echo "Current version: $(VERSION)"
	$(eval NEW_VERSION := $(shell echo "$(VERSION)" | awk -F. '{print $$1 "." $$2 "." $$3+1}'))
	@echo "New version: $(NEW_VERSION)"
	@jq '.version = "$(NEW_VERSION)"' npm/package.json > npm/package.json.tmp && mv npm/package.json.tmp npm/package.json
	@echo "$(NEW_VERSION)" > VERSION
	@echo "Updated version in VERSION file"

jupyterlab_dev:
	cd jupyterlab_extension && \
  pip install -ve . &&  \
  jupyter labextension develop --overwrite .

release:
	cd build/npm && npm publish langforge-$$(cat ../../VERSION | tr -d '[:space:]').tgz
	cd build/pypi && python -m twine upload dist/*

.SUFFIXES:

.PHONY: all jupyterlab go compress copy pypi archive npm clean bump jupyterlab_dev release-npm
