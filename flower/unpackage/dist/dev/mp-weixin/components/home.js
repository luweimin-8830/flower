"use strict";
const common_vendor = require("../common/vendor.js");
const utils_request = require("../utils/request.js");
const common_assets = require("../common/assets.js");
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
      allPlantsList: [],
      loadedImagesMap: {},
      // 用于追踪图片加载状态
      isFirstLoad: true,
      isSelecting: false,
      currentFamilyIndex: 0,
      isFiltering: false
      // 防止过滤期间的并发修改
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
  watch: {
    allPlantsList: {
      handler(newVal, oldVal) {
        if (newVal) {
          newVal.forEach((plant, idx) => {
            if (plant.tags) {
              plant.tags.forEach((tag) => {
              });
            }
          });
        }
      },
      deep: true
    }
  },
  /**
   * 规则：如果没有配置expose，则methods中的方法均对外暴露，如果配置了expose，则以expose的配置为准向外暴露
   * ['publicMethod'] 含义为：只有 `publicMethod` 在实例上可用
   * 
   * 注意：如果在data中声明了一个变量，此时组件配置了 expose字段，但未在expose字段中包含此变量。会导致该变量被标记为`private`：仅能在组件内使用，不能在组件外访问
   */
  //expose: [''],
  methods: {
    async loadFamilyData() {
      try {
        const [familyListResult, familyIdResult] = await Promise.all([
          new Promise((resolve, reject) => {
            common_vendor.index.getStorage({ key: "family", success: resolve, fail: reject });
          }),
          new Promise((resolve, reject) => {
            common_vendor.index.getStorage({ key: "familyId", success: resolve, fail: reject });
          })
        ]);
        const familyList = (familyListResult == null ? void 0 : familyListResult.data) || [];
        const cachedFamilyId = familyIdResult == null ? void 0 : familyIdResult.data;
        common_vendor.index.__f__("log", "at components/home.vue:193", "loadFamilyData - 家庭列表:", familyList, "缓存的家庭ID:", cachedFamilyId);
        if (familyList && Array.isArray(familyList) && familyList.length > 0) {
          this.familyRange = familyList.map((item) => ({
            text: item.name,
            value: item.id,
            disable: false
          }));
          this.value = cachedFamilyId || this.familyRange[0].value;
          this.currentFamilyIndex = this.familyRange.findIndex((item) => item.value === this.value);
          if (this.currentFamilyIndex === -1) {
            this.currentFamilyIndex = 0;
            this.value = this.familyRange[0].value;
          }
        } else {
          this.familyRange = [];
          this.currentFamilyIndex = 0;
          this.value = null;
        }
        await this.$nextTick();
        await this.getTagList();
        await this.getPlantsList();
      } catch (error) {
        common_vendor.index.__f__("error", "at components/home.vue:225", "加载家庭数据失败:", error);
      }
    },
    async refreshFamilyList() {
      try {
        const user = await utils_request.callContainer("/api/login");
        const familyList = user.data.family || [];
        await new Promise((resolve) => {
          common_vendor.index.setStorage({ key: "family", data: familyList, success: resolve });
        });
        if (familyList && Array.isArray(familyList) && familyList.length > 0) {
          this.familyRange = familyList.map((item) => ({
            text: item.name,
            value: item.id,
            disable: false
          }));
          const currentFamilyExists = this.familyRange.some((item) => item.value === this.value);
          if (!currentFamilyExists && this.familyRange.length > 0) {
            this.value = this.familyRange[0].value;
            this.currentFamilyIndex = 0;
            await new Promise((resolve) => {
              common_vendor.index.setStorage({ key: "familyId", data: this.value, success: resolve });
            });
            await this.getTagList();
            await this.getPlantsList();
          }
        } else {
          this.familyRange = [];
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at components/home.vue:267", "刷新家庭列表失败:", error);
      }
    },
    async getPlantsList() {
      const familyId = this.value;
      try {
        const plants = await utils_request.callContainer("/api/plant/list", {
          "familyId": familyId
        });
        common_vendor.index.__f__("log", "at components/home.vue:277", "plants list:", plants);
        const newData = (plants == null ? void 0 : plants.data) || [];
        newData.forEach((item, idx) => {
          if (item.tags && item.tags.length > 0) {
          }
        });
        this.allPlantsList = newData.map((item) => {
          let frozenTags = null;
          if (item.tags && Array.isArray(item.tags)) {
            frozenTags = item.tags.map((tag) => ({ ...tag }));
            frozenTags.forEach((tag) => Object.freeze(tag));
            Object.freeze(frozenTags);
          }
          const newItem = { ...item, tags: frozenTags };
          return Object.freeze(newItem);
        });
        this.filterPlants();
      } catch (error) {
        common_vendor.index.__f__("error", "at components/home.vue:304", "获取植物列表失败:", error);
      }
    },
    filterPlants() {
      const currentTag = this.tagList[this.currentTagIndex];
      const tagId = currentTag ? currentTag.id : 0;
      let filtered = [];
      if (tagId === 0) {
        filtered = this.allPlantsList;
      } else {
        filtered = this.allPlantsList.filter((plant) => {
          return plant.tags && plant.tags.some((t) => t.id === tagId);
        });
      }
      this.plantsList = filtered;
    },
    async handleFamilyChange(e) {
      const selectedIndex = e.detail.value;
      this.currentFamilyIndex = selectedIndex;
      const newFamilyId = this.familyRange[selectedIndex].value;
      try {
        await utils_request.callContainer("/api/family/switch", {
          familyId: newFamilyId
        });
        common_vendor.index.__f__("log", "at components/home.vue:332", "家庭切换成功");
        await new Promise((resolve) => {
          common_vendor.index.setStorage({ key: "familyId", data: newFamilyId, success: resolve });
        });
        common_vendor.index.$emit("familyChanged", newFamilyId);
      } catch (error) {
        common_vendor.index.__f__("error", "at components/home.vue:343", "切换家庭失败:", error);
        const errorMsg = (error == null ? void 0 : error.msg) || (error == null ? void 0 : error.message) || "切换家庭失败，请稍后重试";
        common_vendor.index.showToast({
          title: errorMsg,
          icon: "none",
          duration: 2e3
        });
        this.currentFamilyIndex = this.familyRange.findIndex((item) => item.value === this.value);
        await this.refreshFamilyList();
        return;
      }
      this.value = newFamilyId;
      this.currentTagIndex = 0;
      this.tagList = [];
      this.allPlantsList = [];
      this.plantsList = [];
      this.loadedImagesMap = {};
      await this.$nextTick();
      await this.getTagList();
      await this.getPlantsList();
      this.$nextTick(() => {
        setTimeout(() => {
          this.updateSliderPosition(0);
        }, 200);
      });
      common_vendor.wx$1.vibrateShort({ type: "light" });
    },
    toggleFamilySelect() {
      common_vendor.index.__f__("log", "at components/home.vue:390", "触发家庭选择器");
    },
    onTouchStart() {
      this.isSelecting = true;
    },
    onTouchEnd() {
      setTimeout(() => {
        this.isSelecting = false;
      }, 200);
    },
    async getTagList() {
      const familyId = this.value;
      try {
        const tagList = await utils_request.callContainer("/api/tag/", {
          familyId
        });
        common_vendor.index.__f__("log", "at components/home.vue:409", "tagList:", tagList);
        const apiTags = (tagList == null ? void 0 : tagList.data) || [];
        this.tagList = [
          { name: "全部", id: 0 },
          ...apiTags.map((item) => ({
            name: item.name,
            id: item.id,
            ...item
          }))
        ];
        common_vendor.index.__f__("log", "at components/home.vue:419", "tags:", this.tagList);
        this.$nextTick(() => {
          setTimeout(() => {
            this.updateSliderPosition(0);
          }, 200);
        });
      } catch (error) {
        common_vendor.index.__f__("error", "at components/home.vue:427", "获取标签列表失败:", error);
      }
    },
    searchPlant(e) {
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
          this.$nextTick(() => {
            this.sliderWidth = finalWidth;
            this.sliderLeft = finalLeft;
          });
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
      if (item.id && !this.loadedImagesMap[item.id]) {
        this.$set(this.loadedImagesMap, item.id, true);
      }
    },
    // 确保图片显示的方法（用于页面返回时恢复图片状态）
    ensureImagesVisible() {
      this.$nextTick(() => {
        if (this.plantsList && this.plantsList.length > 0) {
          this.plantsList.forEach((plant) => {
            if (!this.loadedImagesMap[plant.id]) {
              this.$set(this.loadedImagesMap, plant.id, false);
            }
          });
        }
      });
    },
    goAddPage() {
      common_vendor.wx$1.vibrateShort({ type: "medium" });
      common_vendor.index.navigateTo({ url: `/pages/plantEdit/plantEdit?type=add` });
    },
    gotoDetail(item) {
      common_vendor.index.navigateTo({
        url: `/pages/plantDetail/plantDetail?id=${item.id}`
      });
    },
    onPageShow() {
      if (this.isFirstLoad) {
        this.isFirstLoad = false;
        return;
      }
      this.loadFamilyData();
      setTimeout(() => {
        this.ensureImagesVisible();
      }, 500);
    }
  },
  async created() {
    var _a;
    const menuButtonInfo = common_vendor.wx$1.getMenuButtonBoundingClientRect();
    this.menuButtonInfo = menuButtonInfo;
    const systemInfo = common_vendor.index.getWindowInfo();
    this.paddingLeft = systemInfo.screenWidth - menuButtonInfo.right;
    const app = getApp();
    this.topBarHeight = app.globalData.topBarHeight;
    this.windowWidth = app.globalData.windowWidth;
    const user = await utils_request.callContainer("/api/login");
    common_vendor.index.__f__("log", "at components/home.vue:546", "callContainer login:", user);
    const userInfo = user.data.user;
    const familyList = user.data.family;
    await new Promise((resolve) => {
      common_vendor.index.setStorage({ key: "userInfo", data: userInfo, success: resolve });
    });
    await new Promise((resolve) => {
      common_vendor.index.setStorage({ key: "family", data: familyList, success: resolve });
    });
    const defaultFamilyId = (userInfo == null ? void 0 : userInfo.currentFamilyId) || familyList && ((_a = familyList[0]) == null ? void 0 : _a.id);
    if (defaultFamilyId) {
      await new Promise((resolve) => {
        common_vendor.index.setStorage({ key: "familyId", data: defaultFamilyId, success: resolve });
      });
    }
    this.loadFamilyData();
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  const _easycom_uni_search_bar2 = common_vendor.resolveComponent("uni-search-bar");
  const _component_WaterfallBox = common_vendor.resolveComponent("WaterfallBox");
  (_easycom_uni_icons2 + _easycom_uni_search_bar2 + _component_WaterfallBox)();
}
const _easycom_uni_icons = () => "../uni_modules/uni-icons/components/uni-icons/uni-icons.js";
const _easycom_uni_search_bar = () => "../uni_modules/uni-search-bar/components/uni-search-bar/uni-search-bar.js";
if (!Math) {
  (_easycom_uni_icons + _easycom_uni_search_bar)();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  var _a;
  return common_vendor.e({
    a: common_vendor.p({
      type: "home",
      size: "18",
      color: "#6B8857"
    }),
    b: $data.isSelecting ? 1 : "",
    c: common_vendor.t(((_a = $data.familyRange[$data.currentFamilyIndex]) == null ? void 0 : _a.text) || "选择家庭"),
    d: $data.currentFamilyIndex,
    e: $data.familyRange,
    f: common_vendor.o((...args) => $options.handleFamilyChange && $options.handleFamilyChange(...args)),
    g: $data.isSelecting ? 1 : "",
    h: $data.menuButtonInfo.height + "px",
    i: $data.menuButtonInfo.height / 2 + "px",
    j: $data.menuButtonInfo.top + "px",
    k: $data.paddingLeft + "px",
    l: common_vendor.o((...args) => $options.toggleFamilySelect && $options.toggleFamilySelect(...args)),
    m: common_vendor.o((...args) => $options.onTouchStart && $options.onTouchStart(...args)),
    n: common_vendor.o((...args) => $options.onTouchEnd && $options.onTouchEnd(...args)),
    o: $data.topBarHeight + "px",
    p: common_vendor.o($options.searchPlant),
    q: common_vendor.o(($event) => $data.searchValue = $event),
    r: common_vendor.p({
      placeholder: "输入植物名称",
      radius: "20",
      focus: false,
      bgColor: "rgba(255,255,255,0.5)",
      clearButton: "auto",
      cancelButton: "none",
      modelValue: $data.searchValue
    }),
    s: common_vendor.p({
      type: "plusempty",
      size: "22",
      color: "#333"
    }),
    t: common_vendor.o((...args) => $options.goAddPage && $options.goAddPage(...args)),
    v: common_vendor.f($data.tagList, (item, index, i0) => {
      return {
        a: common_vendor.t(item.name),
        b: "tag-text-" + index,
        c: item.id,
        d: "tag-item-" + index,
        e: $data.currentTagIndex === index ? 1 : "",
        f: common_vendor.o(($event) => $options.selectTag(index, item), item.id)
      };
    }),
    w: common_vendor.s($options.sliderStyle),
    x: "tag-item-" + ($data.currentTagIndex > 1 ? $data.currentTagIndex - 1 : 0),
    y: $data.plantsList.length === 0
  }, $data.plantsList.length === 0 ? {
    z: common_assets._imports_0
  } : {
    A: common_vendor.w(({
      item
    }, s0, i0) => {
      return common_vendor.e({
        a: item.cover.url,
        b: $data.loadedImagesMap[item.id] ? 1 : "",
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
      path: "A",
      vueId: "045d88fd-3"
    }),
    B: common_vendor.p({
      list: $data.plantsList,
      idKey: "id",
      cols: "2"
    })
  });
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render], ["__scopeId", "data-v-045d88fd"]]);
wx.createComponent(Component);
//# sourceMappingURL=../../.sourcemap/mp-weixin/components/home.js.map
