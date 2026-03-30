#!/bin/bash

# MCP Bridge Backend Verification Script
# This script checks if the backend is properly configured for MCP

echo "🔍 MCP Bridge Backend Verification"
echo "==================================="
echo ""

# Check if we're in the right directory
if [ ! -f "../../go.mod" ]; then
    echo "❌ Error: Must run from mcp-bridge/test directory"
    echo "   cd mcp-bridge/test && bash verify.sh"
    exit 1
fi

echo "✅ Directory check passed"

# Check if backend can be built
echo ""
echo "🔨 Checking if backend compiles..."
cd ../../
if go build -o /tmp/crosspath-test ./cmd/server > /dev/null 2>&1; then
    echo "✅ Backend compiles successfully"
    rm -f /tmp/crosspath-test
else
    echo "❌ Backend compilation failed"
    echo "   Run: go build ./cmd/server"
    exit 1
fi

# Check if MCP tables exist in database schema
echo ""
echo "📋 Checking database schema..."
if grep -q "mcp_connections" internal/database/database.go; then
    echo "✅ mcp_connections table defined"
else
    echo "❌ mcp_connections table missing"
    exit 1
fi

if grep -q "mcp_tools" internal/database/database.go; then
    echo "✅ mcp_tools table defined"
else
    echo "❌ mcp_tools table missing"
    exit 1
fi

if grep -q "mcp_audit_log" internal/database/database.go; then
    echo "✅ mcp_audit_log table defined"
else
    echo "❌ mcp_audit_log table missing"
    exit 1
fi

# Check if MCP models exist
echo ""
echo "📦 Checking MCP models..."
if [ -f "internal/models/mcp.go" ]; then
    echo "✅ MCP models file exists"
else
    echo "❌ internal/models/mcp.go missing"
    exit 1
fi

# Check if MCP service exists
echo ""
echo "⚙️  Checking MCP service..."
if [ -f "internal/services/mcp_bridge_service.go" ]; then
    echo "✅ MCP bridge service exists"
else
    echo "❌ internal/services/mcp_bridge_service.go missing"
    exit 1
fi

# Check if MCP handler exists
echo ""
echo "🔌 Checking MCP handler..."
if [ -f "internal/handlers/mcp_websocket.go" ]; then
    echo "✅ MCP WebSocket handler exists"
else
    echo "❌ internal/handlers/mcp_websocket.go missing"
    exit 1
fi

# Check if tool registry has user tool support
echo ""
echo "🛠️  Checking tool registry extensions..."
if grep -q "RegisterUserTool" internal/tools/registry.go; then
    echo "✅ RegisterUserTool method exists"
else
    echo "❌ RegisterUserTool method missing"
    exit 1
fi

if grep -q "GetUserTools" internal/tools/registry.go; then
    echo "✅ GetUserTools method exists"
else
    echo "❌ GetUserTools method missing"
    exit 1
fi

# Check if main.go has MCP endpoint
echo ""
echo "🌐 Checking MCP endpoint..."
if grep -q "/mcp/connect" cmd/server/main.go; then
    echo "✅ /mcp/connect endpoint registered"
else
    echo "❌ /mcp/connect endpoint missing"
    exit 1
fi

# Check if mock client dependencies are available
echo ""
echo "📚 Checking mock client dependencies..."
cd mcp-bridge
if go list -m github.com/gorilla/websocket > /dev/null 2>&1; then
    echo "✅ gorilla/websocket dependency installed"
else
    echo "⚠️  gorilla/websocket not installed (run: go get github.com/gorilla/websocket)"
fi

# Final summary
echo ""
echo "=================================="
echo "✅ All checks passed!"
echo ""
echo "📝 Next steps:"
echo "   1. Start backend: cd backend && go run ./cmd/server"
echo "   2. Get JWT token from browser (see QUICKSTART.md)"
echo "   3. Run mock client: cd mcp-bridge/test && go run mock_client.go --token YOUR_TOKEN"
echo ""
echo "📖 Documentation:"
echo "   - Quick test: mcp-bridge/test/QUICKSTART.md"
echo "   - Full guide: mcp-bridge/test/README.md"
echo "   - Architecture: mcp-bridge/README.md"
echo ""
echo "🎉 Backend is ready for testing!"
