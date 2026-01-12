<template>
    <!-- 家庭选择框 -->
    <view class="family-select" :style="{
        width: menuButtonInfo.height + 'px',
        height: menuButtonInfo.height + 'px', 
        top: menuButtonInfo.top + 'px',
        left: paddingLeft + 'px'
    }">
        <uni-data-select 
            class="custom-select" 
            v-model="value" 
            :localdata="range" 
            @change="change" 
            :clear="false"
        ></uni-data-select>
    </view>
    <view :style="{height:topBarHeight+'px'}"></view> 
    <view>
        <!-- <button>aaabsisnis</button> -->
    </view>
</template>

<script>
export default {
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
            topBarHeight:0,
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
    },
}
</script>

<style scoped lang="scss">
.family-select {
    position: fixed;
    z-index: 999;
    
    /* --- 核心毛玻璃样式 --- */
    background-color: rgba(255, 255, 255, 0.6); /* 半透明白底 */
    backdrop-filter: blur(10px); /* 模糊背景 */
    -webkit-backdrop-filter: blur(10px); /* 兼容 iOS */
    border-radius: 50%; /* 圆形 */
    border: 1px solid rgba(0, 0, 0, 0.08); /* 极细的浅色边框 */
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05); /* 轻微阴影 */
    
    /* 布局与过渡 */
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
    box-sizing: border-box; /* 确保边框不撑大尺寸 */
    overflow: hidden; /* 超出圆角部分隐藏 */
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
    font-size: 12px; /* 字体改小一点以适应按钮 */
    color: #333;
}

/* --- 深色模式适配 --- */
@media (prefers-color-scheme: dark) {
    .family-select {
        background-color: rgba(30, 30, 30, 0.5); /* 半透明黑底 */
        border: 1px solid rgba(255, 255, 255, 0.1); /* 浅白边框 */
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
    }
    
    ::v-deep .uni-select__input-text {
        color: #fff;
    }
}
</style>