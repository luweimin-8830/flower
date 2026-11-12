<template>
    <scroll-view scroll-y class="mine-page theme-page" :style="themeStyle">
        <view :style="{ height: topBarHeight + 'px' }"></view>
        <view class="mine-content">
            <view class="family-summary">
                <view>
                    <text class="family-name">{{ currentFamily.name }}</text>
                    <text class="family-meta">{{ currentFamily.memberCount }}位成员 · {{ currentFamily.plantCount }}株植物</text>
                </view>
                <view class="family-mark">
                    <uni-icons type="home-filled" size="24" color="var(--button-text)"></uni-icons>
                </view>
            </view>

            <view v-for="section in menuSections" :key="section.title" class="menu-section">
                <text class="section-title">{{ section.title }}</text>
                <view class="menu-card">
                    <template v-for="(item, index) in section.items" :key="item.type">
                        <picker v-if="item.type === 'reminder'" mode="time" :value="remindTime" @change="handleTimeChange">
                            <view class="menu-item" hover-class="item-hover">
                                <view class="menu-main">
                                    <view class="menu-icon"><uni-icons :type="item.icon" size="20" color="var(--primary-color)"></uni-icons></view>
                                    <text class="menu-text">{{ item.label }}</text>
                                </view>
                                <view class="menu-extra"><text>{{ remindTime }}</text><uni-icons type="right" size="15" color="var(--text-sub)"></uni-icons></view>
                            </view>
                        </picker>
                        <view v-else-if="item.type === 'theme'" class="menu-item" hover-class="item-hover" @click="openThemePopup">
                                <view class="menu-main">
                                    <view class="menu-icon"><uni-icons :type="item.icon" size="20" color="var(--primary-color)"></uni-icons></view>
                                    <text class="menu-text">{{ item.label }}</text>
                                </view>
                                <view class="menu-extra">
                                    <view class="theme-dot" :style="{ background: activeTheme.swatch }"></view>
                                    <text>{{ activeTheme.name }}</text>
                                    <uni-icons type="right" size="15" color="var(--text-sub)"></uni-icons>
                                </view>
                        </view>
                        <button v-else-if="item.type === 'contact'" class="contact-btn" open-type="contact" hover-class="item-hover">
                            <view class="menu-item">
                                <view class="menu-main"><view class="menu-icon"><uni-icons :type="item.icon" size="20" color="var(--primary-color)"></uni-icons></view><text class="menu-text">{{ item.label }}</text></view>
                                <uni-icons type="right" size="15" color="var(--text-sub)"></uni-icons>
                            </view>
                        </button>
                        <view v-else class="menu-item" hover-class="item-hover" @click="handleNav(item.type)">
                            <view class="menu-main">
                                <view class="menu-icon"><uni-icons :type="item.icon" size="20" color="var(--primary-color)"></uni-icons></view>
                                <text class="menu-text">{{ item.label }}</text>
                            </view>
                            <uni-icons type="right" size="15" color="var(--text-sub)"></uni-icons>
                        </view>
                        <view v-if="index < section.items.length - 1" class="divider"></view>
                    </template>
                </view>
            </view>
            <view class="bottom-space"></view>
        </view>

        <uni-popup ref="themePopup" type="bottom">
            <view class="theme-popup">
                <view class="sheet-handle"></view>
                <view class="theme-popup-head">
                    <text class="theme-popup-title">选择主题配色</text>
                    <text class="theme-popup-tip">给花园换一种心情</text>
                </view>
                <view class="theme-options">
                    <view v-for="option in themeOptions" :key="option.key" class="theme-option"
                        :class="{ active: option.key === themeKey }" hover-class="theme-option-hover"
                        @click="selectTheme(option.key)">
                        <view class="theme-swatch" :style="{ background: option.swatch }">
                        </view>
                        <text class="theme-option-name">{{ option.name }}</text>
                        <view class="theme-option-check">
                            <uni-icons v-if="option.key === themeKey" type="checkmarkempty" size="21" color="var(--primary-color)"></uni-icons>
                        </view>
                    </view>
                </view>
                <view class="theme-popup-close" hover-class="item-hover" @click="closeThemePopup">取消</view>
            </view>
        </uni-popup>
    </scroll-view>
</template>

<script>
import { callContainer } from '../utils/request.js'
import { themes } from '../utils/theme.js'

export default {
    name: 'my',
    data() {
        return {
            topBarHeight: 0,
            remindTime: '08:00',
            themeKeys: Object.keys(themes),
            currentFamily: { name: '南窗小花园', memberCount: 3, plantCount: 12 },
            menuSections: [
                { title: '花园', items: [
                    { type: 'family', label: '家庭管理', icon: 'staff-filled' },
                    { type: 'display', label: '陈列室', icon: 'medal-filled' }
                ]},
                { title: '偏好设置', items: [
                    { type: 'reminder', label: '提醒时间', icon: 'notification-filled' },
                    { type: 'theme', label: '主题配色', icon: 'color-filled' }
                ]},
                { title: '数据管理', items: [
                    { type: 'tag', label: '植物标签管理', icon: 'flag-filled' },
                    { type: 'care', label: '日常养护管理', icon: 'settings-filled' },
                    { type: 'charts', label: '图表分析', icon: 'bars' }
                ]},
                { title: '帮助与支持', items: [
                    { type: 'manual', label: '使用手册', icon: 'help-filled' },
                    { type: 'contact', label: '联系客服', icon: 'chat-filled' }
                ]}
            ]
        }
    },
    computed: {
        themeOptions() {
            return this.themeKeys.map(key => ({
                key,
                name: themes[key].name,
                swatch: themes[key].vars['--primary-soft'],
                soft: themes[key].vars['--primary-soft']
            }))
        }
    },
    methods: {
        loadUserData() {
            const userInfo = uni.getStorageSync('userInfo') || {}
            if (userInfo.remindTime) this.remindTime = userInfo.remindTime
            const familyId = uni.getStorageSync('familyId')
            const families = uni.getStorageSync('family') || []
            const family = families.find(item => String(item.id) === String(familyId)) || families[0]
            if (family) {
                this.currentFamily = {
                    name: family.name || '我的花园',
                    memberCount: family.memberCount ?? family.member_count ?? 1,
                    plantCount: family.plantCount ?? family.plant_count ?? 0
                }
            }
        },
        openThemePopup() {
            this.syncGardenTheme()
            wx.vibrateShort({ type: 'light' })
            this.$refs.themePopup.open()
        },
        closeThemePopup() {
            this.$refs.themePopup.close()
            this.$nextTick(() => this.syncGardenTheme())
        },
        selectTheme(key) {
            this.setGardenTheme(key)
            this.closeThemePopup()
            wx.vibrateShort({ type: 'light' })
            uni.showToast({ title: `已切换为${themes[key].name}`, icon: 'none' })
        },
        async handleTimeChange(e) {
            const newTime = e.detail.value
            this.remindTime = newTime
            try {
                uni.showLoading({ title: '保存中...' })
                await callContainer('/api/user/update', { remindTime: newTime })
                const userInfo = uni.getStorageSync('userInfo') || {}
                userInfo.remindTime = newTime
                uni.setStorageSync('userInfo', userInfo)
                uni.showToast({ title: '设置已更新', icon: 'success' })
                wx.vibrateShort({ type: 'medium' })
            } catch (error) {
                console.error('更新提醒时间失败:', error)
                uni.showToast({ title: '更新失败', icon: 'none' })
            } finally {
                uni.hideLoading()
            }
        },
        handleNav(type) {
            wx.vibrateShort({ type: 'medium' })
            const routeMap = {
                tag: '/pages/tagEdit/tagEdit',
                care: '/pages/careEdit/careEdit',
                family: '/pages/familyDetail/familyDetail',
                display: '/pages/displayRoom/displayRoom',
                charts: '/subpackages/analytics/charts/charts',
                manual: '/pages/manual/manual'
            }
            if (routeMap[type]) uni.navigateTo({ url: routeMap[type] })
        }
    },
    created() {
        const app = getApp()
        this.topBarHeight = app.globalData.topBarHeight || 44
        this.loadUserData()
        this.syncGardenTheme()
    },
    activated() {
        this.syncGardenTheme()
    }
}
</script>

<style scoped lang="scss">
.mine-page { width: 100%; height: 100vh; color: var(--text-color); }
.mine-content { min-height: 100%; padding: 8px 16px 0; box-sizing: border-box; }
.family-summary { display: flex; align-items: center; justify-content: space-between; padding: 22px 20px; margin-bottom: 24px; background: var(--primary-soft); border: 1px solid var(--border-color); border-radius: 22px; box-shadow: 0 8px 24px var(--shadow-color); }
.family-name { display: block; margin-bottom: 7px; font-size: 19px; font-weight: 650; }
.family-meta { display: block; font-size: 13px; color: var(--text-sub); }
.family-mark { display: flex; align-items: center; justify-content: center; width: 46px; height: 46px; border-radius: 50%; background: var(--primary-color); }
.menu-section { margin-bottom: 20px; }
.section-title { display: block; margin: 0 5px 8px; font-size: 13px; color: var(--text-sub); }
.menu-card { overflow: hidden; padding: 2px 16px; background: var(--card-bg); border: 1px solid var(--border-color); border-radius: 18px; box-shadow: 0 6px 18px var(--shadow-color); }
.menu-item { display: flex; align-items: center; justify-content: space-between; min-height: 58px; background: transparent; }
.menu-main, .menu-extra { display: flex; align-items: center; }
.menu-icon { display: flex; align-items: center; justify-content: center; width: 31px; height: 31px; margin-right: 10px; border-radius: 10px; background: var(--primary-soft); }
.menu-text { font-size: 15px; color: var(--text-color); }
.menu-extra { gap: 6px; font-size: 13px; color: var(--text-sub); }
.theme-dot { width: 12px; height: 12px; border: 2px solid var(--card-bg); border-radius: 50%; box-shadow: 0 0 0 1px var(--border-color); }
.divider { height: 1px; margin-left: 41px; background: var(--border-color); }
.item-hover { opacity: .68; }
.contact-btn { width: 100%; padding: 0; margin: 0; line-height: inherit; text-align: left; color: inherit; background: transparent; border: 0; }
.contact-btn::after { border: 0; }
.bottom-space { height: calc(140rpx + env(safe-area-inset-bottom)); }
.theme-popup { width: 100%; margin-bottom: 124rpx; padding: 8px 14px 12px; box-sizing: border-box; color: var(--text-color); background: var(--surface-bg); border-radius: 26px; box-shadow: 0 -10px 40px rgba(0, 0, 0, .16); }
.sheet-handle { width: 38px; height: 5px; margin: 0 auto 15px; border-radius: 3px; background: var(--border-color); }
.theme-popup-head { padding: 0 6px 14px; text-align: left; }
.theme-popup-title { display: block; font-size: 21px; line-height: 29px; font-weight: 700; letter-spacing: -.3px; color: var(--text-color); }
.theme-popup-tip { display: block; margin-top: 2px; font-size: 13px; line-height: 19px; color: var(--text-sub); }
.theme-options { overflow: hidden; background: var(--card-bg); border: 1px solid var(--border-color); border-radius: 16px; }
.theme-option { position: relative; display: flex; align-items: center; min-height: 58px; padding: 0 16px; background: var(--card-bg); box-sizing: border-box; }
.theme-option + .theme-option::before { content: ''; position: absolute; top: 0; right: 0; left: 59px; height: 1px; background: var(--border-color); transform: scaleY(.5); }
.theme-option.active { background: var(--primary-soft); }
.theme-option-hover { opacity: .62; }
.theme-swatch { flex: 0 0 auto; width: 28px; height: 28px; margin-right: 14px; border: 2px solid var(--card-bg); border-radius: 50%; box-shadow: 0 0 0 1px var(--border-color); }
.theme-option-name { flex: 1; font-size: 16px; font-weight: 500; letter-spacing: -.1px; color: var(--text-color); }
.theme-option-check { display: flex; align-items: center; justify-content: center; width: 26px; height: 26px; }
.theme-popup-close { height: 52px; margin-top: 9px; line-height: 52px; text-align: center; font-size: 17px; font-weight: 600; color: var(--primary-color); background: var(--card-bg); border: 1px solid var(--border-color); border-radius: 16px; }
</style>
