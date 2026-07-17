package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// GodotClient fala com o plugin CromAI (@tool) rodando dentro do Editor Godot,
// que expõe um servidor WebSocket (padrão :8080) e executa as ações na cena viva.
type GodotClient struct {
	url  string
	conn *websocket.Conn
	mu   sync.Mutex
}

func NewGodotClient(port int) *GodotClient {
	return &GodotClient{url: fmt.Sprintf("ws://127.0.0.1:%d", port)}
}

func (g *GodotClient) Connect() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	dialer := websocket.Dialer{HandshakeTimeout: 5 * time.Second}
	conn, _, err := dialer.Dial(g.url, nil)
	if err != nil {
		return fmt.Errorf("falha ao conectar no Godot (o projeto está aberto com o plugin CromAI ativo em %s?): %w", g.url, err)
	}
	g.conn = conn
	return nil
}

func (g *GodotClient) IsConnected() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.conn != nil
}

func (g *GodotClient) Close() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.conn != nil {
		_ = g.conn.Close()
		g.conn = nil
	}
}

// SendCommand envia {action, params} ao editor e aguarda a resposta JSON.
func (g *GodotClient) SendCommand(action string, params map[string]interface{}) (map[string]interface{}, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.conn == nil {
		return nil, fmt.Errorf("não conectado ao Godot")
	}
	payload, _ := json.Marshal(map[string]interface{}{"action": action, "params": params})
	if err := g.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		_ = g.conn.Close()
		g.conn = nil
		return nil, err
	}
	_ = g.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	_, data, err := g.conn.ReadMessage()
	if err != nil {
		_ = g.conn.Close()
		g.conn = nil
		return nil, fmt.Errorf("erro ao ler resposta do Godot: %w", err)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("resposta inválida do Godot: %w", err)
	}
	return resp, nil
}
