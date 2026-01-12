"use strict";
Object.defineProperty(exports, Symbol.toStringTag, { value: "Module" });
const common_vendor = require("./common/vendor.js");
const utils_request = require("./utils/request.js");
if (!Math) {
  "./pages/index/index.js";
  "./pages/my/my.js";
  "./pages/system/system.js";
}
const _sfc_main = {
  globalData: {
    topBarHeight: 0,
    bottomSafeAreaHeight: 0
  },
  onLaunch: async function() {
    common_vendor.index.__f__("log", "at App.vue:9", "App Launch");
    const systemInfo = common_vendor.wx$1.getWindowInfo();
    const menuButtonInfo = common_vendor.wx$1.getMenuButtonBoundingClientRect();
    const statusBarHeight = systemInfo.statusBarHeight;
    const navBarHeight = (menuButtonInfo.top - statusBarHeight) * 2 + menuButtonInfo.height;
    const barHeight = statusBarHeight + navBarHeight;
    this.globalData.topBarHeight = barHeight;
    const bottomSafeAreaHeight = systemInfo.screenHeight - systemInfo.safeArea.bottom;
    this.globalData.bottomSafeAreaHeight = bottomSafeAreaHeight;
    common_vendor.wx$1.cloud.init();
    const user = await utils_request.callContainer("/api/login");
    common_vendor.index.__f__("log", "at App.vue:22", "callContainer login:", user);
    await new Promise((resolve) => {
      common_vendor.index.setStorage({ key: "family", data: user.data.family, success: resolve });
    });
    await new Promise((resolve) => {
      common_vendor.index.setStorage({ key: "userInfo", data: user.data.user, success: resolve });
    });
  },
  onShow: function() {
  },
  onHide: function() {
  }
};
function createApp() {
  const app = common_vendor.createSSRApp(_sfc_main);
  return {
    app
  };
}
createApp().app.mount("#app");
exports.createApp = createApp;
//# sourceMappingURL=../.sourcemap/mp-weixin/app.js.map
