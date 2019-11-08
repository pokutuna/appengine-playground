const express = require("express");
const winston = require("winston");
const logginWinston = require("@google-cloud/logging-winston");

const main = async () => {
  const app = express();

  const logger = winston.createLogger();
  const middleware = await logginWinston.express.makeMiddleware(logger);

  app.use(middleware);

  app.get("/", (req, res) => {
    req.log.info("hello");
    return res.send("hello").end();
  });

  let counter = 0;
  app.get("/counter", (req, res) => {
    counter += 1;
    req.log.info("counter", { count: counter });
    return res.send(`${counter}`).end();
  });

  const PORT = process.env.PORT || 8080;
  app.listen(PORT, () => {
    console.log(`App listening on port ${PORT}`);
    console.log("Press Ctrl+C to quit.");
  });
};

main();
