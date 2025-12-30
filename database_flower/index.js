import express from "express";
import cors from "cors";
import morgan from "morgan";
import { init, Counter } from "./db.js";
import apiRouter from "./src/api/index.js";
import { errorHandler, notFoundHandler } from "./src/error.js";
import { ok, fail } from "./src/response.js";


const logger = morgan("tiny");

const app = express();
app.use(express.urlencoded({ extended: false }));
app.use(express.json());
app.use(cors());
app.use(logger);
app.use("/api", apiRouter);

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

const server = app.listen(port, async () => {
  console.log(`Web 服务已启动，正在监听端口 ${port}`);

  // 服务启动后，再去尝试连接数据库
  try {
    console.log("正在尝试连接数据库...");
    await init();
    console.log("数据库连接成功！");
  } catch (error) {
    // 这里的报错非常重要，要去云托管日志里看
    console.error("！！！严重错误：数据库连接失败！！！");
    console.error("错误详情:", error);
  }
});
