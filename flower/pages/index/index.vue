<template>
  <view class="app-shell theme-page" :style="themeStyle">
    <loadingPage ref="loading" />

    <view class="page-stage">
      <view v-show="activeTab === 'photos'" class="tab-page"><photoAlbum ref="photoAlbumComponent" /></view>
      <view v-show="activeTab === 'garden'" class="tab-page"><home ref="homeComponent" /></view>
      <view v-show="activeTab === 'calendar'" class="tab-page"><logCalendar ref="calendarComponent" /></view>
      <view v-show="activeTab === 'mine'" class="tab-page"><my /></view>
    </view>

    <view class="bottom-nav">
      <view class="center-add-button" @click="addPlant">
        <uni-icons type="plusempty" size="29" color="var(--button-text)" />
      </view>
      <view
        v-for="item in tabs"
        :key="item.key"
        class="nav-item"
        :class="{ active: activeTab === item.key }"
        @click="switchTab(item.key)"
      >
        <view class="nav-icon-wrap">
          <uni-icons :type="item.icon" size="23" :color="activeTab === item.key ? 'var(--button-text)' : 'var(--text-sub)'" />
        </view>
        <text>{{ item.label }}</text>
      </view>
    </view>
  </view>
</template>

<script>
import home from '@/components/home.vue'
import my from '@/components/my.vue'
import loadingPage from '@/components/loading.vue'
import logCalendar from '@/components/logCalendar.vue'
import photoAlbum from '@/components/photoAlbum.vue'

export default {
  components: { home, my, loadingPage, logCalendar, photoAlbum },
  data() {
    return {
      activeTab: 'garden',
      tabs: [
        { key: 'photos', label: '照片', icon: 'image' },
        { key: 'garden', label: '花园', icon: 'home' },
        { key: 'calendar', label: '日历', icon: 'calendar' },
        { key: 'mine', label: '我的', icon: 'person-filled' }
      ]
    }
  },
  onLoad(options) {
    if (options && options.tab === 'photos') this.activeTab = 'photos'
    this.$nextTick(() => {
      if (this.$refs.loading) {
        this.$refs.loading.open()
        setTimeout(() => this.$refs.loading && this.$refs.loading.close(), 700)
      }
    })
    uni.$on('switchTab', this.handleExternalTab)
  },
  onShow() {
    this.$nextTick(() => {
      if (this.$refs.homeComponent) this.$refs.homeComponent.onPageShow()
      if (this.$refs.calendarComponent) this.$refs.calendarComponent.initData()
    })
  },
  onUnload() {
    uni.$off('switchTab', this.handleExternalTab)
  },
  methods: {
    addPlant() {
      // #ifdef MP-WEIXIN
      wx.vibrateShort({ type: 'medium' })
      // #endif
      uni.navigateTo({ url: '/pages/plantEdit/plantEdit?type=add' })
    },
    switchTab(key) {
      if (this.activeTab === key) return
      this.activeTab = key
      uni.$emit('tabChange', key)
      // #ifdef MP-WEIXIN
      wx.vibrateShort({ type: 'light' })
      // #endif
    },
    handleExternalTab(key) {
      const legacyMap = { a: 'photos', b: 'garden', c: 'garden', d: 'calendar', e: 'mine' }
      this.switchTab(legacyMap[key] || key)
    }
  },
  onShareAppMessage() {
    const albumPage = this.$refs.photoAlbumComponent
    if (this.activeTab === 'photos' && albumPage) {
      const familyId = albumPage.currentFamilyId || uni.getStorageSync('familyId') || 0
      return {
        title: `${albumPage.currentFamilyName}的植物相册`,
        path: `/pages/index/index?tab=photos&familyId=${familyId}`,
        imageUrl: albumPage.shareCoverUrl
      }
    }
    const homePage = this.$refs.homeComponent
    if (this.activeTab === 'garden' && homePage) {
      const familyId = homePage.value || 0
      const tag = homePage.tagList && homePage.tagList[homePage.currentTagIndex]
      const tagId = tag ? tag.id : 0
      const tagName = tag ? tag.name : ''
      return {
        title: tagName && tagName !== '全部' ? `分享我的「${tagName}」植物` : '分享我的小花园',
        path: `/pages/sharePlantList/sharePlantList?familyId=${familyId}&tagId=${tagId}&tagName=${encodeURIComponent(tagName)}`,
        imageUrl: '/static/images/share-garden-decorated.jpg'
      }
    }
    return { title: '我的花园', path: '/pages/index/index', imageUrl: '/static/images/share-garden-decorated.jpg' }
  },
  onShareTimeline() {
    return { title: '爱护绿植，从记录开始', imageUrl: '/static/images/share-garden-decorated.jpg' }
  }
}
</script>

<style lang="scss" scoped>
.app-shell {
  height: 100vh;
  overflow: hidden;
}

.page-stage {
  height: 100%;
}

.tab-page {
  height: 100%;
}

.bottom-nav {
  position: fixed;
  z-index: 1000;
  left: 20rpx;
  right: 20rpx;
  bottom: calc(14rpx + env(safe-area-inset-bottom));
  height: 104rpx;
  padding: 8rpx 18rpx;
  display: flex;
  align-items: center;
  justify-content: space-around;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 999rpx;
  box-shadow: 0 10rpx 34rpx var(--shadow-color);
}

.center-add-button {
  position: absolute;
  z-index: 2;
  left: 50%;
  top: 50%;
  width: 76rpx;
  height: 76rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 8rpx solid var(--card-bg);
  border-radius: 50%;
  background: var(--primary-color);
  box-shadow: 0 8rpx 22rpx var(--shadow-color);
  transform: translate(-50%, -50%);
}

.nav-item {
  width: 112rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2rpx;
  color: var(--text-sub);
  font-size: 21rpx;
}

.nav-item:nth-of-type(3) {
  margin-right: 54rpx;
}

.nav-item:nth-of-type(4) {
  margin-left: 54rpx;
}

.nav-icon-wrap {
  width: 54rpx;
  height: 54rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color .2s ease;
}

.nav-item.active {
  color: var(--primary-strong);
  font-weight: 600;
}

.nav-item.active .nav-icon-wrap {
  background: var(--primary-color);
}
</style>
