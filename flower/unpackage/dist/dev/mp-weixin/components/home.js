"use strict";
const common_vendor = require("../common/vendor.js");
const utils_request = require("../utils/request.js");
const WaterfallBox = () => "./WaterfallBox.js";
const _sfc_main = {
  components: {
    WaterfallBox
  },
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
      topBarHeight: 0,
      windowWidth: 0,
      familyRange: [],
      value: null,
      searchValue: "",
      tagList: [],
      currentTagIndex: 0,
      sliderLeft: 0,
      sliderWidth: 0,
      sliderTimer: null,
      plantsList: [],
      allPlantsList: []
    };
  },
  computed: {
    // 动态生成滑块样式
    sliderStyle() {
      return {
        transform: `translateX(${this.sliderLeft}px)`,
        width: `${this.sliderWidth}px`
      };
    }
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
    async loadFamilyData() {
      var _a;
      try {
        const familyList = await new Promise((resolve, reject) => {
          common_vendor.index.getStorage({ key: "family", success: resolve, fail: reject });
        });
        this.familyRange = [];
        (_a = familyList == null ? void 0 : familyList.data) == null ? void 0 : _a.forEach((item) => {
          this.familyRange = [...this.familyRange, { "text": item.name, "value": item.ID, "disable": false }];
        });
        const exists = this.familyRange.some((item) => item.value === this.value);
        if (!exists && this.familyRange.length > 0) {
          this.value = this.familyRange[0].value;
        }
        this.getTagList();
        this.getPlantsList();
      } catch (error) {
        common_vendor.index.__f__("error", "at components/home.vue:166", error);
      }
    },
    async getPlantsList() {
      try {
        const plants = await utils_request.callContainer("/api/plant/list", {
          "familyId": this.value
        });
        common_vendor.index.__f__("log", "at components/home.vue:174", "plants list:", plants);
        this.allPlantsList = plants == null ? void 0 : plants.data;
        this.filterPlants();
      } catch (error) {
        common_vendor.index.__f__("error", "at components/home.vue:178", error);
      }
    },
    filterPlants() {
      const currentTag = this.tagList[this.currentTagIndex];
      const tagId = currentTag ? currentTag.ID : 0;
      let filtered = [];
      if (tagId === 0) {
        filtered = [...this.allPlantsList];
      } else {
        filtered = this.allPlantsList.filter((plant) => {
          return plant.tags && plant.tags.some((t) => t.ID === tagId);
        });
      }
      this.plantsList = filtered.map((item) => {
        const newItem = { ...item };
        if (Array.isArray(item.tags)) {
          newItem.tags = [...item.tags];
        }
        return newItem;
      });
    },
    changeFamily(e) {
      common_vendor.index.__f__("log", "at components/home.vue:201", e);
    },
    async getTagList() {
      try {
        const tagList = await utils_request.callContainer("/api/tag/", {
          familyId: this.value
        });
        common_vendor.index.__f__("log", "at components/home.vue:208", "tagList:", tagList);
        const apiTags = (tagList == null ? void 0 : tagList.data) || [];
        this.tagList = [
          { name: "全部", ID: 0 },
          ...apiTags.map((item) => ({
            name: item.name,
            ID: item.ID,
            ...item
          }))
        ];
        common_vendor.index.__f__("log", "at components/home.vue:218", "tags:", this.tagList);
        this.$nextTick(() => {
          setTimeout(() => {
            this.updateSliderPosition(0);
          }, 200);
        });
      } catch (error) {
        common_vendor.index.__f__("error", "at components/home.vue:226", error);
      }
    },
    searchPlant(e) {
      common_vendor.index.__f__("log", "at components/home.vue:230", "e", e);
      common_vendor.index.__f__("log", "at components/home.vue:231", "search:", this.searchValue);
    },
    selectTag(index, item) {
      if (this.currentTagIndex === index)
        return;
      common_vendor.wx$1.vibrateShort({ type: "medium" });
      this.currentTagIndex = index;
      this.filterPlants();
      const query = common_vendor.index.createSelectorQuery().in(this);
      query.select("#tag-container").boundingClientRect();
      query.select("#tag-text-" + index).boundingClientRect();
      query.exec((res) => {
        if (res[0] && res[1]) {
          const containerLeft = res[0].left;
          const currentTextLeft = res[1].left;
          const currentTextWidth = res[1].width;
          const ratio = 22 / 18;
          const finalWidth = currentTextWidth * ratio;
          const widthDiff = finalWidth - currentTextWidth;
          const finalLeft = currentTextLeft - containerLeft - widthDiff / 2;
          this.sliderWidth = finalWidth;
          this.sliderLeft = finalLeft;
          if (this.sliderTimer)
            clearTimeout(this.sliderTimer);
          this.sliderTimer = setTimeout(() => {
            this.updateSliderPosition(index);
          }, 350);
        }
      });
    },
    updateSliderPosition(index) {
      const query = common_vendor.index.createSelectorQuery().in(this);
      query.select("#tag-container").boundingClientRect();
      query.select("#tag-text-" + index).boundingClientRect();
      query.exec((res) => {
        if (res[0] && res[1]) {
          const containerLeft = res[0].left;
          const textLeft = res[1].left;
          const textWidth = res[1].width;
          this.sliderLeft = textLeft - containerLeft;
          this.sliderWidth = textWidth;
        }
      });
    },
    onImgLoad(item) {
      item.isLoaded = true;
      if (this.allPlantsList && this.allPlantsList.length > 0) {
        const sourceItem = this.allPlantsList.find((i) => i.ID === item.ID);
        if (sourceItem) {
          sourceItem.isLoaded = true;
        }
      }
    },
    goAddPage() {
      common_vendor.wx$1.vibrateShort({ type: "medium" });
      common_vendor.index.navigateTo({ url: `/pages/addPlant/addPlant?familyId=${this.value}` });
    },
    gotoDetail(item) {
      common_vendor.index.navigateTo({
        url: `/pages/plantDetail/plantDetail?id=${item.ID}`
      });
    }
    /**
    * 内部使用的组件方法
    */
    //privateMethod() {}
  },
  /**
   * [可选实现] 组件被创建，组件第一个生命周期，
   * 在内存中被占用的时候被调用，开发者可以在这里执行一些需要提前执行的初始化逻辑
   */
  async created() {
    const menuButtonInfo = common_vendor.wx$1.getMenuButtonBoundingClientRect();
    this.menuButtonInfo = menuButtonInfo;
    const systemInfo = common_vendor.index.getWindowInfo();
    this.paddingLeft = systemInfo.screenWidth - menuButtonInfo.right;
    const app = getApp();
    this.topBarHeight = app.globalData.topBarHeight;
    this.windowWidth = app.globalData.windowWidth;
    const user = await utils_request.callContainer("/api/login");
    common_vendor.index.__f__("log", "at components/home.vue:316", "callContainer login:", user);
    await new Promise((resolve) => {
      common_vendor.index.setStorage({ key: "family", data: user.data.family, success: resolve });
    });
    await new Promise((resolve) => {
      common_vendor.index.setStorage({ key: "familyId", data: user.data.family[0].ID, success: resolve });
    });
    await new Promise((resolve) => {
      common_vendor.index.setStorage({ key: "userInfo", data: user.data.user, success: resolve });
    });
    this.loadFamilyData();
    common_vendor.index.$off("refreshHomeList");
    common_vendor.index.$on("refreshHomeList", (data) => {
      common_vendor.index.__f__("log", "at components/home.vue:329", "收到刷新通知", data);
      this.getPlantsList();
    });
  },
  beforeDestroy() {
    common_vendor.index.$off("refreshHomeList");
  }
};
if (!Array) {
  const _easycom_uni_data_select2 = common_vendor.resolveComponent("uni-data-select");
  const _easycom_uni_search_bar2 = common_vendor.resolveComponent("uni-search-bar");
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  const _component_WaterfallBox = common_vendor.resolveComponent("WaterfallBox");
  (_easycom_uni_data_select2 + _easycom_uni_search_bar2 + _easycom_uni_icons2 + _component_WaterfallBox)();
}
const _easycom_uni_data_select = () => "../uni_modules/uni-data-select/components/uni-data-select/uni-data-select.js";
const _easycom_uni_search_bar = () => "../uni_modules/uni-search-bar/components/uni-search-bar/uni-search-bar.js";
const _easycom_uni_icons = () => "../uni_modules/uni-icons/components/uni-icons/uni-icons.js";
if (!Math) {
  (_easycom_uni_data_select + _easycom_uni_search_bar + _easycom_uni_icons)();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {
    a: common_vendor.o($options.changeFamily),
    b: common_vendor.o(($event) => $data.value = $event),
    c: common_vendor.p({
      localdata: $data.familyRange,
      clear: false,
      modelValue: $data.value
    }),
    d: $data.menuButtonInfo.height / 2 + "px",
    e: $data.menuButtonInfo.top + "px",
    f: $data.paddingLeft + "px",
    g: $data.topBarHeight + "px",
    h: common_vendor.o($options.searchPlant),
    i: common_vendor.o(($event) => $data.searchValue = $event),
    j: common_vendor.p({
      placeholder: "输入植物名称",
      radius: "20",
      focus: false,
      bgColor: "rgba(255,255,255,0.5)",
      clearButton: "auto",
      cancelButton: "none",
      modelValue: $data.searchValue
    }),
    k: common_vendor.p({
      type: "plusempty",
      size: "22",
      color: "#333"
    }),
    l: common_vendor.o((...args) => $options.goAddPage && $options.goAddPage(...args)),
    m: common_vendor.f($data.tagList, (item, index, i0) => {
      return {
        a: common_vendor.t(item.name),
        b: "tag-text-" + index,
        c: item.ID,
        d: "tag-item-" + index,
        e: $data.currentTagIndex === index ? 1 : "",
        f: common_vendor.o(($event) => $options.selectTag(index, item), item.ID)
      };
    }),
    n: common_vendor.s($options.sliderStyle),
    o: "tag-item-" + ($data.currentTagIndex > 1 ? $data.currentTagIndex - 1 : 0),
    p: common_vendor.w(({
      item
    }, s0, i0) => {
      return common_vendor.e({
        a: item.cover.url,
        b: item.isLoaded ? 1 : "",
        c: common_vendor.o(($event) => $options.onImgLoad(item)),
        d: item.cover.height / item.cover.width * 100 + "%",
        e: common_vendor.t(item.name),
        f: item.tags
      }, item.tags ? {} : {}, {
        g: common_vendor.o(($event) => $options.gotoDetail(item)),
        h: i0,
        i: s0
      });
    }, {
      name: "item",
      path: "p",
      vueId: "045d88fd-3"
    }),
    q: common_vendor.p({
      list: $data.plantsList,
      idKey: "ID",
      cols: "2"
    })
  };
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render], ["__scopeId", "data-v-045d88fd"]]);
wx.createComponent(Component);
//# sourceMappingURL=../../.sourcemap/mp-weixin/components/home.js.map
