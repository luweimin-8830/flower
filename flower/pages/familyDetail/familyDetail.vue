<template>
    <navBar title="家庭管理" />
    <view :style="{ height: topBarHeight + 'px' }"></view>

    <view class="page-container">

        <!-- 头部 -->
        <view class="header-row">
            <text class="page-title">我的家庭</text>
            <view class="sort-btn" @click="toggleSortMode">
                <uni-icons type="list" size="16" color="#666"></uni-icons>
                <text class="sort-text">{{ isSorting ? '完成' : '排序' }}</text>
            </view>
        </view>

        <!-- 列表容器 -->
        <view class="list-card-container">

            <!-- A: 普通模式 -->
            <block v-if="!isSorting">
                <view v-for="(item, index) in familyList" :key="item.ID" class="list-item"
                    :class="{ 'has-border': index !== familyList.length - 1 }">

                    <!-- 左侧：名称 + 编辑 (仅户主) -->
                    <view class="left-content">
                        <text class="item-name">{{ item.name }}</text>
                        <view class="edit-icon-wrap" v-if="item.role === 'owner'" @click.stop="editName(item)">
                            <uni-icons type="compose" size="16" color="#999"></uni-icons>
                        </view>
                    </view>
                    <!-- 右侧：操作图标 + 权限 + 数量 -->
                    <view class="right-content">
                        <!-- 1. 邀请按钮 (户主 & 管理员) -->
                        <view class="icon-btn" v-if="item.role === 'owner' || item.role === 'admin'"
                            @click.stop="handleInvite(item)">
                            <uni-icons type="personadd" size="18" color="#4A6139"></uni-icons>
                        </view>

                        <!-- 2. 🌟 删除按钮 (仅户主) -->
                        <view class="icon-btn delete-btn" v-if="item.role === 'owner'"
                            @click.stop="handleDelete(item, index)">
                            <uni-icons type="trash" size="18" color="#dd524d"></uni-icons>
                        </view>

                        <!-- 3. 权限标签 -->
                        <view class="role-tag" :class="item.role">
                            {{ getRoleName(item.role) }}
                        </view>

                        <!-- 4. 数量徽标 -->
                        <view class="count-badge">
                            {{ item.memberCount }}
                        </view>
                    </view>
                </view>
            </block>

            <!-- B: 排序模式 -->
            <movable-area v-else :style="{ height: areaHeight + 'px' }" class="sort-area">
                <block v-for="(item, index) in familyList" :key="item.ID">
                    <movable-view class="sort-movable-item" :y="item.y" direction="vertical" :damping="40"
                        @change="onDragChange($event, index)" @touchstart="onDragStart(index)" @touchend="onDragEnd"
                        :style="{ zIndex: curDragIndex === index ? 99 : 1 }">

                        <view class="list-item sort-inner-item">
                            <view class="left-content">
                                <uni-icons type="bars" size="20" color="#999" style="margin-right: 8px;"></uni-icons>
                                <text class="item-name">{{ item.name }}</text>
                            </view>
                            <view class="right-content">
                                <view class="role-tag" :class="item.role">
                                    {{ getRoleName(item.role) }}
                                </view>
                                <view class="count-badge">
                                    {{ item.memberCount }}
                                </view>
                            </view>
                        </view>
                        <view v-if="index !== familyList.length - 1" class="divider-line"></view>
                    </movable-view>
                </block>
            </movable-area>

        </view>

        <!-- 底部添加按钮 -->
        <view class="fab-add-btn" @click="createNewFamily" v-if="!isSorting">
            <uni-icons type="plusempty" size="24" color="#fff"></uni-icons>
        </view>

    </view>
</template>

<script>
import navBar from '@/components/navBar.vue'
import { callContainer } from '../../utils/request';

const ITEM_HEIGHT = 56;

export default {
    components: { navBar },
    data() {
        return {
            topBarHeight: 0,
            familyList: [],
            isSorting: false,
            areaHeight: 0,
            curDragIndex: -1,
            tempY: 0,
        }
    },
    methods: {
        async getFamilyList() {
            try {
                const family = await callContainer("/api/family/", {})
                console.log("call container get family", family)
                if (family.data) {
                    let rawList = family.data;
                    rawList.sort((a, b) => a.mySortOrder - b.mySortOrder);
                    this.familyList = rawList.map((item, index) => ({
                        ID: item.ID,
                        name: item.name,
                        role: item.myRole,
                        memberCount: item.memberCount,
                        y: index * ITEM_HEIGHT
                    }));

                    // 4. 设置拖拽区域总高度
                    this.areaHeight = this.familyList.length * ITEM_HEIGHT;
                }

            } catch (e) { console.error(e); }
        },

        getRoleName(role) {
            const map = { 'owner': '户主', 'admin': '管理员', 'member': '成员' };
            return map[role] || '成员';
        },

        // 🌟 删除家庭 (仅户主)
        async handleDelete(item, index) {
            uni.showModal({
                title: '危险操作',
                content: `确定要解散家庭 "${item.name}" 吗？\n解散后所有数据将无法恢复！`,
                confirmText: '解散',
                confirmColor: '#dd524d', // 红色确认按钮
                success: async (res) => {
                    if (res.confirm) {
                        try {
                            // 调用后端删除API
                            await callContainer("/api/family/delete", {
                                familyId: item.ID
                            });
                            
                            // 前端移除该家庭
                            this.familyList.splice(index, 1);
                            // 重新计算高度和位置
                            this.areaHeight = this.familyList.length * ITEM_HEIGHT;
                            this.familyList.forEach((it, idx) => it.y = idx * ITEM_HEIGHT);
                        } catch (e) {
                            console.error("删除家庭失败:", e);
                        }
                    }
                }
            });
        },
        // 邀请成员
        handleInvite(item) {
            uni.vibrateShort();
            const inviteCode = Math.random().toString(36).substring(2, 8).toUpperCase();
            uni.showModal({
                title: `邀请加入 "${item.name}"`,
                content: `邀请码：${inviteCode}`,
                confirmText: '复制',
                success: (res) => {
                    if (res.confirm) {
                        uni.setClipboardData({ data: inviteCode });
                    }
                }
            });
        },
        async editName(item) {
            const res = await new Promise((resolve, reject) => {
                uni.showModal({
                    title: '修改名称',
                    content: item.name,
                    editable: true,
                    placeholderText: '最多10个字符',
                    success: resolve,
                    fail: reject
                });
            })
            if (res.confirm && res.content) {
                const newName = res.content.trim();
                
                // 字符长度校验
                if (newName.length > 10) {
                    uni.showToast({
                        title: '家庭名称不能超过10个字符',
                        icon: 'none',
                        duration: 2000
                    });
                    return;
                }
                
                // 如果名称未改变，直接返回
                if (newName === item.name) {
                    return;
                }
                
                item.name = newName;
                const uploadName = await callContainer("/api/family/update", {
                    familyId: item.ID,
                    name: newName
                })
                console.log("call container upload family:", uploadName)
                // uni.$emit('refreshFamilyList', { needRefresh: true });
                const result = await new Promise((resolve) => {
                    uni.getStorage({ key: "family", success: resolve, fail: resolve });
                })
                let localList = result.data;
                if (Array.isArray(localList)) {
                    const targetIndex = localList.findIndex(f => (f.id === item.ID) || (f.ID === item.ID));
                    if (targetIndex !== -1) {
                        // 修改名称
                        localList[targetIndex].name = item.name;
                        // 写回缓存
                        await new Promise((resolve) => {
                            uni.setStorageSync({ key: "family", data: localList, success: resolve, fail: resolve });
                        })
                    }
                }
                
                uni.showToast({
                    title: '修改成功',
                    icon: 'success'
                });
            }
        },

        createNewFamily() {
            uni.showModal({
                title: '创建新家庭',
                editable: true,
                placeholderText: '请输入家庭名称（最多10个字符）',
                success: async (res) => {
                    if (res.confirm && res.content) {
                        const name = res.content.trim();
                        
                        // 字符长度校验
                        if (name.length > 10) {
                            uni.showToast({
                                title: '家庭名称不能超过10个字符',
                                icon: 'none',
                                duration: 2000
                            });
                            return;
                        }
                        
                        try {
                            // 调用后端创建家庭API
                            const result = await callContainer("/api/family/add", {
                                name: name
                            });
                            console.log("创建家庭成功:", result);
                            
                            // 重新获取家庭列表
                            await this.getFamilyList();
                        } catch (e) {
                            console.error("创建家庭失败:", e);
                            uni.showToast({ 
                                title: '创建失败', 
                                icon: 'none' 
                            });
                        }
                    }
                }
            });
        },

        goDetail(item) {
            if (this.isSorting) return;
            console.log('进入家庭详情', item.name);
        },

        // --- 排序逻辑 ---
        toggleSortMode() {
            if (this.isSorting) {
                this.isSorting = false;
            } else {
                this.familyList.forEach((item, index) => item.y = index * ITEM_HEIGHT);
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
            let target = Math.round(this.tempY / ITEM_HEIGHT);
            if (target < 0) target = 0;
            if (target > this.familyList.length - 1) target = this.familyList.length - 1;

            if (target !== this.curDragIndex) {
                const item = this.familyList[this.curDragIndex];
                this.familyList.splice(this.curDragIndex, 1);
                this.familyList.splice(target, 0, item);
            }
            this.$nextTick(() => {
                this.familyList.forEach((item, index) => item.y = index * ITEM_HEIGHT);
                this.curDragIndex = -1;
            });
        }
    },
    onLoad() {
        const app = getApp()
        this.topBarHeight = app.globalData.topBarHeight
        this.getFamilyList()
    }
}
</script>

<style scoped lang="scss">
.page-container {
    min-height: 100vh;
    padding: 20px 16px;
    box-sizing: border-box;
    background-color: #C1D0B7;
}

.header-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
    padding: 0 4px;
}

.page-title {
    font-size: 15px;
    font-weight: bold;
    color: #4A6139;
}

.sort-btn {
    display: flex;
    align-items: center;
    background-color: rgba(255, 255, 255, 0.6);
    padding: 4px 10px;
    border-radius: 14px;
    gap: 4px;
}

.sort-text {
    font-size: 13px;
    color: #555;
}

.list-card-container {
    background-color: rgba(255, 255, 255, 0.55);
    border-radius: 16px;
    overflow: hidden;
    padding: 0;
}

.list-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 56px;
    padding: 0 16px;
    background-color: transparent;
    transition: background-color 0.2s;
}

.list-item:active {
    background-color: rgba(255, 255, 255, 0.3);
}

.has-border {
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.left-content {
    display: flex;
    align-items: center;
    flex: 1;
}

.item-name {
    font-size: 16px;
    color: #333;
    font-weight: 500;
}

.edit-icon-wrap {
    padding: 8px;
    margin-left: 2px;
    opacity: 0.6;
}

.right-content {
    display: flex;
    align-items: center;
}

/* 图标按钮 */
.icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 30px;
    height: 30px;
    margin-right: 4px;
    border-radius: 50%;
}

.icon-btn:active {
    background-color: rgba(0, 0, 0, 0.05);
}

/* 角色标签 */
.role-tag {
    font-size: 11px;
    padding: 2px 6px;
    border-radius: 4px;
    margin-right: 8px;
    font-weight: bold;
    min-width: 36px;
    text-align: center;
}

.role-tag.owner {
    color: #E6A23C;
    background-color: rgba(230, 162, 60, 0.15);
}

.role-tag.admin {
    color: #409EFF;
    background-color: rgba(64, 158, 255, 0.15);
}

.role-tag.member {
    color: #6B8857;
    background-color: rgba(107, 136, 87, 0.15);
}

/* 数量徽标 */
.count-badge {
    min-width: 24px;
    height: 24px;
    background-color: rgba(0, 0, 0, 0.1);
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 13px;
    color: #555;
    font-weight: 500;
}

/* 排序相关 */
.sort-area {
    width: 100%;
}

.sort-movable-item {
    width: 100%;
    height: 56px;
    z-index: 1;
}

.sort-movable-item[style*="z-index: 99"] .sort-inner-item {
    background-color: rgba(255, 255, 255, 0.9);
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
}

.sort-inner-item {
    background-color: rgba(255, 255, 255, 0.55);
    border-bottom: none;
}

.divider-line {
    height: 1px;
    background-color: rgba(0, 0, 0, 0.05);
    margin: 0 16px;
}

.fab-add-btn {
    position: fixed;
    bottom: 40px;
    right: 20px;
    width: 50px;
    height: 50px;
    background-color: #6B8857;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(107, 136, 87, 0.4);
    z-index: 100;
}

.fab-add-btn:active {
    transform: scale(0.95);
}
</style>