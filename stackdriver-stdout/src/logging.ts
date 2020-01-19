import express = require('express');
import { Request, Response } from 'express';

export const router = express.Router();

const projectId = 'pokutuna-playground';

function getMessage(req: Request): any {
  return {
    headers: req.headers,
    ary: [1, 2, 3],
  };
}

// 1. console.log(obj)
// this outputs multiple log entries as textPayload each.
//   textPayload: "{ headers:"
//   textPayload: "  { host: 'stackdriver-stdout-dot-..."
router.get('/raw', (req: Request, res: Response) => {
  const payload = getMessage(req);
  console.log(payload);
  res.json(payload);
});

// 2. console.log(JSON.stringify(obj))
// this outputs as a single log entry with jsonPayload.
//   jsonPayload: { ary:[3], headers: { ... } }
router.get('/json', (req: Request, res: Response) => {
  const payload = getMessage(req);
  console.log(JSON.stringify(payload));
  res.json(payload);
});

// 3. special fields
// https://cloud.google.com/logging/docs/agent/configuration#special-fields

function getTraceFields(req: Request): any {
  const context = req.header('X-Cloud-Trace-Context') || '';
  const [traceId, spanId] = context.split('/');

  return {
    'logging.googleapis.com/trace': `projects/${projectId}/traces/${traceId}`,
    'logging.googleapis.com/spanId': spanId,
  };
}

// 3-1. log with trace
// this.outputs as a single log entry with jsonPayload & bundled with a parent httpReqeuest entry.
// trace fields are removed from jsonPayload
//   jsonPayload: { ary: [3], headers: { ... } }
//   spanId: ...
//   trace: ...
router.get('/special-fields/trace', (req: Request, res: Response) => {
  const payload = {
    ...getMessage(req),
    ...getTraceFields(req),
  };
  console.log(JSON.stringify(payload));
  res.json(payload);
});

// 3-2. log with severity
// this outputs a log entry marked with "i" icon
//   jsonPayload: { ary: [3], headers: { ... } }
//   severity: "INFO"
router.get('/special-fields/severity', (req: Request, res: Response) => {
  const payload = {
    severity: 'INFO',
    ...getMessage(req),
  };
  console.log(JSON.stringify(payload));
  res.json(payload);
});

// 4. console.error(JSON.stringify(obj))
router.get('/stderr', (req: Request, res: Response) => {
  const payload = getMessage(req);
  console.error(JSON.stringify(payload));
  res.json(payload);
});
