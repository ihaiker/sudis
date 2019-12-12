binout=bin/sudis

ifeq ($(P),release)
Version=$(shell git describe --tags `git rev-list --tags --max-count=1`)
BuildDate=$(shell date +"%F %T")
GitCommit=$(shell git rev-parse HEAD)
debug=-w -s
param=-X main.VERSION=${Version} -X main.GITLOG_VERSION=${GitCommit} -X 'main.BUILD_TIME=${BuildDate}'
else
debug=
param=
endif

gobinddata=$(shell command -v go-bindata)

ifeq ($(gobinddata),'')
	go get -u github.com/shuLhan/go-bindata/cmd/go-bindata
endif

build: bindata
	go mod download
	go build -tags bindata -ldflags "${debug} ${param}" -o ${binout}

bindata: webui
	go generate generator.go

webui:
	make -C webui build

install: clean build
	@mkdir -p bin/conf
	@cp conf/sudis.toml.example bin/conf/sudis.toml

release:
	make -C . -e P=release

clean:
	@rm -rf bin
	@rm -rf webui/dist
	@rm -rf webui/node_modules
	@rm -f master/server/http/http_static_bindata_assets.go

.PHONY: build webui bindata install release clean
