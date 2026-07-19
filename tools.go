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
	{"godot_add_node", "Adiciona um nó novo na cena aberta. As 'properties' podem incluir RECURSOS INLINE: para dar um shape a um CollisionShape2D numa única chamada, use {\"shape\":{\"__resource_type\":\"RectangleShape2D\",\"size\":[32,32]}}. Para textura de Sprite2D use o caminho res:// direto: {\"texture\":\"res://icon.svg\"}. Assim você monta um nó completo (com colisão/visual) sem chamadas extras.", schema(map[string][2]string{
		"node_type":   {"string", "Classe Godot (ex: Sprite2D, CharacterBody2D, Timer, Label)"},
		"node_name":   {"string", "Nome do novo nó"},
		"parent_path": {"string", "Caminho do pai ('.' para a raiz)"},
		"properties":  {"object", "Propriedades iniciais. Vetores [x,y]; cores [r,g,b,a]; recurso inline {\"__resource_type\":\"CircleShape2D\",\"radius\":16}; recurso/textura no disco por caminho res://."},
	}, "node_type", "node_name")},
	{"godot_add_nodes_batch", "Cria VÁRIOS nós de uma vez (uma subárvore inteira) — prefira isto a chamar add_node repetidamente: menos idas-e-vindas, menos erros. Ordene os itens com pais antes dos filhos. Cada item aceita node_type, node_name, parent_path e properties (incl. recursos inline).", schema(map[string][2]string{
		"nodes": {"array", "Array de {node_type, node_name, parent_path, properties}. Ex: [{\"node_type\":\"CharacterBody2D\",\"node_name\":\"Player\",\"parent_path\":\".\"},{\"node_type\":\"CollisionShape2D\",\"node_name\":\"Col\",\"parent_path\":\"Player\",\"properties\":{\"shape\":{\"__resource_type\":\"RectangleShape2D\",\"size\":[32,32]}}}]"},
	}, "nodes")},
	{"godot_remove_node", "Remove um nó da cena aberta.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó"},
	}, "node_path")},
	{"godot_set_node_property", "Altera uma propriedade de um nó (position, scale, text, modulate, visible, autostart, wait_time...). Para propriedades que são RECURSOS: passe um objeto inline {\"__resource_type\":\"RectangleShape2D\",\"size\":[32,32]} (ex: property='shape' de um CollisionShape2D) OU um caminho res:// de um recurso/textura no disco (ex: property='texture', value='res://player.png').", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó"},
		"property":  {"string", "Nome da propriedade (ex: shape, texture, position, text)"},
		"value":     {"string|number|boolean|array|object", "Valor. Vetores [x,y]/[x,y,z]; cores [r,g,b,a]; recurso inline {\"__resource_type\":...}; ou caminho res://"},
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
	{"godot_set_main_scene", "Define a cena principal do projeto (a que roda no F5 e no export). Um jogo PRECISA disso — chame depois de criar a cena principal.", schema(map[string][2]string{
		"scene_path": {"string", "Cena res:// que será a principal"},
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

	// --- Scripts (editar sem recriar) ---
	{"godot_set_script_source", "Edita o código-fonte de um script GDScript JÁ existente (anexado a um nó ou no disco) sem recriar o nó. Use node_path para editar o script de um nó, ou script_path para editar direto pelo caminho res://.", schema(map[string][2]string{
		"node_path":     {"string", "Caminho do nó cujo script será editado ('.' para a raiz). Opcional se script_path for dado."},
		"script_path":   {"string", "Caminho res:// do arquivo .gd a editar. Opcional se node_path for dado."},
		"gdscript_code": {"string", "Código GDScript 4 completo novo"},
	}, "gdscript_code")},
	{"godot_detach_script", "Remove o script de um nó da cena sem excluir o arquivo .gd do disco.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó ('.' para a raiz da cena)"},
	}, "node_path")},

	// --- TileMap ---
	{"godot_set_tilemap_cell", "Define uma célula num TileMapLayer ou TileMap da cena aberta.", schema(map[string][2]string{
		"node_path":        {"string", "Caminho do nó TileMap/TileMapLayer na cena"},
		"coords":           {"array", "Coordenadas da célula [x, y]"},
		"source_id":        {"number", "ID da fonte de tile no TileSet (padrão: 0)"},
		"atlas_coords":     {"array", "Coordenadas no atlas [x, y] (padrão: [0,0])"},
		"alternative_tile": {"number", "ID do tile alternativo (padrão: 0)"},
		"layer":            {"number", "Layer do TileMap legado (só TileMap, padrão: 0)"},
	}, "node_path", "coords")},
	{"godot_get_tilemap_cells", "Retorna as células usadas de um TileMapLayer ou TileMap.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó TileMap/TileMapLayer"},
		"layer":     {"number", "Layer do TileMap legado (só TileMap, padrão: 0)"},
	}, "node_path")},

	// --- Animação ---
	{"godot_list_animations", "Lista as animações de um AnimationPlayer ou os sprite_frames de um AnimatedSprite2D.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó AnimationPlayer ou AnimatedSprite2D"},
	}, "node_path")},
	{"godot_play_animation", "Toca uma animação por nome num AnimationPlayer ou AnimatedSprite2D.", schema(map[string][2]string{
		"node_path":      {"string", "Caminho do nó AnimationPlayer ou AnimatedSprite2D"},
		"animation_name": {"string", "Nome da animação a tocar"},
	}, "node_path", "animation_name")},

	// --- Câmera ---
	{"godot_set_camera_target", "Configura posição, zoom e limites de uma Camera2D da cena.", schema(map[string][2]string{
		"node_path":         {"string", "Caminho do nó Camera2D"},
		"position":          {"array", "Posição [x, y]"},
		"zoom":              {"number|array", "Zoom: número escalar ou [x, y]"},
		"limit_left":        {"number", "Limite esquerdo (pixels)"},
		"limit_top":         {"number", "Limite superior (pixels)"},
		"limit_right":       {"number", "Limite direito (pixels)"},
		"limit_bottom":      {"number", "Limite inferior (pixels)"},
		"smoothing_enabled": {"boolean", "Ativa/desativa suavização de posição"},
	}, "node_path")},

	// --- Documentação offline ---
	{"godot_docs_search", "Busca textual na documentação offline do Godot. Complementa godot_class_reference para tutoriais, guias e conceitos.", schema(map[string][2]string{
		"query":       {"string", "Termo de busca (ex: 'TileMap', 'move_and_slide', 'signals')"},
		"max_results": {"number", "Número máximo de resultados (padrão: 5, máx: 20)"},
	}, "query")},

	// --- Runtime & QA (jogo em execução — requer play_scene rodando) ---
	{"godot_set_game_node_property", "Altera uma propriedade de um nó no JOGO EM EXECUÇÃO (não no editor). Útil para testar reações (ex: forçar position, health). Requer play_scene rodando.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó na cena em execução ('.' para a raiz)"},
		"property":  {"string", "Nome da propriedade"},
		"value":     {"string|number|boolean|array", "Valor. Vetores [x,y]; cores [r,g,b,a]"},
	}, "node_path", "property", "value")},
	{"godot_get_game_node_properties", "Retorna TODAS as propriedades editáveis de um nó no JOGO EM EXECUÇÃO. Requer play_scene rodando.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó na cena em execução ('.' para a raiz)"},
	}, "node_path")},
	{"godot_wait_for_node", "Espera (com timeout) um nó aparecer no JOGO EM EXECUÇÃO. Use antes de asserts sobre nós criados dinamicamente. Requer play_scene rodando.", schema(map[string][2]string{
		"node_path":  {"string", "Caminho do nó a esperar"},
		"timeout_ms": {"number", "Timeout em ms (200-10000, padrão 3000)"},
	}, "node_path")},
	{"godot_assert_node_state", "QA: afirma que uma propriedade de um nó no JOGO EM EXECUÇÃO tem o valor esperado. Devolve passed=true/false. Requer play_scene rodando.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó ('.' para a raiz)"},
		"property":  {"string", "Propriedade a checar (ex: position, text, visible)"},
		"expected":  {"string|number|boolean", "Valor esperado (comparação por igualdade/substring)"},
	}, "node_path", "property", "expected")},
	{"godot_assert_screen_text", "QA: afirma que um texto aparece na tela do jogo (em Labels/Buttons visíveis). Ex: confirmar 'Score: 10' ou 'Game Over'. Requer play_scene rodando.", schema(map[string][2]string{
		"text": {"string", "Texto (ou trecho) que deve aparecer na tela"},
	}, "text")},
	{"godot_find_ui_elements", "Lista nós do JOGO EM EXECUÇÃO por tipo (ex: Button, Label, Control) com caminho — para localizar UI antes de interagir. Requer play_scene rodando.", schema(map[string][2]string{
		"type": {"string", "Classe a filtrar (ex: Button, Label). Vazio = todos"},
	})},
	{"godot_click_button_by_text", "Aciona (emite 'pressed') um Button do JOGO EM EXECUÇÃO pelo seu texto — para testar menus/UI. Requer play_scene rodando.", schema(map[string][2]string{
		"text": {"string", "Texto do botão a clicar"},
	}, "text")},

	// --- Nós & Cena (inspeção/edição fina no editor) ---
	{"godot_get_node_properties", "Retorna TODAS as propriedades editáveis de um nó na cena aberta no editor (nome → valor). Use para saber o estado atual antes de alterar.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó ('.' para a raiz)"},
	}, "node_path")},
	{"godot_disconnect_signal", "Desconecta um sinal de um nó (inverso de connect_signal), removendo a ligação da cena.", schema(map[string][2]string{
		"from_node": {"string", "Nó emissor ('.' para a raiz)"},
		"signal":    {"string", "Nome do sinal"},
		"to_node":   {"string", "Nó receptor ('.' para a raiz)"},
		"method":    {"string", "Método conectado"},
	}, "signal", "method")},
	{"godot_find_nodes_in_group", "Lista os nós da cena aberta que pertencem a um grupo (ex: 'enemies', 'coins').", schema(map[string][2]string{
		"group": {"string", "Nome do grupo"},
	}, "group")},
	{"godot_get_node_groups", "Retorna os grupos aos quais um nó pertence.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó ('.' para a raiz)"},
	}, "node_path")},
	{"godot_delete_scene", "Apaga um arquivo de cena .tscn do disco.", schema(map[string][2]string{
		"scene_path": {"string", "Caminho res:// da cena a apagar"},
	}, "scene_path")},
	{"godot_get_project_info", "Retorna informações do projeto: nome, cena principal, versão do Godot e features.", schema(nil)},
	{"godot_search_files", "Busca arquivos sob res:// por nome (e opcionalmente conteúdo). Útil para achar scripts/cenas.", schema(map[string][2]string{
		"query":          {"string", "Termo a buscar no nome do arquivo"},
		"search_content": {"boolean", "Se true, também busca no conteúdo de .gd/.tscn/.tres (padrão false)"},
	}, "query")},

	// --- Física 2D ---
	{"godot_setup_physics_body", "Cria um corpo físico 2D (CharacterBody2D/RigidBody2D/StaticBody2D/Area2D) JÁ com um CollisionShape2D + shape, num único passo. Evita nó de física sem colisão.", schema(map[string][2]string{
		"parent_path": {"string", "Nó pai ('.' para a raiz)"},
		"node_name":   {"string", "Nome do corpo"},
		"body_type":   {"string", "CharacterBody2D (padrão), RigidBody2D, StaticBody2D, Area2D"},
		"shape_type":  {"string", "RectangleShape2D (padrão) ou CircleShape2D"},
		"size":        {"array", "Tamanho [x,y] p/ RectangleShape2D"},
		"radius":      {"number", "Raio p/ CircleShape2D"},
		"position":    {"array", "Posição [x,y] do corpo (opcional)"},
	}, "parent_path")},
	{"godot_set_physics_layers", "Define collision_layer e collision_mask de um corpo físico (bitmask).", schema(map[string][2]string{
		"node_path": {"string", "Caminho do corpo físico"},
		"layer":     {"number", "collision_layer (bitmask)"},
		"mask":      {"number", "collision_mask (bitmask)"},
	}, "node_path")},
	{"godot_get_physics_layers", "Lê collision_layer e collision_mask de um corpo físico.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do corpo físico"},
	}, "node_path")},
	{"godot_add_raycast", "Adiciona um RayCast2D a um nó pai (para detecção de colisão/linha de visão).", schema(map[string][2]string{
		"parent_path":     {"string", "Nó pai ('.' para a raiz)"},
		"node_name":       {"string", "Nome do RayCast2D"},
		"target_position": {"array", "Direção/alcance [x,y] (relativo ao nó)"},
		"enabled":         {"boolean", "Ativo (padrão true)"},
	}, "parent_path")},

	// --- Animação (AnimationPlayer) ---
	{"godot_create_animation", "Cria uma Animation vazia num AnimationPlayer (biblioteca padrão). Depois use add_animation_track + set_animation_keyframe.", schema(map[string][2]string{
		"node_path":      {"string", "Caminho do AnimationPlayer"},
		"animation_name": {"string", "Nome da nova animação"},
		"length":         {"number", "Duração em segundos (padrão 1.0)"},
		"loop":           {"boolean", "Se a animação repete (padrão false)"},
	}, "node_path", "animation_name")},
	{"godot_add_animation_track", "Adiciona uma track de VALOR (nó:propriedade) a uma animação existente. Ex: track_path 'Sprite2D:position'.", schema(map[string][2]string{
		"node_path":      {"string", "Caminho do AnimationPlayer"},
		"animation_name": {"string", "Nome da animação"},
		"track_path":     {"string", "Caminho da propriedade animada (ex: 'Sprite2D:position', '.:modulate')"},
	}, "node_path", "animation_name", "track_path")},
	{"godot_set_animation_keyframe", "Insere um keyframe numa track de animação (por track_index ou track_path) em um tempo dado.", schema(map[string][2]string{
		"node_path":      {"string", "Caminho do AnimationPlayer"},
		"animation_name": {"string", "Nome da animação"},
		"track_index":    {"number", "Índice da track (ou informe track_path)"},
		"track_path":     {"string", "Caminho da track (alternativa ao índice)"},
		"time":           {"number", "Tempo do keyframe em segundos"},
		"value":          {"string|number|boolean|array", "Valor no keyframe (vetores [x,y])"},
	}, "node_path", "animation_name")},

	// --- BLOCO 5: TileMap extras ---
	{"godot_tilemap_fill_rect", "Preenche um retângulo de células num TileMap/TileMapLayer com um tile (coordenadas 'from' a 'to').", schema(map[string][2]string{
		"node_path":        {"string", "Caminho do nó TileMap/TileMapLayer"},
		"from":             {"array", "Canto inicial [x, y] do retângulo"},
		"to":               {"array", "Canto final [x, y] do retângulo"},
		"source_id":        {"number", "ID da fonte de tile no TileSet (padrão: 0)"},
		"atlas_coords":     {"array", "Coordenadas no atlas [x, y] (padrão: [0,0])"},
		"alternative_tile": {"number", "ID do tile alternativo (padrão: 0)"},
		"layer":            {"number", "Layer do TileMap legado (só TileMap, padrão: 0)"},
	}, "node_path", "from", "to")},
	{"godot_tilemap_clear", "Limpa TODAS as células de um TileMap/TileMapLayer (ou uma camada específica do TileMap).", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó TileMap/TileMapLayer"},
		"layer":     {"number", "Camada a limpar (só TileMap; omitir = todas)"},
	}, "node_path")},
	{"godot_tilemap_get_info", "Retorna metadados de um TileMap/TileMapLayer: TileSet, tile_size, fontes, camadas e contagem de células.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó TileMap/TileMapLayer"},
	}, "node_path")},

	// --- BLOCO 5: Input extras ---
	{"godot_simulate_mouse_click", "Simula um clique do mouse (press + release) nas coordenadas especificadas do jogo em execução.", schema(map[string][2]string{
		"position":     {"array", "Coordenadas [x, y] do clique"},
		"button":       {"number", "Botão: 1=esquerdo (padrão), 2=direito, 3=meio"},
		"double_click": {"boolean", "Se é duplo clique (padrão false)"},
	}, "position")},
	{"godot_simulate_mouse_move", "Simula o movimento do mouse para as coordenadas especificadas do jogo em execução.", schema(map[string][2]string{
		"position": {"array", "Coordenadas destino [x, y]"},
		"relative": {"array", "Movimento relativo [dx, dy] (opcional)"},
	}, "position")},
	{"godot_simulate_sequence", "Executa uma SEQUÊNCIA de inputs com delays entre cada passo (ex: direita, direita, cima, espaço). Cada step pode ser key, action, mouse_click ou wait.", schema(map[string][2]string{
		"steps":       {"array", `Array de passos. Cada passo: {"type":"key","key":"Right"}, {"type":"action","action":"jump"}, {"type":"mouse_click","position":[100,200]}, {"type":"wait"}. Opcional: "delay_ms" e "hold_ms" por passo.`},
		"interval_ms": {"number", "Delay padrão entre passos em ms (padrão: 100, máx: 2000)"},
	}, "steps")},

	// --- BLOCO 5: Editor extras ---
	{"godot_execute_editor_script", "Executa um snippet GDScript no contexto do editor (como EditorScript). O código é envolvido em _run() automaticamente. Poderoso mas use com cuidado.", schema(map[string][2]string{
		"code": {"string", "Código GDScript a executar (será envolvido em @tool extends EditorScript / func _run())"},
	}, "code")},
	{"godot_reload_plugin", "Desabilita e re-habilita um EditorPlugin pelo nome da pasta (recarrega). Útil após alterar scripts do plugin.", schema(map[string][2]string{
		"plugin_name": {"string", "Nome da pasta do addon em addons/ (ex: crom_ai)"},
	}, "plugin_name")},
	{"godot_reload_project", "Reinicia o editor Godot inteiro. ATENÇÃO: a conexão WebSocket será perdida — reconecte após o restart.", schema(nil)},

	// --- BLOCO 5: Animation extras ---
	{"godot_remove_animation", "Remove uma animação de um AnimationPlayer.", schema(map[string][2]string{
		"node_path":      {"string", "Caminho do AnimationPlayer"},
		"animation_name": {"string", "Nome da animação a remover"},
	}, "node_path", "animation_name")},

	// --- BLOCO 5: Physics extras ---
	{"godot_get_collision_info", "Retorna informações de colisão de um corpo físico: collision_layer, collision_mask, shapes filhas (tipo, tamanho/raio, disabled). Para Area2D inclui monitoring/monitorable.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do corpo físico (CharacterBody2D, RigidBody2D, StaticBody2D, Area2D)"},
	}, "node_path")},

	// --- BLOCO 6: Audio ---
	{"godot_add_audio_player", "Adiciona um nó de áudio (AudioStreamPlayer, AudioStreamPlayer2D ou 3D) à cena.", schema(map[string][2]string{
		"parent_path": {"string", "Caminho do nó pai (padrão: root)"},
		"type":        {"string", "Tipo: AudioStreamPlayer, AudioStreamPlayer2D ou AudioStreamPlayer3D"},
		"node_name":   {"string", "Nome do nó"},
		"stream_path": {"string", "Caminho do recurso de áudio (res://...)"},
		"bus":         {"string", "Nome do bus de áudio (ex: Master)"},
		"volume_db":   {"number", "Volume em dB"},
		"autoplay":    {"boolean", "Se deve tocar automaticamente"},
	})},
	{"godot_add_audio_bus", "Adiciona um novo bus de áudio ao AudioServer.", schema(map[string][2]string{
		"bus_name":  {"string", "Nome do novo bus"},
		"send":      {"string", "Nome do bus de destino"},
		"volume_db": {"number", "Volume em dB"},
	}, "bus_name")},
	{"godot_add_audio_bus_effect", "Adiciona um efeito a um bus de áudio.", schema(map[string][2]string{
		"bus_name":    {"string", "Nome do bus"},
		"effect_type": {"string", "Classe do AudioEffect (ex: AudioEffectReverb, AudioEffectFilter)"},
	}, "bus_name", "effect_type")},
	{"godot_set_audio_bus", "Altera configurações de um bus de áudio.", schema(map[string][2]string{
		"bus_name":  {"string", "Nome do bus"},
		"volume_db": {"number", "Volume em dB"},
		"mute":      {"boolean", "Mudo"},
		"solo":      {"boolean", "Solo"},
		"send":      {"string", "Bus de destino"},
	}, "bus_name")},
	{"godot_get_audio_bus_layout", "Retorna a estrutura completa de buses de áudio e seus efeitos.", schema(nil)},
	{"godot_get_audio_info", "Retorna informações sobre o componente de áudio de um nó.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó de áudio"},
	}, "node_path")},

	// --- BLOCO 7: Theme & UI ---
	{"godot_create_theme", "Cria e opcionalmente salva um novo recurso Theme.", schema(map[string][2]string{
		"save_path": {"string", "Caminho para salvar o tema (ex: res://theme.tres)"},
	})},
	{"godot_set_theme_color", "Define uma cor de tema em um Control.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó Control"},
		"item_name": {"string", "Nome do item de cor"},
		"type_name": {"string", "Nome do tipo de Control (ex: Label, Button)"},
		"color":     {"array", "Cor RGBA [r, g, b, a]"},
	}, "node_path", "item_name")},
	{"godot_set_theme_constant", "Define uma constante de tema em um Control.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó Control"},
		"item_name": {"string", "Nome da constante"},
		"type_name": {"string", "Nome do tipo de Control"},
		"value":     {"number", "Valor inteiro da constante"},
	}, "node_path", "item_name", "value")},
	{"godot_set_theme_font_size", "Define um tamanho de fonte no tema de um Control.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó Control"},
		"item_name": {"string", "Nome do item de fonte (padrão: font_size)"},
		"type_name": {"string", "Nome do tipo de Control"},
		"size":      {"number", "Tamanho em pixels"},
	}, "node_path", "size")},
	{"godot_set_theme_stylebox", "Define um StyleBoxFlat no tema de um Control.", schema(map[string][2]string{
		"node_path":     {"string", "Caminho do nó Control"},
		"item_name":     {"string", "Nome da propriedade de stylebox (ex: panel)"},
		"type_name":     {"string", "Nome do tipo de Control"},
		"bg_color":      {"array", "Cor de fundo [r, g, b, a]"},
		"corner_radius": {"number", "Raio dos cantos"},
		"border_width":  {"number", "Largura das bordas"},
	}, "node_path", "item_name")},
	{"godot_get_theme_info", "Retorna informações sobre o tema aplicado a um nó Control.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó Control"},
	}, "node_path")},

	// --- BLOCO 8: Resource/Project ---
	{"godot_edit_resource", "Edita propriedades de um arquivo de recurso (.tres) existente.", schema(map[string][2]string{
		"resource_path": {"string", "Caminho do recurso (res://...)"},
		"properties":    {"object", "Dicionário de propriedades a alterar"},
	}, "resource_path", "properties")},
	{"godot_get_resource_preview", "Retorna informações e metadados de um recurso.", schema(map[string][2]string{
		"resource_path": {"string", "Caminho do recurso"},
	}, "resource_path")},
	{"godot_add_autoload", "Adiciona um script/cena como Autoload (Singleton) do projeto.", schema(map[string][2]string{
		"name": {"string", "Nome do Autoload (Singleton)"},
		"path": {"string", "Caminho do arquivo (res://...)"},
	}, "name", "path")},
	{"godot_remove_autoload", "Remove um Autoload do projeto.", schema(map[string][2]string{
		"name": {"string", "Nome do Autoload"},
	}, "name")},
	{"godot_uid_to_project_path", "Converte um UID do Godot (uid://...) para o caminho relativo do projeto.", schema(map[string][2]string{
		"uid": {"string", "String do UID (ex: uid://...)"},
	}, "uid")},
	{"godot_project_path_to_uid", "Converte um caminho de arquivo no projeto para seu UID.", schema(map[string][2]string{
		"path": {"string", "Caminho res://..."},
	}, "path")},
	{"godot_list_scripts", "Lista todos os scripts (.gd, .cs) no diretório especificado ou todo o projeto.", schema(map[string][2]string{
		"dir_path": {"string", "Diretório inicial (padrão: res://)"},
	})},
	{"godot_search_in_files", "Busca por um texto/string dentro de arquivos de texto (.gd, .tscn, .tres) do projeto.", schema(map[string][2]string{
		"query":      {"string", "Texto a buscar"},
		"dir_path":   {"string", "Diretório (padrão: res://)"},
		"extensions": {"array", "Extensões a incluir (padrão: [.gd, .tscn, .tres, .cfg])"},
	}, "query")},

	// --- BLOCO 9: Node/Selection ---
	{"godot_select_nodes", "Seleciona nós no inspetor/árvore de cena do editor.", schema(map[string][2]string{
		"node_paths": {"array", "Array de caminhos de nós a selecionar"},
	}, "node_paths")},
	{"godot_clear_editor_selection", "Limpa a seleção atual de nós no editor.", schema(nil)},
	{"godot_set_anchor_preset", "Define o preset de âncora/alinhamento de um nó Control (0=top-left, 8=center, 15=full-rect, etc).", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó Control"},
		"preset":    {"number", "Preset enum (0 a 15)"},
		"resize":    {"boolean", "Se deve redimensionar o nó ao aplicar"},
	}, "node_path", "preset")},

	// --- BLOCO 10: Runtime avançado ---
	{"godot_execute_game_script", "Executa um snippet GDScript dentro do jogo em execução via crom_runtime.", schema(map[string][2]string{
		"code": {"string", "Código GDScript a executar no jogo"},
	}, "code")},
	{"godot_start_recording", "Inicia a gravação de alterações de uma propriedade em runtime.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó"},
		"property":  {"string", "Propriedade a monitorar"},
	}, "node_path", "property")},
	{"godot_stop_recording", "Para a gravação e retorna a lista de amostras gravadas.", schema(map[string][2]string{
		"samples": {"number", "Número de amostras finais a coletar"},
	})},
	{"godot_replay_recording", "Aplica as amostras da gravação anterior de volta ao nó no jogo.", schema(nil)},
	{"godot_find_nodes_by_script", "Acha nós no jogo em execução que possuem um determinado script.", schema(map[string][2]string{
		"script_path": {"string", "Caminho do script res://..."},
	}, "script_path")},
	{"godot_get_autoload", "Lista todos os Autoloads configurados e seus estados.", schema(nil)},
	{"godot_batch_get_properties", "Retorna propriedades de múltiplos nós do jogo em execução de uma só vez.", schema(map[string][2]string{
		"node_paths": {"array", "Lista de caminhos de nós"},
		"properties": {"array", "Lista de propriedades"},
	}, "node_paths", "properties")},
	{"godot_find_nearby_nodes", "Localiza nós próximos a uma determinada posição em runtime 2D/3D.", schema(map[string][2]string{
		"position": {"array", "Coordenadas [x, y] ou [x, y, z]"},
		"radius":   {"number", "Raio de busca em pixels/unidades"},
	}, "position", "radius")},
	{"godot_navigate_to", "Solicita que um NavigationAgent2D/3D navegue até uma posição alvo.", schema(map[string][2]string{
		"agent_path": {"string", "Caminho do nó NavigationAgent"},
		"target":     {"array", "Coordenadas de destino"},
	}, "agent_path", "target")},
	{"godot_move_to", "Move suavemente um nó até uma posição alvo em runtime.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó"},
		"target":    {"array", "Coordenadas [x, y] ou [x, y, z]"},
		"speed":     {"number", "Velocidade de movimento"},
	}, "node_path", "target")},

	// --- BLOCO 11: Testing/QA ---
	{"godot_compare_screenshots", "Compara duas imagens PNG e calcula a porcentagem e contagem de pixels diferentes.", schema(map[string][2]string{
		"path_a": {"string", "Caminho da primeira imagem"},
		"path_b": {"string", "Caminho da segunda imagem"},
	}, "path_a", "path_b")},
	{"godot_run_stress_test", "Executa a cena headless por N frames para detectar travamentos e vazamentos.", schema(map[string][2]string{
		"scene_path": {"string", "Caminho da cena (opcional, padrão: main scene)"},
		"frames":     {"number", "Quantidade de frames a rodar (padrão: 300)"},
	})},
	{"godot_get_test_report", "Gera um relatório automatizado de testes para uma cena.", schema(map[string][2]string{
		"scene_path": {"string", "Caminho da cena"},
	})},

	// --- BLOCO 12: Particle + Navigation ---
	{"godot_create_particles", "Cria um nó GPUParticles2D ou GPUParticles3D na cena.", schema(map[string][2]string{
		"parent_path": {"string", "Caminho do pai"},
		"node_name":   {"string", "Nome do nó"},
		"2d":          {"boolean", "Se é 2D (true) ou 3D (false)"},
		"amount":      {"number", "Quantidade de partículas"},
		"lifetime":    {"number", "Tempo de vida em segundos"},
		"emitting":    {"boolean", "Emitindo imediatamente"},
	})},
	{"godot_set_particle_material", "Define propriedades de um ParticleProcessMaterial.", schema(map[string][2]string{
		"node_path":            {"string", "Caminho do nó de partículas"},
		"direction":            {"array", "Vetor de direção [x, y, z]"},
		"spread":               {"number", "Ângulo de espalhamento em graus"},
		"gravity":              {"array", "Gravidade [x, y, z]"},
		"initial_velocity_min": {"number", "Velocidade inicial mínima"},
		"initial_velocity_max": {"number", "Velocidade inicial máxima"},
	}, "node_path")},
	{"godot_set_particle_color_gradient", "Define uma rampa de cores (Gradient) para as partículas.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó de partículas"},
		"colors":    {"array", "Array de cores RGBA [[r,g,b,a], ...]"},
	}, "node_path", "colors")},
	{"godot_apply_particle_preset", "Aplica um preset pronto (fire, smoke, sparks, explosion) às partículas.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó de partículas"},
		"preset":    {"string", "Nome do preset: fire, smoke, sparks, explosion"},
	}, "node_path", "preset")},
	{"godot_get_particle_info", "Retorna a configuração atual de um emissor de partículas.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó de partículas"},
	}, "node_path")},
	{"godot_setup_navigation_region", "Cria uma NavigationRegion2D com NavigationPolygon associado.", schema(map[string][2]string{
		"parent_path": {"string", "Caminho do nó pai"},
		"node_name":   {"string", "Nome da região"},
	})},
	{"godot_setup_navigation_agent", "Cria um nó NavigationAgent2D na cena.", schema(map[string][2]string{
		"parent_path":             {"string", "Caminho do nó pai"},
		"node_name":               {"string", "Nome do agente"},
		"target_desired_distance": {"number", "Distância mínima do alvo"},
		"path_desired_distance":   {"number", "Distância do ponto do caminho"},
	})},
	{"godot_bake_navigation_mesh", "Bake do polígono de navegação de um NavigationRegion2D.", schema(map[string][2]string{
		"node_path": {"string", "Caminho da NavigationRegion2D"},
	}, "node_path")},
	{"godot_set_navigation_layers", "Define as camadas de navegação (navigation_layers bitmask) em um nó.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó"},
		"layers":    {"number", "Bitmask de camadas"},
	}, "node_path", "layers")},
	{"godot_get_navigation_info", "Retorna dados de navegação de um nó.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó de navegação"},
	}, "node_path")},

	// --- BLOCO 13: AnimationTree/StateMachine ---
	{"godot_create_animation_tree", "Cria um nó AnimationTree configurado com nó raiz.", schema(map[string][2]string{
		"parent_path": {"string", "Caminho do nó pai"},
		"node_name":   {"string", "Nome do AnimationTree"},
		"root_type":   {"string", "Tipo da raiz: AnimationNodeStateMachine ou AnimationNodeBlendTree"},
		"anim_player": {"string", "NodePath para o AnimationPlayer"},
	})},
	{"godot_get_animation_tree_structure", "Retorna a estrutura interna e estado de um AnimationTree.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do AnimationTree"},
	}, "node_path")},
	{"godot_set_tree_parameter", "Define o valor de um parâmetro de blend/state no AnimationTree (parameters/...).", schema(map[string][2]string{
		"node_path": {"string", "Caminho do AnimationTree"},
		"parameter": {"string", "Caminho do parâmetro (ex: parameters/playback)"},
		"value":     {"string|number|boolean", "Novo valor do parâmetro"},
	}, "node_path", "parameter")},
	{"godot_add_state_machine_state", "Adiciona um nó de estado em uma AnimationNodeStateMachine.", schema(map[string][2]string{
		"node_path":  {"string", "Caminho do AnimationTree"},
		"state_name": {"string", "Nome do novo estado"},
		"node_type":  {"string", "Tipo de nó (ex: AnimationNodeAnimation)"},
		"animation":  {"string", "Nome da animação no AnimationPlayer"},
		"position_x": {"number", "Posição X no grafo do editor"},
		"position_y": {"number", "Posição Y no grafo do editor"},
	}, "node_path", "state_name")},
	{"godot_remove_state_machine_state", "Remove um estado de uma AnimationNodeStateMachine.", schema(map[string][2]string{
		"node_path":  {"string", "Caminho do AnimationTree"},
		"state_name": {"string", "Nome do estado a remover"},
	}, "node_path", "state_name")},
	{"godot_add_state_machine_transition", "Conecta dois estados com uma transição na AnimationNodeStateMachine.", schema(map[string][2]string{
		"node_path":    {"string", "Caminho do AnimationTree"},
		"from":         {"string", "Estado de origem"},
		"to":           {"string", "Estado de destino"},
		"auto_advance": {"boolean", "Se avança automaticamente"},
	}, "node_path", "from", "to")},
	{"godot_remove_state_machine_transition", "Remove uma transição entre dois estados na AnimationNodeStateMachine.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do AnimationTree"},
		"from":      {"string", "Estado de origem"},
		"to":        {"string", "Estado de destino"},
	}, "node_path", "from", "to")},
	{"godot_set_blend_tree_node", "Adiciona um nó ao grafo AnimationNodeBlendTree.", schema(map[string][2]string{
		"node_path":       {"string", "Caminho do AnimationTree"},
		"blend_node_name": {"string", "Nome do nó no BlendTree"},
		"node_type":       {"string", "Tipo de AnimationNode"},
		"position_x":      {"number", "Posição X no grafo"},
		"position_y":      {"number", "Posição Y no grafo"},
	}, "node_path", "blend_node_name")},

	// --- BLOCO 14: Shader + Export + Profiling ---
	{"godot_create_shader", "Cria um novo arquivo de shader (.gdshader).", schema(map[string][2]string{
		"save_path":   {"string", "Caminho do arquivo (res://...)"},
		"shader_type": {"string", "Tipo: canvas_item, spatial, particles, fog"},
		"code":        {"string", "Código GDShader inicial (opcional)"},
	}, "save_path")},
	{"godot_read_shader", "Lê o código-fonte de um arquivo de shader (.gdshader).", schema(map[string][2]string{
		"shader_path": {"string", "Caminho do shader (res://...)"},
	}, "shader_path")},
	{"godot_edit_shader", "Sobrescreve o código-fonte de um arquivo de shader existente.", schema(map[string][2]string{
		"shader_path": {"string", "Caminho do shader"},
		"code":        {"string", "Novo código do shader"},
	}, "shader_path", "code")},
	{"godot_assign_shader_material", "Cria um ShaderMaterial e o atribui a um nó.", schema(map[string][2]string{
		"node_path":   {"string", "Caminho do nó"},
		"shader_path": {"string", "Caminho do shader (.gdshader)"},
	}, "node_path", "shader_path")},
	{"godot_set_shader_param", "Define o valor de um uniform/parâmetro em um ShaderMaterial.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó"},
		"param":     {"string", "Nome do uniform no shader"},
		"value":     {"string|number|boolean|array", "Valor do uniform"},
	}, "node_path", "param")},
	{"godot_get_shader_params", "Lista os uniforms e parâmetros disponíveis de um ShaderMaterial.", schema(map[string][2]string{
		"node_path": {"string", "Caminho do nó"},
	}, "node_path")},
	{"godot_list_export_presets", "Lista os presets de exportação configurados no export_presets.cfg do projeto.", schema(nil)},
	{"godot_export_project", "Executa a exportação do projeto usando um preset configurado.", schema(map[string][2]string{
		"preset":      {"string", "Nome do preset de exportação"},
		"output_path": {"string", "Caminho do executável de saída"},
	}, "preset", "output_path")},
	{"godot_get_export_info", "Retorna informações sobre presets e versão do engine para exportação.", schema(nil)},
	{"godot_get_performance_monitors", "Coleta métricas de desempenho em tempo real (FPS, tempo de render, memória, nós, draw calls).", schema(nil)},
	{"godot_get_editor_performance", "Retorna as métricas de desempenho no contexto do editor.", schema(nil)},

	// --- BLOCO 15: Batch/Refactoring + Analysis ---
	{"godot_find_nodes_by_type", "Busca todos os nós de uma determinada classe na cena aberta.", schema(map[string][2]string{
		"type": {"string", "Nome da classe (ex: Sprite2D, Area2D, CollisionShape2D)"},
	}, "type")},
	{"godot_find_signal_connections", "Lista todas as conexões de sinais ativas na cena aberta.", schema(map[string][2]string{
		"signal": {"string", "Nome do sinal (opcional)"},
	})},
	{"godot_batch_set_property", "Altera uma propriedade em TODOS os nós de uma classe específica na cena.", schema(map[string][2]string{
		"type":     {"string", "Classe dos nós afetados"},
		"property": {"string", "Nome da propriedade"},
		"value":    {"string|number|boolean|array", "Novo valor"},
	}, "type", "property")},
	{"godot_find_node_references", "Procura referências ao nome de um nó nos arquivos do projeto.", schema(map[string][2]string{
		"node_name": {"string", "Nome do nó"},
	}, "node_name")},
	{"godot_get_scene_dependencies", "Lista os recursos externos (ext_resource) dos quais uma cena depende.", schema(map[string][2]string{
		"scene_path": {"string", "Caminho da cena (res://...)"},
	}, "scene_path")},
	{"godot_cross_scene_set_property", "Modifica a propriedade de um nó dentro de um arquivo de cena (.tscn) sem precisar abri-la no editor.", schema(map[string][2]string{
		"scene_path": {"string", "Caminho do arquivo .tscn"},
		"node_path":  {"string", "Caminho do nó dentro da cena"},
		"property":   {"string", "Nome da propriedade"},
		"value":      {"string|number|boolean|array", "Novo valor"},
	}, "scene_path", "node_path", "property")},
	{"godot_find_script_references", "Busca em quais cenas e scripts um determinado script res:// é utilizado.", schema(map[string][2]string{
		"script_path": {"string", "Caminho do script"},
	}, "script_path")},
	{"godot_detect_circular_dependencies", "Analisa a arvore de dependências entre cenas e detecta ciclos.", schema(map[string][2]string{
		"scene_path": {"string", "Cena inicial (opcional, padrão: main scene)"},
	})},
	{"godot_analyze_scene_complexity", "Calcula métricas de complexidade de uma cena (contagem de nós, profundidade máxima, scripts).", schema(nil)},
	{"godot_analyze_signal_flow", "Mapeia todo o fluxo de sinais e ouvintes na cena aberta.", schema(nil)},
	{"godot_find_unused_resources", "Detecta arquivos de recursos em res:// que não aparecem referenciados em nenhuma cena ou script.", schema(map[string][2]string{
		"dir_path": {"string", "Diretório inicial (padrão: res://)"},
	})},
	{"godot_get_project_statistics", "Gera estatísticas globais do projeto (contagem de scripts, cenas, recursos, cena principal).", schema(nil)},

	// --- BLOCO 16: 3D ---
	{"godot_add_mesh_instance", "Adiciona um nó MeshInstance3D com primitiva (BoxMesh, SphereMesh, etc) à cena.", schema(map[string][2]string{
		"parent_path": {"string", "Caminho do nó pai"},
		"node_name":   {"string", "Nome do nó"},
		"mesh_type":   {"string", "Tipo de mesh: BoxMesh, SphereMesh, CylinderMesh, PlaneMesh, PrismMesh"},
		"size":        {"array", "Tamanho [x, y, z] (para BoxMesh)"},
		"radius":      {"number", "Raio (para SphereMesh)"},
		"position":    {"array", "Posição [x, y, z]"},
	})},
	{"godot_setup_camera_3d", "Adiciona e configura um nó Camera3D.", schema(map[string][2]string{
		"parent_path": {"string", "Caminho do nó pai"},
		"node_name":   {"string", "Nome do nó"},
		"position":    {"array", "Posição [x, y, z]"},
		"fov":         {"number", "Campo de visão em graus (FOV)"},
		"current":     {"boolean", "Se é a câmera ativa"},
	})},
	{"godot_setup_lighting", "Cria um nó de luz 3D (DirectionalLight3D, OmniLight3D, SpotLight3D).", schema(map[string][2]string{
		"parent_path": {"string", "Caminho do nó pai"},
		"light_type":  {"string", "Tipo de luz: DirectionalLight3D, OmniLight3D, SpotLight3D"},
		"node_name":   {"string", "Nome do nó"},
		"energy":      {"number", "Energia da luz"},
		"color":       {"array", "Cor [r, g, b, a]"},
		"position":    {"array", "Posição [x, y, z]"},
	})},
	{"godot_setup_environment", "Adiciona um WorldEnvironment com Environment configurado.", schema(map[string][2]string{
		"parent_path":   {"string", "Caminho do nó pai"},
		"node_name":     {"string", "Nome do nó"},
		"bg_mode":       {"number", "Modo de fundo (0=clear, 1=custom_color, 2=sky)"},
		"bg_color":      {"array", "Cor de fundo [r, g, b, a]"},
		"ambient_color": {"array", "Cor de luz ambiente [r, g, b, a]"},
	})},
	{"godot_add_gridmap", "Adiciona um nó GridMap 3D à cena.", schema(map[string][2]string{
		"parent_path": {"string", "Caminho do nó pai"},
		"node_name":   {"string", "Nome do nó"},
		"cell_size":   {"array", "Tamanho de cada célula [x, y, z]"},
	})},
	{"godot_set_material_3d", "Cria e atribui um StandardMaterial3D a uma superfície de um MeshInstance3D.", schema(map[string][2]string{
		"node_path":      {"string", "Caminho do MeshInstance3D"},
		"surface_index": {"number", "Índice da superfície (padrão: 0)"},
		"albedo_color":  {"array", "Cor Albedo [r, g, b, a]"},
		"metallic":      {"number", "Valor metálico (0 a 1)"},
		"roughness":     {"number", "Rugosidade (0 a 1)"},
		"emission_color": {"array", "Cor de emissão [r, g, b, a]"},
	}, "node_path")},
}

