"use strict";
const common_vendor = require("../common/vendor.js");
const _sfc_main = {
  /**
   * 组件名称，也就是开发者使用的标签
   */
  name: "home",
  /**
   * 组件涉及的事件声明，只有声明过的事件，才能被正常发送
   */
  emits: [],
  /**
   * 属性声明，组件的使用者会传递这些属性值到组件
   */
  props: {
    navText: {
      type: String,
      default: ""
    }
  },
  /**
   * 组件内部变量声明
   */
  data() {
    return {
      menuButtonInfo: {},
      paddingLeft: 0,
      topBarHeight: 0
    };
  },
  /**
   * 属性变化监听器实现
   */
  watch: {},
  /**
   * 规则：如果没有配置expose，则methods中的方法均对外暴露，如果配置了expose，则以expose的配置为准向外暴露
   * ['publicMethod'] 含义为：只有 `publicMethod` 在实例上可用
   * 
   * 注意：如果在data中声明了一个变量，此时组件配置了 expose字段，但未在expose字段中包含此变量。会导致该变量被标记为`private`：仅能在组件内使用，不能在组件外访问
   */
  //expose: [''],
  methods: {
    /**
    * 内部使用的组件方法
    */
    //privateMethod() {}
  },
  /**
   * [可选实现] 组件被创建，组件第一个生命周期，
   * 在内存中被占用的时候被调用，开发者可以在这里执行一些需要提前执行的初始化逻辑
   */
  created() {
    const menuButtonInfo = common_vendor.wx$1.getMenuButtonBoundingClientRect();
    this.menuButtonInfo = menuButtonInfo;
    const systemInfo = common_vendor.index.getWindowInfo();
    this.paddingLeft = systemInfo.screenWidth - menuButtonInfo.right;
    const app = getApp();
    this.topBarHeight = app.globalData.topBarHeight;
  }
};
if (!Array) {
  const _easycom_uni_data_select2 = common_vendor.resolveComponent("uni-data-select");
  _easycom_uni_data_select2();
}
const _easycom_uni_data_select = () => "../uni_modules/uni-data-select/components/uni-data-select/uni-data-select.js";
if (!Math) {
  _easycom_uni_data_select();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {
    a: common_vendor.o(_ctx.change),
    b: common_vendor.o(($event) => _ctx.value = $event),
    c: common_vendor.p({
      localdata: _ctx.range,
      clear: false,
      modelValue: _ctx.value
    }),
    d: $data.menuButtonInfo.height + "px",
    e: $data.menuButtonInfo.height + "px",
    f: $data.menuButtonInfo.top + "px",
    g: $data.paddingLeft + "px",
    h: $data.topBarHeight + "px"
  };
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render], ["__scopeId", "data-v-045d88fd"]]);
wx.createComponent(Component);
//# sourceMappingURL=../../.sourcemap/mp-weixin/components/home.js.map
