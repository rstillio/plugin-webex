package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/mythingies/plugin-webex/internal/webex"
)

func registerDownloadAttachment(s *mcpserver.MCPServer, client *webex.Client) {
	tool := mcp.NewTool("download_attachment",
		mcp.WithDescription("Download a file attachment from a Webex message. Use this when a message includes attachments (indicated by [N attachment(s)] in message output). The file URL comes from the Webex API and requires authenticated access."),
		mcp.WithString("file_url",
			mcp.Required(),
			mcp.Description("The Webex file URL to download (from the message's files array)."),
		),
		mcp.WithString("dest_dir",
			mcp.Description("Directory to save the file to. Defaults to ~/Downloads."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fileURL, err := req.RequireString("file_url")
		if err != nil {
			return mcp.NewToolResultError("file_url is required"), nil
		}

		destDir := req.GetString("dest_dir", "")
		if destDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("cannot determine home directory: %v", err)), nil
			}
			destDir = filepath.Join(home, "Downloads")
		}

		localPath, contentType, err := client.DownloadAttachment(fileURL, destDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("download failed: %v", err)), nil
		}

		text := fmt.Sprintf("Downloaded to: %s\nContent-Type: %s", localPath, contentType)
		return mcp.NewToolResultText(text), nil
	})
}
