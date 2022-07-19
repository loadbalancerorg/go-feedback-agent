.PHONEY: clean get

VERSION=`git describe --tags`
BUILD=`git rev-parse HEAD`
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"
LIGHT=light
CANDLE=candle


default: build

build: windows
windows:
	 env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -v -o ./bin/windows64/LBCPUMon.exe ./src
linux:
	 env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -v -o ./bin/linux64/lbcpumon ./src
FeedbackAgentService.wixobj: FeedbackAgentService.wxs
	$(CANDLE) FeedbackAgentService.wxs

FeedbackAgent.msm: FeedbackAgentService.wixobj bin/windows64/LBCPUMon.exe LICENSE
	$(LIGHT) FeedbackAgent.wixobj
get:
	go mod download
clean:
	go clean -modcache
