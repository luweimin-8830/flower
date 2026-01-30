<template>
    <view class="album-container">
        <!-- 顶部状态栏占位 -->
        <view :style="{ height: topBarHeight + 'px' }"></view>

        <!-- 内容区域 -->
        <view class="content-container">
            <view v-if="loading" class="loading-wrapper">
                <text class="loading-text">加载中...</text>
            </view>

            <view v-else-if="photoList.length === 0" class="empty-wrapper">
                <image src="/static/c2m.svg" class="empty-icon" mode="aspectFit"></image>
                <text class="empty-text">暂无照片数据</text>
            </view>

            <scroll-view v-else scroll-y class="photo-scroll-view" @scrolltolower="onReachBottom">
                <view class="photo-grid">
                    <view v-for="(photo, index) in photoList" :key="index" class="photo-item" @click="previewImage(index)">
                        <image :src="photo.url" mode="aspectFill" class="photo-image" lazy-load></image>
                        <view class="photo-info" v-if="photo.plantName">
                            <text class="plant-name-tag">{{ photo.plantName }}</text>
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
            photoList: [],
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
                    
                    // 按时间倒序排列（如果后端没排好）
                    this.photoList = photos;
                }
            } catch (error) {
                console.error('加载照片失败:', error);
            } finally {
                this.loading = false;
            }
        },
        previewImage(index) {
            const urls = this.photoList.map(p => p.url);
            uni.previewImage({
                current: urls[index],
                urls: urls
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

.photo-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 2rpx;
    padding: 2rpx;
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
