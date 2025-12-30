const express = require("express");
const cors = require("cors");
const morgan = require("morgan");
const { init: initDB, Counter } = require("./db");
const apiRouter = require("./src/api/index.js");
import { errorHandler,notFoundHandler } from "./src/error.js";
import { ok,fail } from "./src/response.js";


const logger = morgan("tiny");

const app = express();
app.use(express.urlencoded({ extended: false }));
app.use(express.json());
app.use(cors());
app.use(logger);
app.use("/api",apiRouter);

// 首页
app.get("/", async (req, res) => {
  res.send("Hello");
});

// 小程序调用，获取微信 Open ID
app.get("/api/wx_openid", async (req, res) => {
  if (req.headers["x-wx-source"]) {
    res.send(req.headers["x-wx-openid"]);
  }
});

app.use(notFoundHandler);//404
app.use(errorHandler);//500

const port = process.env.PORT || 80;

async function bootstrap() {
  await initDB();
  app.listen(port, () => {
    console.log("启动成功", port);
  });
}

bootstrap();
