# crom-godot-mcp — Checklist mestre de ferramentas

> **ATUALIZAÇÃO (2026-07-18): 53 ferramentas implementadas.** O laço de feedback
> está COMPLETO (`verify_playable`, `class_reference`, `gdscript_check`,
> `get_console_errors`, `record_property_over_time`, runtime tree/property). Este
> arquivo é o plano por fases; para a **paridade completa vs godot-mcp-pro (175
> tools)** veja **[CHECKLIST_PARIDADE_GODOT_MCP_PRO.md](CHECKLIST_PARIDADE_GODOT_MCP_PRO.md)**.
> Consolidação: este é o **único** MCP — as tools do antigo `crom-godot-agent/go`
> foram portadas para cá.

Servidor MCP (Go) que expõe o Editor Godot para agentes de IA (crom-agente).
Este documento é o **plano de construção**: cada ferramenta MCP declarada aqui no Go
(`tools.go`) precisa de uma **ação correspondente no plugin** (`command_processor.gd`,
que roda dentro do editor). O Go é só o encaminhador; o trabalho real é GDScript.

## Legenda
- [x] **implementado** e testado no crom-godot-mcp
- [ ] a fazer
- 🟡 parcial (existe mas fraco / precisa melhorar)
- ⭐ **prioridade** — fecha o laço de feedback (o que falta hoje para o agente se auto-corrigir)
- 🆕 **além do godot-mcp-pro** (ideia nossa)

> Referência de escopo: godot-mcp-pro (162 tools, proprietário — **não importável**),
> Coding-Solo/godot-mcp (MIT), IvanMurzak/Godot-MCP (Apache-2.0). Reimplementamos as
> capacidades no nosso próprio servidor, sob nossa licença.

---

## Fase 0 — o que já temos (herdado do bridge atual: 22 tools)

### Cena & Nós
- [x] `godot_get_scene_tree` — árvore de nós da cena aberta
- [x] `godot_add_node` — adicionar nó
- [x] `godot_remove_node` — remover nó
- [x] `godot_set_node_property` — definir propriedade (com coerção Vector2/3, Color)
- [x] `godot_move_node` — mover (posição 2D/3D)
- [x] `godot_rename_node` — renomear
- [x] `godot_reparent_node` — reparentar
- [x] `godot_instantiate_scene` — instanciar cena como filha

### Scripts & Sinais
- [x] `godot_create_and_attach_script` — criar .gd e anexar ao nó
- [x] `godot_connect_signal` — conectar sinal (persistente na cena) 🆕 *adicionado nesta investigação*

### Cenas & Projeto
- [x] `godot_create_scene` — criar .tscn com raiz
- [x] `godot_open_scene` — abrir cena no editor
- [x] `godot_save_scene` — salvar cena
- [x] `godot_set_project_setting` — editar project.godot
- [x] `godot_add_input_action` — criar ação de input mapeada a teclas
- [x] `godot_get_open_editor_context` — scripts abertos, cena, seleção

### Execução & Arquivos
- [x] `godot_play_scene` — executar cena 🟡 *não devolve erros (ver Fase 1)*
- [x] `godot_stop_scene` — parar execução
- [x] `godot_capture_screenshot` — screenshot (1 frame) 🟡 *não detecta movimento*
- [x] `godot_read_project_file` — ler arquivo res://
- [x] `godot_modify_project_file` — escrever arquivo res://
- [x] `godot_list_project_dir` — listar diretório res://

---

## Fase 1 — LAÇO DE FEEDBACK ⭐ (o que destrava tudo)

O agente hoje age mas não vê o resultado. Sem isso, nada mais importa.

- [x] ⭐ `godot_get_console_errors` — lê os `SCRIPT ERROR` / `Parse Error` / `ERROR` recentes do log do Godot (editor + jogo em execução) e devolve ao agente
- [x] ⭐ `godot_get_output` — devolve o conteúdo recente do painel Output (prints, avisos)
- [x] `godot_clear_output` — limpa o log/buffer antes de um novo teste (baseline limpo)
- [ ] ⭐ `godot_play_scene` **retornando erros** — após rodar, coletar e devolver os erros do console (upgrade do atual)
- [x] 🆕 `godot_gdscript_check` — valida a sintaxe de um .gd (parse) ANTES de anexar, sem rodar a cena
- [x] 🆕 ⭐ `godot_verify_playable` — composto: play → coleta erros → (0 erros) → simula input → confirma que o estado mudou → devolve veredito "jogável / não jogável + porquê"

---

## Fase 2 — RUNTIME & INSPEÇÃO DO JOGO EM EXECUÇÃO ⭐

Detectar bugs de gameplay (ex.: "a cobra não se move") que 1 screenshot não pega.

- [x] ⭐ `godot_capture_frames` — captura N frames com intervalo (detecta movimento/animação)
- [x] ⭐ `godot_get_runtime_scene_tree` — árvore de nós do JOGO em execução (não do editor)
- [x] ⭐ `godot_get_runtime_property` — ler propriedade de um nó no jogo rodando (ex.: posição da cobra)
- [x] `godot_record_property_over_time` — grava valores de uma propriedade por N frames (prova de movimento)
- [ ] `godot_get_runtime_nodes_by_type` — achar nós por classe/script no jogo em execução
- [ ] `godot_watch_signal` — registra emissões de um sinal durante a execução

---

## Fase 3 — TESTES & QA ⭐

Transforma "achei que terminei" em "verifiquei que funciona".

- [ ] ⭐ `godot_assert_node_state` — afirma que uma propriedade tem valor esperado (falha = erro claro)
- [ ] `godot_assert_screen_text` — verifica texto na tela (score, "Game Over")
- [ ] `godot_run_test_scenario` — roda uma sequência de input + asserts como um teste
- [ ] `godot_compare_screenshot` — compara com um baseline (regressão visual)
- [ ] 🆕 `godot_smoke_test_scene` — abre a cena principal, roda 2s, retorna erros + se algo se moveu

---

## Fase 4 — INPUT & SIMULAÇÃO

Para "jogar" o jogo e ver se responde.

- [x] ⭐ `godot_simulate_key` — pressiona/solta uma tecla (setas, WASD, espaço, enter)
- [x] `godot_simulate_action` — dispara uma ação do InputMap (ui_accept, jump…)
- [ ] `godot_simulate_mouse` — clique/movimento do mouse em coordenada
- [ ] `godot_input_sequence` — sequência de inputs com frame-delays (ex.: direita, direita, cima)
- [ ] 🟡 `godot_simulate_editor_input` — *existe no plugin, mas fraco; refazer/expor melhor*

---

## Fase 5 — SCRIPTS & CÓDIGO (edição fina)

- [x] `godot_read_script` — ler o script de um nó
- [x] `godot_set_script_source` — reescrever o corpo de um script existente
- [x] `godot_detach_script` — remover script de um nó
- [x] `godot_list_node_methods` — métodos/sinais disponíveis de um nó (ajuda o agente a acertar nomes)
- [x] `godot_list_node_signals` — sinais de um nó (para connect_signal correto)

---

## Fase 6 — NÓS AVANÇADO

- [x] `godot_duplicate_node` — duplicar nó (com filhos)
- [x] `godot_add_to_group` / `godot_remove_from_group` — grupos de nós
- [ ] `godot_set_node_owner` — corrigir owner (evita nós que não salvam)
- [ ] `godot_batch_add_nodes` — adicionar vários nós de uma vez (menos turnos)
- [x] 🆕 `godot_get_node_config_warnings` — avisos de configuração do nó (ex.: CollisionShape2D sem shape)

---

## Fase 7 — RECURSOS & PROJETO

- [x] `godot_create_resource` — criar .tres (ex.: RectangleShape2D, StyleBox)
- [ ] `godot_set_resource_property` — editar um recurso
- [ ] `godot_import_asset` — reimportar/registrar asset
- [x] `godot_get_project_setting` — ler configuração (hoje só escreve)
- [x] `godot_list_input_actions` — listar o InputMap atual

---

## Fase 8 — DOMÍNIOS DE JOGO (2D primeiro)

- [ ] `godot_create_tilemap` / `godot_paint_tiles` — TileMapLayer
- [ ] `godot_create_animation` / `godot_add_animation_track` — AnimationPlayer
- [ ] `godot_setup_physics_body` — CharacterBody2D/RigidBody2D + collision shape num passo
- [ ] `godot_setup_camera` — Camera2D com limites/zoom
- [ ] `godot_add_audio_player` — AudioStreamPlayer + stream

---

## Fase 9 — DOCUMENTAÇÃO & SKILLS 🆕 (evita o erro antes de acontecer)

O modelo deriva para Godot 3 (`update()`, `Color.gray`, `yield`). Duas defesas:

- [x] 🆕 ⭐ **Skill "Godot 4 GDScript"** (arquivo criado em addons/crom_ai/skills/godot4.crom; falta o daemon carregar — ver nota) — arquivo `.crom` (o crom-agente já tem `internal/skills`)
      sempre injetado: migrações 3→4 (`yield`→`await`, `update()`→`queue_redraw()`,
      `.instance()`→`.instantiate()`, `Color.gray`→`Color.GRAY`, `KinematicBody2D`→`CharacterBody2D`),
      "Timer precisa de start()/autostart", "conecte sinais", padrões de nó.
- [x] 🆕 `godot_docs_search` — busca por palavra-chave na doc offline do Godot 4
      (`crom-godot-ai/docs/godot-docs-stable.zip` já existe) → classe/membro/exemplo.
      *Começar com busca estruturada por classe (Timer → start/autostart/wait_time/timeout);
      RAG com embeddings só se a busca simples não bastar.*
- [x] 🆕 `godot_class_reference` — retorna métodos/propriedades/sinais de uma classe Godot

---

## Princípios de arquitetura (para cada ferramenta)
1. **Toda tool no Go tem uma ação no `command_processor.gd`** — senão é stub inútil (pior que não existir).
2. **Não misturar** edição via MCP com `write_file`/`edit_file` no mesmo `.tscn` (causa cenas-lixo).
3. **Toda edição de cena salva** (`save_scene`) para o runtime refletir.
4. **Ferramentas de leitura devem devolver dados úteis ao LLM** (erros exatos, não "ok").

## Ordem de execução recomendada
`Fase 1 (feedback) → Fase 2 (runtime) → Fase 3 (QA) → Fase 4 (input) → Fase 9 (skill Godot 4) → resto`

O par **Fase 1 + Fase 9** (ver o erro + evitar o erro) é o que faz até um modelo médio
entregar jogo jogável. É por onde começamos.
