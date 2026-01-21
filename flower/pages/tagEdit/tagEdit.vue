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

                <view v-for="(item, index) in tagList" :key="item.id">

                    <view class="list-item">
                        <!-- A. 编辑模式 -->
                        <template v-if="item.isEditing">
                            <input class="edit-input" v-model="item.tempName" :focus="true" />
                            <view class="action-group">
                                <view class="mini-btn save-btn" @click="saveEdit(index)">保存</view>
                                <view class="mini-btn cancel-btn" @click="cancelEdit(index)">取消</view>
                            </view>
                        </template>

                        <!-- B. 普通展示模式 -->
                        <template v-else>
                            <view class="left-info">
                                <text class="tag-name">{{ item.name }}</text>
                                <!-- 编辑图标 -->
                                <view class="icon-wrapper" @click="startEdit(index)">
                                    <uni-icons type="compose" size="18" color="#2F3E25"></uni-icons>
                                </view>
                            </view>
                            <!-- 数量徽标 -->
                            <view class="count-badge">
                                <text>{{ item.count }}</text>
                            </view>
                        </template>
                    </view>

                    <!-- 分割线 (最后一项不显示) -->
                    <view v-if="index < tagList.length - 1" class="divider"></view>
                </view>

            </view>
        </view>
    </view>
</template>

<script>
export default {
    data() {
        return {
            newTagName: '',
            // 模拟数据
            tagList: [
                { id: 1, name: '窗台', count: 0, isEditing: false, tempName: '' },
                { id: 2, name: '花架', count: 0, isEditing: false, tempName: '' },
                { id: 3, name: '多肉', count: 12, isEditing: false, tempName: '' },
            ],
            topBarHeight:0
        }
    },
    methods: {
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
            // 先把其他正在编辑的关掉（可选逻辑）
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
}

/* 通用标题样式 */
.section-title {
    font-size: 14px;
    color: #4A6139;
    /* 深绿色标题 */
    font-weight: bold;
    margin-bottom: 10px;
    margin-left: 4px;
    display: block;
}

/* 通用卡片样式 */
.card-box {
    background-color: rgba(255,255,255,0.55);
    /* 豆沙绿背景 */
    border-radius: 16px;
    overflow: hidden;
    margin-bottom: 24px;
}

/* --- 创建新标签区域 --- */
.input-card {
    display: flex;
    align-items: center;
    padding: 6px 6px 6px 16px;
    /* 右侧padding小一点为了放按钮 */
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
    /* 输入框提示文字颜色 */
}

.add-btn {
    background-color: rgba(0, 0, 0, 0.1);
    /* 默认半透明黑 */
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
    padding: 0 16px;
}

.list-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 54px;
    /* 列表项高度 */
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
    /* 增加点击热区 */
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
    /* 加上底部光标线，模拟输入状态 */
    // border-bottom: 1px solid #4A6139; 
    margin-right: 10px;
}

.action-group {
    display: flex;
    align-items: center;
    gap: 8px;
    /* 按钮间距 */
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
}
</style>