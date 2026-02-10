<template>
    <view class="page-container">
        <navBar />
        <!-- 顶部占位 -->
        <view :style="{ height: topBarHeight + 'px' }"></view>
        
        <view class="section-container">
            <text class="section-title">植物分类统计</text>
            <view class="card-box">
                <view class="charts-box">
                    <qiun-data-charts 
                        type="pie" 
                        :opts="opts" 
                        :chartData="chartData" 
                        :canvas2d="true"
                        canvasId="nFJjIDtblJABnYwVKwdOawkzzVcyGbpf" 
                    />
                </view>
            </view>
        </view>

        <view class="section-container">
            <text class="section-title">数据说明</text>
            <view class="card-box description-text">
                展示当前家庭中各类植物的分布比例。
            </view>
        </view>
    </view>
</template>

<script>
import navBar from '@/components/navBar.vue'
import qiunDataCharts from '@/uni_modules/qiun-data-charts/components/qiun-data-charts/qiun-data-charts.vue'
import { callContainer } from '../../utils/request';

export default {
    components: {
        qiunDataCharts,
        navBar
    },
    data() {
        return {
            familyId: 0,
            topBarHeight: 0,
            chartData: {},
            opts: {
                color: ["#778D61", "#A1AF8F", "#B9C7A9", "#D2DBCC", "#8798A5", "#E8D099", "#C9A09B"],
                padding: [5, 5, 5, 5],
                enableScroll: false,
                extra: {
                    pie: {
                        activeOpacity: 0.8,
                        activeRadius: 10,
                        offsetAngle: 0,
                        labelWidth: 15,
                        border: false
                    }
                }
            }
        };
    },
    async onLoad() {
        const app = getApp()
        this.topBarHeight = app.globalData.topBarHeight;
        
        try {
            const familyID = await new Promise((resolve, reject) => {
                uni.getStorage({ key: 'familyId', success: resolve, fail: reject })
            })
            this.familyId = familyID.data
            this.getServerData();
        } catch (error) {
            console.error('获取 familyId 失败:', error)
        }
    },
    /**
 * 页面初次渲染完成时触发的生命周期回调。
 * 此处无需执行数据请求：已在 onLoad 获取 familyId 后调用 getServerData 完成数据加载，避免重复拉取。
 */
onReady() {
        // getServerData 已在 onLoad 获取到 familyId 后执行
    },
    methods: {
        async getServerData() {
            try {
                const res = await callContainer("/api/chart/pie", { familyId: this.familyId });
                if (res.data && res.data.length > 0) {
                    // 后端已经处理好了“其他”分类和数据结构，直接赋值即可
                    this.chartData = {
                        series: [{
                            data: res.data
                        }]
                    };
                } else {
                    this.chartData = { series: [] };
                }
            } catch (error) {
                console.error('获取统计数据失败:', error);
                uni.showToast({ title: '加载失败', icon: 'none' });
            }
        },
    }
};
</script>

<style scoped lang="scss">
.page-container {
    min-height: 100vh;
    padding-bottom: 30px;
}

.section-container {
    padding: 0 16px;
    margin-top: 20px;
}

.section-title {
    font-size: 14px;
    color: var(--text-color);
    font-weight: 500;
    margin-bottom: 10px;
    margin-left: 4px;
    display: block;
}

.card-box {
    background-color: var(--bg-btn-color);
    border-radius: 16px;
    padding: 16px;
    overflow: hidden;
}

.charts-box {
    width: 100%;
    height: 300px;
}

.description-text {
    font-size: 13px;
    color: var(--text-sub);
    line-height: 1.6;
}
</style>