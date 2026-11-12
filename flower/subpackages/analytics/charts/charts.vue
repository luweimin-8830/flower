<template>
    <view class="page-container theme-page" :style="themeStyle">
        <navBar />
        <!-- 顶部占位 -->
        <view :style="{ height: topBarHeight + 'px' }"></view>
        
        <view class="section-container">
            <!-- 🌟 图表类型与时间筛选 (仿 Home.vue 样式) -->
            <view class="selector-row">
                <view class="selector-wrapper">
                    <picker :value="chartIndex" :range="chartOptions" range-key="name" @change="handleChartChange">
                        <view class="chart-selector" hover-class="selector-hover">
                            <text class="selector-text">{{ chartOptions[chartIndex].name }}</text>
                            <uni-icons type="bottom" size="14" color="var(--primary-color)"
                                style="margin-left: 4px;"></uni-icons>
                        </view>
                    </picker>
                </view>

                <!-- 🚀 月份筛选框：仅在柱状图模式下显示 -->
                <view class="selector-wrapper" v-if="chartOptions[chartIndex].type === 'bar'">
                    <picker mode="date" fields="month" :value="currentDate" @change="handleDateChange">
                        <view class="chart-selector" hover-class="selector-hover">
                            <text class="selector-text">{{ currentDate }}</text>
                            <uni-icons type="bottom" size="14" color="var(--primary-color)"
                                style="margin-left: 4px;"></uni-icons>
                        </view>
                    </picker>
                </view>
            </view>

            <view class="card-box">
                <view class="charts-box">
                    <qiun-data-charts v-if="renderChart" :key="chartIndex" :type="chartOptions[chartIndex].type"
                        :opts="currentOpts" :chartData="chartData" :canvas2d="true"
                        :canvasId="'plantChartCanvas' + chartIndex" />
                </view>
            </view>
        </view>

        <view class="section-container">
            <text class="section-title">数据说明</text>
            <view class="card-box description-text">
                {{ chartOptions[chartIndex].desc }}
            </view>
        </view>
    </view>
</template>

<script>
import navBar from '@/components/navBar.vue'
import qiunDataCharts from '../uni_modules/qiun-data-charts/components/qiun-data-charts/qiun-data-charts.vue'
import { callContainer } from '@/utils/request';

// Labubu 配色：采用用户提供色卡中的 HEX 色值，并交错冷暖、明暗排列。
const DISTINCT_CHART_COLORS = [
    '#F67828', // 粉彩橙
    '#4F98CB', // 亮蓝色
    '#CBDD78', // 嫩芽绿
    '#E292B6', // 橡皮糖粉
    '#FEE37F', // 亮黄色
    '#9085C4', // 紫晶灰
    '#EF8E87', // 火烈鸟红
    '#6F787D', // 烟蓝色
    '#FCB937', // 梨黄色
    '#8C6F54', // 灰金色
    '#CBEBE1', // 冰水绿
    '#F4DCD8', // 蜜桃奶油
    '#A1A6D6', // 冰沙淡紫
    '#FDB482', // 杏橙色
    '#B5DBE6', // 极地蓝
    '#C2A585', // 南瓜浅棕
    '#F6F4A2', // 蜂鸣黄
    '#E6CED8', // 兰玫红
    '#D2E6F3', // 石灰淡青
    '#E3C08D', // 黄米色
    '#D4D3E6', // 兰白色
    '#F9C8C2', // 桃红色
    '#F1EAE2', // 香草白
    '#FFFAFA'  // 雪白色
];

export default {
    components: {
        qiunDataCharts,
        navBar
    },
    data() {
        return {
            familyId: 0,
            topBarHeight: 0,
            chartData: { categories: [], series: [] },
            chartIndex: 0,
            currentDate: '', // 格式：YYYY-MM
            renderChart: true,
            chartOptions: [
                { name: '植物分类统计', type: 'pie', desc: '展示当前家庭中各类标签下的植物分布比例。' },
                { name: '养护操作统计', type: 'bar', desc: '展示所选月份内，家庭中各项养护操作的累计执行次数。' }
            ],
            opts: {
                color: DISTINCT_CHART_COLORS,
                padding: [5, 5, 5, 5],
                enableScroll: false,
                fontColor: "#666666", // 默认亮色模式文字颜色
                extra: {
                    pie: {
                        activeOpacity: 0.8,
                        activeRadius: 10,
                        offsetAngle: 0,
                        labelWidth: 15,
                        border: false
                    }
                }
            },
            // 🚀 柱状图专用配置
            barOpts: {
                color: DISTINCT_CHART_COLORS,
                padding: [15, 30, 0, 5],
                enableScroll: false,
                fontColor: "#666666",
                legend: {},
                xAxis: {
                    boundaryGap: "justify",
                    disableGrid: false,
                    min: 0,
                    axisLine: false,
                },
                yAxis: {},
                extra: {
                    bar: {
                        type: "group",
                        width: 25,
                        meterBorde: 1,
                        activeBgColor: "#000000",
                        activeBgOpacity: 0.08,
                        barBorderCircle: true,
                        seriesGap: 2,
                        categoryGap: 2
                    }
                }
            }
        };
    },
    computed: {
        currentOpts() {
            const isPie = this.chartOptions[this.chartIndex].type === 'pie';
            return isPie ? this.opts : this.barOpts;
        }
    },
    async onLoad() {
        const app = getApp()
        this.topBarHeight = app.globalData.topBarHeight;
        
        // 设置默认日期为当前年月
        const now = new Date();
        this.currentDate = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;

        this.updateThemeStyle();

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
        updateThemeStyle() {
            this.opts.fontColor = "#666666";
            this.barOpts.fontColor = "#666666";
            
            // 如果图表已经渲染，需要强制刷新
            if (this.chartData.series) {
                this.chartData = JSON.parse(JSON.stringify(this.chartData));
            }
        },
        async handleChartChange(e) {
            this.renderChart = false; // 彻底销毁当前图表组件
            this.chartIndex = e.detail.value;
            // 🚀 重置为标准空结构
            this.chartData = { categories: [], series: [] }; 
            wx.vibrateShort({ type: 'light' });
            
            await this.getServerData();
            this.$nextTick(() => {
                this.renderChart = true; // 重新挂载新类型的图表组件
            });
        },
        async handleDateChange(e) {
            this.renderChart = false;
            this.currentDate = e.detail.value;
            this.chartData = { categories: [], series: [] }; 
            wx.vibrateShort({ type: 'light' });
            
            await this.getServerData();
            this.$nextTick(() => {
                this.renderChart = true;
            });
        },
        async getServerData() {
            const currentOption = this.chartOptions[this.chartIndex];
            try {
                if (currentOption.type === 'pie') {
                    const res = await callContainer("/api/chart/pie", { familyId: this.familyId });
                    if (res.data && res.data.length > 0) {
                        this.chartData = {
                            categories: [],
                            series: [{
                                data: res.data.map((item, index) => ({
                                    ...item,
                                    color: DISTINCT_CHART_COLORS[index % DISTINCT_CHART_COLORS.length]
                                }))
                            }]
                        };
                    } else {
                        this.chartData = { categories: [], series: [] };
                    }
                } else if (currentOption.type === 'bar') {
                    const [year, month] = this.currentDate.split('-').map(Number);
                    const res = await callContainer("/api/chart/bar", { 
                        familyId: this.familyId,
                        year: year,
                        month: month
                    });
                    if (res.data && res.data.series) {
                        this.chartData = res.data;
                    } else {
                        this.chartData = { categories: [], series: [] };
                    }
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
    padding-bottom: calc(30px + env(safe-area-inset-bottom));
}

.section-container {
    padding: 0 16px;
    margin-top: 20px;
}

.selector-row {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 15px;
}

.selector-wrapper {
    margin-left: 4px;
}

.chart-selector {
    display: inline-flex;
    align-items: center;
    padding: 6px 16px;
    background-color: var(--bg-btn-color);
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    border-radius: 20px;
    border: 1px solid var(--border-color);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    transition: all 0.2s;
}

.selector-hover {
    transform: scale(0.96);
    opacity: 0.8;
}

.selector-text {
    font-size: 14px;
    color: var(--text-color);
    font-weight: 500;
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
    height: 500px;
}

.description-text {
    font-size: 13px;
    color: var(--text-sub);
    line-height: 1.6;
}
</style>
