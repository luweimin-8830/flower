<template>
    <view class="nav-back-button" hover-class="nav-back-button-hover" :style="{width:menuButtonInfo.height+'px',
    height:menuButtonInfo.height+'px',top:menuButtonInfo.top+'px',
    left: paddingLeft + 'px'}" @click="handleBack">
        <uni-icons v-if="isHome" type="home" size="20" color="#333"></uni-icons>
        <view v-else class="iconfont plant-fanhui nav-back-icon"></view>
    </view>
</template>

<script>
export default {
    /**
     * 组件名称，也就是开发者使用的标签
     */
    name: 'navBar',
    /**
     * 属性声明
     */
    props: {
        navText: {
            type: String,
            default: "",
        },
        isHome: {
            type: Boolean,
            default: false
        }
    },
    data() {
        return {
            menuButtonInfo:{},
            paddingLeft:0
        }
    },
    methods: {
        handleBack(){
            wx.vibrateShort({ type: 'medium' })
            if (this.isHome) {
                uni.reLaunch({
                    url: '/pages/index/index'
                })
            } else {
                uni.navigateBack()
            }
        }
    },
    /**
     * [可选实现] 组件被创建，组件第一个生命周期，
     * 在内存中被占用的时候被调用，开发者可以在这里执行一些需要提前执行的初始化逻辑
     */
    created() {
        const menuButtonInfo = wx.getMenuButtonBoundingClientRect()
        this.menuButtonInfo = menuButtonInfo
        const systemInfo = wx.getWindowInfo()
        this.paddingLeft = systemInfo.screenWidth - menuButtonInfo.right
    },
}
</script>

<style scoped lang="scss">
    .nav-back-button{
        position: fixed;
        z-index: 999;
        opacity:0.55;
        background-color: #fff;
        // backdrop-filter: blur(10px);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        border: 1px solid rgba(0,0,0,0.08);
        box-shadow: 0 2px 8px rgba(0,0,0,0.15);
        transition: transform 0.1s,background-color 0.2s;
    }
    .nav-back-icon{
        font-size: 20px;
        // color: var(--text-color,#333);
        color:#333
    
    }
    .nav-back-button-hover {
        transform: scale(0.92);
        background-color: #f5f5f5; 
        opacity:0.65;
        box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
    }
    @media (prefers-color-scheme: dark) {
        .nav-back-button {
            background-color: rgba(30, 30, 30, 0.6);
            border: 1px solid rgba(255, 255, 255, 0.1);
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
        }
        .nav-back-button-hover {
            background-color: rgba(50, 50, 50, 0.8);
        }
    }
</style>