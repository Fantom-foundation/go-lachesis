all: lachesis tx-storm

TAG ?= "latest"

.PHONY: lachesis
lachesis: test_accs.json
	cd ../ && \
	rm -Rf ./build && \
	mkdir ./build && \
	cp -R ../go-ethereum ./build/ && \
	docker build \
	    --network=host \
	    --build-arg GOPROXY=$(GOPROXY) \
	    -f docker/Dockerfile.lachesis -t "lachesis:$(TAG)" .


.PHONY: tx-storm
tx-storm:
	cd ../ && \
	rm -Rf ./build && \
	mkdir ./build && \
	cp -R ../go-ethereum ./build/ && \
	docker build \
	    --network=host \
	    --build-arg GOPROXY=$(GOPROXY) \
	    -f docker/Dockerfile.tx-storm -t "tx-storm:$(TAG)" .

test_accs.json:
	./_test_accs.sh
