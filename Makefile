make: 
	go build -o bin/mnemo mnemo.go

run:
	bin/mnemo

publish:
	GOPROXY=proxy.golang.org go list -m github.com/kitkitchen/mnemo@v0.0.1-beta.1