<template>
    <view :style="{height:topBarHeight+'px'}"></view>
    <view class="page-container">

        <!-- 1. 创建新标签区域 -->
        <view class="section">
            <text class="section-title">创建新标签</text>
            <view class="card-box input-card">
                <input class="input-field" type="text" v-model="newTagName" placeholder="输入标签名称"
                    placeholder-class="placeholder-style" />
                <view class="add-btn" :class="{ 'disabled': !newTagName }" @click="handleAddTag">
                    <text class="btn-text">添加</text>
                </view>
            </view>
        </view>

        <!-- 2. 可用标签列表区域 -->
        <view class="section">
            <text class="section-title">可用标签</text>
            <view class="card-box list-card">
                <!-- 🌟 1. 外层包裹 uni-swipe-action -->
                <uni-swipe-action>
                    <block v-for="(item, index) in tagList" :key="item.id">
                        
                        <!-- 🌟 2. 每一个项包裹 uni-swipe-action-item -->
                        <!-- right-options 定义左滑出现的按钮 -->
                        <uni-swipe-action-item 
                            :right-options="swipeOptions" 
                            @click="swipeClick($event, index)"
                            :auto-close="true"
                        >
                            <view class="list-item">
                                <!-- A. 编辑模式 -->
                                <template v-if="item.isEditing">
                                    <input class="edit-input" v-model="item.tempName" :focus="true" />
                                    <view class="action-group">
                                        <!-- 注意：这里加了 @click.stop 防止触发其他点击事件 -->
                                        <view class="mini-btn save-btn" @click.stop="saveEdit(index)">保存</view>
                                        <view class="mini-btn cancel-btn" @click.stop="cancelEdit(index)">取消</view>
                                    </view>
                                </template>

                                <!-- B. 普通展示模式 -->
                                <template v-else>
                                    <view class="left-info">
                                        <text class="tag-name">{{ item.name }}</text>
                                        <!-- 编辑图标 -->
                                        <view class="icon-wrapper" @click.stop="startEdit(index)">
                                            <uni-icons type="compose" size="18" color="#2F3E25"></uni-icons>
                                        </view>
                                    </view>
                                    <!-- 数量徽标 -->
                                    <view class="count-badge">
                                        <text>{{ item.count }}</text>
                                    </view>
                                </template>
                            </view>
                        </uni-swipe-action-item>

                        <!-- 分割线 (最后一项不显示) -->
                        <view v-if="index < tagList.length - 1" class="divider"></view>
                    </block>
                </uni-swipe-action>

            </view>
        </view>
    </view>
</template>

<script>
export default {
    data() {
        return {
            newTagName: '',
            // 🌟 3. 定义左滑按钮样式
            swipeOptions: [
                {
                    text: '删除',
                    style: {
                        backgroundColor: '#dd524d', // 红色背景
                        color: '#fff',
                        fontSize: '14px'
                    }
                }
            ],
            // 模拟数据
            tagList: [
                { id: 1, name: '窗台', count: 0, isEditing: false, tempName: '' },
                { id: 2, name: '花架', count: 0, isEditing: false, tempName: '' },
                { id: 3, name: '多肉', count: 12, isEditing: false, tempName: '' },
            ],
            topBarHeight: 0
        }
    },
    methods: {
        // 🌟 4. 处理左滑按钮点击
        swipeClick(e, index) {
            // e.content.text 是按钮的文字，比如 '删除'
            if (e.content.text === '删除') {
                // 可以在这里加一个弹窗确认
                uni.showModal({
                    title: '提示',
                    content: '确定要删除这个标签吗？',
                    success: (res) => {
                        if (res.confirm) {
                            this.tagList.splice(index, 1);
                            uni.showToast({ title: '已删除', icon: 'none' });
                        }
                    }
                });
            }
        },

        // 添加标签
        handleAddTag() {
            if (!this.newTagName.trim()) return;

            this.tagList.push({
                id: Date.now(),
                name: this.newTagName,
                count: 0,
                isEditing: false,
                tempName: ''
            });

            this.newTagName = ''; // 清空输入框
            uni.showToast({ title: '添加成功', icon: 'none' });
        },

        // 开始编辑
        startEdit(index) {
            // 先把其他正在编辑的关掉
            this.tagList.forEach(item => item.isEditing = false);

            const item = this.tagList[index];
            item.tempName = item.name; // 备份当前名字
            item.isEditing = true;
        },

        // 保存编辑
        saveEdit(index) {
            const item = this.tagList[index];
            if (!item.tempName.trim()) {
                uni.showToast({ title: '名称不能为空', icon: 'none' });
                return;
            }
            item.name = item.tempName;
            item.isEditing = false;
        },

        // 取消编辑
        cancelEdit(index) {
            this.tagList[index].isEditing = false;
        }
    },
    onLoad() {
        const app = getApp()
        this.topBarHeight = app.globalData.topBarHeight;
    }
}
</script>

<style scoped lang="scss">
/* 页面背景 */
.page-container {
    min-height: 100vh;
    padding: 20px 16px;
    box-sizing: border-box;
    /* 这里的背景色很重要，防止滑动时露出奇怪的底色 */
    // background-color: #F7F9F5; 
}

/* 通用标题样式 */
.section-title {
    font-size: 14px;
    color: #4A6139;
    font-weight: bold;
    margin-bottom: 10px;
    margin-left: 4px;
    display: block;
}

/* 通用卡片样式 */
.card-box {
    background-color: rgba(255, 255, 255, 0.55);
    border-radius: 16px;
    overflow: hidden; /* 关键：保证圆角不被子元素撑破 */
    margin-bottom: 24px;
}

/* 🌟 修复 uni-swipe-action 可能自带的背景色问题 */
::v-deep .uni-swipe_button-group {
    /* 确保删除按钮高度撑满 */
    height: 100%; 
}

::v-deep .uni-swipe {
    /* 让滑动组件背景透明，透出 card-box 的豆沙绿 */
    background-color: transparent !important; 
}

/* --- 创建新标签区域 --- */
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
    transition: all 0.2s;

    &.disabled {
        opacity: 0.5;
    }

    &:active {
        background-color: rgba(0, 0, 0, 0.2);
    }
}

.btn-text {
    font-size: 13px;
    color: #4A6139;
    font-weight: 500;
}

/* --- 列表区域 --- */
.list-card {
    /* 列表卡片不需要 padding，因为 padding 会导致滑动条也被挤进去 */
    /* 我们把 padding 移到 list-item 内部 */
    padding: 0; 
}

.list-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 54px;
    /* 🌟 关键：因为 list-card 去掉了 padding，这里要补回来 */
    padding: 0 16px; 
    width: 100%;
    box-sizing: border-box;
    background-color: transparent; /* 确保透明 */
}

/* 左侧：名字+图标 */
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

/* 右侧：数量徽标 */
.count-badge {
    background-color: rgba(0, 0, 0, 0.08);
    padding: 2px 8px;
    border-radius: 6px;
    min-width: 16px;
    text-align: center;

    text {
        font-size: 12px;
        color: #4A6139;
    }
}

/* --- 编辑模式样式 --- */
.edit-input {
    flex: 1;
    font-size: 16px;
    color: #333;
    height: 100%;
    margin-right: 10px;
}

.action-group {
    display: flex;
    align-items: center;
    gap: 8px;
}

.mini-btn {
    font-size: 13px;
    padding: 4px 12px;
    border-radius: 12px;
}

.save-btn {
    background-color: #6B8857;
    color: #fff;
}

.cancel-btn {
    background-color: rgba(0, 0, 0, 0.1);
    color: #555;
}

/* 分割线 */
.divider {
    height: 1px;
    background-color: rgba(0, 0, 0, 0.05);
    width: 100%;
    /* 稍微缩进一点，看起来更像 iOS 风格 */
    margin-left: 16px; 
}
</style>
