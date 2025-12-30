import { Sequelize, DataTypes } from "sequelize";
import process from "process";

// 从环境变量中读取数据库配置
const { MYSQL_USERNAME, MYSQL_PASSWORD, MYSQL_ADDRESS = "" } = process.env;

const [host, port] = MYSQL_ADDRESS.split(":");

const sequelize = new Sequelize("nodejs_demo", MYSQL_USERNAME, MYSQL_PASSWORD, {
  host,
  port,
  dialect: "mysql" /* one of 'mysql' | 'mariadb' | 'postgres' | 'mssql' */,
  logging: false,
  timezone: '+08:00'
});

// 定义数据模型
const Counter = sequelize.define("Counter", {
  count: {
    type: DataTypes.INTEGER,
    allowNull: false,
    defaultValue: 1, 
  },
});

const User = sequelize.define("User",{
  openId: {
    type: DataTypes.STRING,
    allowNull: false, //不允许为空
    unique: true //唯一
  },
  name:{
    type: DataTypes.STRING
  },
  isAdmin:{
    type:DataTypes.BOOLEAN,
    defaultValue:false //默认值
  },
  age:{
    type:DataTypes.INTEGER
  },
  lastLogin:{
    type: DataTypes.DATE
  }
})

// 数据库初始化方法
async function init() {
  await Counter.sync({ alter: true });
}

// 导出初始化方法和模型
export {
  init,
  Counter,
  User,
};
