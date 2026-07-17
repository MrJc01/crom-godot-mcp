package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const protocolVersion = "2024-11-05"

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      *int64          `json:"id,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type rpcResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
	ID      int64       `json:"id"`
}

// RunStdio processa JSON-RPC 2.0 do stdin e responde no stdout.
// Todo log vai para stderr para não corromper o canal MCP.
func RunStdio(port int) {
	log.SetOutput(os.Stderr)
	log.Printf("[crom-godot-mcp] iniciado (Godot esperado em ws://127.0.0.1:%d)", port)

	godot := NewGodotClient(port)
	writer := bufio.NewWriter(os.Stdout)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 1024*1024), 16*1024*1024)

	respond := func(r rpcResponse) {
		r.JSONRPC = "2.0"
		data, err := json.Marshal(r)
		if err != nil {
			log.Printf("[crom-godot-mcp] erro ao serializar: %v", err)
			return
		}
		writer.Write(data)
		writer.WriteByte('\n')
		writer.Flush()
	}

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var req rpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			log.Printf("[crom-godot-mcp] linha inválida: %v", err)
			continue
		}
		if req.ID == nil { // notificação: sem resposta
			continue
		}
		id := *req.ID

		switch req.Method {
		case "initialize":
			respond(rpcResponse{ID: id, Result: map[string]interface{}{
				"protocolVersion": protocolVersion,
				"capabilities":    map[string]interface{}{"tools": map[string]interface{}{}},
				"serverInfo":      map[string]string{"name": "crom-godot-mcp", "version": version},
			}})
		case "ping":
			respond(rpcResponse{ID: id, Result: map[string]interface{}{}})
		case "tools/list":
			respond(rpcResponse{ID: id, Result: map[string]interface{}{"tools": catalog}})
		case "tools/call":
			var p struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments"`
			}
			if err := json.Unmarshal(req.Params, &p); err != nil {
				respond(rpcResponse{ID: id, Error: &rpcError{Code: -32602, Message: "parâmetros inválidos: " + err.Error()}})
				continue
			}
			text := callTool(godot, p.Name, p.Arguments)
			respond(rpcResponse{ID: id, Result: map[string]interface{}{
				"content": []map[string]string{{"type": "text", "text": text}},
			}})
		default:
			respond(rpcResponse{ID: id, Error: &rpcError{Code: -32601, Message: "método não suportado: " + req.Method}})
		}
	}
	log.Printf("[crom-godot-mcp] stdin encerrado.")
}

// callTool encaminha a chamada ao editor (WS), com conexão lazy e 1 retry.
func callTool(godot *GodotClient, name string, args map[string]interface{}) string {
	known := false
	for _, t := range catalog {
		if t.Name == name {
			known = true
			break
		}
	}
	if !known {
		return fmt.Sprintf(`{"status":"error","message":"Ferramenta desconhecida: %s"}`, name)
	}
	action := strings.TrimPrefix(name, "godot_")
	if args == nil {
		args = map[string]interface{}{}
	}

	var resp map[string]interface{}
	var err error
	for attempt := 0; attempt < 2; attempt++ {
		if !godot.IsConnected() {
			if err = godot.Connect(); err != nil {
				continue
			}
		}
		resp, err = godot.SendCommand(action, args)
		if err == nil {
			break
		}
	}
	if err != nil {
		return fmt.Sprintf(`{"status":"error","message":"Sem conexão com o Editor Godot (o projeto está aberto com o plugin CromAI ativo?): %s"}`, strings.ReplaceAll(err.Error(), `"`, `'`))
	}

	// Screenshot: salva o PNG e devolve o caminho (evita estourar o contexto do LLM).
	if action == "capture_screenshot" {
		if b64, ok := resp["image_base64"].(string); ok {
			if path, e := saveScreenshot(b64); e == nil {
				return fmt.Sprintf(`{"status":"success","message":"Screenshot salvo em %s. Analise a imagem.","file_path":"%s"}`, path, path)
			}
		}
	}
	data, mErr := json.Marshal(resp)
	if mErr != nil {
		return fmt.Sprintf(`{"status":"error","message":"Falha ao serializar resposta: %s"}`, mErr.Error())
	}
	return string(data)
}

func saveScreenshot(b64 string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	dir := os.Getenv("GODOT_PROJECT_DIR")
	if dir == "" {
		dir = os.TempDir()
	} else {
		dir = filepath.Join(dir, ".crom", "screenshots")
	}
	_ = os.MkdirAll(dir, 0o755)
	path := filepath.Join(dir, fmt.Sprintf("godot_%s.png", time.Now().Format("20060102_150405")))
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return "", err
	}
	return path, nil
}
