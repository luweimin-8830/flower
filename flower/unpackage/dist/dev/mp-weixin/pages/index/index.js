"use strict";
const common_vendor = require("../../common/vendor.js");
const common_assets = require("../../common/assets.js");
let flag;
const {
  windowWidth
} = common_vendor.index.getSystemInfoSync();
const navBar = () => "../../components/navBar.js";
const home = () => "../../components/home.js";
const my = () => "../../components/my.js";
const loadingPage = () => "../../components/loading.js";
const plantList = () => "../../components/plantList.js";
const logCalendar = () => "../../components/logCalendar.js";
const photoAlbum = () => "../../components/photoAlbum.js";
const _sfc_main = {
  components: {
    navBar,
    home,
    my,
    loadingPage,
    plantList,
    logCalendar,
    photoAlbum
  },
  data() {
    return {
      animation01: false,
      animation02: false,
      ballStyleLeft: 0,
      liquidStyleLeft: 0,
      onKey: "c",
      list: {
        a: {
          is: false,
          title: "01",
          iconOff: common_assets.a1,
          iconOn: common_assets.a2
        },
        b: {
          is: false,
          title: "02",
          iconOff: common_assets.b1,
          iconOn: common_assets.b2
        },
        c: {
          is: true,
          title: "03",
          iconOff: common_assets.c1,
          iconOn: common_assets._imports_0
        },
        d: {
          is: false,
          title: "04",
          iconOff: common_assets.d1,
          iconOn: common_assets.d2
        },
        e: {
          is: false,
          title: "05",
          iconOff: common_assets.e1,
          iconOn: common_assets.e2
        }
      }
    };
  },
  computed: {
    ballStyle() {
      let style = {
        left: `${this.ballStyleLeft}px`
      };
      return style;
    },
    liquidStyle() {
      let style = {
        left: `${this.liquidStyleLeft}px`
      };
      return style;
    }
  },
  methods: {
    tabbarPageScrollLower() {
    },
    init(e) {
      common_vendor.index.__f__("log", "at pages/index/index.vue:175", "Tabbar init elements:", e);
      if (!e || e.length === 0)
        return;
      const keys = Object.keys(this.list);
      e.forEach((item2, index) => {
        const key = keys[index];
        if (key && this.list[key]) {
          item2.left = item2.left + 20;
          this.list[key]["style"] = item2;
        }
      });
      this.onKey = "c";
      const item = this.list[this.onKey];
      if (item && item.style) {
        this.onTabbar(item, this.onKey);
      } else {
        common_vendor.index.__f__("warn", "at pages/index/index.vue:208", "初始化匹配失败，尝试使用默认值");
        const defaultWidth = windowWidth / 5;
        this.ballStyleLeft = defaultWidth * 2 + defaultWidth / 2 - 22;
        this.liquidStyleLeft = defaultWidth * 2 + defaultWidth / 2 - windowWidth / 2;
      }
    },
    onTabbar(item, key) {
      this.throttle(() => {
        if (!item || !item.style) {
          common_vendor.index.__f__("error", "at pages/index/index.vue:218", "onTabbar error: item or item.style is undefined");
          return;
        }
        this.switchTabbarPage(key);
        this.onKey = key;
        let left = item.style.left;
        this.ballStyleLeft = left - 22;
        this.animation01 = true;
        this.animation02 = true;
        setTimeout(() => {
          this.liquidStyleLeft = left - windowWidth / 2;
          this.animation01 = false;
        }, 300);
        setTimeout(() => {
          this.animation02 = false;
        }, 610);
      }, 610);
    },
    switchTabbarPage(key) {
      if (this.onKey === key)
        return;
      const selectPageFlag = this.list[key]["is"];
      if (selectPageFlag === void 0) {
        return;
      }
      if (selectPageFlag === false) {
        this.list[key]["is"] = true;
      }
    },
    // 获取元素位置
    getDemRefAll(e) {
      const {
        selector,
        success = () => {
        }
      } = e;
      return new Promise((r) => {
        let view = common_vendor.index.createSelectorQuery().in(this).selectAll(selector);
        view.boundingClientRect((ref) => {
          success(ref);
          r(ref);
        }).exec();
      });
    },
    // 节流函数
    throttle(func, wait = 500, immediate = true) {
      if (immediate) {
        if (!flag) {
          flag = true;
          typeof func === "function" && func();
          setTimeout(() => {
            flag = false;
          }, wait);
        }
      } else if (!flag) {
        flag = true;
        setTimeout(() => {
          flag = false;
          typeof func === "function" && func();
        }, wait);
      }
    }
  },
  async onLoad() {
    const app = getApp();
    await app.globalData.initPromise;
    this.$refs.loading.open();
    setTimeout(() => {
      this.$refs.loading.close();
    }, 1500);
    this.$nextTick(() => {
      this.getDemRefAll({
        selector: ".pan-tabbar-item",
        success: (e) => {
          this.init(e);
        }
      });
    });
    common_vendor.index.$on("switchTab", (key) => {
      if (this.list[key]) {
        const item = this.list[key];
        if (item.style) {
          this.onTabbar(item, key);
        } else {
          this.onKey = key;
          this.switchTabbarPage(key);
        }
      }
    });
  },
  onShow() {
    this.$nextTick(() => {
      if (this.$refs.homeComponent)
        this.$refs.homeComponent.onPageShow();
      if (this.$refs.calendarComponent)
        this.$refs.calendarComponent.initData();
    });
  },
  onUnload() {
    common_vendor.index.$off("switchTab");
  },
  onShareAppMessage() {
    return {
      title: "分享我的植物",
      path: "/pages/index/index"
    };
  }
};
if (!Array) {
  const _component_loadingPage = common_vendor.resolveComponent("loadingPage");
  const _component_photoAlbum = common_vendor.resolveComponent("photoAlbum");
  const _component_plantList = common_vendor.resolveComponent("plantList");
  const _component_home = common_vendor.resolveComponent("home");
  const _component_logCalendar = common_vendor.resolveComponent("logCalendar");
  const _component_my = common_vendor.resolveComponent("my");
  (_component_loadingPage + _component_photoAlbum + _component_plantList + _component_home + _component_logCalendar + _component_my)();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: common_vendor.sr("loading", "24547230-0"),
    b: $data.list["a"].is
  }, $data.list["a"].is ? {
    c: common_vendor.o((...args) => $options.tabbarPageScrollLower && $options.tabbarPageScrollLower(...args)),
    d: $data.onKey === "a" ? "" : "none"
  } : {}, {
    e: $data.list["b"].is
  }, $data.list["b"].is ? {
    f: common_vendor.o((...args) => $options.tabbarPageScrollLower && $options.tabbarPageScrollLower(...args)),
    g: $data.onKey === "b" ? "" : "none"
  } : {}, {
    h: $data.list["c"].is
  }, $data.list["c"].is ? {
    i: common_vendor.sr("homeComponent", "24547230-3"),
    j: common_vendor.o((...args) => $options.tabbarPageScrollLower && $options.tabbarPageScrollLower(...args)),
    k: $data.onKey === "c" ? "" : "none"
  } : {}, {
    l: $data.list["d"].is
  }, $data.list["d"].is ? {
    m: common_vendor.sr("calendarComponent", "24547230-4"),
    n: common_vendor.o((...args) => $options.tabbarPageScrollLower && $options.tabbarPageScrollLower(...args)),
    o: $data.onKey === "d" ? "" : "none"
  } : {}, {
    p: $data.list["e"].is
  }, $data.list["e"].is ? {
    q: common_vendor.o((...args) => $options.tabbarPageScrollLower && $options.tabbarPageScrollLower(...args)),
    r: $data.onKey === "e" ? "" : "none"
  } : {}, {
    s: $data.animation01 ? 1 : "",
    t: common_vendor.s($options.liquidStyle),
    v: $data.animation02 ? 1 : "",
    w: common_vendor.s($options.ballStyle),
    x: common_vendor.f($data.list, (item, key, i0) => {
      return {
        a: item.iconOff,
        b: item.iconOn,
        c: key,
        d: key,
        e: key === $data.onKey ? 1 : "",
        f: key,
        g: common_vendor.o(($event) => $options.onTabbar(item, key), key)
      };
    })
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
_sfc_main.__runtimeHooks = 2;
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/index/index.js.map
