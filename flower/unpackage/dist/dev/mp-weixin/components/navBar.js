"use strict";
const common_vendor = require("../common/vendor.js");
const _sfc_main = {
  /**
   * 组件名称，也就是开发者使用的标签
   */
  name: "navBar",
  /**
   * 属性声明
   */
  props: {
    navText: {
      type: String,
      default: ""
    },
    isHome: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      menuButtonInfo: {},
      paddingLeft: 0
    };
  },
  methods: {
    handleBack() {
      common_vendor.wx$1.vibrateShort({ type: "medium" });
      if (this.isHome) {
        common_vendor.index.reLaunch({
          url: "/pages/index/index"
        });
      } else {
        common_vendor.index.navigateBack();
      }
    }
  },
  /**
   * [可选实现] 组件被创建，组件第一个生命周期，
   * 在内存中被占用的时候被调用，开发者可以在这里执行一些需要提前执行的初始化逻辑
   */
  created() {
    const menuButtonInfo = common_vendor.wx$1.getMenuButtonBoundingClientRect();
    this.menuButtonInfo = menuButtonInfo;
    const systemInfo = common_vendor.wx$1.getWindowInfo();
    this.paddingLeft = systemInfo.screenWidth - menuButtonInfo.right;
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  _easycom_uni_icons2();
}
const _easycom_uni_icons = () => "../uni_modules/uni-icons/components/uni-icons/uni-icons.js";
if (!Math) {
  _easycom_uni_icons();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: $props.isHome
  }, $props.isHome ? {
    b: common_vendor.p({
      type: "home",
      size: "20",
      color: "var(--primary-color)"
    })
  } : {}, {
    c: $data.menuButtonInfo.height + "px",
    d: $data.menuButtonInfo.height + "px",
    e: $data.menuButtonInfo.top + "px",
    f: $data.paddingLeft + "px",
    g: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args))
  });
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render], ["__scopeId", "data-v-3075e18a"]]);
wx.createComponent(Component);
//# sourceMappingURL=../../.sourcemap/mp-weixin/components/navBar.js.map
