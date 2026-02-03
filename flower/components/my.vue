<template>
    <view :style="{ height: topBarHeight + 'px' }"></view>
    <view class="section-container">
        <!-- 标题 -->
        <text class="section-title">个人信息</text>
        <!-- 卡片容器 -->
        <view class="card-box">
            <view class="menu-item" hover-class="item-hover" @click="handleNav('family')">
                <view class="left-content">
                    <uni-icons type="staff" size="22" color="var(--primary-color)" class="menu-icon"></uni-icons>
                    <text class="menu-text">家庭管理</text>
                </view>
                <uni-icons type="right" size="16" color="var(--text-sub)"></uni-icons>
            </view>
        </view>
    </view>
    <view class="section-container">
        <!-- 标题 -->
        <text class="section-title">设置</text>
        <!-- 卡片容器 -->
        <view class="card-box">
            <picker mode="time" :value="remindTime" @change="handleTimeChange">
                <view class="menu-item" hover-class="item-hover">
                    <view class="left-content">
                        <uni-icons type="notification" size="22" color="var(--primary-color)" class="menu-icon"></uni-icons>
                        <text class="menu-text">提醒时间</text>
                    </view>
                    <view class="right-content">
                        <text class="time-text">{{ remindTime }}</text>
                        <uni-icons type="right" size="16" color="var(--text-sub)"></uni-icons>
                    </view>
                </view>
            </picker>
        </view>
    </view>
    <view class="section-container">
        <!-- 标题 -->
        <text class="section-title">数据管理</text>
        <!-- 卡片容器 -->
        <view class="card-box">
            <view class="menu-item" hover-class="item-hover" @click="handleNav('tag')">
                <view class="left-content">
                    <uni-icons type="flag" size="22" color="var(--primary-color)" class="menu-icon"></uni-icons>
                    <text class="menu-text">植物标签管理</text>
                </view>
                <uni-icons type="right" size="16" color="var(--text-sub)"></uni-icons>
            </view>
            <!-- 分割线 -->
            <view class="divider"></view>
            <view class="menu-item" hover-class="item-hover" @click="handleNav('care')">
                <view class="left-content">
                    <uni-icons type="settings" size="22" color="var(--primary-color)" class="menu-icon"></uni-icons>
                    <text class="menu-text">日常养护管理</text>
                </view>
                <uni-icons type="right" size="16" color="var(--text-sub)"></uni-icons>
            </view>
        </view>
    </view>
    <view class="section-container">
        <text class="section-title">帮助与支持</text>
        <view class="card-box">
            <button class="contact-btn" open-type="contact" hover-class="item-hover">
                <view class="menu-item">
                    <view class="left-content">
                        <uni-icons type="chat" size="22" color="var(--primary-color)" class="menu-icon"></uni-icons>
                        <text class="menu-text">联系客服</text>
                    </view>
                    <uni-icons type="right" size="16" color="var(--text-sub)"></uni-icons>
                </view>
            </button>
        </view>
    </view>
</template>

<script>
import { callContainer } from '../utils/request.js';

export default {
    name: 'my',
    data() {
        return {
            topBarHeight: 0,
            remindTime: '08:00'
        }
    },
    methods: {
        async loadUserData() {
            const userInfo = uni.getStorageSync('userInfo');
            if (userInfo && userInfo.remindTime) {
                this.remindTime = userInfo.remindTime;
            }
        },
        async handleTimeChange(e) {
            const newTime = e.detail.value;
            this.remindTime = newTime;
            
            try {
                uni.showLoading({ title: '保存中...' });
                await callContainer('/api/user/update', {
                    remindTime: newTime
                });
                
                // 更新本地存储
                const userInfo = uni.getStorageSync('userInfo') || {};
                userInfo.remindTime = newTime;
                uni.setStorageSync('userInfo', userInfo);
                
                uni.showToast({ title: '设置已更新', icon: 'success' });
                wx.vibrateShort({ type: "medium" });
            } catch (error) {
                console.error('更新提醒时间失败:', error);
                uni.showToast({ title: '更新失败', icon: 'none' });
            } finally {
                uni.hideLoading();
            }
        },
        handleNav(type) {
            // 这里处理点击事件
            console.log('点击了:', type);
            wx.vibrateShort({ type: "medium" })
            if (type === 'tag') {
                uni.navigateTo({ url: '/pages/tagEdit/tagEdit' })
            } else if (type === 'care') {
                uni.navigateTo({ url: '/pages/careEdit/careEdit' })
            } else if (type === 'family') {
                uni.navigateTo({ url: '/pages/familyDetail/familyDetail' })
            }
        }
    },
    created() {
        const app = getApp()
        this.topBarHeight = app.globalData.topBarHeight;
        this.loadUserData();
    }
}
</script>

<style scoped lang="scss">
/* 容器整体样式 */
.section-container {
    padding: 0 16px;
    /* 页面左右边距 */
    margin-bottom: 20px;
}

/* 标题样式 */
.section-title {
    font-size: 14px;
    color: var(--text-color);
    /* 深色文字 */
    font-weight: 500;
    margin-bottom: 8px;
    margin-left: 4px;
    /* 稍微对齐卡片内部文字 */
    display: block;
}

/* 卡片主体样式 */
.card-box {
    background-color: var(--bg-btn-color);
    /* 截图中的浅绿色背景 */
    /* 如果觉得颜色太深，可以试试 rgba(214, 232, 208, 0.6) */

    border-radius: 16px;
    /* 圆角 */
    padding: 4px 16px;
    /* 内部上下留白 */
    overflow: hidden;
}

/* 单个菜单项 */
.menu-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 0;
    /* 行高 */
    background-color: transparent;
}

/* 点击时的按压效果 */
.item-hover {
    opacity: 0.7;
}

/* 左侧图标和文字的容器 */
.left-content {
    display: flex;
    align-items: center;
}

/* 图标微调 */
.menu-icon {
    margin-right: 10px;
    /* 确保图标和文字垂直对齐 */
    display: flex;
    align-items: center;
}

/* 文字样式 */
.menu-text {
    font-size: 15px;
    color: var(--text-color);
    /* 深橄榄绿文字 */
    font-weight: 400;
}

.time-text {
    font-size: 14px;
    color: var(--primary-color);
    margin-right: 4px;
}

/* 分割线 */
.divider {
    height: 1px;
    background-color: var(--border-color);
    /* 极淡的分割线 */
    width: 100%;
    /* 如果想让分割线不顶头，可以加 margin-left: 32px; */
    margin-left: 32px;
}
/* 客服按钮重置样式 */
.contact-btn {
    width: 100%;
    background: transparent;
    padding: 0;
    margin: 0;
    line-height: inherit;
    text-align: left;
    border: none;
    border-radius: 0;
    display: block;

    &::after {
        border: none;
    }
}
</style>