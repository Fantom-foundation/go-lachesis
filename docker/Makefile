all: lachesis tx-storm grafana

TAG ?= "latest"

.PHONY: lachesis
lachesis: test_accs.json
	cd ../ && \
	docker build \
	    --network=host \
	    --build-arg GOPROXY=$(GOPROXY) \
	    -f docker/Dockerfile.lachesis -t "lachesis:$(TAG)" .


.PHONY: tx-storm
tx-storm:
	cd ../ && \
	docker build \
	    --network=host \
	    --build-arg GOPROXY=$(GOPROXY) \
	    -f docker/Dockerfile.tx-storm -t "tx-storm:$(TAG)" .

grafana:
	./grafana/configs-generate.sh && \
	docker build \
		--network==host \
		-f Dockerfile.grafana -t "grafana-lachesis:${TAG}"

test_accs.json:
	./_test_accs.sh
