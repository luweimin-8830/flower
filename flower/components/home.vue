<template>
    <!-- 🌟 1. 最外层包裹一个全屏容器 -->
    <view class="home-container">

        <!-- 🌟 2. 优化后的家庭选择按钮 (毛玻璃效果 + 交互增强) - 使用微信原生 picker -->
        <view class="family-select" :class="{ 'selecting': isSelecting }" :style="{
            width: 'auto',
            height: menuButtonInfo.height + 'px',
            borderRadius: menuButtonInfo.height / 2 + 'px',
            top: menuButtonInfo.top + 'px',
            left: paddingLeft + 'px'
        }" @click="toggleFamilySelect" @touchstart="onTouchStart" @touchend="onTouchEnd">
            <view class="family-select-icon" :class="{ 'selecting': isSelecting }">
                <uni-icons type="home" size="18" color="#6B8857"></uni-icons>
            </view>
            <picker class="custom-select" :value="currentFamilyIndex" :range="familyRange" :range-key="'text'"
                @change="handleFamilyChange">
                <text class="family-select-text">{{ familyRange[currentFamilyIndex]?.text || '选择家庭' }}</text>
            </picker>
        </view>

        <!-- 🌟 3. 顶部固定区域 (包含占位符、搜索框、标签栏) -->
        <view class="fixed-header-group">
            <!-- 顶部占位 (状态栏高度) -->
            <view :style="{ height: topBarHeight + 'px' }"></view>

            <!-- 搜索框 -->
            <view class="header-action-container">
                <view class="search-box-wrapper">
                    <uni-search-bar @confirm="searchPlant" placeholder="输入植物名称" radius="20" :focus="false"
                        v-model="searchValue" bgColor="rgba(255,255,255,0.5)" clearButton="auto" cancelButton="none">
                    </uni-search-bar>
                </view>
                <view class="add-btn" @click="goAddPage">
                    <uni-icons type="plusempty" size="22" color="#333"></uni-icons>
                </view>
            </view>

            <!-- 横向滚动标签 -->
            <view class="tag-scroll-container">
                <scroll-view scroll-x="true" class="tag-scroll-view" :show-scrollbar="false"
                    :scroll-into-view="'tag-item-' + (currentTagIndex > 1 ? currentTagIndex - 1 : 0)"
                    scroll-with-animation>
                    <view class="tag-flex-box" id="tag-container">
                        <view v-for="(item, index) in tagList" :key="item.ID" :id="'tag-item-' + index" class="tag-item"
                            :class="{ 'active': currentTagIndex === index }" @click="selectTag(index, item)">
                            <text :id="'tag-text-' + index" class="tag-text">{{ item.name }}</text>
                        </view>
                        <view class="slider-bar" :style="sliderStyle"></view>
                    </view>
                </scroll-view>
            </view>
        </view>

        <!-- 🌟 4. 中间独立滚动区域 -->
        <!-- flex:1 让它自动填满剩余空间，height:0 防止被内容撑大 -->
        <scroll-view scroll-y class="content-scroll-view">
            <!-- 空状态 -->
            <view v-if="plantsList.length === 0" class="empty-wrapper">
                <image src="/static/icon/c2m.svg" class="empty-icon" mode="aspectFit"></image>
                <text class="empty-text">点击右上方添加第一颗植物吧</text>
            </view>
            
            <!-- 瀑布流列表 -->
            <view v-else class="waterfall-wrapper">
                <WaterfallBox :list="plantsList" idKey="ID" cols="2">
                    <template #item="{ item }">
                        <view class="plant-card" @click="gotoDetail(item)">
                            <view class="image-wrapper"
                                :style="{ paddingBottom: (item.cover.height / item.cover.width * 100) + '%' }">
                                <image :src="item.cover.url" mode="aspectFill" class="plant-image" lazy-load
                                    :class="{ 'show': loadedImagesMap[item.ID] }" @load="onImgLoad(item)"></image>
                            </view>
                            <view class="plant-info">
                                <text class="plant-name">{{ item.name }}</text>
                                <view class="plant-tags" v-if="item.tags"></view>
                            </view>
                        </view>
                    </template>
                </WaterfallBox>

                <!-- 底部垫片：防止内容被 TabBar 遮挡 -->
                <view style="height: 20px;"></view>
            </view>
        </scroll-view>

    </view>
</template>


<script>
import { callContainer } from '../utils/request';
import WaterfallBox from './WaterfallBox.vue';

export default {
    components: {
        WaterfallBox,
    },
    /**
     * 组件名称，也就是开发者使用的标签
     */
    name: 'home',
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
            default: "",

        },
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
            loadedImagesMap: {}, // 用于追踪图片加载状态
            isFirstLoad: true,
            isSelecting: false,
            currentFamilyIndex: 0,
            isFiltering: false, // 防止过滤期间的并发修改
        }
    },
    computed: {
        // 动态生成滑块样式
        sliderStyle() {
            return {
                transform: `translateX(${this.sliderLeft}px)`,
                width: `${this.sliderWidth}px`
            }
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
                            plant.tags.forEach(tag => {
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
                // 同时获取家庭列表和当前选中的家庭ID
                const [familyListResult, familyIdResult] = await Promise.all([
                    new Promise((resolve, reject) => {
                        uni.getStorage({ key: 'family', success: resolve, fail: reject })
                    }),
                    new Promise((resolve, reject) => {
                        uni.getStorage({ key: 'familyId', success: resolve, fail: reject })
                    })
                ]);

                const familyList = familyListResult?.data || [];
                const cachedFamilyId = familyIdResult?.data;

                console.log("loadFamilyData - 家庭列表:", familyList, "缓存的家庭ID:", cachedFamilyId);

                if (familyList && Array.isArray(familyList) && familyList.length > 0) {
                    this.familyRange = familyList.map(item => ({
                        text: item.name,
                        value: item.ID || item.id,
                        disable: false
                    }));

                    // 使用缓存的家庭ID，如果缓存不存在则使用第一个家庭
                    this.value = cachedFamilyId || this.familyRange[0].value;

                    // 设置当前选中的家庭索引
                    this.currentFamilyIndex = this.familyRange.findIndex(item => item.value === this.value);
                    if (this.currentFamilyIndex === -1) {
                        this.currentFamilyIndex = 0;
                        this.value = this.familyRange[0].value;
                    }

                } else {
                    this.familyRange = [];
                    this.currentFamilyIndex = 0;
                    this.value = null;
                }

                // 确保 this.value 设置完成后再加载数据
                await this.$nextTick();

                // 使用正确的 familyId 加载数据
                await this.getTagList();
                await this.getPlantsList();
            } catch (error) {
                console.error("加载家庭数据失败:", error);
            }
        },
        async refreshFamilyList() {
            try {
                const user = await callContainer("/api/login");
                const familyList = user.data.family || [];

                // 更新缓存
                await new Promise((resolve) => {
                    uni.setStorage({ key: "family", data: familyList, success: resolve })
                });

                // 更新家庭选择器的选项
                if (familyList && Array.isArray(familyList) && familyList.length > 0) {
                    this.familyRange = familyList.map(item => ({
                        text: item.name,
                        value: item.ID || item.id,
                        disable: false
                    }));

                    // 确保当前选择的家庭ID在新列表中
                    const currentFamilyExists = this.familyRange.some(item => item.value === this.value);
                    if (!currentFamilyExists && this.familyRange.length > 0) {
                        // 如果当前家庭不在新列表中，切换到第一个
                        this.value = this.familyRange[0].value;
                        this.currentFamilyIndex = 0;

                        // 更新缓存
                        await new Promise((resolve) => {
                            uni.setStorage({ key: "familyId", data: this.value, success: resolve })
                        });

                        // 刷新数据
                        await this.getTagList();
                        await this.getPlantsList();
                    }
                } else {
                    this.familyRange = [];
                }

            } catch (error) {
                console.error("刷新家庭列表失败:", error);
            }
        },
        async getPlantsList() {
            const familyId = this.value;

            try {
                const plants = await callContainer("/api/plant/list", {
                    "familyId": familyId
                })
                console.log("plants list:", plants)
                const newData = plants?.data || [];

                // 打印原始数据的 tags
                newData.forEach((item, idx) => {
                    if (item.tags && item.tags.length > 0) {
                    }
                });

                // 🔥 关键修复：完全冻结所有对象，防止被修改
                this.allPlantsList = newData.map(item => {
                    // 深拷贝 tags 数组
                    let frozenTags = null;
                    if (item.tags && Array.isArray(item.tags)) {
                        frozenTags = item.tags.map(tag => ({ ...tag }));
                        // 冻结每个 tag 对象
                        frozenTags.forEach(tag => Object.freeze(tag));
                        // 冻结 tags 数组
                        Object.freeze(frozenTags);
                    }

                    // 创建新对象并冻结（不再需要 isLoaded 属性）
                    const newItem = { ...item, tags: frozenTags };
                    return Object.freeze(newItem);
                });
                this.filterPlants();
            } catch (error) {
                console.error("获取植物列表失败:", error)
            }
        },
        filterPlants() {
            const currentTag = this.tagList[this.currentTagIndex];
            const tagId = currentTag ? currentTag.ID : 0;
            let filtered = [];
            if (tagId === 0) {

                filtered = this.allPlantsList;
            } else {
                // 🟢 只过滤，不 map
                filtered = this.allPlantsList.filter(plant => {
                    return plant.tags && plant.tags.some(t => t.ID === tagId);
                });
            }
            this.plantsList = filtered;
        },
        async handleFamilyChange(e) {
            const selectedIndex = e.detail.value;
            this.currentFamilyIndex = selectedIndex;
            const newFamilyId = this.familyRange[selectedIndex].value;

            try {
                // 调用后端切换家庭接口
                await callContainer("/api/family/switch", {
                    familyId: newFamilyId
                });
                console.log("家庭切换成功");

                // 更新storage中的familyId
                await new Promise((resolve) => {
                    uni.setStorage({ key: "familyId", data: newFamilyId, success: resolve })
                });
                
                // 触发家庭切换事件，通知其他组件
                uni.$emit('familyChanged', newFamilyId);
                
            } catch (error) {
                console.error("切换家庭失败:", error);

                // 显示错误提示
                const errorMsg = error?.msg || error?.message || "切换家庭失败，请稍后重试";
                uni.showToast({
                    title: errorMsg,
                    icon: 'none',
                    duration: 2000
                });

                // 恢复之前的选择
                this.currentFamilyIndex = this.familyRange.findIndex(item => item.value === this.value);

                // 刷新家庭列表，移除无权限的家庭
                await this.refreshFamilyList();
                return;
            }

            // 使用新familyId更新数据
            this.value = newFamilyId;
            this.currentTagIndex = 0;


            // 清空旧数据
            this.tagList = [];
            this.allPlantsList = [];
            this.plantsList = [];
            this.loadedImagesMap = {}; // 清空图片加载状态

            // 等待 DOM 更新
            await this.$nextTick();

            // 直接获取新家庭的标签和植物列表
            await this.getTagList();
            await this.getPlantsList();

            this.$nextTick(() => {
                setTimeout(() => {
                    this.updateSliderPosition(0);
                }, 200);
            });

            wx.vibrateShort({ type: "light" });
        },
        toggleFamilySelect() {
            // 微信原生 picker 会自动展开，无需额外触发
            // 这里可以添加一些额外的逻辑，比如聚焦或高亮
            console.log("触发家庭选择器");
        },
        onTouchStart() {
            // 按钮按下时的样式变化
            this.isSelecting = true;
        },
        onTouchEnd() {
            // 按钮释放时的样式恢复
            setTimeout(() => {
                this.isSelecting = false;
            }, 200);
        },
        async getTagList() {
            const familyId = this.value;

            try {
                const tagList = await callContainer("/api/tag/", {
                    familyId: familyId
                })
                console.log("tagList:", tagList)
                const apiTags = tagList?.data || []
                this.tagList = [
                    { name: "全部", ID: 0 },
                    ...apiTags.map(item => ({
                        name: item.name,
                        ID: item.ID,
                        ...item
                    }))
                ]
                console.log("tags:", this.tagList)
                this.$nextTick(() => {
                    // 稍微延迟一点，确保 DOM 渲染完成
                    setTimeout(() => {
                        this.updateSliderPosition(0);
                    }, 200);
                });
            } catch (error) {
                console.error("获取标签列表失败:", error)
            }
        },
        searchPlant(e) {

        
        },
        selectTag(index, item) {
            if (this.currentTagIndex === index) return;
            wx.vibrateShort({ type: "medium" })

            // 优化标签切换逻辑：先更新索引，再过滤，避免中间状态
            this.currentTagIndex = index;
            this.filterPlants();

            // 滑块动画逻辑
            const query = uni.createSelectorQuery().in(this);
            query.select('#tag-container').boundingClientRect();
            query.select('#tag-text-' + index).boundingClientRect();
            query.exec((res) => {
                if (res[0] && res[1]) {
                    const containerLeft = res[0].left;
                    const currentTextLeft = res[1].left;
                    const currentTextWidth = res[1].width;
                    const ratio = 22 / 18;
                    const finalWidth = currentTextWidth * ratio;
                    const widthDiff = finalWidth - currentTextWidth;
                    const finalLeft = (currentTextLeft - containerLeft) - (widthDiff / 2);

                    // 使用 Vue 的 nextTick 确保响应式更新
                    this.$nextTick(() => {
                        this.sliderWidth = finalWidth;
                        this.sliderLeft = finalLeft;
                    });

                    if (this.sliderTimer) clearTimeout(this.sliderTimer);
                    this.sliderTimer = setTimeout(() => {
                        this.updateSliderPosition(index);
                    }, 350);
                }
            });
        },


        updateSliderPosition(index) {
            const query = uni.createSelectorQuery().in(this);
            query.select('#tag-container').boundingClientRect();
            query.select('#tag-text-' + index).boundingClientRect();

            query.exec((res) => {
                if (res[0] && res[1]) {
                    const containerLeft = res[0].left; // 容器距离屏幕左边的距离
                    const textLeft = res[1].left;      // 文字距离屏幕左边的距离
                    const textWidth = res[1].width;    // 文字宽度

                    // 计算相对位置：文字位置 - 容器位置 = 滑块在容器内的 left
                    // 注意：因为是在 scroll-view 内部，这种相对计算方式即使在滚动后也是正确的
                    this.sliderLeft = textLeft - containerLeft;
                    this.sliderWidth = textWidth;
                }
            });
        },
        onImgLoad(item) {
            // 使用独立的映射对象来追踪加载状态，避免修改冻结对象
            if (item.ID && !this.loadedImagesMap[item.ID]) {
                this.$set(this.loadedImagesMap, item.ID, true);
            }

        },
        // 确保图片显示的方法（用于页面返回时恢复图片状态）
        ensureImagesVisible() {
            this.$nextTick(() => {
                if (this.plantsList && this.plantsList.length > 0) {
                    this.plantsList.forEach(plant => {
                        if (!this.loadedImagesMap[plant.ID]) {
                            // 重置状态，触发重新加载
                            this.$set(this.loadedImagesMap, plant.ID, false);
                        }
                    });
                }
            });
        },
        goAddPage() {
            wx.vibrateShort({ type: "medium" })
            // 传入当前家庭ID
            uni.navigateTo({ url: `/pages/plantEdit/plantEdit?type=add` });
        },
        gotoDetail(item) {
            uni.navigateTo({
                url: `/pages/plantDetail/plantDetail?id=${item.ID}`
            })
        },
        onPageShow() {
            // 首次加载时不调用 loadFamilyData，避免重复加载
            if (this.isFirstLoad) {
                
                this.isFirstLoad = false;
                return;
            }

        
            this.loadFamilyData();

            // 确保图片显示
            setTimeout(() => {
                this.ensureImagesVisible();
            }, 500);
        },
    },
    async created() {

        const menuButtonInfo = wx.getMenuButtonBoundingClientRect()
        this.menuButtonInfo = menuButtonInfo
        const systemInfo = uni.getWindowInfo()
        this.paddingLeft = systemInfo.screenWidth - menuButtonInfo.right
        const app = getApp()
        this.topBarHeight = app.globalData.topBarHeight;
        this.windowWidth = app.globalData.windowWidth;
        const user = await callContainer("/api/login")
        console.log("callContainer login:", user)

        const userInfo = user.data.user;
        const familyList = user.data.family;

        // 保存用户信息
        await new Promise((resolve) => {
            uni.setStorage({ key: "userInfo", data: userInfo, success: resolve })
        })

        // 保存家庭列表
        await new Promise((resolve) => {
            uni.setStorage({ key: "family", data: familyList, success: resolve })
        })

        // 确定默认家庭ID：优先使用 userInfo 中的 currentFamilyId，否则使用家庭列表的第一个
        const defaultFamilyId = userInfo?.currentFamilyId || (familyList && familyList[0]?.ID);
        

        // 保存默认家庭ID到缓存
        if (defaultFamilyId) {
            await new Promise((resolve) => {
                uni.setStorage({ key: "familyId", data: defaultFamilyId, success: resolve })
            })
        }

        this.loadFamilyData()
        // uni.$off('refreshFamilyList');
        // uni.$on('refreshFamilyList', (data) => {
        //     const user = await callContainer("/api/login")
        //     await new Promise((resolve) => {
        //         uni.setStorage({ key: "family", data: user.data.family, success: resolve })
        //     })
        //     this.loadFamilyData()
        // })

    }
}
</script>

<style scoped lang="scss">
/* 1. 页面容器：占满全屏，垂直排列 */
.home-container {
    height: 100vh;
    /* 关键：锁定高度 */
    display: flex;
    flex-direction: column;
    overflow: hidden;
    /* 禁止整个页面拖动 */
    box-sizing: border-box;
    background-color: #C1D0B7;
    /* 建议加个背景色，防止列表滚动到底部露白 */
}

/* 2. 头部固定区域组 */
.fixed-header-group {
    flex-shrink: 0;
    /* 禁止压缩 */
    z-index: 10;
    background-color: #C1D0B7;
    /* 必须给背景色，否则列表滚动时会透过文字看到下面 */
    /* 如果你的设计是背景图通铺，这里可以用 transparent，但要注意视觉重叠 */
}

/* 3. 滚动区域：自动填满剩余空间 */
.content-scroll-view {
    flex: 1;
    /* 占据剩余高度 */
    height: 0;
    /* 🌟 关键：强制触发 Flex 计算，防止 scroll-view 被内容撑开导致失效 */
    overflow: hidden;
    margin-bottom: 160rpx;
}

.waterfall-wrapper {
    padding-bottom: env(safe-area-inset-bottom);
    /* 适配 iPhone 底部安全区 */
}

.family-select {
    position: fixed;
    z-index: 999;

    /* --- 核心毛玻璃样式 --- */
    background-color: rgba(255, 255, 255, 0.5);
    /* 半透明白底 */
    backdrop-filter: blur(10px);
    /* 模糊背景 */
    -webkit-backdrop-filter: blur(10px);
    /* 兼容 iOS */
    border-radius: 20px;
    border: 1px solid rgba(0, 0, 0, 0.08);
    /* 极细的浅色边框 */
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    /* 增强阴影效果 */

    /* 布局与过渡 */
    display: flex;
    align-items: center;
    justify-content: flex-start;
    padding: 0 12px;
    gap: 6px;
    transition: all 0.3s ease;
    box-sizing: border-box;
    /* 确保边框不撑大尺寸 */

    /* 交互状态 */
    &:active,
    &.selecting {
        transform: scale(0.95);
        box-shadow: 0 1px 4px rgba(0, 0, 0, 0.15);
    }

    &.selecting .family-select-icon {
        transform: scale(0.9);
        opacity: 1;
    }
}

.family-select-icon {
    flex-shrink: 0;
    opacity: 0.8;
    transition: opacity 0.2s;
}

.family-select-icon uni-icons {
    font-size: 16px;
}

/* 微信 picker 样式调整 */
.custom-select {
    flex: 1;
    height: 100%;
    background: transparent;
    border: none;
    padding: 0;
    margin: 0;
    display: flex;
    align-items: center;
}

.family-select-text {
    font-size: 14px;
    color: #333;
    opacity: 0.9;
    font-weight: 500;
}



.header-action-container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 10px;
    /* 左右留白 */
    margin-bottom: 5px;
    margin-top: 10px;
    /* 和下方 Tag 保持一点距离 */
}

.search-box-wrapper {
    /* 核心：搜索框占 82% (稍微多一点看起来更协调，留 18% 给按钮) */
    width: 86%;
}

.add-btn {
    width: 74rpx;
    /* 稍微加大一点点，更易点击 */
    height: 74rpx;

    /* 1. 微弱的线性渐变，模拟光照（上亮下暗） */
    // background: linear-gradient(145deg, #7da066, #607a4e);
    background: rgba(255, 255, 255, 0.55);
    /* 如果不支持渐变回退到纯色 */
    // background-color: #6B8857; 

    border-radius: 50%;

    display: flex;
    align-items: center;
    justify-content: center;

    transition: all 0.2s cubic-bezier(0.25, 0.8, 0.25, 1);

    /* 3. 增加一点边框让轮廓更清晰 */
    border: 1px solid rgba(255, 255, 255, 0.1);

    &:active {
        transform: scale(0.92) translateY(2px);
        /* 点击时下沉 */
    }
}

::v-deep .uni-searchbar {
    padding: 10px 0 !important;
    /* 去掉左右默认 padding */
}

/* 这是一个深度选择器，用于去除 uni-data-select 自带的边框，使其融入毛玻璃按钮 */
::v-deep .uni-select {
    border: none !important;
    background-color: transparent !important;
    padding: 0 !important;
    height: 100%;
    justify-content: center;
}

::v-deep .uni-select__input-text {
    font-size: 12px;
    /* 字体改小一点以适应按钮 */
    color: #333;
}

::v-deep .uni-select__selector-item {
    /* 这一行是为了覆盖原生样式，确保我们自定义的 slot 充满整行 */
    padding: 0 !important;
}

// tag css
.tag-scroll-container {
    width: 100%;
    background-color: transparent;
    padding: 5px 0;
}

.tag-scroll-view {
    width: 100%;
    white-space: nowrap;
    /* 关键：禁止换行 */
}

.tag-flex-box {
    display: flex;
    align-items: center;
    padding: 0 10px;
    position: relative;
}

.tag-item {
    position: relative;
    /* 为了定位下划线 */
    display: inline-flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 8px;
    margin-right: 10px;
    font-size: 16px;
    color: #666;
    transition: all 0.3s;

    &.active {
        // color:#BC3823;
        color: #6B8857;
        /* 选中颜色 */
        font-weight: bold;
        font-size: 20px;
        /* 选中稍微变大 */
    }
}

/* 下划线动画样式 */
.slider-bar {
    position: absolute;
    bottom: 8px;
    /* 距离底部的位置 */
    left: 0;
    /* 初始位置，由 JS 控制 translateX */
    height: 3px;
    // background-color: #BC3823;
    background-color: #6B8857;
    /* 下划线颜色 */
    border-radius: 2px;

    /* 动画配置 */
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    /* width 和 transform 都会平滑过渡 */

    z-index: 1;
    /* 确保在文字下方或上方 */
    pointer-events: none;
    /* 不影响点击 */
}

.active-line {
    position: absolute;
    bottom: 0;
    width: 20px;
    /* 下划线宽度 */
    height: 3px;
    background-color: #6B8857;
    border-radius: 2px;
    animation: scaleIn 0.2s ease-out;
}

.plant-card {
    background-color: rgba(255, 255, 255, 0.5);
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
    // transform: translateY(0);
}

.image-wrapper {
    position: relative;
    width: 100%;
    height: 0;
    /* 背景色作为占位时的颜色 */
    background-color: rgba(255, 255, 255, 0.6);
    overflow: hidden;
}

.plant-image {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    // background-color: rgba(255,255,255,1); // 加载时的背景色
    // opacity: 0;
    transition: opacity 0.4s ease-in-out;
}

.plant-image.show {
    opacity: 1;
}

.plant-info {
    padding: 8px;
}

.plant-name {
    font-size: 14px;
    color: #333;
    font-weight: bold;
}

/* 空状态样式 */
.empty-wrapper {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding: 200rpx 0;
}

.empty-icon {
    width: 120rpx;
    height: 120rpx;
    margin-bottom: 20rpx;
    opacity: 0.4;
    filter: grayscale(100%);
}

.empty-text {
    font-size: 14px;
    color: #999;
}

@keyframes scaleIn {
    from {
        transform: scaleX(0);
    }

    to {
        transform: scaleX(1);
    }
}

/* --- 深色模式适配 --- */
@media (prefers-color-scheme: dark) {
    .family-select {
        background-color: rgba(30, 30, 30, 0.5);
        /* 半透明黑底 */
        border: 1px solid rgba(255, 255, 255, 0.1);
        /* 浅白边框 */
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
    }

    ::v-deep .uni-select__input-text {
        color: #fff;
    }

    .tag-item {
        color: #aaa;

        &.active {
            color: #409eff;
        }
    }

    .active-line {
        background-color: #409eff;
    }

    .slider-bar {
        background-color: #8bb374;
    }
}

.waterfall-wrapper {
    padding-bottom: env(safe-area-inset-bottom);
    /* 适配 iPhone 底部安全区 */
}
</style>