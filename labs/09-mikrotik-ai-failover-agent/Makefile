.PHONY: all build-amd64 build-arm64 deploy-lab destroy-lab

CLAB_TOPO := topology.yaml
IMAGE_NAME := mikrotik-ai-failover

all: build-amd64

build-amd64:
	@echo "🔨 Compilando agente para x86_64 (CHR / PC)..."
	docker buildx build --platform linux/amd64 -t $(IMAGE_NAME)-amd64 --load .
	docker save $(IMAGE_NAME)-amd64 > ai-agent-amd64.tar
	@echo "✅ Imagen guardada como ai-agent-amd64.tar"

build-arm64:
	@echo "🔨 Compilando agente para ARM64 (RB5009 / CCR)..."
	docker buildx build --platform linux/arm64 -t $(IMAGE_NAME)-arm64 --load .
	docker save $(IMAGE_NAME)-arm64 > ai-agent-arm64.tar
	@echo "✅ Imagen guardada como ai-agent-arm64.tar"

deploy-lab:
	@echo "🚀 Desplegando topología virtual de MikroTik (Containerlab)..."
	sudo clab deploy -t $(CLAB_TOPO)
	@echo "✅ Laboratorio desplegado. Conectate a router1 para cargar el .tar (vía SFTP o Winbox)."

destroy-lab:
	@echo "🧹 Destruyendo laboratorio virtual..."
	sudo clab destroy -t $(CLAB_TOPO) --cleanup
	@echo "✅ Laboratorio limpiado."
