package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/mythingies/plugin-webex/internal/webex"
)

func registerGetRecordingTranscript(s *mcpserver.MCPServer, client *webex.Client) {
	tool := mcp.NewTool("get_recording_transcript",
		mcp.WithDescription("Get the transcript from a Webex recording. Works for recordings you hosted OR that were shared with you. Use list_recordings to find recording IDs."),
		mcp.WithString("recording_id",
			mcp.Description("The recording ID (from list_recordings)."),
			mcp.Required(),
		),
		mcp.WithBoolean("download",
			mcp.Description("If true, download and return the full transcript content. If false (default), return recording details and download links only."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		recordingID := req.GetString("recording_id", "")
		if recordingID == "" {
			return mcp.NewToolResultError("recording_id is required"), nil
		}
		download := req.GetBool("download", false)

		details, err := client.GetRecordingDetails(recordingID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get recording details: %v", err)), nil
		}

		text := fmt.Sprintf("**%s**\n", details.Topic)
		text += fmt.Sprintf("  Recording ID: %s\n", details.ID)
		text += fmt.Sprintf("  Meeting ID: %s\n", details.MeetingID)
		text += fmt.Sprintf("  Host: %s\n", details.HostEmail)
		text += fmt.Sprintf("  Recorded: %s\n", details.TimeRecorded)
		text += fmt.Sprintf("  Status: %s\n", details.Status)

		if details.TemporaryDirectDownloadLinks == nil {
			text += "\nNo download links available for this recording.\n"
			return mcp.NewToolResultText(text), nil
		}

		links := details.TemporaryDirectDownloadLinks
		text += fmt.Sprintf("  Download links expire: %s\n", links.Expiration)

		hasTranscript := links.TranscriptDownloadLink != ""
		text += fmt.Sprintf("  Transcript available: %v\n", hasTranscript)

		if !download {
			return mcp.NewToolResultText(text), nil
		}

		if !hasTranscript {
			text += "\nNo transcript download link available for this recording.\n"
			return mcp.NewToolResultText(text), nil
		}

		content, err := client.DownloadRecordingTranscript(links.TranscriptDownloadLink)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to download transcript: %v", err)), nil
		}

		text += fmt.Sprintf("\n---\n\n%s", content)
		return mcp.NewToolResultText(text), nil
	})
}
