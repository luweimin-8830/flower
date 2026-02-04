<template>
    <!-- 🌟 1. 最外层包裹一个全屏容器 -->
    <view class="home-container">

        <!-- 🌟 2. 优化后的家庭选择按钮 (毛玻璃效果 + 交互增强) - 使用微信原生 picker -->
        <view class="family-select" :class="{ 'selecting': isSelecting }" :style="{
            width: 'auto',
            height: menuButtonInfo.height + 'px',
            borderRadius: menuButtonInfo.height / 2 + 'px',
            top: menuButtonInfo.top + 'px',
            left: paddingLeft + 'px'
        }" @click="toggleFamilySelect" @touchstart="onTouchStart" @touchend="onTouchEnd">
            <view class="family-select-icon" :class="{ 'selecting': isSelecting }">
                <uni-icons type="home" size="18" color="var(--primary-color)"></uni-icons>
            </view>
            <picker class="custom-select" :value="currentFamilyIndex" :range="familyRange" :range-key="'text'"
                @change="handleFamilyChange">
                <view class="family-name-wrapper">
                    <text class="family-select-text">{{ familyRange[currentFamilyIndex]?.text || '选择家庭' }}</text>
                    <uni-icons type="bottom" size="14" color="var(--primary-color)"
                        style="margin-left: 4px;"></uni-icons>
                </view>
            </picker>
        </view>

        <!-- 🌟 3. 顶部固定区域 (包含占位符、搜索框、标签栏) -->
        <view class="fixed-header-group">
            <!-- 顶部占位 (状态栏高度) -->
            <view :style="{ height: topBarHeight + 'px' }"></view>

            <!-- 搜索框与编辑模式头部 -->
            <view class="header-action-container">
                <template v-if="!isEditMode">
                    <view class="search-box-wrapper">
                        <uni-search-bar @confirm="searchPlant" placeholder="输入植物名称或标签" radius="20" :focus="false"
                            v-model="searchValue" bgColor="rgba(255,255,255,0.5)" clearButton="auto"
                            cancelButton="none">
                        </uni-search-bar>
                    </view>
                    <view class="add-btn" @click="goAddPage">
                        <uni-icons type="plusempty" size="22" color="#333"></uni-icons>
                    </view>
                </template>
                <template v-else>
                    <view class="edit-mode-header">
                        <uni-row :gutter="20" style="width: 100%;">
                            <uni-col :span="6">
                                <view class="clean-btn" @click="exitEditMode">
                                    <text>取消</text>
                                </view>
                            </uni-col>
                            <uni-col :span="12">
                                <view class="selection-info">
                                    <text class="count">已选 {{ selectedPlantIds.length }} 项</text>
                                    <text v-if="selectedPlantIds.length > 0 && selectedPlantIds.length < total"
                                        class="select-all-tag" @click="selectAllInTag">选中全部 {{ total }} 项</text>
                                </view>
                            </uni-col>
                            <uni-col :span="6">
                                <view class="save-btn-rect" @click="handleBatchDone">
                                    <text>完成</text>
                                </view>
                            </uni-col>
                        </uni-row>
                    </view>
                </template>
            </view>

            <!-- 横向滚动标签 / 操作选择 -->
            <view class="tag-scroll-container">
                <template v-if="!isEditMode">
                    <scroll-view scroll-x="true" class="tag-scroll-view" :show-scrollbar="false"
                        :scroll-into-view="'tag-item-' + (currentTagIndex > 1 ? currentTagIndex - 1 : 0)"
                        scroll-with-animation>
                        <view class="tag-flex-box" id="tag-container">
                            <view v-for="(item, index) in tagList" :key="item.id" :id="'tag-item-' + index"
                                class="tag-item" :class="{ 'active': currentTagIndex === index }"
                                @click="selectTag(index, item)">
                                <text :id="'tag-text-' + index" class="tag-text">{{ item.name }}</text>
                            </view>
                            <view class="slider-bar" :style="sliderStyle"></view>
                        </view>
                    </scroll-view>
                </template>
                <template v-else>
                    <view class="action-selection-bar">
                        <view class="select-all-wrapper" @click="toggleSelectAll">
                            <uni-icons :type="isSelectAll ? 'checkbox-filled' : 'circle'" size="20"
                                color="#6B8857"></uni-icons>
                            <text class="select-all-text">全选</text>
                        </view>
                        <scroll-view scroll-x class="action-scroll-view" :show-scrollbar="false">
                            <view class="action-list">
                                <view v-for="item in careOptions" :key="item.type" class="action-item"
                                    :class="{ 'active': batchActionType === item.type }"
                                    @click="batchActionType = item.type">
                                    <view class="iconfont" :class="item.icon" style="font-size: 18px;"></view>
                                    <text class="action-name">{{ item.name }}</text>
                                </view>
                            </view>
                        </scroll-view>
                    </view>
                </template>
            </view>
        </view>

        <!-- 🌟 4. 中间独立滚动区域 -->
        <!-- flex:1 让它自动填满剩余空间，height:0 防止被内容撑大 -->
        <scroll-view scroll-y class="content-scroll-view" @scrolltolower="onScrollToLower">
            <!-- 空状态 -->
            <view v-if="plantsList.length === 0" class="empty-wrapper">
                <image src="/static/icon/c2m.svg" class="empty-icon" mode="aspectFit"></image>
                <text class="empty-text">点击右上方添加第一颗植物吧</text>
            </view>

            <!-- 瀑布流列表 -->
            <view v-else class="waterfall-wrapper">
                <WaterfallBox :list="plantsList" idKey="id" cols="2">
                    <template #item="{ item }">
                        <view class="plant-card" @click="handleCardClick(item)" @longpress="enterEditMode(item)">
                            <view class="image-wrapper"
                                :style="{ paddingBottom: (item.cover.height / item.cover.width * 100) + '%' }">
                                <image :src="item.cover.url" mode="aspectFill" class="plant-image" lazy-load
                                    :class="{ 'show': loadedImagesMap[item.id] }" @load="onImgLoad(item)"></image>
                                <!-- 选择框 -->
                                <view v-if="isEditMode" class="checkbox-wrapper" @click.stop="toggleSelect(item.id)">
                                    <uni-icons :type="selectedPlantIds.includes(item.id) ? 'checkbox-filled' : 'circle'"
                                        size="22"
                                        :color="selectedPlantIds.includes(item.id) ? '#6B8857' : 'rgba(255,255,255,0.8)'"></uni-icons>
                                </view>
                            </view>
                            <view class="plant-info">
                                <text class="plant-name">{{ item.name }}</text>
                                <view class="plant-tags" v-if="item.tags"></view>
                            </view>
                            <!-- 卡片右下角删除图标 -->
                            <view v-if="isEditMode" class="card-delete-btn" @click.stop="handleSingleDelete(item)">
                                <uni-icons type="trash-filled" size="14" color="#dd524d"></uni-icons>
                            </view>
                        </view>
                    </template>
                </WaterfallBox>

                <!-- 底部垫片：防止内容被 TabBar 遮挡 -->
                <view style="height: 20px;"></view>
            </view>
        </scroll-view>

    </view>
</template>


<script>
import { callContainer } from '../utils/request';
import WaterfallBox from './WaterfallBox.vue';

export default {
    components: {
        WaterfallBox,
    },
    /**
     * 组件名称，也就是开发者使用的标签
     */
    name: 'home',
    /**
     * 组件涉及的事件声明，只有声明过的事件，才能被正常发送
     */
    emits: [],
    /**
     * 属性声明，组件的使用者会传递这些属性值到组件
     */
    props: {
        navText: {
            type: String,
            default: "",

        },
    },
    /**
     * 组件内部变量声明
     */
    data() {
        return {
            menuButtonInfo: {},
            paddingLeft: 0,
            topBarHeight: 0,
            windowWidth: 0,
            familyRange: [],
            value: null,
            searchValue: "",
            tagList: [],
            currentTagIndex: 0,
            sliderLeft: 0,
            sliderWidth: 0,
            sliderTimer: null,
            plantsList: [],
            allPlantsList: [],
            loadedImagesMap: {}, // 用于追踪图片加载状态
            isFirstLoad: true,
            isSelecting: false,
            currentFamilyIndex: 0,
            isFiltering: false, // 防止过滤期间的并发修改
            isEditMode: false,
            selectedPlantIds: [],
            isSelectAll: false,
            careOptions: [],
            batchActionType: '',
            // 分页相关
            page: 1,
            pageSize: 20,
            total: 0,
            isLoading: false,
            isNoMore: false,
            searchTimer: null,
            // 🚀 新增：防止重复初始化的锁
            isInitializing: false,
        }
    },
    computed: {
        // 动态生成滑块样式
        sliderStyle() {
            return {
                transform: `translateX(${this.sliderLeft}px)`,
                width: `${this.sliderWidth}px`
            }
        }
    },
    /**
     * 属性变化监听器实现
     */
    watch: {
        searchValue() {
            // 搜索内容变化时增加防抖处理
            if (this.searchTimer) clearTimeout(this.searchTimer);
            this.searchTimer = setTimeout(() => {
                this.getPlantsList();
            }, 500);
        },
        allPlantsList: {
            handler(newVal, oldVal) {
                if (newVal) {
                    newVal.forEach((plant, idx) => {
                        if (plant.tags) {
                            plant.tags.forEach(tag => {
                            });
                        }
                    });
                }
            },
            deep: true
        }
    },
    /**
     * 规则：如果没有配置expose，则methods中的方法均对外暴露，如果配置了expose，则以expose的配置为准向外暴露
     * ['publicMethod'] 含义为：只有 `publicMethod` 在实例上可用
     * 
     * 注意：如果在data中声明了一个变量，此时组件配置了 expose字段，但未在expose字段中包含此变量。会导致该变量被标记为`private`：仅能在组件内使用，不能在组件外访问
     */
    //expose: [''],
    methods: {
        async loadFamilyData() {
            if (this.isInitializing) return;

            try {
                this.isInitializing = true;
                // 🚀 使用同步获取，确保校验逻辑即时执行，防止异步间隙导致的重复触发
                const cachedFamilyId = uni.getStorageSync('familyId');
                const familyList = uni.getStorageSync('family') || [];

                // 🚀 核心拦截：如果当前已有数据且家庭 ID 没变，绝对不执行后续加载
                if (this.value && cachedFamilyId && String(this.value) === String(cachedFamilyId) && this.plantsList.length > 0) {
                    this.isInitializing = false;
                    return;
                }

                if (familyList && Array.isArray(familyList) && familyList.length > 0) {
                    this.familyRange = familyList.map(item => ({
                        text: item.name,
                        value: item.id,
                        disable: false
                    }));

                    const targetId = cachedFamilyId || this.familyRange[0].value;

                    // 只有在 ID 真正变化时（如切换家庭）才重置状态
                    if (String(this.value) !== String(targetId)) {
                        this.value = targetId;
                        this.currentFamilyIndex = this.familyRange.findIndex(item => String(item.value) === String(targetId));
                        if (this.currentFamilyIndex === -1) this.currentFamilyIndex = 0;

                        // ID 变化了，需要重置分页和列表
                        this.plantsList = [];
                        this.page = 1;
                        this.isNoMore = false;
                        this.currentTagIndex = 0;
                    }
                } else {
                    this.familyRange = [];
                    this.currentFamilyIndex = 0;
                    this.value = null;
                }

                // 只有通过了上方拦截（即 ID 变化或初始加载）才会执行到这里
                await this.$nextTick();
                await this.getTagList();
                await this.getPlantsList();
            } catch (error) {
                console.error("加载家庭数据失败:", error);
            } finally {
                this.isInitializing = false;
            }
        },
        async refreshFamilyList() {
            try {
                const user = await callContainer("/api/login");
                const familyList = user.data.family || [];

                // 更新缓存
                await new Promise((resolve) => {
                    uni.setStorage({ key: "family", data: familyList, success: resolve })
                });

                // 更新家庭选择器的选项
                if (familyList && Array.isArray(familyList) && familyList.length > 0) {
                    this.familyRange = familyList.map(item => ({
                        text: item.name,
                        value: item.id,
                        disable: false
                    }));

                    // 确保当前选择的家庭ID在新列表中
                    const currentFamilyExists = this.familyRange.some(item => item.value === this.value);
                    if (!currentFamilyExists && this.familyRange.length > 0) {
                        // 如果当前家庭不在新列表中，切换到第一个
                        this.value = this.familyRange[0].value;
                        this.currentFamilyIndex = 0;

                        // 更新缓存
                        await new Promise((resolve) => {
                            uni.setStorage({ key: "familyId", data: this.value, success: resolve })
                        });

                        // 刷新数据
                        await this.getTagList();
                        await this.getPlantsList();
                    }
                } else {
                    this.familyRange = [];
                }

            } catch (error) {
                console.error("刷新家庭列表失败:", error);
            }
        },
        async getPlantsList(isLoadMore = false) {
            if (this.isLoading || (isLoadMore && this.isNoMore)) return;

            if (!isLoadMore) {
                this.page = 1;
                this.isNoMore = false;
                // 不再立即清空列表，防止页面跳动
            }

            this.isLoading = true;
            const familyId = this.value;
            const currentTag = this.tagList[this.currentTagIndex];

            try {
                const res = await callContainer("/api/plant/list", {
                    familyId: familyId,
                    page: this.page,
                    pageSize: this.pageSize,
                    tagId: currentTag ? currentTag.id : 0,
                    keyword: this.searchValue
                });

                console.log("plants list res:", res);

                // 处理分页返回的数据格式
                let rawData = [];
                if (res?.data?.list) {
                    rawData = res.data.list;
                    this.total = res.data.total;
                } else if (Array.isArray(res?.data)) {
                    // 兼容老接口
                    rawData = res.data;
                    this.total = rawData.length;
                }

                // 处理数据冻结
                const frozenData = rawData.map(item => {
                    let frozenTags = null;
                    if (item.tags && Array.isArray(item.tags)) {
                        frozenTags = item.tags.map(tag => Object.freeze({ ...tag }));
                        Object.freeze(frozenTags);
                    }
                    const newItem = { ...item, tags: frozenTags };
                    return Object.freeze(newItem);
                });

                if (isLoadMore) {
                    this.plantsList = [...this.plantsList, ...frozenData];
                } else {
                    // 🚀 优化：对比新旧数据第一条的 ID。
                    // 如果 ID 没变且当前已经在第一页，说明内容没有实质更新，不替换引用以保持滚动位置
                    const isSameFirstItem = this.plantsList.length > 0 &&
                        frozenData.length > 0 &&
                        String(this.plantsList[0].id) === String(frozenData[0].id);

                    if (isSameFirstItem && this.page === 1) {
                        console.log("首屏数据未发生导致位移的变化，保持当前滚动位置");
                    } else {
                        // 只有数据变了（如新增或删除），才替换列表
                        this.plantsList = frozenData;
                    }
                }

                this.isNoMore = this.plantsList.length >= this.total;
            } catch (error) {
                console.error("获取植物列表失败:", error);
            } finally {
                this.isLoading = false;
            }
        },
        onScrollToLower() {
            if (!this.isNoMore && !this.isLoading) {
                this.page++;
                this.getPlantsList(true);
            }
        },
        async handleFamilyChange(e) {
            const selectedIndex = e.detail.value;
            this.currentFamilyIndex = selectedIndex;
            const newFamilyId = this.familyRange[selectedIndex].value;

            try {
                // 调用后端切换家庭接口
                await callContainer("/api/family/switch", {
                    familyId: newFamilyId
                });
                console.log("家庭切换成功");

                // 更新storage中的familyId
                await new Promise((resolve) => {
                    uni.setStorage({ key: "familyId", data: newFamilyId, success: resolve })
                });

                // 触发家庭切换事件，通知其他组件
                uni.$emit('familyChanged', newFamilyId);

            } catch (error) {
                console.error("切换家庭失败:", error);

                // 显示错误提示
                const errorMsg = error?.msg || error?.message || "切换家庭失败，请稍后重试";
                uni.showToast({
                    title: errorMsg,
                    icon: 'none',
                    duration: 2000
                });

                // 恢复之前的选择
                this.currentFamilyIndex = this.familyRange.findIndex(item => item.value === this.value);

                // 刷新家庭列表，移除无权限的家庭
                await this.refreshFamilyList();
                return;
            }

            // 使用新familyId更新数据
            this.value = newFamilyId;
            this.currentTagIndex = 0;


            // 清空旧数据
            this.tagList = [];
            this.allPlantsList = [];
            this.plantsList = [];
            this.loadedImagesMap = {}; // 清空图片加载状态

            // 等待 DOM 更新
            await this.$nextTick();

            // 直接获取新家庭的标签和植物列表
            await this.getTagList();
            await this.getPlantsList();

            this.$nextTick(() => {
                setTimeout(() => {
                    this.updateSliderPosition(0);
                }, 200);
            });

            wx.vibrateShort({ type: "light" });
        },
        toggleFamilySelect() {
            // 微信原生 picker 会自动展开，无需额外触发
            // 这里可以添加一些额外的逻辑，比如聚焦或高亮
            console.log("触发家庭选择器");
        },
        onTouchStart() {
            // 按钮按下时的样式变化
            this.isSelecting = true;
        },
        onTouchEnd() {
            // 按钮释放时的样式恢复
            setTimeout(() => {
                this.isSelecting = false;
            }, 200);
        },
        async getTagList() {
            const familyId = this.value;

            try {
                const [tagList, careList] = await Promise.all([
                    callContainer("/api/tag/", { familyId: familyId }),
                    callContainer("/api/care/", { familyId: Number(familyId) })
                ]);

                console.log("tagList:", tagList)
                const apiTags = tagList?.data || []
                this.tagList = [
                    { name: "全部", id: 0 },
                    ...apiTags.map(item => ({
                        name: item.name,
                        id: item.id,
                        ...item
                    }))
                ]

                // 批量操作项
                const apiCares = careList?.data || [];
                // 首页不需要“成长记录”操作，过滤掉
                this.careOptions = apiCares.filter(c => c.type !== 'record');
                if (this.careOptions.length > 0) {
                    this.batchActionType = this.careOptions[0].type;
                }

                console.log("tags:", this.tagList)
                this.$nextTick(() => {
                    // 稍微延迟一点，确保 DOM 渲染完成
                    setTimeout(() => {
                        this.updateSliderPosition(0);
                    }, 200);
                });
            } catch (error) {
                console.error("获取标签列表失败:", error)
            }
        },
        enterEditMode(item) {
            if (this.isEditMode) return;
            wx.vibrateShort({ type: "medium" });
            this.isEditMode = true;
            this.selectedPlantIds = [item.id];
            this.checkSelectAll();
        },
        async handleSingleDelete(item) {
            uni.showModal({
                title: '提示',
                content: `确定要删除“${item.name}”吗？删除后相关的记录也会被清除。`,
                confirmColor: '#dd524d',
                success: async (res) => {
                    if (res.confirm) {
                        uni.showLoading({ title: '正在删除...' });
                        try {
                            await callContainer("/api/plant/delete", { id: item.id });
                            uni.showToast({ title: '已删除', icon: 'success' });
                            await this.getPlantsList();
                        } catch (e) {
                            console.error("删除失败:", e);
                            uni.showToast({ title: '删除失败', icon: 'none' });
                        } finally {
                            uni.hideLoading();
                        }
                    }
                }
            });
        },
        exitEditMode() {
            this.isEditMode = false;
            this.selectedPlantIds = [];
            this.isSelectAll = false;
        },
        handleCardClick(item) {
            if (this.isEditMode) {
                this.toggleSelect(item.id);
            } else {
                this.gotoDetail(item);
            }
        },
        toggleSelect(id) {
            const index = this.selectedPlantIds.indexOf(id);
            if (index > -1) {
                this.selectedPlantIds.splice(index, 1);
            } else {
                this.selectedPlantIds.push(id);
            }
            this.checkSelectAll();
            wx.vibrateShort({ type: "light" });
        },
        toggleSelectAll() {
            if (this.isSelectAll) {
                this.selectedPlantIds = [];
                this.isSelectAll = false;
            } else {
                this.selectedPlantIds = this.plantsList.map(p => p.id);
                this.isSelectAll = true;
            }
            wx.vibrateShort({ type: "light" });
        },
        async selectAllInTag() {
            uni.showLoading({ title: '正在获取全部 ID...' });
            try {
                const currentTag = this.tagList[this.currentTagIndex];
                const res = await callContainer("/api/plant/list", {
                    familyId: this.value,
                    page: 1,
                    pageSize: this.total, // 一次性获取所有 ID
                    tagId: currentTag ? currentTag.id : 0,
                    keyword: this.searchValue
                });

                if (res?.data?.list) {
                    this.selectedPlantIds = res.data.list.map(p => p.id);
                    this.isSelectAll = true;
                    uni.showToast({ title: `已选中全部 ${this.selectedPlantIds.length} 项`, icon: 'none' });
                }
            } catch (e) {
                console.error("获取全部 ID 失败:", e);
            } finally {
                uni.hideLoading();
            }
        },
        checkSelectAll() {
            this.isSelectAll = this.plantsList.length > 0 && this.selectedPlantIds.length === this.plantsList.length;
        },
        async handleBatchDone() {
            if (this.selectedPlantIds.length === 0) {
                uni.showToast({ title: '请先选择植物', icon: 'none' });
                return;
            }
            if (!this.batchActionType) {
                uni.showToast({ title: '请选择操作', icon: 'none' });
                return;
            }

            uni.showLoading({ title: '正在处理...' });
            try {
                const now = new Date();
                const logTime = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;

                await callContainer("/api/plant/log/add", {
                    plantIds: this.selectedPlantIds,
                    actionType: this.batchActionType,
                    content: `批量执行了${this.careOptions.find(c => c.type === this.batchActionType)?.name || ''}操作`,
                    logTime: logTime,
                    imageIds: []
                });

                uni.showToast({ title: '操作成功', icon: 'success' });
                this.exitEditMode();

                // 触发日历记录刷新
                uni.$emit('refreshLogCalendar');

                // 🚀 优化刷新逻辑：刷新当前已加载的所有数据，而不是跳回第一页
                const currentLoadedCount = this.plantsList.length;
                const res = await callContainer("/api/plant/list", {
                    familyId: this.value,
                    page: 1,
                    pageSize: Math.max(currentLoadedCount, 20),
                    tagId: this.tagList[this.currentTagIndex]?.id || 0,
                    keyword: this.searchValue
                });

                if (res?.data?.list) {
                    this.plantsList = res.data.list.map(item => {
                        let frozenTags = null;
                        if (item.tags && Array.isArray(item.tags)) {
                            frozenTags = item.tags.map(tag => Object.freeze({ ...tag }));
                            Object.freeze(frozenTags);
                        }
                        return Object.freeze({ ...item, tags: frozenTags });
                    });
                }
            } catch (error) {
                console.error("批量操作失败:", error);
                uni.showToast({ title: '操作失败', icon: 'none' });
            } finally {
                uni.hideLoading();
            }
        },
        searchPlant(e) {


        },
        selectTag(index, item) {
            if (this.currentTagIndex === index) return;
            wx.vibrateShort({ type: "medium" })

            // 优化标签切换逻辑：先更新索引，再请求，由后端过滤
            this.currentTagIndex = index;
            this.getPlantsList();

            // 滑块动画逻辑
            const query = uni.createSelectorQuery().in(this);
            query.select('#tag-container').boundingClientRect();
            query.select('#tag-text-' + index).boundingClientRect();
            query.exec((res) => {
                if (res[0] && res[1]) {
                    const containerLeft = res[0].left;
                    const currentTextLeft = res[1].left;
                    const currentTextWidth = res[1].width;
                    const ratio = 22 / 18;
                    const finalWidth = currentTextWidth * ratio;
                    const widthDiff = finalWidth - currentTextWidth;
                    const finalLeft = (currentTextLeft - containerLeft) - (widthDiff / 2);

                    // 使用 Vue 的 nextTick 确保响应式更新
                    this.$nextTick(() => {
                        this.sliderWidth = finalWidth;
                        this.sliderLeft = finalLeft;
                    });

                    if (this.sliderTimer) clearTimeout(this.sliderTimer);
                    this.sliderTimer = setTimeout(() => {
                        this.updateSliderPosition(index);
                    }, 350);
                }
            });
        },


        updateSliderPosition(index) {
            const query = uni.createSelectorQuery().in(this);
            query.select('#tag-container').boundingClientRect();
            query.select('#tag-text-' + index).boundingClientRect();

            query.exec((res) => {
                if (res[0] && res[1]) {
                    const containerLeft = res[0].left; // 容器距离屏幕左边的距离
                    const textLeft = res[1].left;      // 文字距离屏幕左边的距离
                    const textWidth = res[1].width;    // 文字宽度

                    // 计算相对位置：文字位置 - 容器位置 = 滑块在容器内的 left
                    // 注意：因为是在 scroll-view 内部，这种相对计算方式即使在滚动后也是正确的
                    this.sliderLeft = textLeft - containerLeft;
                    this.sliderWidth = textWidth;
                }
            });
        },
        onImgLoad(item) {
            // 使用独立的映射对象来追踪加载状态，避免修改冻结对象
            if (item.id && !this.loadedImagesMap[item.id]) {
                this.$set(this.loadedImagesMap, item.id, true);
            }

        },
        // 确保图片显示的方法（用于页面返回时恢复图片状态）
        ensureImagesVisible() {
            this.$nextTick(() => {
                if (this.plantsList && this.plantsList.length > 0) {
                    this.plantsList.forEach(plant => {
                        if (!this.loadedImagesMap[plant.id]) {
                            // 重置状态，触发重新加载
                            this.$set(this.loadedImagesMap, plant.id, false);
                        }
                    });
                }
            });
        },
        goAddPage() {
            wx.vibrateShort({ type: "medium" })
            // 传入当前家庭ID
            uni.navigateTo({ url: `/pages/plantEdit/plantEdit?type=add` });
        },
        gotoDetail(item) {
            uni.navigateTo({
                url: `/pages/plantDetail/plantDetail?id=${item.id}`
            })
        },
        onPageShow() {
            // 首次加载时不调用，避免与 created 冲突
            if (this.isFirstLoad) {
                this.isFirstLoad = false;
                return;
            }

            // 🚀 核心逻辑：不再盲目刷新。
            // 只有当缓存中的家庭ID与当前不一致时（说明用户在设置或其他页面改了家庭）才刷新
            const cachedFamilyId = uni.getStorageSync('familyId');
            if (this.value && cachedFamilyId && String(this.value) !== String(cachedFamilyId)) {
                console.log("检测到家庭 ID 变更，执行强制同步");
                this.loadFamilyData();
            }
        },
    },
    beforeDestroy() {
        // 🚀 销毁时移除监听
        uni.$off('refreshHomeList');
    },
    async created() {
        // 🚀 监听全局刷新事件
        uni.$on('refreshHomeList', () => {
            console.log("收到全局刷新指令，正在重置列表...");
            this.page = 1;
            this.isNoMore = false;
            this.getPlantsList(false);
        });

        const menuButtonInfo = wx.getMenuButtonBoundingClientRect()
        this.menuButtonInfo = menuButtonInfo
        const systemInfo = uni.getWindowInfo()
        this.paddingLeft = systemInfo.screenWidth - menuButtonInfo.right
        const app = getApp()
        this.topBarHeight = app.globalData.topBarHeight;
        this.windowWidth = app.globalData.windowWidth;
        const user = await callContainer("/api/login")
        console.log("callContainer login:", user)

        const userInfo = user.data.user;
        const familyList = user.data.family;

        // 保存用户信息
        await new Promise((resolve) => {
            uni.setStorage({ key: "userInfo", data: userInfo, success: resolve })
        })

        // 保存家庭列表
        await new Promise((resolve) => {
            uni.setStorage({ key: "family", data: familyList, success: resolve })
        })

        // 确定默认家庭ID：优先使用 userInfo 中的 currentFamilyId，否则使用家庭列表的第一个
        const defaultFamilyId = userInfo?.currentFamilyId || (familyList && familyList[0]?.id);


        // 保存默认家庭ID到缓存
        if (defaultFamilyId) {
            await new Promise((resolve) => {
                uni.setStorage({ key: "familyId", data: defaultFamilyId, success: resolve })
            })
        }

        this.loadFamilyData()
        // uni.$off('refreshFamilyList');
        // uni.$on('refreshFamilyList', (data) => {
        //     const user = await callContainer("/api/login")
        //     await new Promise((resolve) => {
        //         uni.setStorage({ key: "family", data: user.data.family, success: resolve })
        //     })
        //     this.loadFamilyData()
        // })

    }
}
</script>

<style scoped lang="scss">
/* 1. 页面容器：占满全屏，垂直排列 */
.home-container {
    height: 100vh;
    /* 关键：锁定高度 */
    display: flex;
    flex-direction: column;
    overflow: hidden;
    /* 禁止整个页面拖动 */
    box-sizing: border-box;
    background-color: var(--bg-color);
    /* 建议加个背景色，防止列表滚动到底部露白 */
}

/* 2. 头部固定区域组 */
.fixed-header-group {
    flex-shrink: 0;
    /* 禁止压缩 */
    z-index: 10;
    background-color: var(--bg-color);
    /* 必须给背景色，否则列表滚动时会透过文字看到下面 */
    /* 如果你的设计是背景图通铺，这里可以用 transparent，但要注意视觉重叠 */
}

.selection-info {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;

    .count {
        font-size: 26rpx;
        font-weight: bold;
        color: var(--text-color);
    }

    .select-all-tag {
        font-size: 22rpx;
        color: var(--primary-color);
        text-decoration: underline;
        margin-top: 4rpx;
    }
}

.edit-mode-header {
    width: 100%;
    padding: 0 10px;
}

.clean-btn,
.save-btn-rect {
    width: 100%;
    height: 74rpx;
    background: var(--bg-btn-color);
    border-radius: 60rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s cubic-bezier(0.25, 0.8, 0.25, 1);
    border: 1px solid var(--border-color);
    font-size: 28rpx;
    color: var(--text-color);

    &:active {
        transform: scale(0.92) translateY(2px);
    }
}

.save-btn-rect {
    color: var(--primary-color);
    font-weight: bold;
}

.action-selection-bar {
    display: flex;
    align-items: center;
    padding: 0 0 0 10px;
    height: 140rpx;
}

.select-all-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    margin-right: 30rpx;
    flex-shrink: 0;

    .select-all-text {
        font-size: 10px;
        color: var(--primary-color);
        margin-top: 4rpx;
    }
}

.action-scroll-view {
    flex: 1;
    width: 0;
    white-space: nowrap;
    height: 100%;
}

.action-list {
    display: inline-flex;
    align-items: center;
    padding-right: 60rpx;
    height: 100%;
}

.action-item {
    display: inline-flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    margin-right: 20rpx;
    width: 120rpx;
    height: 120rpx;
    background: var(--bg-btn-color);
    border-radius: 24rpx;
    transition: all 0.2s;
    flex-shrink: 0;

    .iconfont {
        color: #666;

        @media (prefers-color-scheme: dark) {
            color: rgba(245, 245, 245, 0.6);
        }
    }

    &.active {
        background: var(--primary-color);

        .action-name {
            color: #fff;
        }

        .iconfont {
            color: #fff;
        }

        @media (prefers-color-scheme: dark) {
            background: #FAF2CB;

            .action-name {
                color: #0A3323;
            }

            .iconfont {
                color: #0A3323;
            }
        }
    }

    .action-name {
        font-size: 11px;
        color: var(--text-sub);
        margin-top: 8rpx;
        width: 100%;
        text-align: center;
        line-height: 1.2;
        display: -webkit-box;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        overflow: hidden;

        @media (prefers-color-scheme: dark) {
            color: rgba(245, 245, 245, 0.6);
        }
    }
}


.checkbox-wrapper {
    position: absolute;
    top: 10rpx;
    right: 10rpx;
    z-index: 5;
    background: var(--bg-btn-color);
    border-radius: 50%;
    width: 44rpx;
    height: 44rpx;
    display: flex;
    align-items: center;
    justify-content: center;
}

/* 3. 滚动区域：自动填满剩余空间 */
.content-scroll-view {
    flex: 1;
    /* 占据剩余高度 */
    height: 0;
    /* 🌟 关键：强制触发 Flex 计算，防止 scroll-view 被内容撑开导致失效 */
    overflow: hidden;
    margin-bottom: 160rpx;
}

.waterfall-wrapper {
    padding-bottom: env(safe-area-inset-bottom);
    /* 适配 iPhone 底部安全区 */
}

.family-select {
    position: fixed;
    z-index: 999;

    /* --- 核心毛玻璃样式 --- */
    background-color: var(--bg-btn-color);
    /* 半透明白底 */
    backdrop-filter: blur(10px);
    /* 模糊背景 */
    -webkit-backdrop-filter: blur(10px);
    /* 兼容 iOS */
    border-radius: 20px;
    border: 1px solid var(--border-color);
    /* 极细的浅色边框 */
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    /* 增强阴影效果 */

    /* 布局与过渡 */
    display: flex;
    align-items: center;
    justify-content: flex-start;
    padding: 0 12px;
    gap: 6px;
    transition: all 0.3s ease;
    box-sizing: border-box;
    /* 确保边框不撑大尺寸 */

    /* 交互状态 */
    &:active,
    &.selecting {
        transform: scale(0.95);
        box-shadow: 0 1px 4px rgba(0, 0, 0, 0.15);
    }

    &.selecting .family-select-icon {
        transform: scale(0.9);
        opacity: 1;
    }
}

.family-select-icon {
    flex-shrink: 0;
    opacity: 0.8;
    transition: opacity 0.2s;
}

.family-select-icon uni-icons {
    font-size: 16px;
}

/* 微信 picker 样式调整 */
.custom-select {
    flex: 1;
    height: 100%;
    background: transparent;
    border: none;
    padding: 0;
    margin: 0;
    display: flex;
    align-items: center;
}

.family-name-wrapper {
    display: flex;
    align-items: center;
    flex-direction: row;
}

.family-select-text {
    font-size: 14px;
    color: #333;
    opacity: 0.9;
    font-weight: 500;
}

@media (prefers-color-scheme: dark) {
    .family-select-text {
        color: var(--primary-color);
    }

    .family-select-icon ::v-deep .uni-icons {
        color: var(--primary-color) !important;
    }
}



.header-action-container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 10px;
    /* 左右留白 */
    margin-bottom: 5px;
    margin-top: 10px;
    /* 和下方 Tag 保持一点距离 */
}

.search-box-wrapper {
    /* 核心：搜索框占 82% (稍微多一点看起来更协调，留 18% 给按钮) */
    width: 86%;
}

.add-btn {
    width: 74rpx;
    /* 稍微加大一点点，更易点击 */
    height: 74rpx;

    /* 1. 微弱的线性渐变，模拟光照（上亮下暗） */
    // background: linear-gradient(145deg, #7da066, #607a4e);
    background: rgba(255, 255, 255, 0.55);
    /* 如果不支持渐变回退到纯色 */
    // background-color: #6B8857; 

    border-radius: 50%;

    display: flex;
    align-items: center;
    justify-content: center;

    transition: all 0.2s cubic-bezier(0.25, 0.8, 0.25, 1);

    /* 3. 增加一点边框让轮廓更清晰 */
    border: 1px solid rgba(255, 255, 255, 0.1);

    &:active {
        transform: scale(0.92) translateY(2px);
        /* 点击时下沉 */
    }
}

::v-deep .uni-searchbar {
    padding: 10px 0 !important;
    /* 去掉左右默认 padding */

    @media (prefers-color-scheme: dark) {
        .uni-searchbar__box-search-input {
            color: #f5f5f5 !important;
        }

        .uni-searchbar__text-placeholder {
            color: rgba(245, 245, 245, 0.7) !important;
        }
    }
}

/* 这是一个深度选择器，用于去除 uni-data-select 自带的边框，使其融入毛玻璃按钮 */
::v-deep .uni-select {
    border: none !important;
    background-color: transparent !important;
    padding: 0 !important;
    height: 100%;
    justify-content: center;
}

::v-deep .uni-select__input-text {
    font-size: 12px;
    /* 字体改小一点以适应按钮 */
    color: #333;
}

::v-deep .uni-select__selector-item {
    /* 这一行是为了覆盖原生样式，确保我们自定义的 slot 充满整行 */
    padding: 0 !important;
}

// tag css
.tag-scroll-container {
    width: 100%;
    background-color: transparent;
    padding: 5px 0;
}

.tag-scroll-view {
    width: 100%;
    white-space: nowrap;
    /* 关键：禁止换行 */
}

.tag-flex-box {
    display: flex;
    align-items: center;
    padding: 0 40rpx 0 20rpx;
    position: relative;
}

.tag-item {
    position: relative;
    /* 为了定位下划线 */
    display: inline-flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 8px;
    margin-right: 10px;
    font-size: 16px;
    color: #666;
    transition: all 0.3s;

    @media (prefers-color-scheme: dark) {
        color: #f5f5f5;
    }

    &.active {
        // color:#BC3823;
        color: #6B8857;
        /* 选中颜色 */
        font-weight: bold;
        font-size: 20px;
        /* 选中稍微变大 */

        @media (prefers-color-scheme: dark) {
            color: #FAF2CB;
        }
    }
}

/* 下划线动画样式 */
.slider-bar {
    position: absolute;
    bottom: 8px;
    /* 距离底部的位置 */
    left: 0;
    /* 初始位置，由 JS 控制 translateX */
    height: 3px;
    // background-color: #BC3823;
    background-color: #6B8857;
    /* 下划线颜色 */

    @media (prefers-color-scheme: dark) {
        background-color: #FAF2CB;
    }

    border-radius: 2px;

    /* 动画配置 */
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    /* width 和 transform 都会平滑过渡 */

    z-index: 1;
    /* 确保在文字下方或上方 */
    pointer-events: none;
    /* 不影响点击 */
}

.active-line {
    position: absolute;
    bottom: 0;
    width: 20px;
    /* 下划线宽度 */
    height: 3px;
    background-color: #6B8857;
    border-radius: 2px;
    animation: scaleIn 0.2s ease-out;
}

.plant-card {
    position: relative;
    background-color: rgba(255, 255, 255, 0.5);
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
    // transform: translateY(0);
}

.card-delete-btn {
    position: absolute;
    right: 12rpx;
    bottom: 12rpx;
    z-index: 10;
    background: rgba(255, 255, 255, 0.85);
    border-radius: 50%;
    width: 44rpx;
    height: 44rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    transition: all 0.2s;

    &:active {
        transform: scale(0.9);
        background: #fff;
    }
}

.image-wrapper {
    position: relative;
    width: 100%;
    height: 0;
    /* 背景色作为占位时的颜色 */
    background-color: rgba(255, 255, 255, 0.6);
    overflow: hidden;
}

.plant-image {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    // background-color: rgba(255,255,255,1); // 加载时的背景色
    // opacity: 0;
    transition: opacity 0.4s ease-in-out;
}

.plant-image.show {
    opacity: 1;
}

.plant-info {
    padding: 8px;
}

.plant-name {
    font-size: 14px;
    color: #333;
    font-weight: bold;
}

/* 空状态样式 */
.empty-wrapper {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding: 200rpx 0;
}

.empty-icon {
    width: 120rpx;
    height: 120rpx;
    margin-bottom: 20rpx;
    opacity: 0.4;
    filter: grayscale(100%);
}

.empty-text {
    font-size: 14px;
    color: #999;
}

@keyframes scaleIn {
    from {
        transform: scaleX(0);
    }

    to {
        transform: scaleX(1);
    }
}

.waterfall-wrapper {
    padding-bottom: env(safe-area-inset-bottom);
    /* 适配 iPhone 底部安全区 */
}
</style>