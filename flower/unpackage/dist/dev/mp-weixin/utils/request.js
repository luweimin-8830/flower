"use strict";
const common_vendor = require("../common/vendor.js");
const CLOUD_ENV = "prod-0gr2o3qpe533f1fb";
const SERVICE_NAME = "golang-17j9";
const callContainer = (path, data = {}, method = "POST") => {
  return new Promise((resolve, reject) => {
    const formatPath = path.startsWith("/") ? path : `/${path}`;
    common_vendor.wx$1.cloud.callContainer({
      config: {
        env: CLOUD_ENV
      },
      path: formatPath,
      method,
      header: {
        "X-WX-SERVICE": SERVICE_NAME,
        "content-type": "application/json"
      },
      data,
      success: (res) => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data);
        } else {
          reject(res);
        }
      },
      fail: (err) => {
        common_vendor.index.__f__("error", "at utils/request.js:37", `请求 ${path} 失败`, err);
        reject(err);
      }
    });
  });
};
exports.callContainer = callContainer;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/request.js.map
