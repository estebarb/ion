get-deps:
	go get github.com/gorilla/mux
	go get github.com/justinas/alice
	go get github.com/gorilla/sessions
fmt:
	go fmt ./context
	go fmt ./examples
	go fmt ./session
	go fmt .