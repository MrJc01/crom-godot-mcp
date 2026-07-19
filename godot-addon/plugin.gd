@tool
extends EditorPlugin

# ==============================================================================
# Plugin MÍNIMO do crom-godot-mcp (lado Godot do MCP).
# Expõe o Editor Godot num servidor WebSocket (porta 8080) e registra o autoload
# CromRuntime (inspeção do jogo em execução, porta 8091). O servidor Go
# (godot-mcp) encaminha as ações godot_* para cá; o command_processor.gd executa.
# NÃO inclui chat/hub/menus — isso é do app que consome este addon.
# ==============================================================================

var websocket_server: Node = null
var command_processor: Node = null

func _addon_dir() -> String:
	var s: Script = get_script()
	return s.resource_path.get_base_dir()

func _enter_tree() -> void:
	var d := _addon_dir()
	var ProcessorClass: GDScript = load(d + "/command_processor.gd")
	if ProcessorClass:
		command_processor = ProcessorClass.new(self)
		command_processor.name = "CommandProcessor"
		add_child(command_processor)
	var ServerClass: GDScript = load(d + "/websocket_server.gd")
	if ServerClass and command_processor:
		websocket_server = ServerClass.new(command_processor, 8080)
		websocket_server.name = "WebSocketServer"
		add_child(websocket_server)
		var ok: bool = websocket_server.start_server()
		if ok:
			print("[crom-godot-mcp] Editor exposto na porta 8080 (ws://127.0.0.1:8080).")
	add_autoload_singleton("CromRuntime", d + "/crom_runtime.gd")

func _process(_delta: float) -> void:
	if websocket_server and websocket_server.has_method("process_network"):
		websocket_server.process_network()

func _exit_tree() -> void:
	remove_autoload_singleton("CromRuntime")
	if websocket_server:
		if websocket_server.has_method("stop_server"):
			websocket_server.stop_server()
		websocket_server.queue_free()
		websocket_server = null
	if command_processor:
		command_processor.queue_free()
		command_processor = null
