// crom-godot-mcp — servidor MCP (Model Context Protocol) para o Editor Godot 4.
//
// Expõe as ferramentas do editor (via o plugin CromAI, WebSocket na porta 8080)
// para agentes de IA como o crom-agente, por JSON-RPC 2.0 em stdio.
//
//	crom-godot-mcp --mcp-stdio [--port 8080]
package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	stdio := flag.Bool("mcp-stdio", false, "Roda como servidor MCP stdio (JSON-RPC 2.0)")
	port := flag.Int("port", 8080, "Porta WebSocket do plugin CromAI no Editor Godot")
	showVersion := flag.Bool("version", false, "Mostra a versão e sai")
	flag.Parse()

	if *showVersion {
		fmt.Printf("crom-godot-mcp %s\n", version)
		return
	}
	if *stdio {
		RunStdio(*port)
		return
	}
	fmt.Fprintln(os.Stderr, "crom-godot-mcp: use --mcp-stdio para rodar como servidor MCP. Veja --help.")
	os.Exit(1)
}
