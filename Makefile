LDFLAGS=-w -s

all:
	make mac
	make linux
mac:
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/notifyme_darwin64
	GOOS=darwin GOARCH=386 go build -ldflags "$(LDFLAGS)" -o bin/notifyme_darwin32
linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/notifyme_linux64
	GOOS=linux GOARCH=386 go build -ldflags "$(LDFLAGS)" -o bin/notifyme_linux32
