default:
	@export GOPATH=$$GOPATH:$$(pwd) && go install runner
edit:
	@export GOPATH=$$GOPATH:$$(pwd) && atom .
edit2:
	@export GOPATH=$$GOPATH:$$(pwd) && code .
run: default
	@rm -rf dumps
	@mkdir dumps
	@bin/runner
	@echo ""
clean:
	@rm -rf bin
test:
	@export GOPATH=$$GOPATH:$$(pwd) && go test ./...
test_ver:
	@export GOPATH=$$GOPATH:$$(pwd) && go test -v ./...
setup:
	@if [ ! -d $$GOPATH/pkg/linux_amd64/gopkg.in/mgo.v2 ]; then go get gopkg.in/mgo.v2; else echo "mgo installed."; fi
	@if [ ! -d $$GOPATH/pkg/linux_amd64/github.com/aws/ ]; then go get -u github.com/aws/aws-sdk-go/...; else echo "aws installed."; fi
	@if [ ! -d $$GOPATH/pkg/linux_amd64/github.com/mxk/go-sqlite/ ]; then go get github.com/mxk/go-sqlite/sqlite3; else echo "sqlite installed."; fi
	@if [ ! -d $$GOPATH/pkg/linux_amd64/golang.org/x/crypto/ ]; then go get golang.org/x/crypto/ssh; else echo "ssl installed."; fi
	@if [ ! -f $$GOPATH/pkg/linux_amd64/github.com/bitly/go-simplejson.a ]; then go get github.com/bitly/go-simplejson; else echo "json installed."; fi
