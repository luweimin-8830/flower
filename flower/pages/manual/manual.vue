<template>
    <view class="manual-page theme-page" :style="themeStyle">
        <navBar title="使用手册" />
        <view :style="{ height: topBarHeight + 'px' }"></view>
        <view class="manual-tabs">
            <view v-for="(tab, index) in tabs" :key="tab" class="manual-tab" :class="{ active: activeTab === index }" @click="activeTab = index">{{ tab }}</view>
        </view>

        <scroll-view scroll-y class="manual-scroll">
            <view class="manual-content">
                <template v-if="activeTab === 0">
                    <view class="welcome-card">
                        <text class="welcome-emoji">🌱</text>
                        <view><text class="content-title">欢迎来到我的花园</text><text class="content-desc">记录植物成长，与家人一起照顾每一株植物</text></view>
                    </view>
                    <text class="group-title">三步开始</text>
                    <view v-for="(item, index) in startSteps" :key="item.title" class="guide-card step-card">
                        <text class="step-number">{{ index + 1 }}</text>
                        <view class="guide-icon"><uni-icons :type="item.icon" size="25" color="var(--primary-color)"></uni-icons></view>
                        <view class="guide-copy"><text class="guide-title">{{ item.title }}</text><text class="guide-desc">{{ item.desc }}</text></view>
                    </view>
                    <view class="tip-card"><uni-icons type="info-filled" size="17" color="var(--primary-color)"></uni-icons><text>上传照片时，可像微信头像一样拖动和缩放后完成方形截取</text></view>
                </template>

                <template v-else-if="activeTab === 1">
                    <text class="panel-title">🌿 管理你的植物</text>
                    <view v-for="item in gardenGuides" :key="item.title" class="guide-card">
                        <view class="guide-icons"><uni-icons v-for="icon in item.icons" :key="icon" :type="icon" size="25" color="var(--primary-color)"></uni-icons></view>
                        <view class="guide-copy"><text class="guide-title">{{ item.title }}</text><text class="guide-desc">{{ item.desc }}</text></view>
                    </view>
                    <view class="tip-card"><uni-icons type="info-filled" size="18" color="var(--primary-color)"></uni-icons><text>右上角“+”可添加植物，分享图标可生成花园长图</text></view>
                </template>

                <template v-else-if="activeTab === 2">
                    <text class="panel-title">🌿 记录每一次养护</text>
                    <view class="care-icons">
                        <view v-for="item in careTypes" :key="item.name"><view class="care-icon"><uni-icons :type="item.icon" size="24" color="var(--primary-color)"></uni-icons></view><text>{{ item.name }}</text></view>
                    </view>
                    <view v-for="item in careGuides" :key="item.title" class="guide-card">
                        <view class="guide-icon"><uni-icons :type="item.icon" size="25" color="var(--primary-color)"></uni-icons></view>
                        <view class="guide-copy"><text class="guide-title">{{ item.title }}</text><text class="guide-desc">{{ item.desc }}</text></view>
                    </view>
                </template>

                <template v-else>
                    <text class="panel-title">🌿 一起分享成长</text>
                    <view v-for="item in shareGuides" :key="item.title" class="guide-card clickable" @click="handleGuideAction(item.action)">
                        <view class="guide-icon"><uni-icons :type="item.icon" size="25" color="var(--primary-color)"></uni-icons></view>
                        <view class="guide-copy"><text class="guide-title">{{ item.title }}</text><text class="guide-desc">{{ item.desc }}</text></view>
                        <uni-icons type="right" size="16" color="var(--text-sub)"></uni-icons>
                    </view>
                    <button class="contact-button" open-type="contact">联系客服</button>
                </template>
                <view class="safe-space"></view>
            </view>
        </scroll-view>
    </view>
</template>

<script>
import navBar from '@/components/navBar.vue'

export default {
    components: { navBar },
    data() {
        return {
            topBarHeight: 0,
            activeTab: 0,
            tabs: ['快速开始', '花园', '养护', '分享'],
            startSteps: [
                { title: '创建花园', desc: '设置家庭花园，邀请家人共同管理', icon: 'home-filled' },
                { title: '添加植物', desc: '上传照片，填写名称、日期和标签', icon: 'camera-filled' },
                { title: '开始记录', desc: '浇水、施肥、换盆都能轻松记下', icon: 'compose' }
            ],
            gardenGuides: [
                { title: '两种展示方式', desc: '在卡片模式与索引模式之间自由切换', icons: ['image-filled', 'list'] },
                { title: '搜索与标签', desc: '按名称搜索，用标签快速筛选', icons: ['search', 'flag-filled'] },
                { title: '添加植物', desc: '选择照片后，像微信头像一样移动和缩放', icons: ['camera-filled'] },
                { title: '批量操作', desc: '长按植物卡片，可批量记录养护', icons: ['checkbox-filled'] }
            ],
            careTypes: [
                { name: '浇水', icon: 'fire-filled' }, { name: '施肥', icon: 'wallet-filled' },
                { name: '换盆', icon: 'trash-filled' }, { name: '除虫', icon: 'settings-filled' }, { name: '修剪', icon: 'tune' }
            ],
            careGuides: [
                { title: '快捷养护', desc: '在植物详情点击一次即可记录', icon: 'redo-filled' },
                { title: '成长日志', desc: '添加文字和照片，留下变化过程', icon: 'compose' },
                { title: '养护日历', desc: '按日期查看记录、天气和提醒', icon: 'calendar-filled' },
                { title: '预约提醒', desc: '设置日期与时间，不错过下一次养护', icon: 'notification-filled' }
            ],
            shareGuides: [
                { title: '家庭协作', desc: '邀请家人加入，共同查看和记录植物', icon: 'staff-filled', action: 'family' },
                { title: '长图分享', desc: '选择植物和内容，生成长图保存或分享', icon: 'upload-filled', action: 'share' },
                { title: '图表分析', desc: '查看植物分类和每月养护次数', icon: 'bars', action: 'charts' },
                { title: '陈列室', desc: '纪念曾经陪伴过的植物', icon: 'images-filled', action: 'display' },
                { title: '主题配色', desc: '在“我的”中选择喜欢的界面颜色', icon: 'color-filled', action: 'back' }
            ]
        }
    },
    onLoad(options) {
        const app = getApp()
        this.topBarHeight = app.globalData.topBarHeight || 44
        if (options && options.tab) this.activeTab = Math.max(0, Math.min(3, Number(options.tab)))
    },
    methods: {
        handleGuideAction(action) {
            const routes = { family: '/pages/familyDetail/familyDetail', charts: '/subpackages/analytics/charts/charts', display: '/pages/displayRoom/displayRoom', share: '/pages/sharePlantList/sharePlantList' }
            if (routes[action]) uni.navigateTo({ url: routes[action] })
            else if (action === 'back') uni.navigateBack()
        }
    }
}
</script>

<style scoped lang="scss">
.manual-page { min-height: 100vh; color: var(--text-color); }
.manual-tabs { display: flex; gap: 8px; padding: 10px 16px 12px; }
.manual-tab { flex: 1; padding: 8px 3px; border: 1px solid var(--border-color); border-radius: 19px; text-align: center; font-size: 13px; background: var(--surface-bg); }
.manual-tab.active { color: var(--button-text); background: var(--primary-color); border-color: var(--primary-color); }
.manual-scroll { height: calc(100vh - 115px); }
.manual-content { padding: 4px 16px calc(30px + env(safe-area-inset-bottom)); }
.welcome-card, .guide-card, .tip-card { border: 1px solid var(--border-color); background: var(--card-bg); box-shadow: 0 5px 16px var(--shadow-color); }
.welcome-card { display: flex; align-items: center; gap: 18px; padding: 22px 20px; margin-bottom: 20px; border-radius: 18px; background: var(--primary-soft); }
.welcome-emoji { font-size: 52px; }
.content-title, .content-desc, .guide-title, .guide-desc, .panel-title, .group-title { display: block; }
.content-title { margin-bottom: 7px; font-size: 19px; font-weight: 650; }
.content-desc, .guide-desc { font-size: 13px; line-height: 1.55; color: var(--text-sub); }
.group-title, .panel-title { margin: 10px 3px 13px; font-size: 17px; font-weight: 650; }
.guide-card { display: flex; align-items: center; gap: 13px; min-height: 76px; padding: 12px 14px; margin-bottom: 11px; border-radius: 15px; }
.step-card { position: relative; padding-left: 17px; }
.step-number { display: flex; align-items: center; justify-content: center; position: absolute; left: -7px; top: -7px; width: 25px; height: 25px; border-radius: 50%; background: var(--primary-color); color: var(--button-text); font-size: 13px; }
.guide-icon { display: flex; align-items: center; justify-content: center; width: 47px; height: 47px; flex-shrink: 0; border-radius: 14px; background: var(--primary-soft); }
.guide-icons { display: flex; align-items: center; justify-content: space-evenly; width: 88px; height: 50px; flex-shrink: 0; border-radius: 14px; background: var(--primary-soft); }
.guide-copy { flex: 1; min-width: 0; }
.guide-title { margin-bottom: 5px; font-size: 16px; font-weight: 600; }
.tip-card { display: flex; align-items: flex-start; gap: 9px; padding: 14px; margin-top: 16px; border-radius: 14px; background: var(--primary-soft); font-size: 12px; line-height: 1.5; color: var(--text-sub); }
.care-icons { display: flex; justify-content: space-between; padding: 12px 8px 20px; }
.care-icons > view { display: flex; flex-direction: column; align-items: center; gap: 6px; font-size: 12px; }
.care-icon { display: flex; align-items: center; justify-content: center; width: 48px; height: 48px; border: 1px solid var(--border-color); border-radius: 13px; background: var(--card-bg); }
.clickable:active { opacity: .7; }
.contact-button { width: 180px; margin: 24px auto 0; color: var(--button-text); background: var(--primary-color); border-radius: 22px; font-size: 15px; }
.contact-button::after { border: 0; }
.safe-space { height: env(safe-area-inset-bottom); }
</style>
