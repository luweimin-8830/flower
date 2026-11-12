<template>
    <view class="album-container">
        <!-- 顶部状态栏占位 -->
        <view :style="{ height: topBarHeight + 'px' }"></view>

        <view class="album-top-controls" :style="{
            top: menuButtonInfo.top + 'px',
            left: paddingLeft + 'px',
            height: menuButtonInfo.height + 'px'
        }">
            <picker class="family-picker-wrap" :style="{ height: menuButtonInfo.height + 'px' }"
                :value="currentFamilyIndex" :range="familyRange" :range-key="'name'" @change="handleFamilyChange">
                <view class="family-picker" :style="{ height: menuButtonInfo.height + 'px' }">
                    <uni-icons type="home" size="18" color="var(--primary-color)"></uni-icons>
                    <text>{{ familyRange[currentFamilyIndex]?.name || '选择花园' }}</text>
                    <uni-icons type="bottom" size="14" color="var(--primary-color)"></uni-icons>
                </view>
            </picker>
            <button class="album-share-button" :style="{ width: menuButtonInfo.height + 'px' }" open-type="share">
                <uni-icons type="upload" size="19" color="var(--text-color)"></uni-icons>
            </button>
            <view class="filter-button" :class="{ active: hasPhotoFilter }"
                :style="{ width: menuButtonInfo.height + 'px' }" @click="openFilterPopup">
                <uni-icons type="tune" size="19" :color="hasPhotoFilter ? 'var(--button-text)' : 'var(--text-color)'"></uni-icons>
                <view v-if="hasPhotoFilter" class="filter-active-dot"></view>
            </view>
        </view>

        <!-- 内容区域 -->
        <view class="content-container">
            <view v-if="loading" class="loading-wrapper">
                <text class="loading-text">加载中...</text>
            </view>

            <view v-else-if="filteredPhotoGroups.length === 0" class="empty-wrapper">
                <image src="/static/icon/c2m.svg" class="empty-icon" mode="aspectFit"></image>
                <text class="empty-text">{{ hasPhotoFilter ? '没有符合条件的照片' : '暂无照片，请去成长记录中添加' }}</text>
            </view>

            <scroll-view v-else scroll-y class="photo-scroll-view" @scrolltolower="onReachBottom">
                <view v-for="group in filteredPhotoGroups" :key="group.date" class="date-group">
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
            </scroll-view>
        </view>

        <uni-popup ref="filterPopup" type="bottom">
            <view class="filter-sheet">
                <view class="sheet-handle"></view>
                <view class="filter-sheet-header">
                    <view>
                        <text class="filter-sheet-title">搜索与筛选</text>
                        <text class="filter-sheet-subtitle">按植物名称和照片日期查找</text>
                    </view>
                    <view class="filter-reset" :class="{ disabled: !hasPhotoFilter }" @click="resetPhotoFilter">重置</view>
                </view>

                <view class="filter-search-box">
                    <uni-icons type="search" size="18" color="var(--text-sub)"></uni-icons>
                    <input v-model="photoKeyword" class="filter-search-input" placeholder="搜索植物名称"
                        placeholder-class="filter-placeholder" confirm-type="search" />
                    <view v-if="photoKeyword" class="search-clear" @click="photoKeyword = ''">
                        <uni-icons type="clear" size="18" color="var(--text-sub)"></uni-icons>
                    </view>
                </view>

                <text class="filter-section-label">照片日期</text>
                <view class="date-filter-row">
                    <picker class="date-filter-picker" mode="date" :value="filterStartDate" :end="filterEndDate || today"
                        @change="filterStartDate = $event.detail.value">
                        <view class="date-filter-value" :class="{ placeholder: !filterStartDate }">
                            {{ filterStartDate || '开始日期' }}
                        </view>
                    </picker>
                    <text class="date-filter-separator">至</text>
                    <picker class="date-filter-picker" mode="date" :value="filterEndDate" :start="filterStartDate" :end="today"
                        @change="filterEndDate = $event.detail.value">
                        <view class="date-filter-value" :class="{ placeholder: !filterEndDate }">
                            {{ filterEndDate || '结束日期' }}
                        </view>
                    </picker>
                </view>

                <view class="filter-result-tip">当前显示 {{ filteredPhotoCount }} 张照片</view>
                <button class="filter-done-button" @click="closeFilterPopup">完成</button>
            </view>
        </uni-popup>
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
            menuButtonInfo: {},
            paddingLeft: 0,
            currentFamilyId: null,
            familyRange: [],
            currentFamilyIndex: 0,
            photoKeyword: '',
            filterStartDate: '',
            filterEndDate: '',
            today: ''
        };
    },
    computed: {
        hasPhotoFilter() {
            return Boolean(this.photoKeyword.trim() || this.filterStartDate || this.filterEndDate);
        },
        filteredPhotoGroups() {
            const keyword = this.photoKeyword.trim().toLowerCase();
            return this.groupedPhotos.map(group => ({
                ...group,
                photos: group.photos.filter(photo => {
                    const matchesName = !keyword || String(photo.plantName || '').toLowerCase().includes(keyword);
                    const matchesStart = !this.filterStartDate || group.date >= this.filterStartDate;
                    const matchesEnd = !this.filterEndDate || group.date <= this.filterEndDate;
                    return matchesName && matchesStart && matchesEnd;
                })
            })).filter(group => group.photos.length > 0);
        },
        filteredPhotoCount() {
            return this.filteredPhotoGroups.reduce((count, group) => count + group.photos.length, 0);
        },
        currentFamilyName() {
            return this.familyRange[this.currentFamilyIndex]?.name || '我的花园';
        },
        shareCoverUrl() {
            return this.filteredPhotoGroups[0]?.photos?.[0]?.url || '/static/images/share-garden-decorated.jpg';
        }
    },
    mounted() {
        this.initPageInfo();
        this.loadPhotos();
        
        // 监听家庭切换
        uni.$on('familyChanged', this.handleFamilyChanged);
        
        // 监听数据刷新事件（与首页、详情页保持一致）
        uni.$on('refreshHomeList', () => {
            this.loadPhotos();
        });
        
        // 监听页面显示事件，每次页面显示时刷新数据
        uni.$on('tabChange', (tabIndex) => {
            // 假设相册页是特定的tab索引，根据实际情况调整
            if (tabIndex === 'photos') {
                this.loadPhotos();
            }
        });
    },
    beforeUnmount() {
        uni.$off('familyChanged', this.handleFamilyChanged);
        uni.$off('refreshHomeList');
        uni.$off('tabChange');
    },
    methods: {
        initPageInfo() {
            const now = new Date();
            this.today = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
            try {
                const systemInfo = uni.getSystemInfoSync();
                this.menuButtonInfo = uni.getMenuButtonBoundingClientRect();
                this.paddingLeft = systemInfo.screenWidth - this.menuButtonInfo.right;
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
                this.menuButtonInfo = { top: 8, width: 87, height: 32 };
                this.paddingLeft = 12;
            }
            const families = uni.getStorageSync('family') || [];
            this.familyRange = families;
            const familyId = uni.getStorageSync('familyId');
            const familyIndex = families.findIndex(item => String(item.id) === String(familyId));
            this.currentFamilyIndex = familyIndex < 0 ? 0 : familyIndex;
        },
        async handleFamilyChange(e) {
            const index = Number(e.detail.value);
            const family = this.familyRange[index];
            if (!family) return;
            this.currentFamilyIndex = index;
            this.currentFamilyId = family.id;
            uni.setStorageSync('familyId', family.id);
            try { await callContainer('/api/family/switch', { familyId: family.id }); } catch (error) { console.error('切换家庭失败:', error); }
            uni.$emit('familyChanged', family.id);
            this.loadPhotos();
        },
        async handleFamilyChanged(newFamilyId) {
            this.currentFamilyId = newFamilyId;
            await this.loadPhotos();
        },
        // 处理图片上传完成事件
        async handlePhotoUploaded() {
            console.log('检测到新图片上传，刷新相册');
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
            this.filteredPhotoGroups.forEach(group => {
                group.photos.forEach(p => allUrls.push(p.url));
            });
            uni.previewImage({
                current: currentUrl,
                urls: allUrls
            });
        },
        onReachBottom() {
            // 分页逻辑可以在此添加
        },
        openFilterPopup() {
            wx.vibrateShort({ type: 'light' });
            this.$refs.filterPopup.open();
        },
        closeFilterPopup() {
            this.$refs.filterPopup.close();
        },
        resetPhotoFilter() {
            if (!this.hasPhotoFilter) return;
            this.photoKeyword = '';
            this.filterStartDate = '';
            this.filterEndDate = '';
            wx.vibrateShort({ type: 'light' });
        }
    },
    // 添加页面显示时的生命周期钩子
    onShow() {
        // 页面显示时刷新数据
        this.loadPhotos();
    }
};
</script>

<style lang="scss" scoped>
.album-container {
    width: 100%;
    height: 100vh;
    display: flex;
    flex-direction: column;
    padding-bottom: calc(100rpx + env(safe-area-inset-bottom));
    box-sizing: border-box;
}

.content-container {
    flex: 1;
    overflow: hidden;
    padding-top: 12px;
}

.album-top-controls { position: fixed; z-index: 100; display: flex; align-items: center; gap: 8px; }
.family-picker-wrap { display: block; width: auto; flex-shrink: 0; box-sizing: border-box; }
.album-share-button { display: flex; flex-shrink: 0; align-items: center; justify-content: center; height: 100%; padding: 0; margin: 0; line-height: 1; border: 1px solid var(--border-color); border-radius: 50%; background: var(--card-bg); box-shadow: 0 2px 8px var(--shadow-color); }
.album-share-button::after { border: 0; }
.family-picker { display: flex; align-items: center; justify-content: flex-start; gap: 6px; width: auto; min-height: 32px; box-sizing: border-box; padding: 0 12px; border: 1px solid var(--border-color); border-radius: 999px; background: var(--card-bg); color: var(--text-color); font-size: 14px; font-weight: 500; line-height: 1; box-shadow: 0 2px 8px var(--shadow-color); }
.family-picker text { overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
.filter-button { position: relative; display: flex; flex-shrink: 0; align-items: center; justify-content: center; height: 100%; box-sizing: border-box; border-radius: 50%; background: var(--card-bg); border: 1px solid var(--border-color); box-shadow: 0 2px 8px var(--shadow-color); }
.filter-button.active { background: var(--primary-color); border-color: var(--primary-color); }
.filter-active-dot { position: absolute; top: -2px; right: -2px; width: 8px; height: 8px; border: 2px solid var(--page-bg); border-radius: 50%; background: #ff5a5f; }

.photo-scroll-view {
    height: 100%;
    box-sizing: border-box;
}

.date-group {
    margin-bottom: 10rpx;
}

.date-group:last-child {
    padding-bottom: 16px;
}

.date-header {
    padding: 20rpx 30rpx 10rpx;
}

.date-text {
    font-size: 32rpx;
    font-weight: bold;
    color: var(--text-color);
}

.photo-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 6rpx;
    padding: 0 20rpx;
}

.photo-item {
    aspect-ratio: 1 / 1;
    position: relative;
    overflow: hidden;
    border-radius: 9px;
    background: var(--surface-bg);
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
    color: var(--text-sub);
}

.empty-icon {
    width: 120rpx;
    height: 120rpx;
    opacity: 0.4;
    filter: grayscale(100%);

}

.filter-sheet { width: 100%; margin-bottom: 124rpx; padding: 8px 16px 16px; box-sizing: border-box; color: var(--text-color); background: var(--surface-bg); border-radius: 26px; box-shadow: 0 -10px 40px rgba(0, 0, 0, .16); }
.sheet-handle { width: 38px; height: 5px; margin: 0 auto 16px; border-radius: 3px; background: var(--border-color); }
.filter-sheet-header { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 18px; }
.filter-sheet-title { display: block; font-size: 21px; line-height: 29px; font-weight: 700; letter-spacing: -.3px; color: var(--text-color); }
.filter-sheet-subtitle { display: block; margin-top: 2px; font-size: 13px; color: var(--text-sub); }
.filter-reset { padding: 5px 2px 5px 12px; font-size: 15px; color: var(--primary-color); }
.filter-reset.disabled { opacity: .35; }
.filter-search-box { display: flex; align-items: center; height: 48px; padding: 0 14px; margin-bottom: 20px; background: var(--card-bg); border: 1px solid var(--border-color); border-radius: 14px; }
.filter-search-input { flex: 1; height: 48px; padding: 0 10px; font-size: 16px; color: var(--text-color); }
.filter-placeholder { color: var(--text-sub); }
.search-clear { display: flex; align-items: center; justify-content: center; padding: 5px 0 5px 8px; }
.filter-section-label { display: block; margin: 0 4px 8px; font-size: 13px; color: var(--text-sub); }
.date-filter-row { display: flex; align-items: center; gap: 9px; }
.date-filter-picker { flex: 1; }
.date-filter-value { height: 48px; line-height: 48px; text-align: center; font-size: 15px; color: var(--text-color); background: var(--card-bg); border: 1px solid var(--border-color); border-radius: 14px; }
.date-filter-value.placeholder { color: var(--text-sub); }
.date-filter-separator { flex-shrink: 0; font-size: 13px; color: var(--text-sub); }
.filter-result-tip { padding: 15px 2px 11px; font-size: 13px; color: var(--text-sub); }
.filter-done-button { height: 50px; padding: 0; margin: 0; line-height: 50px; font-size: 17px; font-weight: 600; color: var(--button-text); background: var(--primary-color); border: 0; border-radius: 15px; }
.filter-done-button::after { border: 0; }
</style>
