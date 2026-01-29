<template>
    <view class="page-wrapper">
        <navBar title="养护管理" />
        <view :style="{ height: topBarHeight + 'px' }"></view>
        
        <scroll-view scroll-y class="main-scroll" :show-scrollbar="false" :enhanced="true">
            <view class="page-container">
                <!-- 1. 创建/编辑养护项区域 -->
                <view class="section" v-if="!isSorting">
                    <text class="section-title">{{ isEditing ? '编辑养护项' : '创建新养护项' }}</text>
                    <view class="card-box form-card">
                        <view class="form-item">
                            <text class="label">名称</text>
                            <input class="input-field" type="text" v-model="formData.name" placeholder="如: 浇水, 施肥" />
                        </view>
                        <view class="form-item">
                            <text class="label">标识 (Type)</text>
                            <input class="input-field" type="text" v-model="formData.type" placeholder="如: Watering" />
                        </view>
                        <view class="form-item">
                            <text class="label">选择图标</text>
                            <view class="options-grid">
                                <view v-for="icon in iconOptions" :key="icon" 
                                    class="option-item icon-item" 
                                    :class="{ 'active': formData.icon === icon }"
                                    @click="formData.icon = icon">
                                    <uni-icons :type="icon" size="20" :color="formData.icon === icon ? '#fff' : '#666'"></uni-icons>
                                </view>
                            </view>
                        </view>
                        <view class="form-item">
                            <text class="label">选择颜色</text>
                            <view class="options-grid">
                                <view v-for="color in colorOptions" :key="color" 
                                    class="option-item color-item" 
                                    :class="{ 'active': formData.color === color }"
                                    :style="{ backgroundColor: color }"
                                    @click="formData.color = color">
                                    <uni-icons v-if="formData.color === color" type="checkmarkempty" size="16" color="#fff"></uni-icons>
                                </view>
                            </view>
                        </view>
                        <view class="btn-group">
                            <view class="action-btn cancel" v-if="isEditing" @click="resetForm">取消</view>
                            <view class="action-btn save" :class="{ 'disabled': !formData.name || !formData.type }" @click="handleSubmit">
                                {{ isEditing ? '保存修改' : '立即添加' }}
                            </view>
                        </view>
                    </view>
                </view>

                <!-- 2. 养护项列表区域 -->
                <view class="section">
                    <view class="section-header">
                        <text class="section-title">现有养护项</text>
                        <view class="sort-btn" @click="toggleSortMode">
                            <uni-icons type="list" size="16" :color="isSorting ? '#4A6139' : '#666'"></uni-icons>
                            <text :class="{ 'active-text': isSorting }">{{ isSorting ? '完成' : '排序' }}</text>
                        </view>
                    </view>

                    <view class="card-box list-card" :class="{ 'sorting-mode': isSorting }">
                        <!-- 模式 A: 普通模式 (左滑删除) -->
                        <uni-swipe-action v-if="!isSorting">
                            <block v-for="(item, index) in careList" :key="item.ID">
                                <uni-swipe-action-item :right-options="swipeOptions" @click="swipeClick($event, index)">
                                    <view class="list-item" @click="startEdit(item)">
                                        <view class="left-info">
                                            <view class="care-icon-preview" :style="{ backgroundColor: item.color }">
                                                <uni-icons :type="item.icon" size="20" color="#fff"></uni-icons>
                                            </view>
                                            <view class="care-text-info">
                                                <text class="care-name">{{ item.name }}</text>
                                                <text class="care-type">{{ item.type }}</text>
                                            </view>
                                        </view>
                                        <uni-icons type="compose" size="18" color="#999"></uni-icons>
                                    </view>
                                </uni-swipe-action-item>
                                <view v-if="index < careList.length - 1" class="divider"></view>
                            </block>
                        </uni-swipe-action>

                        <!-- 模式 B: 排序模式 -->
                        <movable-area v-else :style="{ height: areaHeight + 'px' }" class="sort-area">
                            <block v-for="(item, index) in careList" :key="item.ID">
                                <movable-view class="sort-item" :y="item.y" direction="vertical" :damping="40"
                                    @change="onDragChange($event, index)" @touchstart="onDragStart(index)" @touchend="onDragEnd"
                                    :style="{ zIndex: curDragIndex === index ? 99 : 1 }">
                                    <view class="list-item sort-inner">
                                        <view class="left-info">
                                            <uni-icons type="bars" size="20" color="#8FA385" class="drag-handle"></uni-icons>
                                            <view class="care-icon-preview" :style="{ backgroundColor: item.color }">
                                                <uni-icons :type="item.icon" size="20" color="#fff"></uni-icons>
                                            </view>
                                            <text class="care-name">{{ item.name }}</text>
                                        </view>
                                    </view>
                                    <view class="divider"></view>
                                </movable-view>
                            </block>
                        </movable-area>
                    </view>
                </view>
                
                <!-- 底部安全区域占位 -->
                <view class="safe-area-bottom"></view>
            </view>
        </scroll-view>
    </view>
</template>

<script>
import navBar from '@/components/navBar.vue'
import { callContainer } from '../../utils/request';

const ROW_HEIGHT = 60;

export default {
    components: { navBar },
    data() {
        return {
            familyId: 0,
            formData: {
                id: null,
                name: '',
                type: '',
                icon: 'checkbox-filled',
                color: '#D6EAF8'
            },
            isEditing: false,
            swipeOptions: [{ text: '删除', style: { backgroundColor: '#dd524d' } }],
            careList: [],
            iconOptions: [
                'checkbox-filled', 'flask', 'scissors', 'download', 
                'sun', 'flag', 'heart', 'calendar', 'fire', 'medal',
                'camera', 'image', 'chat'
            ],
            colorOptions: [
                '#D6EAF8', '#DCECC9', '#F2D7D5', '#E8E0D5', 
                '#FEF5E7', '#EBDEF0', '#E5E8E8', '#A2D9CE',
                '#FAD7A0', '#D2B4DE', '#F1948A', '#85C1E9'
            ],
            topBarHeight: 0,
            isSorting: false,
            areaHeight: 0,
            curDragIndex: -1,
            tempY: 0,
        }
    },
    methods: {
        async loadData() {
            try {
                // 优先从 storage 获取，如果没有则尝试获取家庭列表并使用第一个
                let fId = uni.getStorageSync('familyId');
                console.log("Current familyId from storage:", fId);
                
                if (!fId) {
                    console.log("Storage 中无 familyId，尝试从缓存的家庭列表获取");
                    const families = uni.getStorageSync('family');
                    if (families && families.length > 0) {
                        fId = families[0].ID || families[0].id;
                        uni.setStorageSync('familyId', fId);
                    }
                }
                
                // 如果还是没有，可能需要重新登录获取
                if (!fId) {
                    console.log("仍然没有 familyId，调用登录接口获取");
                    uni.showLoading({ title: '同步家庭信息...' });
                    const user = await callContainer("/api/login");
                    if (user.data && user.data.family && user.data.family.length > 0) {
                        fId = user.data.family[0].ID;
                        uni.setStorageSync('familyId', fId);
                        uni.setStorageSync('family', user.data.family);
                    }
                    uni.hideLoading();
                }
                
                if (fId) {
                    this.familyId = Number(fId);
                    await this.getCareList();
                } else {
                    uni.showToast({ title: '未找到家庭信息', icon: 'none' });
                }
            } catch (e) { 
                console.error("loadData error:", e); 
                uni.hideLoading();
            }
        },
        async getCareList() {
            if (!this.familyId) return;
            try {
                const res = await callContainer("/api/care/", { familyId: this.familyId });
                if (res.data) {
                    this.careList = res.data.map((item, index) => ({
                        ...item,
                        y: index * ROW_HEIGHT
                    }));
                    this.areaHeight = this.careList.length * ROW_HEIGHT;
                }
            } catch (e) { 
                console.error("getCareList error:", e);
                uni.showToast({ title: '获取列表失败', icon: 'none' });
            }
        },
        resetForm() {
            this.formData = { 
                id: null, 
                name: '', 
                type: '', 
                icon: this.iconOptions[0], 
                color: this.colorOptions[0] 
            };
            this.isEditing = false;
        },
        startEdit(item) {
            this.formData = { ...item, id: item.ID };
            this.isEditing = true;
            uni.pageScrollTo({ scrollTop: 0, duration: 300 });
        },
        async handleSubmit() {
            if (!this.formData.name || !this.formData.type) return;
            try {
                if (this.isEditing) {
                    await callContainer("/api/care/update", this.formData);
                } else {
                    await callContainer("/api/care/add", { ...this.formData, familyId: this.familyId });
                }
                uni.showToast({ title: '操作成功' });
                this.resetForm();
                this.getCareList();
            } catch (e) { console.error(e); }
        },
        async swipeClick(e, index) {
            if (e.content.text === '删除') {
                const item = this.careList[index];
                uni.showModal({
                    title: '提示',
                    content: '确定删除该养护项吗？',
                    success: async (res) => {
                        if (res.confirm) {
                            await callContainer("/api/care/delete", { id: item.ID });
                            this.getCareList();
                        }
                    }
                });
            }
        },
        // --- 排序逻辑 ---
        async toggleSortMode() {
            if (this.isSorting) {
                this.isSorting = false;
                const sortedIds = this.careList.map(item => item.ID);
                await callContainer("/api/care/sort", { careIds: sortedIds });
            } else {
                this.careList.forEach((item, index) => item.y = index * ROW_HEIGHT);
                this.areaHeight = this.careList.length * ROW_HEIGHT;
                this.isSorting = true;
            }
        },
        onDragStart(index) { this.curDragIndex = index; },
        onDragChange(e, index) {
            if (e.detail.source === 'touch' && index === this.curDragIndex) {
                this.tempY = e.detail.y;
            }
        },
        onDragEnd() {
            if (this.curDragIndex === -1) return;
            let target = Math.round(this.tempY / ROW_HEIGHT);
            if (target < 0) target = 0;
            if (target > this.careList.length - 1) target = this.careList.length - 1;

            if (target !== this.curDragIndex) {
                const item = this.careList[this.curDragIndex];
                this.careList.splice(this.curDragIndex, 1);
                this.careList.splice(target, 0, item);
            }
            this.$nextTick(() => {
                this.careList.forEach((item, index) => item.y = index * ROW_HEIGHT);
                this.curDragIndex = -1;
            });
        }
    },
    onLoad() {
        const app = getApp();
        this.topBarHeight = app.globalData.topBarHeight;
        this.loadData();
    }
}
</script>

<style scoped lang="scss">
.page-wrapper {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background-color: #C1D0B7;
}

.main-scroll {
    flex: 1;
    overflow: hidden;
}

.page-container {
    padding: 20px 16px;
}

.safe-area-bottom {
    height: calc(30px + env(safe-area-inset-bottom));
}

.section {
    margin-bottom: 24px;
}

.section-title {
    font-size: 14px;
    color: #4A6139;
    font-weight: bold;
    margin-bottom: 12px;
    margin-left: 4px;
    display: block;
}

.card-box {
    background-color: rgba(255, 255, 255, 0.55);
    border-radius: 16px;
    overflow: hidden;
}

.form-card {
    padding: 16px;
}

.form-item {
    margin-bottom: 12px;
    .label {
        font-size: 12px;
        color: #666;
        margin-bottom: 4px;
        display: block;
    }
    .input-field {
        height: 40px;
        background-color: rgba(255,255,255,0.7);
        border-radius: 8px;
        padding: 0 12px;
        font-size: 14px;
    }
}

.options-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 4px;
}

.option-item {
    width: 36px;
    height: 36px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #fff;
    transition: all 0.2s;
    border: 1px solid transparent;
    
    &.active {
        border-color: #4A6139;
        transform: scale(1.1);
    }
}

.icon-item {
    &.active {
        background-color: #4A6139;
    }
}

.color-item {
    &.active {
        box-shadow: 0 0 0 2px #fff, 0 0 0 4px #4A6139;
    }
}

.btn-group {
    display: flex;
    gap: 10px;
    margin-top: 20px;
}

.action-btn {
    flex: 1;
    height: 44px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 22px;
    font-size: 14px;
    font-weight: bold;
    &.save { background-color: #4A6139; color: #fff; }
    &.cancel { background-color: #eee; color: #666; }
    &.disabled { opacity: 0.5; }
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
}

.sort-btn {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 4px 10px;
    background-color: rgba(255, 255, 255, 0.6);
    border-radius: 12px;
    text { font-size: 12px; color: #666; }
    .active-text { color: #4A6139; font-weight: bold; }
}

.list-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 60px;
    padding: 0 16px;
    background-color: transparent;
}

.left-info {
    display: flex;
    align-items: center;
    gap: 12px;
}

.care-icon-preview {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
}

.care-text-info {
    display: flex;
    flex-direction: column;
    .care-name { font-size: 15px; color: #333; font-weight: 500; }
    .care-type { font-size: 11px; color: #999; }
}

.divider {
    height: 1px;
    background-color: rgba(0, 0, 0, 0.05);
    margin-left: 64px;
}

.sort-area { width: 100%; }
.sort-item { width: 100%; height: 60px; z-index: 1; background-color: rgba(255, 255, 255, 0.7); }
.sort-item[style*="z-index: 99"] { box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1); background-color: #fff; }
.drag-handle { margin-right: 8px; }
</style>
