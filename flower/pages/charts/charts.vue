<template>
    <view class="page-wrapper">
        <nav-bar />
        <view :style="{ height: topBarHeight + 'px' }"></view>

        <view class="container">
            <view class="chart-card">
                <text class="chart-title">班级比例</text>
                <view class="charts-box">
                    <!-- 🌟 关键：使用 v-if 确保数据就绪后再渲染组件，解决 Canvas 初始化宽高为 0 的问题 -->
                    <qiun-data-charts v-if="isDataReady" type="pie" :opts="opts" :chartData="chartData"
                        :animation="true" :canvas2d="true" canvasId="BKdzxMyvrniCHuVpFDZeLStxHyvnBHxt" />
                    <!-- 加载中状态占位 -->
                    <view v-else class="loading-box">
                        <uni-load-more status="loading" />
                    </view>
                </view>
            </view>
        </view>
    </view>
</template>

<script>
import navBar from '@/components/navBar.vue'
export default {
    components: { navBar },
    data() {
        return {
            topBarHeight: 0,
            isDataReady: false,
            chartData: {},
            opts: {
                color: ["#1890FF", "#91CB74", "#FAC858", "#EE6666", "#73C0DE", "#3CA272", "#FC8452", "#9A60B4", "#ea7ccc"],
                padding: [5, 5, 5, 5],
                legend: {
                    show: true,
                    position: "bottom",
                    lineHeight: 25
                },
                extra: {
                    pie: {
                        activeOpacity: 0.5,
                        activeRadius: 10,
                        offsetAngle: 0,
                        labelWidth: 15,
                        border: true,
                        borderWidth: 3,
                        borderColor: "#FFFFFF"
                    }
                }
            }
        };
    },
    onLoad() {
        const app = getApp();
        this.topBarHeight = app.globalData.topBarHeight;
    },
    onReady() {
        this.getServerData();
    },
    methods: {
        getServerData() {
            // 模拟从服务器获取数据
            setTimeout(() => {
                let res = {
                    series: [
                        {
                            data: [
                                { "name": "一班", "value": 50 },
                                { "name": "二班", "value": 30 },
                                { "name": "三班", "value": 20 },
                                { "name": "四班", "value": 18 },
                                { "name": "五班", "value": 8 }
                            ]
                        }
                    ]
                };
                this.chartData = JSON.parse(JSON.stringify(res));
                this.isDataReady = true; // 🌟 数据就绪，触发渲染
            }, 800);
        },
    }
};
</script>

<style scoped lang="scss">
.page-wrapper {
    min-height: 100vh;
    background-color: var(--bg-color);
}

.container {
    padding: 20px 16px;
}

.chart-card {
    background-color: var(--bg-btn-color);
    backdrop-filter: blur(10px);
    border-radius: 20px;
    padding: 20px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
}

.chart-title {
    font-size: 16px;
    font-weight: bold;
    color: var(--text-color);
    margin-bottom: 20px;
    display: block;
    text-align: center;
}

.charts-box {
    width: 100%;
    height: 320px;
}

.loading-box {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>