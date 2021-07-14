build:
	go build -ldflags "-s -w"
	GOOS=windows go build -ldflags "-s -w -H=windowsgui"
