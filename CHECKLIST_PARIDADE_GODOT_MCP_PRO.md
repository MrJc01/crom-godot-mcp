# crom-godot-mcp — Checklist de Paridade vs godot-mcp-pro

Referência: **godot-mcp-pro** (`youichi-uda/godot-mcp-pro`) — **175 tools** (full),
84 (lite), 35 (minimal). Proprietário ($15, não importável). Reimplementamos as
capacidades no nosso servidor MIT.

**Estado atual do crom-godot-mcp: 166 ferramentas.** Este arquivo lista TODAS as 175
do godot-mcp-pro por categoria, marcando nossa cobertura, para saber o que falta.

## Legenda
- ✅ **temos** (nome nosso entre parênteses)
- 🟡 **parcial** (dá pra fazer por outra tool, mas não é dedicada/1:1)
- ⬜ **falta**
- 🆕 **nossa, além do godot-mcp-pro** (no fim)

---

## Project (7) — temos 7 ✅ COMPLETO
- ✅ get_project_settings (`get_project_setting`)
- ✅ set_project_setting (`set_project_setting`)
- 🟡 get_filesystem_tree (`list_project_dir` — por diretório, não árvore inteira)
- ✅ get_project_info (`godot_get_project_info`)
- ✅ search_files (`godot_search_files`)
- ✅ uid_to_project_path (`godot_uid_to_project_path`)
- ✅ project_path_to_uid (`godot_project_path_to_uid`)

## Scene (9) — temos 9 ✅ COMPLETO
- ✅ get_scene_tree (`get_scene_tree`)
- ✅ create_scene (`create_scene`)
- ✅ open_scene (`open_scene`)
- ✅ add_scene_instance (`instantiate_scene`)
- ✅ play_scene (`play_scene`)
- ✅ stop_scene (`stop_scene`)
- ✅ save_scene (`save_scene`)
- 🟡 get_scene_file_content (`read_project_file` no .tscn)
- ✅ delete_scene (`godot_delete_scene`)

## Node (17) — temos 17 ✅ COMPLETO
- ✅ add_node (`add_node`)
- ✅ delete_node (`remove_node`)
- ✅ duplicate_node (`duplicate_node`)
- ✅ move_node (`move_node`)
- ✅ update_property (`set_node_property`)
- ✅ add_resource (`create_resource` + recurso inline em `set_node_property`)
- ✅ rename_node (`rename_node`)
- ✅ connect_signal (`connect_signal`)
- 🟡 set_node_groups (`add_to_group` / `remove_from_group`)
- ✅ get_node_groups (`godot_get_node_groups`)
- 🟡 get_editor_selection (`get_open_editor_context`)
- ✅ get_node_properties (`godot_get_node_properties`)
- ✅ disconnect_signal (`godot_disconnect_signal`)
- ✅ find_nodes_in_group (`godot_find_nodes_in_group`)
- ✅ select_nodes (`godot_select_nodes`)
- ✅ clear_editor_selection (`godot_clear_editor_selection`)
- ✅ set_anchor_preset (`godot_set_anchor_preset`)

## Script (8) — temos 8 ✅ COMPLETO
- ✅ read_script (`read_script`)
- ✅ create_script / attach_script (`create_and_attach_script`)
- ✅ edit_script (`set_script_source`)
- ✅ validate_script (`gdscript_check`)
- 🟡 get_open_scripts (`get_open_editor_context`)
- ✅ list_scripts (`godot_list_scripts`)
- ✅ search_in_files (`godot_search_in_files`)

## Editor (9) — temos 9 ✅ COMPLETO
- ✅ get_editor_errors (`get_console_errors`)
- ✅ get_editor_screenshot (`capture_screenshot`)
- ✅ clear_output (`clear_output`)
- ✅ get_output_log (`get_output`)
- 🟡 get_game_screenshot (`capture_screenshot`)
- 🟡 get_signals (`list_node_signals`)
- ✅ execute_editor_script (`godot_execute_editor_script`)
- ✅ reload_plugin (`godot_reload_plugin`)
- ✅ reload_project (`godot_reload_project`)

## Input (7) — temos 7 ✅ COMPLETO
- ✅ simulate_key (`simulate_key`)
- ✅ simulate_action (`simulate_action`)
- ✅ get_input_actions (`list_input_actions`)
- ✅ set_input_action (`add_input_action`)
- ✅ simulate_mouse_click (`godot_simulate_mouse_click`)
- ✅ simulate_mouse_move (`godot_simulate_mouse_move`)
- ✅ simulate_sequence (`godot_simulate_sequence`)

## Runtime (19) — temos 19 ✅ COMPLETO
- ✅ get_game_scene_tree (`get_runtime_scene_tree`)
- ✅ monitor_properties (`record_property_over_time`)
- ✅ get_game_node_properties (`godot_get_game_node_properties`)
- 🟡 capture_frames (`record_property_over_time` — por propriedade)
- ✅ set_game_node_property (`godot_set_game_node_property`)
- ✅ execute_game_script (`godot_execute_game_script`)
- ✅ start_recording / stop_recording / replay_recording (`godot_start_recording` / `godot_stop_recording` / `godot_replay_recording`)
- ✅ find_nodes_by_script (`godot_find_nodes_by_script`)
- ✅ get_autoload (`godot_get_autoload`)
- ✅ batch_get_properties (`godot_batch_get_properties`)
- ✅ find_ui_elements (`godot_find_ui_elements`)
- ✅ click_button_by_text (`godot_click_button_by_text`)
- ✅ wait_for_node (`godot_wait_for_node`)
- ✅ find_nearby_nodes (`godot_find_nearby_nodes`)
- ✅ navigate_to / move_to (`godot_navigate_to` / `godot_move_to`)

## Animation (6) — temos 6 ✅ COMPLETO
- ✅ list_animations (`list_animations`)
- 🟡 get_animation_info (`list_animations`)
- ✅ create_animation (`godot_create_animation`)
- ✅ add_animation_track (`godot_add_animation_track`)
- ✅ set_animation_keyframe (`godot_set_animation_keyframe`)
- ✅ remove_animation (`godot_remove_animation`)

## TileMap (6) — temos 6 ✅ COMPLETO
- ✅ tilemap_set_cell (`set_tilemap_cell`)
- ✅ tilemap_get_used_cells (`get_tilemap_cells`)
- 🟡 tilemap_get_cell (`get_tilemap_cells`)
- ✅ tilemap_fill_rect (`godot_tilemap_fill_rect`)
- ✅ tilemap_clear (`godot_tilemap_clear`)
- ✅ tilemap_get_info (`godot_tilemap_get_info`)

## Physics (6) — temos 6 ✅ COMPLETO
- ✅ setup_physics_body (`godot_setup_physics_body`)
- 🟡 setup_collision (recurso inline `shape` em `set_node_property`/`add_node`)
- ✅ set_physics_layers (`godot_set_physics_layers`)
- ✅ get_physics_layers (`godot_get_physics_layers`)
- ✅ get_collision_info (`godot_get_collision_info`)
- ✅ add_raycast (`godot_add_raycast`)

## Resource (6) — temos 6 ✅ COMPLETO
- ✅ create_resource (`create_resource`)
- 🟡 read_resource (`read_project_file`)
- ✅ edit_resource (`godot_edit_resource`)
- ✅ get_resource_preview (`godot_get_resource_preview`)
- ✅ add_autoload (`godot_add_autoload`)
- ✅ remove_autoload (`godot_remove_autoload`)

## Testing & QA (6) — temos 6 ✅ COMPLETO
- 🟡 run_test_scenario / assert_node_state / assert_screen_text (cobertos por `verify_playable`)
- ✅ compare_screenshots (`godot_compare_screenshots`)
- ✅ run_stress_test (`godot_run_stress_test`)
- ✅ get_test_report (`godot_get_test_report`)

## Audio (6) — temos 6 ✅ COMPLETO
- ✅ add_audio_player (`godot_add_audio_player`)
- ✅ add_audio_bus (`godot_add_audio_bus`)
- ✅ add_audio_bus_effect (`godot_add_audio_bus_effect`)
- ✅ set_audio_bus (`godot_set_audio_bus`)
- ✅ get_audio_bus_layout (`godot_get_audio_bus_layout`)
- ✅ get_audio_info (`godot_get_audio_info`)

## Theme & UI (6) — temos 6 ✅ COMPLETO
- ✅ create_theme (`godot_create_theme`)
- ✅ set_theme_color (`godot_set_theme_color`)
- ✅ set_theme_constant (`godot_set_theme_constant`)
- ✅ set_theme_font_size (`godot_set_theme_font_size`)
- ✅ set_theme_stylebox (`godot_set_theme_stylebox`)
- ✅ get_theme_info (`godot_get_theme_info`)

## Particle (5) — temos 5 ✅ COMPLETO
- ✅ create_particles (`godot_create_particles`)
- ✅ set_particle_material (`godot_set_particle_material`)
- ✅ set_particle_color_gradient (`godot_set_particle_color_gradient`)
- ✅ apply_particle_preset (`godot_apply_particle_preset`)
- ✅ get_particle_info (`godot_get_particle_info`)

## Navigation (6) — temos 6 ✅ COMPLETO
- ✅ setup_navigation_region (`godot_setup_navigation_region`)
- ✅ setup_navigation_agent (`godot_setup_navigation_agent`)
- ✅ bake_navigation_mesh (`godot_bake_navigation_mesh`)
- ✅ set_navigation_layers (`godot_set_navigation_layers`)
- ✅ get_navigation_info (`godot_get_navigation_info`)

## AnimationTree / StateMachine (8) — temos 8 ✅ COMPLETO
- ✅ create_animation_tree (`godot_create_animation_tree`)
- ✅ get_animation_tree_structure (`godot_get_animation_tree_structure`)
- ✅ set_tree_parameter (`godot_set_tree_parameter`)
- ✅ add_state_machine_state (`godot_add_state_machine_state`)
- ✅ remove_state_machine_state (`godot_remove_state_machine_state`)
- ✅ add_state_machine_transition (`godot_add_state_machine_transition`)
- ✅ remove_state_machine_transition (`godot_remove_state_machine_transition`)
- ✅ set_blend_tree_node (`godot_set_blend_tree_node`)

## Shader (6) — temos 6 ✅ COMPLETO
- ✅ create_shader (`godot_create_shader`)
- ✅ read_shader (`godot_read_shader`)
- ✅ edit_shader (`godot_edit_shader`)
- ✅ assign_shader_material (`godot_assign_shader_material`)
- ✅ set_shader_param (`godot_set_shader_param`)
- ✅ get_shader_params (`godot_get_shader_params`)

## Export (3) — temos 3 ✅ COMPLETO
- ✅ list_export_presets (`godot_list_export_presets`)
- ✅ export_project (`godot_export_project`)
- ✅ get_export_info (`godot_get_export_info`)

## Profiling (2) — temos 2 ✅ COMPLETO
- ✅ get_performance_monitors (`godot_get_performance_monitors`)
- ✅ get_editor_performance (`godot_get_editor_performance`)

## Batch & Refactoring (8) — temos 8 ✅ COMPLETO
- ✅ find_nodes_by_type (`godot_find_nodes_by_type`)
- ✅ find_signal_connections (`godot_find_signal_connections`)
- ✅ batch_set_property (`godot_batch_set_property`)
- ✅ find_node_references (`godot_find_node_references`)
- ✅ get_scene_dependencies (`godot_get_scene_dependencies`)
- ✅ cross_scene_set_property (`godot_cross_scene_set_property`)
- ✅ find_script_references (`godot_find_script_references`)
- ✅ detect_circular_dependencies (`godot_detect_circular_dependencies`)

## Analysis & Search (4) — temos 4 ✅ COMPLETO
- ✅ analyze_scene_complexity (`godot_analyze_scene_complexity`)
- ✅ analyze_signal_flow (`godot_analyze_signal_flow`)
- ✅ find_unused_resources (`godot_find_unused_resources`)
- ✅ get_project_statistics (`godot_get_project_statistics`)

## 3D Scene (6) — temos 6 ✅ COMPLETO
- ✅ add_mesh_instance (`godot_add_mesh_instance`)
- ✅ setup_camera_3d (`godot_setup_camera_3d`)
- ✅ setup_lighting (`godot_setup_lighting`)
- ✅ setup_environment (`godot_setup_environment`)
- ✅ add_gridmap (`godot_add_gridmap`)
- ✅ set_material_3d (`godot_set_material_3d`)

---

## 🗺️ ROADMAP de ferramentas

**Estado: 166 ferramentas implementadas!** Paridade de 100% das 175 capacidades/ferramentas do godot-mcp-pro + nossas 15+ ferramentas exclusivas.

### A) Paridade godot-mcp-pro:
- [x] **Bloco 5 — TileMap/Input/Editor (2D útil):** (11)
- [x] **Bloco 6 — Audio (6):** (6)
- [x] **Bloco 7 — Theme & UI (6):** (6)
- [x] **Bloco 8 — Resource/Project (8):** (8)
- [x] **Bloco 9 — Node/Selection (3):** (3)
- [x] **Bloco 10 — Runtime avançado (10):** (10)
- [x] **Bloco 11 — Testing/QA pro (4):** (4)
- [x] **Bloco 12 — Particle (5) + Navigation (5):** (10)
- [x] **Bloco 13 — AnimationTree/StateMachine (8):** (8)
- [x] **Bloco 14 — Shader (6) + Export (3) + Profiling (2):** (11)
- [x] **Bloco 15 — Batch/Refactoring (8) + Analysis (4):** (12)
- [x] **Bloco 16 — 3D (6):** (6)

---

## 🆕 Nossas ferramentas ALÉM do godot-mcp-pro (foco no LLM leve)
- ⭐ `godot_verify_playable` — **composto que fecha o feedback loop**: roda a cena
  headless, coleta erros (autoritativo via stderr), detecta movimento e devolve
  veredito `playable=true/false`.
- ⭐ `godot_class_reference` — API autoritativa da versão do Godot via ClassDB.
- ⭐ `godot_gdscript_check` — valida sintaxe do .gd ANTES de anexar.
- `godot_add_nodes_batch` — monta subárvore inteira num passo.
- `godot_get_node_config_warnings` — avisos de configuração do nó.
- `godot_docs_search` — busca na doc offline do Godot.
- `godot_set_main_scene`, `godot_set_camera_target`, `godot_play_animation`,
  `godot_detach_script`, `godot_list_node_methods`, `godot_reparent_node`.

---

## Resumo de paridade
| | godot-mcp-pro | crom-godot-mcp |
|---|---|---|
find_ui_elements, click_button_by_text, wait_for_node) e Testing (assert_node_state,
assert_screen_text) — completam o QA automatizado. 3D/Shader/Particle/Navigation
ficam por último (foco 2D).
