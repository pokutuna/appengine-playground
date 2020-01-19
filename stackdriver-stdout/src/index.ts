import express = require('express');
import { Request, Response } from 'express';

import { router as loggingRoute } from './logging';
import { router as errorReportingRoute } from './errorreporting';

const app = express();

app.get('/', (req: Request, res: Response) => {
  res.status(200).send('ok');
});

app.use('/logging', loggingRoute);
app.use('/errorreporting', errorReportingRoute);

const port = process.env.PORT || 3000;
app.listen(port, () => {
  console.log(`App listening on port ${port}`);
});
