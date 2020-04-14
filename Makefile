.PHONY: build
build:
	GIT_COMMIT=`git rev-list -1 HEAD` && \
	GIT_DATE=`git log -1 --date=short --pretty=format:%ct` && \
	go build -ldflags "-s -w -X main.gitCommit=$${GIT_COMMIT} -X main.gitDate=$${GIT_DATE}" \
	-o build/lachesis \
	./cmd/lachesis


.PHONY: txstorm
txstorm:
	go build -o build/tx-storm ./cmd/tx-storm


.PHONY: test
test:
	go test ./...


.PHONY: clean
clean:
	rm ./build/lachesis ./build/tx-storm


.PHONY: xxx
xxx:
	echo HERE
	XXX=`echo 111` && echo $${XXX}