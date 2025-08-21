import type { CallToolResult } from '@modelcontextprotocol/sdk/types.js';

export function formatTextToolResult(text: string, toolName: string): CallToolResult {
  return {
    content: [
      {
        type: 'text',
        text
      }
    ]
  };
}

export function formatToolError(error: unknown, toolName: string): CallToolResult {
    return {
    content: [
        {
        type: 'text',
        text: `Error ${toolName}:\n${error instanceof Error ? error.message : String(error)}`
        }
    ],
    isError: true
    };
}
