# crom-godot-mcp — Godot Addon (lado Godot)

Este é o **lado Godot** do crom-godot-mcp: o plugin do Editor que **executa** as
ferramentas `godot_*`. O servidor Go (`../` — `godot-mcp`) só encaminha JSON-RPC;
o trabalho real (criar nós, física, animação, rodar/verificar cena) acontece aqui.

## Conteúdo
- `command_processor.gd` — handlers de TODAS as ferramentas `godot_*` (fonte única).
- `crom_runtime.gd` — autoload que roda dentro do jogo (porta 8091) p/ inspeção em runtime.
- `websocket_server.gd` — servidor WS na porta 8080 que recebe as ações do godot-mcp.
- `plugin.gd` / `plugin.cfg` — plugin mínimo do editor (sobe o WS + registra o autoload).

## Instalar num projeto Godot
1. Copie esta pasta para `res://addons/crom_mcp/` (ou qualquer nome — o plugin é
   auto-localizável).
2. Ative o plugin em **Projeto → Ajustes → Plugins**.
3. Rode o servidor Go (`godot-mcp --mcp-stdio --port 8080`) e conecte seu agente/MCP.

## Fonte única
Este addon é a **fonte única** do lado Godot do MCP. Apps que consomem o CromAI
(ex.: crom-godot-ai) **sincronizam** estes `.gd` daqui — não editam cópias.
