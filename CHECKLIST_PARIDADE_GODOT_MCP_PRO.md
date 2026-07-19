# crom-godot-mcp — Checklist de Paridade vs godot-mcp-pro

Referência: **godot-mcp-pro** (`youichi-uda/godot-mcp-pro`) — **175 tools** (full),
84 (lite), 35 (minimal). Proprietário ($15, não importável). Reimplementamos as
capacidades no nosso servidor MIT.

**Estado atual do crom-godot-mcp: 74 ferramentas.** Este arquivo lista TODAS as 175
do godot-mcp-pro por categoria, marcando nossa cobertura, para saber o que falta.

## Legenda
- ✅ **temos** (nome nosso entre parênteses)
- 🟡 **parcial** (dá pra fazer por outra tool, mas não é dedicada/1:1)
- ⬜ **falta**
- 🆕 **nossa, além do godot-mcp-pro** (no fim)

---

## Project (7) — temos 2, parcial 1
- ✅ get_project_settings (`get_project_setting`)
- ✅ set_project_setting (`set_project_setting`)
- 🟡 get_filesystem_tree (`list_project_dir` — por diretório, não árvore inteira)
- ✅ get_project_info (`godot_get_project_info`)
- ✅ search_files (`godot_search_files`)
- ⬜ uid_to_project_path
- ⬜ project_path_to_uid

## Scene (9) — temos 7, parcial 1
- ✅ get_scene_tree (`get_scene_tree`)
- ✅ create_scene (`create_scene`)
- ✅ open_scene (`open_scene`)
- ✅ add_scene_instance (`instantiate_scene`)
- ✅ play_scene (`play_scene`)
- ✅ stop_scene (`stop_scene`)
- ✅ save_scene (`save_scene`)
- 🟡 get_scene_file_content (`read_project_file` no .tscn)
- ✅ delete_scene (`godot_delete_scene`)

## Node (17) — temos 8, parcial 4
- ✅ add_node (`add_node`)
- ✅ delete_node (`remove_node`)
- ✅ duplicate_node (`duplicate_node`)
- ✅ move_node (`move_node`)
- ✅ update_property (`set_node_property`)
- ✅ add_resource (`create_resource` + recurso inline em `set_node_property`)
- ✅ rename_node (`rename_node`)
- ✅ connect_signal (`connect_signal`)
- 🟡 set_node_groups (`add_to_group` / `remove_from_group`)
- ✅ get_node_groups (`godot_get_node_groups`) *(agora dedicada)*
- 🟡 get_editor_selection (`get_open_editor_context`)
- ✅ get_node_properties (`godot_get_node_properties`) *(agora dedicada)*
- ✅ disconnect_signal (`godot_disconnect_signal`)
- ✅ find_nodes_in_group (`godot_find_nodes_in_group`)
- ⬜ select_nodes
- ⬜ clear_editor_selection
- ⬜ set_anchor_preset

## Script (8) — temos 6, parcial 1
- ✅ read_script (`read_script`)
- ✅ create_script / attach_script (`create_and_attach_script`)
- ✅ edit_script (`set_script_source`)
- ✅ validate_script (`gdscript_check`)
- 🟡 get_open_scripts (`get_open_editor_context`)
- ⬜ list_scripts
- ⬜ search_in_files

## Editor (9) — temos 4, parcial 2
- ✅ get_editor_errors (`get_console_errors`)
- ✅ get_editor_screenshot (`capture_screenshot`)
- ✅ clear_output (`clear_output`)
- ✅ get_output_log (`get_output`)
- 🟡 get_game_screenshot (`capture_screenshot`)
- 🟡 get_signals (`list_node_signals`)
- ⬜ execute_editor_script
- ⬜ reload_plugin
- ⬜ reload_project

## Input (7) — temos 4
- ✅ simulate_key (`simulate_key`)
- ✅ simulate_action (`simulate_action`)
- ✅ get_input_actions (`list_input_actions`)
- ✅ set_input_action (`add_input_action`)
- ⬜ simulate_mouse_click
- ⬜ simulate_mouse_move
- ⬜ simulate_sequence

## Runtime (19) — temos 2, parcial 2  ⭐ (área mais fraca vs pro)
- ✅ get_game_scene_tree (`get_runtime_scene_tree`)
- ✅ monitor_properties (`record_property_over_time`)
- ✅ get_game_node_properties (`godot_get_game_node_properties`) *(agora dedicada)*
- 🟡 capture_frames (`record_property_over_time` — por propriedade, não frames)
- ✅ set_game_node_property (`godot_set_game_node_property`)
- ⬜ execute_game_script
- ⬜ start_recording / stop_recording / replay_recording
- ⬜ find_nodes_by_script
- ⬜ get_autoload
- ⬜ batch_get_properties
- ✅ find_ui_elements (`godot_find_ui_elements`)
- ✅ click_button_by_text (`godot_click_button_by_text`)
- ✅ wait_for_node (`godot_wait_for_node`)
- ⬜ find_nearby_nodes
- ⬜ navigate_to / move_to

## Animation (6) — temos 1, parcial 1
- ✅ list_animations (`list_animations`)
- 🟡 get_animation_info (`list_animations`)
- ✅ create_animation (`godot_create_animation`)
- ✅ add_animation_track (`godot_add_animation_track`)
- ✅ set_animation_keyframe (`godot_set_animation_keyframe`)
- ⬜ remove_animation

## TileMap (6) — temos 2, parcial 1
- ✅ tilemap_set_cell (`set_tilemap_cell`)
- ✅ tilemap_get_used_cells (`get_tilemap_cells`)
- 🟡 tilemap_get_cell (`get_tilemap_cells`)
- ⬜ tilemap_fill_rect
- ⬜ tilemap_clear
- ⬜ tilemap_get_info

## Physics (6) — parcial 2
- ✅ setup_physics_body (`godot_setup_physics_body`) *(agora dedicada)*
- 🟡 setup_collision (recurso inline `shape` em `set_node_property`/`add_node`)
- ✅ set_physics_layers (`godot_set_physics_layers`)
- ✅ get_physics_layers (`godot_get_physics_layers`)
- ⬜ get_collision_info
- ✅ add_raycast (`godot_add_raycast`)

## Resource (6) — temos 1, parcial 1
- ✅ create_resource (`create_resource`)
- 🟡 read_resource (`read_project_file`)
- ⬜ edit_resource
- ⬜ get_resource_preview
- ⬜ add_autoload
- ⬜ remove_autoload

## Testing & QA (6) — parcial (coberto por nossa `verify_playable` 🆕)
- 🟡 run_test_scenario / assert_node_state / assert_screen_text — cobertos parcialmente por `verify_playable` (veredito jogável + erros + movimento)
- ⬜ compare_screenshots
- ⬜ run_stress_test
- ⬜ get_test_report

## Categorias ainda NÃO cobertas (0 tools nossas)
- **3D Scene (6)**: add_mesh_instance, setup_camera_3d, setup_lighting, setup_environment, add_gridmap, set_material_3d — *foco 2D primeiro*
- **Particle (5)**: create_particles, set_particle_material, set_particle_color_gradient, apply_particle_preset, get_particle_info
- **Navigation (6)**: setup_navigation_region, setup_navigation_agent, bake_navigation_mesh, set_navigation_layers, get_navigation_info
- **Audio (6)**: add_audio_player, add_audio_bus, add_audio_bus_effect, set_audio_bus, get_audio_bus_layout, get_audio_info
- **AnimationTree/State Machine/Blend (8)**: create_animation_tree, get_animation_tree_structure, set_tree_parameter, add_state_machine_state, remove_state_machine_state, add/remove_state_machine_transition, set_blend_tree_node
- **Theme & UI (6)**: create_theme, set_theme_color, set_theme_constant, set_theme_font_size, set_theme_stylebox, get_theme_info
- **Shader (6)**: create_shader, read_shader, edit_shader, assign_shader_material, set_shader_param, get_shader_params
- **Export (3)**: list_export_presets, export_project, get_export_info
- **Profiling (2)**: get_performance_monitors, get_editor_performance
- **Batch & Refactoring (8)**: find_nodes_by_type, find_signal_connections, batch_set_property, find_node_references, get_scene_dependencies, cross_scene_set_property, find_script_references, detect_circular_dependencies
- **Analysis & Search (4)**: analyze_scene_complexity, analyze_signal_flow, find_unused_resources, get_project_statistics

---

## 🆕 Nossas ferramentas ALÉM do godot-mcp-pro (foco no LLM leve)
- ⭐ `godot_verify_playable` — **composto que fecha o feedback loop**: roda a cena
  headless, coleta erros (autoritativo via stderr), detecta movimento e devolve
  veredito `playable=true/false`. O godot-mcp-pro tem asserts avulsos; nós temos o
  passo único "é jogável?".
- ⭐ `godot_class_reference` — API autoritativa da versão do Godot via ClassDB
  (métodos com assinatura, props, sinais). Combate o drift Godot 3→4 na raiz.
- ⭐ `godot_gdscript_check` — valida sintaxe do .gd ANTES de anexar.
- `godot_add_nodes_batch` — monta subárvore inteira num passo (menos round-trips).
- `godot_get_node_config_warnings` — avisos do editor (ex.: CollisionShape2D sem shape).
- `godot_docs_search` — busca na doc offline do Godot.
- `godot_set_main_scene`, `godot_set_camera_target`, `godot_play_animation`,
  `godot_detach_script`, `godot_list_node_methods`, `godot_reparent_node`.
- Recurso **inline** em properties (`{"__resource_type":...}`) e **caminhos de nó relativos**.

---

## Resumo de paridade
| | godot-mcp-pro | crom-godot-mcp |
|---|---|---|
| Full | 175 | **74** |
| Cobertura direta (✅) | — | ~55 |
| Parcial (🟡) | — | ~9 |
| Nossas extras (🆕) | — | ~8 |

**Prioridade para paridade útil (2D + feedback):** Runtime (set_game_node_property,
find_ui_elements, click_button_by_text, wait_for_node) e Testing (assert_node_state,
assert_screen_text) — completam o QA automatizado. 3D/Shader/Particle/Navigation
ficam por último (foco 2D).
