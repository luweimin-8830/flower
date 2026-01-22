<template>
    <view class="loading-mask" v-if="show" @touchmove.stop.prevent="preventTouch">
        <view class="loading-container">
            <!-- 动画主体 -->
            <view class="pulse-ring"></view>
            <view class="icon-box">
                <!-- 这里用了 uni-icons 的 leaf 图标，也可以换成你的 logo 图片 -->
                <view class="sprout-icon iconfont plant-zhiwuzhiyuan-duorouzhiwuyuan"></view>
                <!-- <uni-icons type="checkbox-filled" size="32" color="#fff" class="sprout-icon"></uni-icons> -->
            </view>
            <text class="loading-text">加载中...</text>
        </view>
    </view>
</template>

<script>
export default {
    name: 'loadingPage',
    data() {
        return {
            show: false
        }
    },
    methods: {
        open() {
            this.show = true;
        },
        close() {
            // 给个小延时，让动画不至于一闪而过，体验更好
            setTimeout(() => {
                this.show = false;
            }, 300);
        },
        preventTouch() {
            // 阻止遮罩层下的页面滚动
            return;
        }
    }
}
</script>

<style scoped lang="scss">
.loading-mask {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(255, 255, 255, 0.9);
    /* 90%透明度的白色背景 */
    z-index: 9999;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
}

.loading-container {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}

/* 核心动画：脉冲光环 */
.pulse-ring {
    width: 60px;
    height: 60px;
    background-color: #D6E8D0;
    /* 豆沙绿 */
    border-radius: 50%;
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    animation: pulse 2s infinite ease-in-out;
    opacity: 0.6;
}

/* 核心动画：图标容器 */
.icon-box {
    width: 50px;
    height: 50px;
    background-color: #4A6139;
    /* 深绿色实心圆 */
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    z-index: 2;
    box-shadow: 0 4px 10px rgba(74, 97, 57, 0.3);
    animation: float 3s infinite ease-in-out;
}

/* 图标微调 */
.sprout-icon {
    margin-top: 2px;
    color: #fff;
    font-size: 32px;
    /* 视觉修正 */
}

.loading-text {
    margin-top: 20px;
    font-size: 14px;
    color: #4A6139;
    font-weight: 500;
    letter-spacing: 1px;
    animation: fadeIn 1s infinite alternate;
}

/* 动画定义 */
@keyframes pulse {
    0% {
        transform: translate(-50%, -50%) scale(1);
        opacity: 0.6;
    }

    50% {
        transform: translate(-50%, -50%) scale(1.6);
        /* 放大光环 */
        opacity: 0;
    }

    100% {
        transform: translate(-50%, -50%) scale(1);
        opacity: 0;
    }
}

@keyframes float {

    0%,
    100% {
        transform: translateY(0);
    }

    50% {
        transform: translateY(-6px);
        /* 轻轻上浮 */
    }
}

@keyframes fadeIn {
    from {
        opacity: 0.6;
    }

    to {
        opacity: 1;
    }
}
</style>