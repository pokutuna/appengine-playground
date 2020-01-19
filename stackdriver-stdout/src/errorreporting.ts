import express = require('express');
import { Request, Response } from 'express';

export const router = express.Router();

const projectId = 'pokutuna-playground';

// 1. throw Error object (collected as a error)
router.get('/raw', (req: Request, res: Response) => {
  throw new Error('jsut throw a error!');
});

// 2. print trace (NOT collected)
// I guess the stacktrace ErrorReporting collectable must start with "Error: "
router.get('/trace', (req: Request, res: Response) => {
  console.trace('console.trace()!');
  res.status(500).send('trace');
});

// 3-1. print stacktrace to stdout (collected)
router.get('/stacktrace/stdout', (req: Request, res: Response) => {
  const err = new Error('error to write stdout');
  console.log(err.stack);
  res.status(500).send(err);
});

// 3-2. print stacktrace to stderr (collected)
router.get('/stacktrace/stderr', (req: Request, res: Response) => {
  const err = new Error('error to write stderr');
  console.error(err.stack);
  res.status(500).send(err);
});

// https://cloud.google.com/error-reporting/docs/formatting-error-messages?hl=ja
// https://cloud.google.com/error-reporting/reference/rest/v1beta1/projects.events/report?hl=ja#ReportedErrorEvent

// 4. print a json with @type field (collected)
router.get('/formatting/json', (req: Request, res: Response) => {
  // stack_trace, exception, message
  const prop = req.query.prop || 'message';

  const err = new Error('error to write with a json');
  const payload: any = {
    '@type':
      'type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent',
    [prop]: err.stack,
  };
  console.error(JSON.stringify(payload));
  res.status(500).json(payload);
});

// 4-ex. print a json with @type field without stacktrace (collected)
// On the console of ErrorReporting, this causes a error but "No parsed stack trace available"
router.get('/formatting/json/without-stacktrace', (req: Request, res: Response) => {
  // stack_trace, exception, message
  const prop = req.query.prop || 'message';

  const payload: any = {
    '@type': 'type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent',
    [prop]: "this is not a stacktrace",
    headers: req.headers,
  };
  console.error(JSON.stringify(payload));
  res.status(500).json(payload);
});

// 5. print a json with logName field (collected)
router.get('/formatting/logName', (req: Request, res: Response) => {
  const prop = req.query.prop || 'message';

  const err = new Error('error to write with logName field');
  const payload: any = {
    logName: `projects/${projectId}/clouderrorreporting.googleapis.com%2Freported_errors`,
    [prop]: err.stack,
  };
  console.error(JSON.stringify(payload));
  res.status(500).json(payload);
});
