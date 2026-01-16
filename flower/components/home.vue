<template>
    <!-- 家庭选择框 -->
    <view class="family-select" :style="{
        width: 'auto',
        // height: menuButtonInfo.height + 'px', 
        borderRadius: menuButtonInfo.height / 2 + 'px',
        top: menuButtonInfo.top + 'px',
        left: paddingLeft + 'px'
    }">
        <uni-data-select class="custom-select" v-model="value" :localdata="familyRange" @change="changeFamily"
            :clear="false">
        </uni-data-select>
    </view>
    <view :style="{ height: topBarHeight + 'px' }"></view>
    <!-- 搜索框 -->
    <view class="header-action-container">
        <view class="search-box-wrapper">
            <uni-search-bar @confirm="searchPlant" placeholder="输入植物名称" radius="20" :focus="true" v-model="searchValue"
            bgColor="rgba(255,255,255,0.5)" clearButton="auto" cancelButton="none">
        </uni-search-bar>
        </view>
        <view class="add-btn" @click="goAddPage">
            <uni-icons type="plusempty" size="22" color="#333"></uni-icons>
        </view>
    </view>
    <!-- 横向滚动列表 -->
    <view class="tag-scroll-container">
        <scroll-view scroll-x="true" class="tag-scroll-view" :show-scrollbar="false"
            :scroll-into-view="'tag-item-' + (currentTagIndex > 1 ? currentTagIndex - 1 : 0)" scroll-with-animation>
            <!-- 必须给 flex 容器一个 id，用于计算相对位置 -->
            <view class="tag-flex-box" id="tag-container">
                <view v-for="(item, index) in tagList" :key="index" :id="'tag-item-' + index" class="tag-item"
                    :class="{ 'active': currentTagIndex === index }" @click="selectTag(index, item)">
                    <text :id="'tag-text-' + index" class="tag-text">{{ item.name }}</text>
                </view>
                <!-- 独立的滑动下划线 -->
                <view class="slider-bar" :style="sliderStyle"></view>

            </view>
        </scroll-view>
    </view>
    <!-- 植物列表瀑布流 -->
    <view>
        <WaterfallBox :list="plantsList" idKey="ID" cols="2">
            <template #item="{ item }">
                <view class="plant-card">
                    <view class="image-wrapper" :style="{ paddingBottom: (item.height / item.width * 100) + '%' }">
                        <image :src="item.cover" mode="aspectFill" class="plant-image" lazy-load
                            :class="{ 'show': item.isLoaded }" @load="onImgLoad(item)"></image>
                    </view>
                    <!-- 文字内容 -->
                    <view class="plant-info">
                        <text class="plant-name">{{ item.name }}</text>
                        <view class="plant-tags" v-if="item.tags"></view>
                    </view>
                </view>
            </template>

        </WaterfallBox>
    </view>
</template>

<script>
import { callContainer } from '../utils/request';
import WaterfallBox from './WaterfallBox.vue'
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
                const familyList = await new Promise((resolve) => {
                    uni.getStorage({ key: 'family', success: resolve })
                })
                this.familyRange = []
                familyList?.data?.forEach(item => {
                    this.familyRange = [...this.familyRange, { "text": item.name, "value": item.ID, "disable": false }]
                });
                const exists = this.familyRange.some(item => item.value === this.value);
                if (!exists && this.familyRange.length > 0) {
                    this.value = this.familyRange[0].value;
                }
                this.getTagList()
                this.getPlantsList()
            } catch (error) {
                console.error(error)
            }
        },
        async getPlantsList() {
            try {
                const plants = await callContainer("/api/plant/list", {
                    "familyId": this.value
                })
                console.log("plants list:", plants)
                this.plantsList = plants?.data

                //测试,删除
                this.plantsList.forEach((item,index)=>{
                    item.width = 256;
                    item.height = 256;
                })
            } catch (error) {

            }
        },
        changeFamily(e) {
            console.log(e)
        },
        async getTagList() {
            try {
                const tagList = await callContainer("/api/tag/", {
                    familyId: this.value
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
                console.error(error)
            }
        },
        searchPlant(e) {
            console.log("e", e)
            console.log("search:", this.searchValue)
        },
        selectTag(index, item) {
            this.currentTagIndex = index;
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
                    this.sliderWidth = finalWidth;
                    this.sliderLeft = finalLeft;
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
            item.isLoaded = true
        },
        goAddPage() {
            wx.vibrateShort({type:"medium"})
            // 传入当前家庭ID
            uni.navigateTo({ url: `/pages/addPlant/addPlant?familyId=${this.value}` });
        },
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
        const menuButtonInfo = wx.getMenuButtonBoundingClientRect()
        this.menuButtonInfo = menuButtonInfo
        const systemInfo = uni.getWindowInfo()
        this.paddingLeft = systemInfo.screenWidth - menuButtonInfo.right
        const app = getApp()
        this.topBarHeight = app.globalData.topBarHeight;
        this.windowWidth = app.globalData.windowWidth;
        this.loadFamilyData()

    },
}
</script>

<style scoped lang="scss">
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
    // border-radius: 50%; /* 圆形 */
    border: 1px solid rgba(0, 0, 0, 0.08);
    /* 极细的浅色边框 */
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
    /* 轻微阴影 */

    /* 布局与过渡 */
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
    box-sizing: border-box;
    /* 确保边框不撑大尺寸 */
    //overflow: hidden; /* 超出圆角部分隐藏 */
}



.header-action-container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 10px; /* 左右留白 */
    margin-bottom: 5px; /* 和下方 Tag 保持一点距离 */
}

.search-box-wrapper {
    /* 核心：搜索框占 82% (稍微多一点看起来更协调，留 18% 给按钮) */
    width: 86%; 
}

.add-btn {
    width: 74rpx; /* 稍微加大一点点，更易点击 */
    height: 74rpx;
    
    /* 1. 微弱的线性渐变，模拟光照（上亮下暗） */
    // background: linear-gradient(145deg, #7da066, #607a4e);
    background: rgba(255,255,255,0.55);
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
        transform: scale(0.92) translateY(2px); /* 点击时下沉 */
    }
}

::v-deep .uni-searchbar {
    padding: 10px 0 !important; /* 去掉左右默认 padding */
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
    opacity: 0;
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
</style>