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
	{"godot_create_and_attach_script", "Cria um GDScript e anexa ao nó indicado.", schema(map[string][2]string{
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
}
