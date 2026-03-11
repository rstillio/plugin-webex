package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/mythingies/plugin-webex/internal/webex"
)

func registerGetMeetingTranscript(s *mcpserver.MCPServer, client *webex.Client) {
	tool := mcp.NewTool("get_meeting_transcript",
		mcp.WithDescription("Pull the transcript from a past Webex meeting. If meeting_id is omitted, lists all available transcripts without downloading."),
		mcp.WithString("meeting_id",
			mcp.Description("The instance ID of the meeting (use list_meetings with meeting_type=\"meeting\" to get instance IDs). Omit to list all available transcripts."),
		),
		mcp.WithString("format",
			mcp.Description("Transcript format: \"txt\" (plain text) or \"vtt\" (WebVTT with timestamps). Default: txt."),
			mcp.Enum("txt", "vtt"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		meetingID := req.GetString("meeting_id", "")
		format := req.GetString("format", "txt")

		// Find transcripts (for a specific meeting or all available).
		transcripts, err := client.ListTranscripts(meetingID, 10)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to list transcripts: %v", err)), nil
		}

		if len(transcripts) == 0 {
			if meetingID == "" {
				return mcp.NewToolResultText("No transcripts available in the system."), nil
			}
			return mcp.NewToolResultText("No transcript available for this meeting."), nil
		}

		// If no meeting_id, list available transcripts without downloading.
		if meetingID == "" {
			text := fmt.Sprintf("%d transcript(s) available:\n\n", len(transcripts))
			for _, t := range transcripts {
				text += fmt.Sprintf("- **%s** (started: %s, meeting_id: %s, transcript_id: %s)\n",
					t.MeetingTopic, t.StartTime, t.MeetingID, t.ID)
			}
			return mcp.NewToolResultText(text), nil
		}

		// Download the first available transcript.
		transcript := transcripts[0]
		content, err := client.DownloadTranscript(transcript.ID, format)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to download transcript: %v", err)), nil
		}

		header := fmt.Sprintf("**Transcript: %s** (started: %s, format: %s)\n\n",
			transcript.MeetingTopic, transcript.StartTime, format)

		return mcp.NewToolResultText(header + content), nil
	})
}
