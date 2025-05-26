mcp-build:
	go build -o docs-mcp mcp/main.go

m-agent:
	go run multi_agent/main.go