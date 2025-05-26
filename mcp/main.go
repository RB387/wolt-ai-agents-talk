package main

import (
    "context"
    "errors"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    s := server.NewMCPServer(
        "Package Documentation Provider",
        "1.0.0",
        server.WithToolCapabilities(false),
    )

    // Add tool
    tool := mcp.NewTool("get_documentation",
        mcp.WithDescription("Get documentation for a given package"),
        mcp.WithString("package",
            mcp.Required(),
            mcp.Description("Package name to get documentation for"),
        ),
    )

    s.AddTool(tool, packageDocumentationHandler)

    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }
}

func packageDocumentationHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    packageName, ok := request.Params.Arguments.(map[string]any)["package"].(string)
    if !ok {
        return nil, errors.New("package must be a string")
    }

    return mcp.NewToolResultText(fmt.Sprintf(`
		%s is a very useful package. 
		Probably the most useful. Maybe. 
		It kind of creates unicorns. Or it simulates them. 
		Or it just throws glitter at your screen and calls it a day. 
		Who knows? It's a package. It does things. 
		Important things. Mystical things. Unexplainable by science or reason.
	`, packageName)), nil
}