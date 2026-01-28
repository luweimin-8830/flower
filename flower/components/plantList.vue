<template>
    <!-- 最外层包裹一个全屏容器 -->
    <view class="plantlist-container">

        <!-- 顶部固定区域 -->
        <view class="fixed-header-group">
            <!-- 顶部占位 + 标题栏合并 -->
            <view class="header-title-container" :style="{ paddingTop: statusBarHeight + 'px' }">
                <text class="page-title">植物列表</text>
            </view>
        </view>

        <!-- 中间独立滚动区域 -->
        <view class="content-container">
            <!-- 加载中状态 -->
            <view v-if="loading" class="loading-wrapper">
                <text class="loading-text">加载中...</text>
            </view>

            <!-- 空状态 -->
            <view v-else-if="indexedOptions.length === 0" class="empty-wrapper">
                <image src="/static/icon/c2m.svg" class="empty-icon" mode="aspectFit"></image>
                <text class="empty-text">暂无植物</text>
            </view>

            <!-- 索引列表 -->
            <view v-else class="indexed-list-wrapper">
                <scroll-view scroll-y class="plant-scroll-view" :scroll-into-view="scrollIntoView">
                    <view v-for="(group, index) in indexedOptions" :key="index" :id="'group-' + group.key" class="plant-group">
                        <view class="group-title">{{ group.key }}</view>
                        <view v-for="item in group.data" :key="item.id" class="plant-item" @click="onPlantItemClick(item)">
                            <image 
                                v-if="item.data.cover && item.data.cover.url"
                                class="plant-image" 
                                :src="item.data.cover.url" 
                                mode="aspectFill"
                            ></image>
                            <view v-else class="plant-image placeholder">
                                <uni-icons type="image" size="30" color="#ccc"></uni-icons>
                            </view>
                            <view class="plant-info">
                                <text class="plant-name">{{ item.name }}</text>
                                <text v-if="item.data.desc" class="plant-desc">{{ item.data.desc }}</text>
                            </view>
                        </view>
                    </view>
                </scroll-view>
                
                <!-- 右侧索引条 -->
                <view class="index-bar">
                    <view 
                        v-for="(group, index) in indexedOptions" 
                        :key="index" 
                        class="index-item"
                        @click="scrollToGroup(group.key)"
                    >
                        {{ group.key }}
                    </view>
                </view>
            </view>
        </view>

    </view>
</template>

<script>
import { callContainer } from '../utils/request.js';

export default {
    name: 'plantList',
    emits: [],
    data() {
        return {
            loading: true,
            plantsList: [],
            indexedOptions: [],
            scrollIntoView: '',
            currentFamilyId: null,
            
            // 系统信息
            statusBarHeight: 0
        };
    },

    mounted() {
        this.initPageInfo();
        this.loadPlantsList();
        
        // 监听家庭切换事件
        uni.$on('familyChanged', this.handleFamilyChanged);
    },
    
    beforeUnmount() {
        // 移除监听
        uni.$off('familyChanged', this.handleFamilyChanged);
    },

    // 页面显示时检查家庭是否变化
    onShow() {
        this.checkFamilyChange();
    },

    methods: {
        // 初始化页面信息
        initPageInfo() {
            try {
                const systemInfo = uni.getSystemInfoSync();
                this.statusBarHeight = systemInfo.statusBarHeight || 44;
            } catch (error) {
                console.error("获取系统信息失败:", error);
                this.statusBarHeight = 44;
            }
        },
        
        // 检查家庭是否变化
        async checkFamilyChange() {
            try {
                const familyIdResult = await new Promise((resolve, reject) => {
                    uni.getStorage({ key: 'familyId', success: resolve, fail: reject })
                });
                
                const newFamilyId = familyIdResult?.data;
                
                // 如果家庭ID变化了，重新加载
                if (newFamilyId && newFamilyId !== this.currentFamilyId) {
                    console.log('检测到家庭切换，重新加载植物列表');
                    this.currentFamilyId = newFamilyId;
                    await this.loadPlantsList();
                }
            } catch (error) {
                console.error('检查家庭变化失败:', error);
            }
        },
        
        // 处理家庭切换事件
        async handleFamilyChanged(newFamilyId) {
            console.log('收到家庭切换事件:', newFamilyId);
            this.currentFamilyId = newFamilyId;
            this.plantsList = [];
            this.indexedOptions = [];
            await this.loadPlantsList();
        },

        // 加载植物列表
        async loadPlantsList() {
            try {
                this.loading = true;

                // 获取当前家庭ID
                const familyIdResult = await new Promise((resolve, reject) => {
                    uni.getStorage({ key: 'familyId', success: resolve, fail: reject })
                });
                
                const familyId = familyIdResult?.data;
                
                if (!familyId) {
                    uni.showToast({
                        title: '请先选择家庭',
                        icon: 'none'
                    });
                    this.loading = false;
                    return;
                }
                
                // 更新当前家庭ID
                this.currentFamilyId = familyId;

                const result = await callContainer('/api/plant/list', {
                    familyId: familyId
                });

                console.log('植物列表数据:', result);
                this.plantsList = result?.data || [];
                this.processIndexedData();

            } catch (error) {
                console.error('加载植物列表失败:', error);
                uni.showToast({
                    title: '加载失败，请重试',
                    icon: 'none'
                });
            } finally {
                this.loading = false;
            }
        },

        // 处理数据为索引列表格式
        processIndexedData() {
            if (!this.plantsList || this.plantsList.length === 0) {
                this.indexedOptions = [];
                return;
            }

            const groupMap = {};

            this.plantsList.forEach(plant => {
                if (!plant.name) return;

                let firstChar = this.getFirstLetter(plant.name);

                if (!groupMap[firstChar]) {
                    groupMap[firstChar] = [];
                }

                // uni-indexed-list 默认显示 name 字段
                groupMap[firstChar].push({
                    name: plant.name,
                    id: plant.ID,
                    data: plant
                });
            });

            this.indexedOptions = Object.keys(groupMap)
                .sort()
                .map(key => ({
                    key: key.toUpperCase(),
                    data: groupMap[key]
                }));

            console.log('处理后的索引数据:', this.indexedOptions);
        },

        // 获取字符串首字母
        getFirstLetter(str) {
            if (!str) return '#';

            const firstChar = str.charAt(0).toUpperCase();

            // 如果已经是英文字母，直接返回
            if (/^[A-Z]$/.test(firstChar)) {
                return firstChar;
            }

            // 如果是数字，返回 #
            if (/^[0-9]$/.test(firstChar)) {
                return '#';
            }

            // 使用 Unicode 范围判断并转换中文到拼音首字母
            const code = str.charCodeAt(0);
            
            // 常用汉字 Unicode 范围判断
            if (code >= 19968 && code <= 40869) {
                // 使用更完整的拼音首字母映射
                return this.getChinesePinyin(str.charAt(0));
            }

            return '#';
        },
        
        // 汉字转拼音首字母（使用边界汉字对比法，最准确）
        getChinesePinyin(char) {
            // 预设每个字母开头的代表性边界汉字
            const letters = "ABCDEFGHJKLMNOPQRSTWXYZ".split("");
            const boundary = "阿八嚓哒妸发旮哈讥咔垃妈拿喔妑七然仨他哇夕丫匝".split("");
            
            // 兜底：处理你提到的特殊字或多音字
            const overrides = {
                '梦': 'M', '重': 'C', '长': 'C', '行': 'X'
            };
            if (overrides[char]) return overrides[char];

            // 使用 localeCompare 比较拼音权重
            // zh 语言环境会强制按拼音排序
            for (let i = 0; i < boundary.length; i++) {
                if (char.localeCompare(boundary[i], 'zh') < 0) {
                    return i === 0 ? 'A' : letters[i - 1];
                }
            }
            
            return 'Z';
        },

        // 点击植物项
        onPlantClick(e) {
            console.log('点击植物:', e);
            const plantData = e.item?.data;

            if (plantData && plantData.ID) {
                uni.navigateTo({
                    url: `/pages/plantDetail/plantDetail?id=${plantData.ID}`
                });
            }
        },
        
        // 自定义列表项点击
        onPlantItemClick(item) {
            console.log('点击植物项:', item);
            if (item && item.data && item.data.ID) {
                uni.navigateTo({
                    url: `/pages/plantDetail/plantDetail?id=${item.data.ID}`
                });
            }
        },
        
        // 滚动到指定分组
        scrollToGroup(key) {
            this.scrollIntoView = 'group-' + key;
            setTimeout(() => {
                this.scrollIntoView = '';
            }, 300);
        }
    }
};
</script>

<style lang="scss" scoped>
.plantlist-container {
    width: 100%;
    height: 100vh;
    display: flex;
    flex-direction: column;
    padding-bottom: calc(100rpx + env(safe-area-inset-bottom));
    box-sizing: border-box;
}

// 固定头部
.fixed-header-group {
    position: sticky;
    top: 0;
    z-index: 998;
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
}

.header-title-container {
    padding: 16px 20px 12px;
    text-align: center;
}

.page-title {
    font-size: 18px;
    font-weight: 600;
    color: #333;
}

// 内容区域
.content-container {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

.loading-wrapper {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 100rpx 0;
}

.loading-text {
    font-size: 14px;
    color: #999;
}

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

.indexed-list-wrapper {
    flex: 1;
    overflow: hidden;
    position: relative;
}

// 滚动视图
.plant-scroll-view {
    height: 100%;
    width: 100%;
}

// 分组容器
.plant-group {
    margin-bottom: 20rpx;
}

// 分组标题
.group-title {
    padding: 16rpx 32rpx;
    font-size: 14px;
    font-weight: 600;
    color: #6B8857;
    background: rgba(107, 136, 87, 0.1);
    position: sticky;
    top: 0;
    z-index: 10;
}

// 自定义植物项样式
.plant-item {
    padding: 24rpx 32rpx;
    background: rgba(255, 255, 255, 0.55);
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);
    transition: background 0.2s;
    display: flex;
    align-items: center;
    gap: 24rpx;
    
    &:active {
        background: rgba(107, 136, 87, 0.1);
    }
}

.plant-image {
    width: 100rpx;
    height: 100rpx;
    border-radius: 12rpx;
    flex-shrink: 0;
    background: #f5f5f5;
    
    &.placeholder {
        display: flex;
        align-items: center;
        justify-content: center;
        background: #f5f5f5;
    }
}

.plant-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8rpx;
    overflow: hidden;
}

.plant-name {
    font-size: 15px;
    color: #333;
    font-weight: 500;
}

.plant-desc {
    font-size: 12px;
    color: #999;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

// 右侧索引条
.index-bar {
    position: absolute;
    right: 10rpx;
    top: 50%;
    transform: translateY(-50%);
    display: flex;
    flex-direction: column;
    align-items: center;
    z-index: 100;
}

.index-item {
    padding: 4rpx 8rpx;
    font-size: 12px;
    color: #6B8857;
    font-weight: 500;
    line-height: 1.2;
    
    &:active {
        color: #fff;
        background: #6B8857;
        border-radius: 50%;
    }
}

// 索引列表样式优化
::v-deep .uni-indexed-list {
    height: 100%;
}

::v-deep .uni-indexed-list__scroll {
    height: 100%;
}

::v-deep .uni-indexed-list__item {
    padding: 12px 20px;
    background: rgba(255, 255, 255, 0.5);
    border-radius: 8px;
    margin: 8px 16px;
    transition: all 0.2s;

    &:active {
        background: rgba(255, 255, 255, 0.8);
        transform: scale(0.98);
    }
}

::v-deep .uni-indexed-list__item-content {
    font-size: 15px;
    color: #333;
}

::v-deep .uni-indexed-list__menu {
    right: 5px;
}
</style>
