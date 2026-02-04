"use strict";
Object.defineProperty(exports, Symbol.toStringTag, { value: "Module" });
const common_vendor = require("./common/vendor.js");
if (!Math) {
  "./pages/index/index.js";
  "./pages/my/my.js";
  "./pages/system/system.js";
  "./pages/plantEdit/plantEdit.js";
  "./pages/plantDetail/plantDetail.js";
  "./pages/tagEdit/tagEdit.js";
  "./pages/familyDetail/familyDetail.js";
  "./pages/shareMember/shareMember.js";
  "./pages/careEdit/careEdit.js";
  "./pages/logEdit/logEdit.js";
  "./pages/shareDetail/shareDetail.js";
  "./pages/displayRoom/displayRoom.js";
}
const _sfc_main = {
  globalData: {
    topBarHeight: 0,
    bottomSafeAreaHeight: 0,
    windowWidth: 0,
    windowHeight: 0
  },
  onLaunch: async function() {
    common_vendor.index.__f__("log", "at App.vue:11", "App Launch");
    if (!common_vendor.wx$1.cloud) {
      common_vendor.index.__f__("error", "at App.vue:13", "请使用 2.2.3 或以上的基础库以使用云能力");
    } else {
      common_vendor.wx$1.cloud.init({
        env: "prod-0gr2o3qpe533f1fb",
        traceUser: true
      });
    }
    const systemInfo = common_vendor.wx$1.getWindowInfo();
    this.globalData.windowWidth = systemInfo.windowWidth;
    this.globalData.windowHeight = systemInfo.windowHeight;
    const menuButtonInfo = common_vendor.wx$1.getMenuButtonBoundingClientRect();
    const statusBarHeight = systemInfo.statusBarHeight;
    const navBarHeight = (menuButtonInfo.top - statusBarHeight) * 2 + menuButtonInfo.height;
    const barHeight = statusBarHeight + navBarHeight;
    this.globalData.topBarHeight = barHeight;
    const bottomSafeAreaHeight = systemInfo.screenHeight - systemInfo.safeArea.bottom;
    this.globalData.bottomSafeAreaHeight = bottomSafeAreaHeight;
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
