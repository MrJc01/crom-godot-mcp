# crom-godot-mcp

Servidor **MCP (Model Context Protocol)** em Go que expõe o **Editor Godot 4** para
agentes de IA (como o [crom-agente](https://github.com/MrJc01/crom-agente)).

O agente ganha ferramentas para **criar cenas, editar nós, anexar scripts, conectar
sinais, rodar o jogo e ler os erros do console** — operando de dentro do editor, na
cena viva, e não só escrevendo arquivos no escuro.

## Arquitetura

```
Agente (crom-agente)  ──stdio/JSON-RPC──►  crom-godot-mcp  ──WebSocket :8080──►  Plugin CromAI no Editor Godot
```

O Go é um **encaminhador fino**: cada ferramenta `godot_*` vira uma ação JSON enviada
ao plugin `addons/crom_ai` (GDScript), que executa a operação real no editor e devolve
o resultado. O "poder" está no plugin; este servidor traduz MCP ↔ WebSocket.

## Uso

```bash
go build -o crom-godot-mcp .
./crom-godot-mcp --mcp-stdio --port 8080
```

Registro no crom-agente (`~/.crom/global.json`):

```json
{
  "mcp_servers": [
    {
      "name": "godot-editor",
      "command": "/caminho/para/crom-godot-mcp",
      "args": ["--mcp-stdio", "--port", "8080"],
      "env": ["GODOT_PROJECT_DIR=/caminho/do/projeto"]
    }
  ]
}
```

Requer o projeto Godot **aberto no editor com o plugin CromAI ativo** (que sobe o
servidor WebSocket na porta 8080).

## Estado

v0.1 — 22 ferramentas base (cena, nós, scripts, sinais, projeto, execução).
O roteiro completo (laço de feedback, runtime, testes, input, docs) está em
[`CHECKLIST.md`](CHECKLIST.md).

## Licença

MIT — veja [`LICENSE`](LICENSE).
