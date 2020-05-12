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

.PHONY: grafana
grafana:
	cd ../ && \
	./docker/grafana/configs-generate.sh && \
	docker build \
	    --network==host \
	    --build-arg GOPROXY=$(GOPROXY) \
	    -f docker/Dockerfile.grafana -t "grafana-lachesis:${TAG}" .

.PHONY: prometheus
prometheus:
	cd ../ && \
	./docker/prometheus/configs-generate.sh && \
	docker build \
	    --network==host \
	    --build-arg GOPROXY=$(GOPROXY) \
	    -f docker/Dockerfile.prometheus -t "prometheus-lachesis:${TAG}" .

.PHONY: alertmanager
alertmanager:
	cd ../ && \
	./docker/alertmanager/configs-generate.sh && \
	docker build \
	    --network==host \
	    --build-arg GOPROXY=$(GOPROXY) \
	    -f docker/Dockerfile.alertmanager -t "alertmanager-lachesis:${TAG}" .

test_accs.json:
	./_test_accs.sh
