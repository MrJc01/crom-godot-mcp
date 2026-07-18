package main

import "strings"

// toolDef descreve uma ferramenta MCP exposta ao agente.
type toolDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// schema monta um JSON Schema compacto. props: nome -> [tipo, descrição].
func schema(props map[string][2]string, required ...string) map[string]interface{} {
	p := map[string]interface{}{}
	for name, td := range props {
		entry := map[string]interface{}{"description": td[1]}
		if td[0] == "array" {
			entry["type"] = "array"
			entry["items"] = map[string]interface{}{}
		} else if strings.Contains(td[0], "|") {
			entry["type"] = strings.Split(td[0], "|")
		} else {
			entry["type"] = td[0]
		}
		p[name] = entry
	}
	s := map[string]interface{}{"type": "object", "properties": p}
	if len(required) > 0 {
		s["required"] = required
	}
	return s
}

// catalog é o inventário de ferramentas. FASE 0: as 22 herdadas (todas com ação
// correspondente no command_processor.gd do plugin). Novas fases são adicionadas
// aqui + a ação no plugin (ver CHECKLIST.md).
var catalog = []toolDef{
	// --- Cena & Nós ---
	{"godot_get_scene_tree", "Lê a árvore de nós da cena aberta no Editor (nomes, tipos, caminhos, filhos).", schema(nil)},
	{"godot_get_open_editor_context", "Contexto do Editor: scripts abertos, cena em edição e nós selecionados.", schema(nil)},
	{"godot_add_node", "Adiciona um nó novo na cena aberta.", schema(map[string][2]string{
		"node_type":   {"string", "Classe Godot (ex: Sprite2D, CharacterBody2D, Timer, Label)"},
		"node_name":   {"string", "Nome do novo nó"},
		"parent_path": {"string", "Caminho do pai ('.' para a raiz)"},
		"properties":  {"object", "Propriedades iniciais (ex: {\"position\":[100,200]})"},
	}, "node_type", "node_name")},
	{"godot_remove_node", "Remove um nó da cena aberta.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó"},
	}, "node_path")},
	{"godot_set_node_property", "Altera uma propriedade de um nó (position, scale, text, modulate, visible, autostart, wait_time...).", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó"},
		"property":  {"string", "Nome da propriedade"},
		"value":     {"string|number|boolean|array|object", "Valor. Vetores [x,y]/[x,y,z]; cores [r,g,b,a]"},
	}, "node_path", "property", "value")},
	{"godot_move_node", "Move um nó 2D/3D/Control para uma posição.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó"},
		"position":  {"array", "[x,y] para 2D/Control ou [x,y,z] para 3D"},
	}, "node_path", "position")},
	{"godot_rename_node", "Renomeia um nó.", schema(map[string][2]string{
		"node_path": {"string", "Caminho atual"},
		"new_name":  {"string", "Novo nome"},
	}, "node_path", "new_name")},
	{"godot_reparent_node", "Move um nó para debaixo de outro pai (reparenting).", schema(map[string][2]string{
		"node_path":       {"string", "Caminho do nó"},
		"new_parent_path": {"string", "Caminho do novo pai ('.' para a raiz)"},
	}, "node_path", "new_parent_path")},
	{"godot_connect_signal", "Conecta um sinal de um nó a um método de outro nó, salvando na cena. USE sempre que criar um handler _on_<no>_<sinal> (ex: timeout de um Timer).", schema(map[string][2]string{
		"from_node": {"string", "Nó emissor ('.' para a raiz)"},
		"signal":    {"string", "Sinal (ex: timeout, pressed, body_entered)"},
		"to_node":   {"string", "Nó receptor ('.' para a raiz)"},
		"method":    {"string", "Método a chamar (ex: _on_timer_timeout)"},
	}, "signal", "method")},

	// --- Scripts & Cenas ---
	{"godot_create_and_attach_script", "Cria um GDScript e anexa ao nó indicado. ESCREVA GODOT 4 (nunca Godot 3): use queue_redraw() (não update()); Color.GRAY/Color.RED em MAIÚSCULO (não Color.gray); await (não yield); packed.instantiate() (não .instance()); CharacterBody2D (não KinematicBody2D). Um Timer só funciona com autostart=true OU start(), e o sinal timeout precisa ser conectado (godot_connect_signal). Valide com godot_gdscript_check antes.", schema(map[string][2]string{
		"node_path":     {"string", "Nó que recebe o script ('.' para a raiz)"},
		"script_path":   {"string", "Caminho res:// do .gd"},
		"gdscript_code": {"string", "Código GDScript 4 completo"},
	}, "script_path", "gdscript_code")},
	{"godot_create_scene", "Cria um .tscn novo com um nó raiz.", schema(map[string][2]string{
		"scene_path": {"string", "Caminho res:// (ex: res://main.tscn)"},
		"root_type":  {"string", "Classe do nó raiz (padrão: Node2D)"},
		"root_name":  {"string", "Nome do nó raiz"},
	}, "scene_path")},
	{"godot_instantiate_scene", "Instancia uma cena existente como filha de um nó.", schema(map[string][2]string{
		"scene_path":  {"string", "Cena res:// a instanciar"},
		"parent_path": {"string", "Nó pai ('.' para a raiz)"},
		"node_name":   {"string", "Nome opcional"},
	}, "scene_path")},
	{"godot_save_scene", "Salva a cena aberta no disco.", schema(nil)},
	{"godot_open_scene", "Abre uma cena .tscn no Editor.", schema(map[string][2]string{
		"scene_path": {"string", "Cena res://"},
	}, "scene_path")},

	// --- Projeto & Arquivos ---
	{"godot_set_project_setting", "Define uma configuração no project.godot.", schema(map[string][2]string{
		"setting": {"string", "Ex: display/window/size/viewport_width"},
		"value":   {"string|number|boolean|array|object", "Valor"},
	}, "setting", "value")},
	{"godot_add_input_action", "Cria uma ação de input mapeada a teclas.", schema(map[string][2]string{
		"action_name": {"string", "Ex: jump, move_left"},
		"keys":        {"array", "Teclas (ex: [\"Space\"], [\"W\",\"Up\"])"},
	}, "action_name", "keys")},
	{"godot_read_project_file", "Lê um arquivo de texto do projeto (res://).", schema(map[string][2]string{
		"file_path": {"string", "Caminho res://"},
	}, "file_path")},
	{"godot_modify_project_file", "Escreve um arquivo de texto do projeto (res://). NÃO use em .tscn quando estiver editando a cena por ferramentas godot_* (causa conflito).", schema(map[string][2]string{
		"file_path":   {"string", "Caminho res://"},
		"new_content": {"string", "Conteúdo novo"},
	}, "file_path", "new_content")},
	{"godot_list_project_dir", "Lista arquivos e subpastas de um diretório res://.", schema(map[string][2]string{
		"dir_path": {"string", "Diretório res:// (padrão res://)"},
	})},

	// --- Execução & Visão ---
	{"godot_play_scene", "Executa uma cena (ou a principal) para testar o jogo.", schema(map[string][2]string{
		"scene_path": {"string", "Cena res:// (opcional)"},
	})},
	{"godot_stop_scene", "Para a execução da cena em teste.", schema(nil)},
	{"godot_capture_screenshot", "Captura um screenshot e salva como PNG, retornando o caminho.", schema(nil)},

	// --- FASE 1: laço de feedback (ver os próprios erros) ---
	{"godot_get_console_errors", "Retorna os erros recentes do console do Godot (SCRIPT ERROR, Parse Error, ERROR) — do editor e do jogo executado por play_scene. Chame DEPOIS de play_scene para verificar se o jogo roda sem erro.", schema(nil)},
	{"godot_get_output", "Retorna as últimas linhas do painel Output (prints, avisos).", schema(map[string][2]string{
		"lines": {"number", "Quantas linhas (padrão 60)"},
	})},
	{"godot_clear_output", "Marca o ponto atual do console como baseline: get_console_errors/get_output passam a olhar só o que vier depois. Chame antes de um novo teste.", schema(nil)},
	{"godot_gdscript_check", "Valida a sintaxe de um GDScript SEM rodar a cena (parse). Use antes de anexar para pegar erros de sintaxe cedo.", schema(map[string][2]string{
		"script_path":   {"string", "Caminho res:// de um .gd existente (ou use gdscript_code)"},
		"gdscript_code": {"string", "Código a validar diretamente"},
	})},

	// --- FASE 5/6: inspeção de scripts e nós (o agente acerta os nomes) ---
	{"godot_read_script", "Lê o código-fonte do script anexado a um nó.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó ('.' para a raiz)"},
	}, "node_path")},
	{"godot_list_node_methods", "Lista os métodos disponíveis de um nó (ajuda a chamar/conectar com nomes corretos).", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó ('.' para a raiz)"},
	}, "node_path")},
	{"godot_list_node_signals", "Lista os sinais de um nó (para usar connect_signal corretamente).", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó ('.' para a raiz)"},
	}, "node_path")},
	{"godot_get_node_config_warnings", "Retorna avisos de configuração de um nó (ex: CollisionShape2D sem shape, Timer mal configurado).", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó ('.' para a raiz)"},
	}, "node_path")},
	{"godot_duplicate_node", "Duplica um nó (com filhos) na cena aberta.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó a duplicar"},
		"new_name":  {"string", "Nome do duplicado (opcional)"},
	}, "node_path")},
	{"godot_add_to_group", "Adiciona um nó a um grupo (persistente na cena).", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó"},
		"group":     {"string", "Nome do grupo (ex: enemies, food)"},
	}, "node_path", "group")},
	{"godot_remove_from_group", "Remove um nó de um grupo.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó"},
		"group":     {"string", "Nome do grupo"},
	}, "node_path", "group")},

	// --- FASE 7: recursos e projeto (leitura) ---
	{"godot_get_project_setting", "Lê uma configuração do project.godot.", schema(map[string][2]string{
		"setting": {"string", "Ex: application/run/main_scene"},
	}, "setting")},
	{"godot_list_input_actions", "Lista as ações de input definidas no projeto (InputMap).", schema(nil)},
	{"godot_create_resource", "Cria (e opcionalmente salva) um recurso .tres (ex: RectangleShape2D, CircleShape2D, StyleBoxFlat).", schema(map[string][2]string{
		"resource_type": {"string", "Classe do Resource (ex: RectangleShape2D)"},
		"save_path":     {"string", "Caminho res:// para salvar (.tres) — opcional"},
		"properties":    {"object", "Propriedades iniciais (ex: {\"size\":[32,32]})"},
	}, "resource_type")},

	// --- FASE 4: simulação de input (testar jogabilidade) ---
	{"godot_simulate_key", "Simula pressionar/soltar uma tecla física (para testar o jogo em execução).", schema(map[string][2]string{
		"key":     {"string", "Tecla (ex: Up, Down, W, Space, Enter)"},
		"pressed": {"boolean", "true = pressiona, false = solta (padrão true)"},
	}, "key")},
	{"godot_simulate_action", "Simula uma ação do InputMap (ex: ui_accept, jump).", schema(map[string][2]string{
		"action":  {"string", "Nome da ação"},
		"pressed": {"boolean", "true = press, false = release (padrão true)"},
	}, "action")},

	// --- FASE 2: inspeção do JOGO EM EXECUÇÃO (verificar gameplay de verdade) ---
	{"godot_get_runtime_scene_tree", "Lê a árvore de nós do JOGO EM EXECUÇÃO (não do editor) — requer play_scene rodando. Mostra a cena viva com posições.", schema(nil)},
	{"godot_get_runtime_property", "Lê o valor de uma propriedade de um nó no JOGO EM EXECUÇÃO (ex: a posição da cobra). Requer play_scene rodando.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó na cena em execução ('.' para a raiz)"},
		"property":  {"string", "Propriedade (ex: position, visible, text)"},
	}, "node_path", "property")},
	{"godot_record_property_over_time", "Amostra uma propriedade de um nó várias vezes ao longo do tempo e diz se ela MUDOU — é como você confirma que 'a cobra se move' (gameplay funciona), não só que não deu erro. Requer play_scene rodando.", schema(map[string][2]string{
		"node_path":   {"string", "Caminho do nó ('.' para a raiz)"},
		"property":    {"string", "Propriedade a observar (padrão: position)"},
		"samples":     {"number", "Quantas leituras (2-20, padrão 5)"},
		"interval_ms": {"number", "Intervalo entre leituras em ms (padrão 250)"},
	}, "node_path")},

	{"godot_class_reference", "DOCUMENTAÇÃO VIVA: consulta a API autoritativa da versão do Godot em uso (via ClassDB) para uma classe — métodos com assinatura, propriedades, sinais e constantes. Use ANTES de escrever código para uma classe que não domina (ex.: CharacterBody2D, Timer, TileMap) e garanta a sintaxe Godot 4 correta. Se a classe não existir, devolve nomes parecidos.", schema(map[string][2]string{
		"class_name": {"string", "Nome da classe Godot (ex.: CharacterBody2D). Parcial faz busca."},
	}, "class_name")},

	{"godot_verify_playable", "FECHA O FEEDBACK LOOP num único passo: roda a cena, espera o boot, checa erros de console E detecta se algo se move em runtime, devolvendo um veredito 'jogável' (playable=true/false). Use como verificação final antes de dizer que terminou — só finalize com playable=true.", schema(map[string][2]string{
		"scene_path":    {"string", "Cena a rodar (vazio = cena principal do projeto)"},
		"node_path":     {"string", "Nó para checar movimento ('.' = cena atual, padrão)"},
		"property":      {"string", "Propriedade observada para movimento (padrão: position)"},
		"check_movement": {"boolean", "Se deve verificar movimento em runtime (padrão true)"},
		"boot_wait_ms":  {"number", "Espera de boot em ms antes de checar (500-8000, padrão 2000)"},
	})},
}
