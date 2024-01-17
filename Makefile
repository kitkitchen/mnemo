make: 
	go build -o bin/mnemo mnemo.go

run:
	bin/mnemo

publish:
	GOPROXY=proxy.golang.org go list -m github.com/kitkitchen/mnemo@v0.0.1-beta.1

lookup:
	curl https://sum.golang.org/lookup/github.com/kitkitchen/mnemo@v0.0.1-beta.1