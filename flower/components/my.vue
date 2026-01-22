<template>
    <view :style="{height:topBarHeight+'px'}"></view>
    <view class="section-container">
        <!-- 标题 -->
        <text class="section-title">数据管理</text>

        <!-- 卡片容器 -->
        <view class="card-box">

            <!-- 列表项 1: 植物标签管理 -->
            <view class="menu-item" hover-class="item-hover" @click="handleNav('tag')">
                <view class="left-content">
                    <!-- 图标：使用了 uni-icons 的 tag 图标 -->
                    <uni-icons type="flag" size="22" color="#4A6139" class="menu-icon"></uni-icons>
                    <text class="menu-text">植物标签管理</text>
                </view>
                <uni-icons type="right" size="16" color="#999"></uni-icons>
            </view>

            <!-- 分割线 -->
            <view class="divider"></view>

            <!-- 列表项 2: 日常养护管理 -->
            <view class="menu-item" hover-class="item-hover" @click="handleNav('care')">
                <view class="left-content">
                    <!-- 图标：使用了 uni-icons 的 hand-up 图标作为替代 -->
                    <uni-icons type="settings" size="22" color="#4A6139" class="menu-icon"></uni-icons>
                    <text class="menu-text">日常养护管理</text>
                </view>
                <uni-icons type="right" size="16" color="#999"></uni-icons>
            </view>

        </view>
    </view>
</template>

<script>
export default {
    name: 'my',
    data() {
        return {
            topBarHeight:0,
        }
    },
    methods: {
        handleNav(type) {
            // 这里处理点击事件
            console.log('点击了:', type);
            wx.vibrateShort({ type: "medium" })
            if (type === 'tag') {
                uni.navigateTo({ url: '/pages/tagEdit/tagEdit' })
            } else if (type === 'care') {
                // uni.navigateTo({ url: '/pages/careManager/careManager' })
            }
        }
    },
    created() {
        const app = getApp()
        this.topBarHeight = app.globalData.topBarHeight;
    }
}
</script>

<style scoped lang="scss">
/* 容器整体样式 */
.section-container {
    padding: 0 16px;
    /* 页面左右边距 */
    margin-bottom: 20px;
}

/* 标题样式 */
.section-title {
    font-size: 14px;
    color: #333;
    /* 深色文字 */
    font-weight: 500;
    margin-bottom: 8px;
    margin-left: 4px;
    /* 稍微对齐卡片内部文字 */
    display: block;
}

/* 卡片主体样式 */
.card-box {
    background-color: rgba(255,255,255,0.55);
    /* 截图中的浅绿色背景 */
    /* 如果觉得颜色太深，可以试试 rgba(214, 232, 208, 0.6) */

    border-radius: 16px;
    /* 圆角 */
    padding: 4px 16px;
    /* 内部上下留白 */
    overflow: hidden;
}

/* 单个菜单项 */
.menu-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 0;
    /* 行高 */
    background-color: transparent;
}

/* 点击时的按压效果 */
.item-hover {
    opacity: 0.7;
}

/* 左侧图标和文字的容器 */
.left-content {
    display: flex;
    align-items: center;
}

/* 图标微调 */
.menu-icon {
    margin-right: 10px;
    /* 确保图标和文字垂直对齐 */
    display: flex;
    align-items: center;
}

/* 文字样式 */
.menu-text {
    font-size: 15px;
    color: #2F3E25;
    /* 深橄榄绿文字 */
    font-weight: 400;
}

/* 分割线 */
.divider {
    height: 1px;
    background-color: rgba(0, 0, 0, 0.05);
    /* 极淡的分割线 */
    width: 100%;
    /* 如果想让分割线不顶头，可以加 margin-left: 32px; */
    margin-left: 32px;
}
</style>