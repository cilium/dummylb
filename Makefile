.PHONY: clean

clean:
	rm -f main

compile: dummylb

dummylb: main.go
	CGO_ENABLED=0 go build -o dummylb main.go
