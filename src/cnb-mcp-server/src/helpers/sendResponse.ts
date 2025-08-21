import type { Response } from 'express';

export function stopWithWrongTransport(res: Response) {
  res.status(400).json({
    jsonrpc: '2.0',
    error: {
      code: -32000,
      message: 'Bad Request: Session exists but uses a different transport protocol'
    },
    id: null
  });
}
