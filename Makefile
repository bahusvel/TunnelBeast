TB_BIN = $(shell which TunnelBeast)
DEPLOY_NODE = 192.168.1.91

assets:
	go-bindata html

build: assets
	go install

deploy: build
	ssh root@$(DEPLOY_NODE) "killall TunnelBeast" || true
	scp $(TB_BIN) root@$(DEPLOY_NODE):/usr/local/bin/
	scp config.yml root@$(DEPLOY_NODE):./
	ssh root@$(DEPLOY_NODE) "TunnelBeast config.yml"

run: build
	sudo TunnelBeast config.yml
