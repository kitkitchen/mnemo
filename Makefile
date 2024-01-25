make: 
	go build -o bin/mnemo mnemo.go

run:
	bin/mnemo

tag:
	git tag -af v0.0.1-beta.5 -m "mnemo v0.0.1-beta.5"

publish-proxy:
	GOPROXY=proxy.golang.org go list -m github.com/kitkitchen/mnemo@v0.0.1-beta.5

lookup:
	curl https://sum.golang.org/lookup/github.com/kitkitchen/mnemo@v0.0.1-beta.5