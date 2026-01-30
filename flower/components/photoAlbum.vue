<template>
    <view class="album-container">
        <!-- 顶部状态栏占位 -->
        <view :style="{ height: topBarHeight + 'px' }"></view>

        <!-- 内容区域 -->
        <view class="content-container">
            <view v-if="loading" class="loading-wrapper">
                <text class="loading-text">加载中...</text>
            </view>

            <view v-else-if="groupedPhotos.length === 0" class="empty-wrapper">
                <image src="/static/c2m.svg" class="empty-icon" mode="aspectFit"></image>
                <text class="empty-text">暂无照片数据</text>
            </view>

            <scroll-view v-else scroll-y class="photo-scroll-view" @scrolltolower="onReachBottom">
                <view v-for="(group, gIndex) in groupedPhotos" :key="group.date" class="date-group">
                    <view class="date-header">
                        <text class="date-text">{{ group.displayDate }}</text>
                    </view>
                    <view class="photo-grid">
                        <view v-for="(photo, pIndex) in group.photos" :key="photo.id" class="photo-item" @click="previewImage(photo.url)">
                            <image :src="photo.url" mode="aspectFill" class="photo-image" lazy-load></image>
                            <view class="photo-info" v-if="photo.plantName">
                                <text class="plant-name-tag">{{ photo.plantName }}</text>
                            </view>
                        </view>
                    </view>
                </view>
                <!-- 底部垫片 -->
                <view class="safe-area-bottom"></view>
            </scroll-view>
        </view>
    </view>
</template>

<script>
import { callContainer } from '../utils/request.js';

export default {
    name: 'photoAlbum',
    data() {
        return {
            loading: true,
            groupedPhotos: [],
            statusBarHeight: 0,
            topBarHeight: 0,
            currentFamilyId: null
        };
    },
    mounted() {
        this.initPageInfo();
        this.loadPhotos();
        
        // 监听家庭切换
        uni.$on('familyChanged', this.handleFamilyChanged);
    },
    beforeUnmount() {
        uni.$off('familyChanged', this.handleFamilyChanged);
    },
    methods: {
        initPageInfo() {
            try {
                const systemInfo = uni.getSystemInfoSync();
                this.statusBarHeight = systemInfo.statusBarHeight || 44;
                const app = getApp();
                if (app && app.globalData) {
                    this.topBarHeight = app.globalData.topBarHeight || this.statusBarHeight;
                } else {
                    this.topBarHeight = this.statusBarHeight;
                }
            } catch (error) {
                this.statusBarHeight = 44;
                this.topBarHeight = 44;
            }
        },
        async handleFamilyChanged(newFamilyId) {
            this.currentFamilyId = newFamilyId;
            await this.loadPhotos();
        },
        async loadPhotos() {
            try {
                this.loading = true;
                const familyId = uni.getStorageSync('familyId');
                if (!familyId) {
                    this.loading = false;
                    return;
                }
                this.currentFamilyId = familyId;

                // 获取家庭日志，日志中包含图片
                const result = await callContainer('/api/plant/log/list', {
                    familyId: parseInt(familyId)
                });

                if (result && result.data) {
                    const photos = [];
                    result.data.forEach(log => {
                        if (log.images && log.images.length > 0) {
                            log.images.forEach(img => {
                                photos.push({
                                    url: img.url,
                                    id: img.id,
                                    plantName: log.plant ? log.plant.name : '',
                                    logTime: log.logTime
                                });
                            });
                        }
                    });
                    
                    // 按日期分组
                    const groups = {};
                    photos.forEach(photo => {
                        // 提取 YYYY-MM-DD
                        const dateStr = photo.logTime ? photo.logTime.split('T')[0] : '未知日期';
                        if (!groups[dateStr]) {
                            groups[dateStr] = [];
                        }
                        groups[dateStr].push(photo);
                    });

                    // 排序并格式化展示日期
                    const sortedDates = Object.keys(groups).sort((a, b) => b.localeCompare(a));
                    this.groupedPhotos = sortedDates.map(date => {
                        let displayDate = date;
                        if (date !== '未知日期') {
                            const d = new Date(date);
                            displayDate = `${d.getMonth() + 1}月${d.getDate()}日`;
                            // 如果不是今年，增加年份
                            if (d.getFullYear() !== new Date().getFullYear()) {
                                displayDate = `${d.getFullYear()}年${displayDate}`;
                            }
                        }
                        return {
                            date,
                            displayDate,
                            photos: groups[date]
                        };
                    });
                }
            } catch (error) {
                console.error('加载照片失败:', error);
            } finally {
                this.loading = false;
            }
        },
        previewImage(currentUrl) {
            const allUrls = [];
            this.groupedPhotos.forEach(group => {
                group.photos.forEach(p => allUrls.push(p.url));
            });
            uni.previewImage({
                current: currentUrl,
                urls: allUrls
            });
        },
        onReachBottom() {
            // 分页逻辑可以在此添加
        }
    }
};
</script>

<style lang="scss" scoped>
.album-container {
    width: 100%;
    height: 100vh;
    display: flex;
    flex-direction: column;
    background-color: #C1D0B7;
}

.content-container {
    flex: 1;
    overflow: hidden;
}

.photo-scroll-view {
    height: 100%;
}

.date-group {
    margin-bottom: 10rpx;
}

.date-header {
    padding: 20rpx 30rpx 10rpx;
}

.date-text {
    font-size: 32rpx;
    font-weight: bold;
    color: #566C44;
}

.photo-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 2rpx;
    padding: 0 2rpx;
}

.photo-item {
    aspect-ratio: 1 / 1;
    position: relative;
    overflow: hidden;
}

.photo-image {
    width: 100%;
    height: 100%;
}

.photo-info {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 4rpx 8rpx;
    background: linear-gradient(transparent, rgba(0,0,0,0.4));
}

.plant-name-tag {
    font-size: 10px;
    color: #fff;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    display: block;
}

.loading-wrapper, .empty-wrapper {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding: 100rpx 0;
}

.loading-text, .empty-text {
    margin-top: 20rpx;
    font-size: 14px;
    color: #999;
}

.empty-icon {
    width: 120rpx;
    height: 120rpx;
    opacity: 0.4;
    filter: grayscale(100%);
}

.safe-area-bottom {
    height: calc(120rpx + env(safe-area-inset-bottom));
}
</style>
