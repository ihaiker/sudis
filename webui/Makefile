npm=$(shell command -v cnpm)

ifeq ("${npm}","")
  npm=npm
endif

build: dep
	$(npm) run build

clean:
	rm -rf dist

dep:
	$(npm) i
