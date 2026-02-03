<template>
    <view class="page-wrapper">
        <navBar />
        <view :style="{ height: topBarHeight + 'px' }"></view>

        <scroll-view scroll-y class="main-scroll" :show-scrollbar="false" :enhanced="true">
            <view class="page-container">
                <!-- 1. 创建新标签区域 (排序模式下隐藏，避免干扰) -->
                <view class="section" v-if="!isSorting">
                    <text class="section-title">创建新标签</text>
                    <view class="card-box input-card">
                        <input class="input-field" type="text" v-model="newTagName" placeholder="输入标签名称"
                            placeholder-class="placeholder-style" />
                        <view class="add-btn" :class="{ 'disabled': !newTagName }" @click="handleAddTag">
                            <text class="btn-text">添加</text>
                        </view>
                    </view>
                </view>

                <!-- 2. 标签列表区域 -->
                <view class="section">
                    <view class="section-header">
                        <text class="section-title">可用标签</text>
                        <!-- 🌟 切换排序模式按钮 -->
                        <view class="sort-btn" @click="toggleSortMode">
                            <!-- 图标：根据状态变色 -->
                            <view class="iconfont plant-paixu" :class="{ 'active-text': isSorting }"></view>
                            <text :class="{ 'active-text': isSorting }">{{ isSorting ? '完成' : '排序' }}</text>
                        </view>
                    </view>

                    <view class="card-box list-card" :class="{ 'sorting-mode': isSorting }">

                        <!-- 🌟 模式 A: 普通模式 (支持左滑删除) -->
                        <uni-swipe-action v-if="!isSorting">
                            <block v-for="(item, index) in tagsList" :key="item.id">
                                <uni-swipe-action-item :right-options="swipeOptions" @click="swipeClick($event, index)"
                                    :auto-close="true">
                                    <view class="list-item">
                                        <template v-if="item.isEditing">
                                            <view class="edit-container">
                                                <input class="edit-input" v-model="item.tempName" :focus="true"
                                                    @confirm="saveEdit(index)" />
                                                <view class="edit-actions">
                                                    <text class="edit-btn cancel"
                                                        @click.stop="cancelEdit(index)">取消</text>
                                                    <view class="edit-divider"></view>
                                                    <text class="edit-btn save" @click.stop="saveEdit(index)">保存</text>
                                                </view>
                                            </view>
                                        </template>
                                        <template v-else>
                                            <view class="left-info">
                                                <text class="tag-name">{{ item.name }}</text>
                                                <view class="icon-wrapper" @click.stop="startEdit(index)">
                                                    <uni-icons type="compose" size="18" color="#2F3E25"></uni-icons>
                                                </view>
                                            </view>
                                            <view class="count-badge">
                                                <text>{{ item.plantCount }}</text>
                                            </view>
                                        </template>
                                    </view>
                                </uni-swipe-action-item>
                                <view v-if="index < tagsList.length - 1" class="divider"></view>
                            </block>
                        </uni-swipe-action>

                        <!-- 🌟 模式 B: 排序模式 (支持拖拽) -->
                        <movable-area v-else :style="{ height: areaHeight + 'px' }" class="sort-area">
                            <block v-for="(item, index) in tagsList" :key="item.id">
                                <movable-view class="sort-item" :y="item.y" direction="vertical" :damping="40"
                                    :disabled="false" @change="onDragChange($event, index)"
                                    @touchstart="onDragStart(index)" @touchend="onDragEnd"
                                    :style="{ zIndex: curDragIndex === index ? 99 : 1 }">
                                    <view class="list-item sort-inner">
                                        <view class="left-info">
                                            <uni-icons type="bars" size="20" color="#8FA385"
                                                class="drag-handle"></uni-icons>
                                            <text class="tag-name">{{ item.name }}</text>
                                        </view>
                                        <view class="count-badge">
                                            <text>{{ item.plantCount }}</text>
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
import navBar from '@/components/navBar.vue'// 假设你已经修复了组件名
import { callContainer } from '../../utils/request';

// 每一行的高度 (54px 内容 + 1px 分割线)
// ⚠️ 注意：这个值必须和 CSS 里的 .list-item 高度 + .divider 高度完全一致
const ROW_HEIGHT = 55;

export default {
    components: {
        navBar,
    },
    data() {
        return {
            familyId: 0,
            newTagName: '',
            swipeOptions: [
                { text: '删除', style: { backgroundColor: '#dd524d', color: '#fff', fontSize: '14px' } }
            ],
            tagsList: [],
            topBarHeight: 0,

            // --- 排序相关数据 ---
            isSorting: false, // 是否处于排序模式
            areaHeight: 0,    // 拖拽区域总高度
            curDragIndex: -1, // 当前正在拖拽的索引
            tempY: 0,         // 记录当前拖拽的 Y 值
        }
    },
    methods: {
        // ... (loadData 等原有方法保持不变) ...
        async loadData() {
            // ... 你的原有代码 ...
            try {
                const familyID = await new Promise((resolve, reject) => {
                    uni.getStorage({ key: 'familyId', success: resolve, fail: reject })
                })
                this.familyId = familyID.data
                this.getTagsList()
            } catch (error) {
                console.error(error)
            }
        },

        async getTagsList() {
            try {
                const res = await callContainer("/api/tag/", { familyId: this.familyId });
                if (res.data) {
                    // 🌟 初始化数据时，必须给每个 item 加上 y 坐标
                    this.tagsList = res.data.map((item, index) => {
                        return {
                            ...item,
                            y: index * ROW_HEIGHT, // 初始位置
                            isEditing: false,
                            tempName: ''
                        }
                    });
                    // 计算区域总高度
                    this.areaHeight = this.tagsList.length * ROW_HEIGHT;
                }
            } catch (error) {
                console.error(error)
            }
        },

        // --- 🌟 排序核心逻辑 ---

        // 1. 切换模式
        async toggleSortMode() {
            if (this.isSorting) {
                // 点击“完成”：保存顺序
                this.isSorting = false;
                await this.saveSortOrder();
            } else {
                // 点击“排序”：进入排序模式
                // 重新计算一次位置，防止数据错乱
                this.tagsList.forEach((item, index) => {
                    item.y = index * ROW_HEIGHT;
                });
                this.areaHeight = this.tagsList.length * ROW_HEIGHT;
                this.isSorting = true;
            }
        },

        // 2. 开始拖拽
        onDragStart(index) {
            this.curDragIndex = index;
        },

        // 3. 拖拽过程 (记录当前 Y 值)
        onDragChange(e, index) {
            // 只有当前拖拽的项才记录，防止其他项抖动
            if (e.detail.source === 'touch' && index === this.curDragIndex) {
                this.tempY = e.detail.y;
            }
        },

        // 4. 拖拽结束 (计算新位置并重排)
        onDragEnd() {
            if (this.curDragIndex === -1) return;

            // 根据最终的 Y 值，计算应该在第几行
            // Math.round 四舍五入，实现吸附效果
            let targetIndex = Math.round(this.tempY / ROW_HEIGHT);

            // 边界限制
            if (targetIndex < 0) targetIndex = 0;
            if (targetIndex > this.tagsList.length - 1) targetIndex = this.tagsList.length - 1;

            if (targetIndex !== this.curDragIndex) {
                // 移动数组元素
                const item = this.tagsList[this.curDragIndex];
                this.tagsList.splice(this.curDragIndex, 1); // 移除旧位置
                this.tagsList.splice(targetIndex, 0, item); // 插入新位置
            }

            // 🌟 关键：重置所有项的 Y 坐标，让它们归位
            // 使用 $nextTick 确保视图更新
            this.$nextTick(() => {
                this.tagsList.forEach((item, index) => {
                    item.y = index * ROW_HEIGHT;
                });
                this.curDragIndex = -1;
            });
        },

        // 5. 保存顺序到服务器
        async saveSortOrder() {
            // 提取 ID 数组
            const sortedIds = this.tagsList.map(item => item.id);
            console.log('新的顺序 ID:', sortedIds);

            try {

                // 这里调用你的后端接口保存顺序
                const tagSort = await callContainer("/api/tag/sort", { tagIds: sortedIds });
                console.log("call container tag sort:", tagSort)
                // uni.showToast({ title: '顺序已保存', icon: 'success' });
            } catch (e) {
                console.error(e);
            }
        },

        // ... (handleAddTag, swipeClick 等原有方法保持不变) ...
        // 注意：handleAddTag 里添加新数据后，也要记得计算 y: this.tagsList.length * ROW_HEIGHT
        async handleAddTag() {
            if (!this.newTagName.trim()) return;
            try {
                const newTag = await callContainer("/api/tag/add", {
                    name: this.newTagName,
                    familyId: this.familyId
                })
                if (newTag.code === 0) {
                    // 重新获取列表，会触发 getTagsList 里的位置初始化
                    this.getTagsList()
                }
            } catch (error) {
                console.error(error)
            }
            this.newTagName = '';
        },

        // ... 其他方法保持不变 ...
        // 记得把 swipeClick, startEdit, saveEdit, cancelEdit 补全
        async swipeClick(e, index) {
            if (e.content.text === '删除') {
                const result = await new Promise((resolve) => uni.showModal({
                    title: '提示', content: '确定删除?', success: resolve
                }));
                if (result.confirm) {
                    const item = this.tagsList[index];
                    await callContainer("/api/tag/delete", { id: item.id });
                    this.tagsList.splice(index, 1);
                    // 删除后要更新高度
                    this.areaHeight = this.tagsList.length * ROW_HEIGHT;
                }
            }
        },
        startEdit(index) {
            this.tagsList.forEach(item => item.isEditing = false);
            const item = this.tagsList[index];
            item.tempName = item.name;
            item.isEditing = true;
        },
        cancelEdit(index) {
            this.tagsList[index].isEditing = false;
        },
        async saveEdit(index) {
            const item = this.tagsList[index];
            if (!item.tempName.trim()) return;
            item.name = item.tempName;
            item.isEditing = false;
            await callContainer("/api/tag/update", { id: item.id, name: item.name });
        }
    },
    onLoad() {
        const app = getApp()
        this.topBarHeight = app.globalData.topBarHeight;
        this.loadData()
    },
}
</script>

<style scoped lang="scss">
/* ... 原有样式保持不变 ... */
.page-wrapper {
    display: flex;
    flex-direction: column;
    height: 100vh;
}

.main-scroll {
    flex: 1;
    overflow: hidden;
}

.page-container {
    padding: 20px 16px;
    box-sizing: border-box;
}

.safe-area-bottom {
    height: calc(30px + env(safe-area-inset-bottom));
}

.section-title {
    font-size: 14px;
    color: #4A6139;
    font-weight: bold;
    margin-bottom: 10px;
    margin-left: 4px;
    display: block;

    @media (prefers-color-scheme: dark) {
        color: #FAF2CB;
    }
}

.card-box {
    background-color: rgba(255, 255, 255, 0.55);
    border-radius: 16px;
    overflow: hidden;
    margin-bottom: 24px;
    transition: all 0.3s;
}

/* 🌟 头部布局：包含标题和排序按钮 */
.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
}

.sort-btn {
    display: flex;
    /* 开启 Flex 布局 */
    align-items: center;
    /* 垂直居中 */
    gap: 4px;
    /* 图标和文字之间的间距 */
    padding: 6px 12px;
    /* 稍微调整一下内边距 */
    background-color: rgba(255, 255, 255, 0.6);
    border-radius: 12px;
    transition: all 0.2s;
    /* 加个过渡动画更丝滑 */
}

/* 点击时的反馈效果 */
.sort-btn:active {
    background-color: rgba(255, 255, 255, 0.8);
}

.sort-btn text {
    font-size: 13px;
    color: #666;
    line-height: 1;
    /* 防止文字垂直方向有偏差 */
}

/* 图标基础样式 */
.plant-paixu {
    font-size: 14px;
    /* 大小跟文字匹配 */
    color: #666;
}

/* 激活状态（绿色） */
.active-text {
    color: #4A6139 !important;
    font-weight: bold;
}

/* --- 排序模式样式 --- */
.sort-area {
    width: 100%;
    /* 背景透明 */
    background-color: transparent;
}

.sort-item {
    width: 100%;
    height: 55px;
    /* 必须与 ROW_HEIGHT 一致 */
    background-color: rgba(255, 255, 255, 0.55);
    /* 保持卡片背景色 */
    z-index: 1;
}

/* 拖动时的样式 */
.sort-item[style*="z-index: 99"] {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    transform: scale(1.02);
    background-color: rgba(255, 255, 255, 0.9);
}

.sort-inner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 54px;
    padding: 0 16px;
}

.drag-handle {
    margin-right: 12px;
}

/* 修复原有的样式引用 */
.input-card {
    display: flex;
    align-items: center;
    padding: 6px 6px 6px 16px;
    height: 50px;
    box-sizing: border-box;
}

.input-field {
    flex: 1;
    font-size: 15px;
    color: #333;
    height: 100%;
}

.placeholder-style {
    color: #8FA385;
}

.add-btn {
    background-color: rgba(0, 0, 0, 0.1);
    padding: 6px 16px;
    border-radius: 14px;
    margin-left: 10px;
}

.btn-text {
    font-size: 13px;
    color: #4A6139;
    font-weight: 500;
}

.list-card {
    padding: 0;
}

.list-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 54px;
    padding: 0 16px;
    width: 100%;
    box-sizing: border-box;
    background-color: transparent;
}

/* --- 编辑模式样式 --- */
.edit-container {
    display: flex;
    align-items: center;
    flex: 1;
    background-color: rgba(255, 255, 255, 0.4);
    border-radius: 10px;
    padding: 0 12px;
    height: 38px;
    margin-right: 12px;
    border: 1px solid rgba(74, 97, 57, 0.1);
}

.edit-input {
    flex: 1;
    font-size: 14px;
    color: #2F3E25;
}

.edit-actions {
    display: flex;
    align-items: center;
    gap: 0;
}

.edit-btn {
    font-size: 13px;
    padding: 6px 8px;
    transition: all 0.2s;
}

.edit-btn:active {
    opacity: 0.6;
}

.edit-btn.cancel {
    color: #888;
}

.edit-btn.save {
    color: #4A6139;
    font-weight: bold;
}

.edit-divider {
    width: 1px;
    height: 12px;
    background-color: rgba(0, 0, 0, 0.08);
}

.left-info {
    display: flex;
    align-items: center;
    flex: 1;
}

.tag-name {
    font-size: 16px;
    color: #2F3E25;
    margin-right: 8px;
}

.icon-wrapper {
    padding: 4px;
    display: flex;
}

.count-badge {
    background-color: rgba(0, 0, 0, 0.08);
    padding: 2px 8px;
    border-radius: 6px;
    min-width: 16px;
    text-align: center;
}

.count-badge text {
    font-size: 12px;
    color: #4A6139;
}

.divider {
    height: 1px;
    background-color: rgba(0, 0, 0, 0.05);
    width: 100%;
    margin-left: 16px;
}

/* 修复 uni-swipe */
::v-deep .uni-swipe_button-group {
    height: 100%;
}

::v-deep .uni-swipe {
    background-color: transparent !important;
}
</style>
