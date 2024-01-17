make: 
	go build -o bin/mnemo mnemo.go

run:
	bin/mnemo

publish:
	GOPROXY=proxy.golang.org go list -m github.com/kitkitchen/mnemo@v.0.0.2-beta