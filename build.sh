#!/usr/bin/env bash
# ==============================================================================
# build.sh — FONTE ÚNICA do binário godot-mcp.
# Compila o crom-godot-mcp para todos os alvos e (opcional) implanta no addon
# do app Godot. Este é o único lugar de onde o godot-mcp deve ser buildado.
# ==============================================================================
set -euo pipefail
cd "$(dirname "$0")"

# Destino padrão: o bin/ do addon CromAI (ajuste via 1º argumento).
DEST="${1:-/home/j/Documentos/GitHub/crom-godot-ai/addons/crom_ai/bin}"

targets=(
  "linux-amd64:linux:amd64"
  "linux-arm64:linux:arm64"
  "darwin-amd64:darwin:amd64"
  "darwin-arm64:darwin:arm64"
  "windows-amd64.exe:windows:amd64"
)

echo "Compilando godot-mcp (crom-godot-mcp) -> ${DEST}"
mkdir -p "${DEST}"
for t in "${targets[@]}"; do
  name="${t%%:*}"; rest="${t#*:}"; os="${rest%%:*}"; arch="${rest##*:}"
  GOOS="${os}" GOARCH="${arch}" go build -ldflags="-s -w" -o "${DEST}/godot-mcp-${name}" .
  echo "  ✓ godot-mcp-${name}"
done

count=$(grep -oE '"godot_[a-z_]+"' tools.go | sort -u | wc -l)
echo "Pronto. ${count} ferramentas expostas."

# --- Sync do lado Godot (FONTE ÚNICA) -----------------------------------------
# O crom-godot-mcp é dono do addon Godot. Apps que o consomem sincronizam os .gd
# core daqui (não editam cópias). Ajuste APP_ADDON via 2º argumento.
APP_ADDON="${2:-/home/j/Documentos/GitHub/crom-godot-ai/addons/crom_ai}"
if [[ -d "${APP_ADDON}" ]]; then
  echo "Sincronizando lado Godot (godot-addon -> ${APP_ADDON})"
  for f in command_processor.gd crom_runtime.gd websocket_server.gd; do
    cp "godot-addon/${f}" "${APP_ADDON}/${f}"
    echo "  ✓ ${f}"
  done
  echo "Sync do addon Godot concluído."
fi
