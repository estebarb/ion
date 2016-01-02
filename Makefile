get-deps:
	go get github.com/gorilla/mux
	go get github.com/justinas/alice
	go get github.com/gorilla/sessions
fmt:
	go fmt ./context
	go fmt ./examples/context
	go fmt ./examples/hello
	go fmt ./examples/ionvc
	go fmt ./examples/restful
	go fmt ./examples/session
	go fmt ./examples/template
	go fmt ./session
	go fmt ./middleware
	go fmt .
test:
	go test -v .
	go test -v ./context
	go test -v ./ionvc
	go test -v ./session
	go test -v ./examples/context/context.go
	go test -v ./examples/hello/hello.go
	go test -v ./examples/restful/restful.go
	go test -v ./examples/session/session.go
	go test -v ./examples/template/template.go
