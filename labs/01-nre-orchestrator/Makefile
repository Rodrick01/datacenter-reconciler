.PHONY: all build deploy destroy logs

# Variables
CLAB_TOPO := topology.yaml
DOCKER_COMPOSE := docker-compose.yml

all: deploy

# Construir las imágenes locales si es necesario
build:
	docker build -t nre-agents:latest .

# Despliegue completo: Levantar NetBox/DBs -> Desplegar Fabric -> Levantar Agentes
deploy: build
	@echo "🚀 Iniciando despliegue Zero-Touch del Laboratorio NRE..."
	docker compose -f $(DOCKER_COMPOSE) up -d postgres redis redis-cache netbox
	@echo "⏳ Esperando a que NetBox inicialice..."
	sleep 5
	sudo clab deploy -t $(CLAB_TOPO)
	@echo "🤖 Levantando agentes de AI y Reconciliación..."
	docker compose -f $(DOCKER_COMPOSE) up -d aico reconciler
	@echo "✅ Despliegue completado con éxito."

# Destrucción limpia del laboratorio
destroy:
	@echo "🧹 Destruyendo laboratorio..."
	sudo clab destroy -t $(CLAB_TOPO) --cleanup
	docker compose -f $(DOCKER_COMPOSE) down -v
	@echo "✅ Laboratorio limpiado."

logs:
	docker compose -f $(DOCKER_COMPOSE) logs -f aico reconciler
