# MCP Transport

This module implements the MCP server in Kratos based on mcp-go.

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/go-kratos/kratos/contrib/transport/mcp/v2)

## Quick start
```go
import(
    tm "github.com/go-kratos/kratos/contrib/transport/mcp/v2"
    mcp "github.com/mark3labs/mcp-go/mcp"
)

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    name, ok := request.Params.Arguments["name"].(string)
    if !ok {
        return nil, errors.New("name must be a string")
    }
    return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}

func Health(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health/ready" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
    srv := tm.NewServer("kratos-mcp", "v1.0.0", tm.Address(":8000"), tm.Middleware(Health))
    tool := mcp.NewTool("hello_world",
        mcp.WithDescription("Say hello to someone"),
        mcp.WithString("name",
            mcp.Required(),
            mcp.Description("Name of the person to greet"),
        ),
    )
    // Add tool handler
    srv.AddTool(tool, helloHandler)
    // creates a kratos application
    app := kratos.New(
        kratos.Name("kratos-app"),
        kratos.Server(srv),
    )
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```
