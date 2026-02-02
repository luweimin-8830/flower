<template>
    <navBar :isHome="true" title="加入家庭" />
    <view :style="{ height: topBarHeight + 'px' }"></view>

    <view class="container">
        <view class="card">
            <view class="header">
                <uni-icons type="home-filled" size="48" color="#4A6139"></uni-icons>
                <text class="title">您收到了一个家庭邀请</text>
            </view>
            
            <view class="family-info">
                <text class="label">家庭名称</text>
                <text class="family-name">{{ familyName }}</text>
            </view>

            <button class="join-btn" @click="handleJoin" :loading="loading">
                加入家庭
            </button>
            
            <text class="tip">加入后，您可以与家人一起管理植物</text>
        </view>
    </view>
</template>


<script>
import navBar from '@/components/navBar.vue'
import { callContainer } from '../../utils/request';

export default {
    components: { navBar },
    data() {
        return {
            topBarHeight: 0,
            familyId: null,
            familyName: '',
            loading: false
        }
    },
    onLoad(options) {
        const app = getApp()
        this.topBarHeight = app.globalData.topBarHeight
        
        if (options.familyId) {
            this.familyId = options.familyId
            this.familyName = options.familyName || '未知家庭'
        } else {
            uni.showToast({
                title: '邀请链接无效',
                icon: 'none'
            })
            setTimeout(() => uni.reLaunch({ url: '/pages/index/index' }), 1500)
        }
    },
    methods: {
        async handleJoin() {
            if (!this.familyId) return
            
            this.loading = true
            try {
                await callContainer("/api/family/join", {
                    familyId: parseInt(this.familyId)
                });
                
                uni.showToast({
                    title: '加入成功',
                    icon: 'success'
                });
                
                // 1. 自动切换到该家庭
                const targetFamilyId = parseInt(this.familyId);
                await callContainer("/api/family/switch", {
                    familyId: targetFamilyId
                });
                
                // 2. 刷新并同步本地缓存
                const user = await callContainer("/api/login")
                uni.setStorageSync('family', user.data.family)
                uni.setStorageSync('familyId', targetFamilyId)
                
                // 3. 通知首页刷新
                uni.$emit('familyChanged', targetFamilyId);
                
                setTimeout(() => {
                    uni.reLaunch({ url: '/pages/index/index' })
                }, 1500);
            } catch (e) {
                console.error("加入家庭失败:", e);
                uni.showToast({
                    title: e.message || '加入失败',
                    icon: 'none'
                });
            } finally {
                this.loading = false
            }
        }
    }
}
</script>

<style scoped lang="scss">
.container {
    min-height: 100vh;
    background-color: #C1D0B7;
    padding: 40px 20px;
    display: flex;
    justify-content: center;
    align-items: flex-start;
}

.card {
    width: 100%;
    background-color: rgba(255, 255, 255, 0.55);
    border-radius: 20px;
    padding: 30px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
    display: flex;
    flex-direction: column;
    align-items: center;
}

.header {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-bottom: 30px;
}

.title {
    font-size: 18px;
    font-weight: bold;
    color: #333;
    margin-top: 15px;
}

.family-info {
    width: 100%;
    background-color: rgba(74, 97, 57, 0.05);
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 30px;
    display: flex;
    flex-direction: column;
    align-items: center;
}

.label {
    font-size: 14px;
    color: #666;
    margin-bottom: 8px;
}

.family-name {
    font-size: 24px;
    font-weight: bold;
    color: #4A6139;
}

.join-btn {
    width: 100%;
    height: 50px;
    background-color: #4A6139;
    color: white;
    border-radius: 25px;
    font-size: 16px;
    font-weight: bold;
    display: flex;
    justify-content: center;
    align-items: center;
    border: none;
    margin-bottom: 20px;
}

.join-btn:active {
    opacity: 0.8;
    transform: scale(0.98);
}

.tip {
    font-size: 12px;
    color: #999;
}
</style>
