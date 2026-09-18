#!/usr/bin/env bash
set -euo pipefail

# Script ejecutado por cron cada 2 minutos para crear contenedores de prueba.
# Requiere Docker instalado y autenticado en el host.

CARNET="202300668"
LOG_FILE="/tmp/cron-proyecto2-${CARNET}.log"

log() {
  echo "[$(date -Iseconds)] $*" >> "$LOG_FILE"
}

ensure_one() {
  local name="$1"
  local image="$2"
  local cmd=("${@:3}")

  if ! docker ps --format '{{.Names}}' | grep -Fxq "$name"; then
    log "Creando contenedor $name usando $image"
    docker run -d --name "$name" "$image" "${cmd[@]}" >/dev/null 2>&1 || true
  fi
}

ensure_one "low-${CARNET}-1" "alpine" "sleep" "240"
ensure_one "low-${CARNET}-2" "alpine" "sleep" "240"
ensure_one "low-${CARNET}-3" "alpine" "sleep" "240"
ensure_one "high-${CARNET}-1" "roldyoran/go-client" ""
ensure_one "high-${CARNET}-2" "roldyoran/go-client" ""

log "Estado final de contenedores:"
docker ps --format 'table {{.Names}}\t{{.Image}}' >> "$LOG_FILE" || true
