.PHONY: build webui bindata release clean

binout=bin/sudis

Version=$(shell git describe --tags `git rev-list --tags --max-count=1`)
BuildDate=$(shell date +"%F %T")
GitCommit=$(shell git rev-parse --short HEAD)
debug=-w -s
param=-X main.VERSION=${Version} -X main.GITLOG_VERSION=${GitCommit} -X 'main.BUILD_TIME=${BuildDate}'

gobinddata=$(shell command -v go-bindata)

ifeq ($(gobinddata),'')
	go get -u github.com/shuLhan/go-bindata/cmd/go-bindata
endif

build: bindata
	go mod download
	go build -tags bindata -ldflags "${debug} ${param}" -o ${binout}

docker:
	docker build --build-arg LDFLAGS="${debug} ${param}" -t xhaiker/sudis:${Version} .

bindata:
	go generate generator.go

webui:
	make -C webui build

clean:
	@rm -rf bin
	@rm -rf webui/dist
	@rm -rf webui/node_modules
	@rm -f nodes/http/http_static_bindata_assets.go


