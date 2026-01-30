<template>
	<view class="container">
		<navBar />

		<!-- 🌟 新增：右上角操作按钮 (悬浮在导航栏之上) -->
		<!-- topBarHeight 是状态栏高度，我们需要往下一点点 -->
		<view class="top-actions" :style="{ top: (topBarHeight + 8) + 'px' }">
			<view class="action-btn icon-only">
				<uni-icons type="upload" size="20" color="#333"></uni-icons>
			</view>
			<view class="action-btn pill-btn" @click="goEdit">
				<text>编辑</text>
			</view>
		</view>

		<!-- 占位符 -->
		<view :style="{ height: topBarHeight + 44 + 'px' }"></view>

		<scroll-view scroll-y class="main-scroll" :enable-back-to-top="true">

			<!-- 🌟 修改：植物基本信息卡片 -->
			<view class="plant-header-card">
				<!-- 上半部分：左图右文 -->
				<view class="card-top">
					<!-- 左侧图片 -->
					<image v-if="plant.cover" :src="plant.cover.url" mode="aspectFill" class="plant-avatar"></image>
					<view v-else class="plant-avatar placeholder"></view>

					<!-- 右侧信息 -->
					<view class="plant-info">
						<text class="plant-name">{{ plant.name }}</text>
						<text class="plant-desc">{{ plant.desc || '暂无描述...' }}</text>

						<!-- 标签列表 -->
						<view class="tag-list">
							<view class="tag-item" v-for="(tag, index) in plant.tags" :key="index">
								<text>{{ tag.name }}</text>
							</view>
							<!-- 如果没有标签显示添加按钮 -->
							<!-- <view class="tag-item add-tag" v-if="!plant.tags || plant.tags.length === 0">
								<uni-icons type="plusempty" size="12" color="#fff"></uni-icons>
							</view> -->
						</view>
					</view>
				</view>

				<!-- 分割线 -->
				<view class="card-divider"></view>

				<!-- 下半部分：统计数据 -->
				<view class="stats-row">
					<view class="stat-item">
						<text class="stat-val">{{ plant.birthday }}</text>
						<text class="stat-label">到家日期</text>
					</view>
					<view class="vertical-line"></view>
					<view class="stat-item">
						<text class="stat-val">{{ plant.daysCount }}</text>
						<text class="stat-label">养护天数</text>
					</view>
					<view class="vertical-line"></view>
					<view class="stat-item">
						<text class="stat-val">{{ totalLogCount }}</text>
						<text class="stat-label">植物日志</text>
					</view>
				</view>
			</view>

			<!-- ... 下面的日常护理、快捷入口、日志部分保持不变 ... -->
			<!-- 3. 日常护理 (横向滚动) -->
			<view class="section-container">
				<!-- 代码省略，保持原样 -->
				<view class="section-header">
					<text class="section-title">日常护理</text>
					<!-- <uni-icons type="notification" size="20" color="#6B8857"></uni-icons> -->
				</view>
				<scroll-view scroll-x class="care-scroll-view" :show-scrollbar="false">
					<view class="care-list">
						<view class="care-item" v-for="(item, index) in careActions" :key="index"
							@click="handleCare(item)">
							<view class="care-icon-box" :style="{ backgroundColor: item.color }">
								<uni-icons :type="item.icon" size="28" color="#fff" v-if="item.icon"></uni-icons>
								<image v-else :src="item.img" class="care-img"></image>
							</view>
							<text class="care-name">{{ item.name }}</text>
						</view>
					</view>
				</scroll-view>
			</view>

			<!-- 4. 快捷入口 -->
			<!-- <view class="quick-links">
				<view class="link-card" @click="goCalendar">
					<view class="link-left">
						<uni-icons type="calendar-filled" size="20" color="#6B8857"></uni-icons>
						<text class="link-text">植物日历</text>
					</view>
					<uni-icons type="right" size="14" color="#999"></uni-icons>
				</view>
				<view class="link-card" @click="goAlbum">
					<view class="link-left">
						<uni-icons type="image-filled" size="20" color="#6B8857"></uni-icons>
						<text class="link-text">植物相册 ({{ plant.photoCount || 0 }})</text>
					</view>
					<uni-icons type="right" size="14" color="#999"></uni-icons>
				</view>
			</view> -->

			<!-- 5. 植物日志 -->
			<view class="log-section">
				<view class="section-header">
					<view class="header-left">
						<uni-icons type="compose" size="20" color="#333"></uni-icons>
						<text class="section-title" style="margin-left: 6px;">植物日志</text>
					</view>
					<view class="header-right">
						<text class="edit-btn" @click="isManageMode = !isManageMode">{{ isManageMode ? '完成' : '编辑' }}</text>
						<uni-icons type="plusempty" size="24" color="#333" style="margin-left: 15px;" @click="goAddLog"></uni-icons>
					</view>
				</view>
				<!-- 日志列表 -->
				<view class="log-list">
					<view v-for="(group, gIndex) in logList" :key="gIndex" class="log-group">
						<text class="log-date-header">{{ group.dateStr }}</text>
						<view v-for="(log, lIndex) in group.items" :key="lIndex" class="log-item" @click="goEditLog(log)">
							<view class="log-time-col">
								<text class="log-time">{{ log.time }}</text>
								<text class="log-date-mini">{{ log.dateMini }}</text>
							</view>
							<view class="log-timeline">
								<view class="dot"></view>
								<view class="line" v-if="lIndex !== group.items.length - 1"></view>
							</view>
							<view class="log-content-box">
								<view class="log-tag pill" :style="{ backgroundColor: log.color + '33' || '#D6EAF8' }">
									<uni-icons :type="log.icon || 'checkbox-filled'" size="16" :color="log.color || '#4A90E2'"></uni-icons>
									<text class="log-text">{{ log.actionName }}</text>
								</view>
								<view class="log-content" v-if="log.content">{{ log.content }}</view>
								
								<!-- 图片展示 -->
								<view class="log-images" v-if="log.images && log.images.length > 0">
									<image v-for="(img, iIndex) in log.images" :key="iIndex" 
										:src="img.url" mode="aspectFill" class="log-img"
										@click.stop="previewLogImages(log.images, iIndex)"></image>
								</view>
							</view>
							
							<!-- 管理模式下的删除按钮 -->
							<view class="delete-icon" v-if="isManageMode" @click.stop="handleDeleteLog(log.id)">
								<uni-icons type="minus-filled" size="20" color="#dd524d"></uni-icons>
							</view>
						</view>
					</view>
				</view>
			</view>

			<view style="height: 100px;"></view>
		</scroll-view>
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
			plantId: 0,
			plant: {
				tags: [],
				desc: ''
			},
			careActions: [],
			logList: [],
			totalLogCount: 0,
			isManageMode: false
		};
	},
	async onLoad(options) {
		const app = getApp()
		const menuButtonInfo = uni.getMenuButtonBoundingClientRect();
		this.topBarHeight = app.globalData.topBarHeight || menuButtonInfo.top;

		if (options && options.id) {
			this.plantId = options.id
			await Promise.all([
				this.getPlant(),
				this.getCareActions()
			]);
			// 确保有了 careActions 后再获取并格式化日志
			this.getLogs();
		}
		uni.$off('refreshHomeList');
        uni.$on('refreshHomeList', async (data) => {
            await Promise.all([
				this.getPlant(),
				this.getCareActions()
			]);
			this.getLogs();
        })
	},
	onUnload() {
		uni.$off('refreshHomeList');
	},
	methods: {
		async getCareActions() {
			try {
				const familyId = uni.getStorageSync('familyId');
				const result = await callContainer("/api/care/", { familyId: Number(familyId) });
				if (result.data) {
					this.careActions = result.data;
				}
			} catch (error) {
				console.error("获取养护项失败:", error);
			}
		},
		async getPlant() {
			try {
				const result = await callContainer("/api/plant/", {
					id: Number(this.plantId)
				})
				this.plant = result.data
				this.plant.birthday = this.plant.birthday ? this.plant.birthday.split('T')[0] : '';
			} catch (error) {
				console.error(error)
			}
		},
		async getLogs() {
			try {
				const result = await callContainer("/api/plant/log/list", {
					plantId: Number(this.plantId)
				})
				if (result.data) {
					this.totalLogCount = result.data.length
					this.formatLogs(result.data)
				}
			} catch (error) {
				console.error(error)
			}
		},
		formatLogs(rawLogs) {
			const groups = {};
			rawLogs.forEach(log => {
				const date = new Date(log.logTime);
				const dateStr = `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日`;
				const time = `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`;
				const dateMini = `${String(date.getMonth() + 1).padStart(2, '0')}/${String(date.getDate()).padStart(2, '0')}`;
				
				if (!groups[dateStr]) {
					groups[dateStr] = {
						dateStr,
						items: []
					};
				}
				
				// 查找对应的操作名称
				const action = this.careActions.find(a => a.type === log.actionType);
				
				groups[dateStr].items.push({
					id: log.id,
					time,
					dateMini,
					type: log.actionType,
					actionName: action ? action.name : log.actionType,
					content: log.content,
					images: log.images,
					icon: action ? action.icon : '',
					color: action ? action.color : ''
				});
			});
			this.logList = Object.values(groups);
		},
		goEdit() {
			uni.navigateTo({ url: `/pages/plantEdit/plantEdit?id=${this.plant.id}&type=edit` })
		},
		goAddLog() {
			uni.navigateTo({ url: `/pages/logEdit/logEdit?plantId=${this.plantId}` })
		},
		goEditLog(log) {
			if (this.isManageMode) return;
			uni.navigateTo({ url: `/pages/logEdit/logEdit?id=${log.id}&plantId=${this.plantId}` })
		},
		async handleDeleteLog(id) {
			uni.showModal({
				title: '提示',
				content: '确定要删除这条日志吗？',
				success: async (res) => {
					if (res.confirm) {
						try {
							await callContainer("/api/plant/log/delete", { id: id });
							uni.showToast({ title: '已删除' });
							this.getLogs();
						} catch (e) { console.error(e); }
					}
				}
			});
		},
		previewLogImages(images, index) {
			uni.previewImage({
				current: index,
				urls: images.map(img => img.url)
			});
		},
		async handleCare(item) {
			uni.showLoading({ title: '正在记录...' });
			try {
				await callContainer("/api/plant/log/add", {
					plantIds: [Number(this.plantId)],
					actionType: item.type,
					content: `完成了一次${item.name}`
				});
				uni.hideLoading();
				uni.showToast({ title: '记录成功', icon: 'success' });
				// 刷新列表
				this.getLogs();
			} catch (error) {
				uni.hideLoading();
				// 后端返回“今日已记录过该操作”会进入这里
				uni.showModal({
					title: '提示',
					content: error.message || '今日已记录过该操作',
					showCancel: false
				});
			}
		},
		goAlbum() { }
	}
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
}

/* 🌟 右上角悬浮按钮组 */
.top-actions {
	position: fixed;
	right: 16px;
	z-index: 100;
	display: flex;
	align-items: center;
	gap: 12px;
}

.action-btn {
	background-color: rgba(255, 255, 255, 0.55);
	backdrop-filter: blur(5px);
	border-radius: 20px;
	display: flex;
	align-items: center;
	justify-content: center;
	border: 1px solid rgba(0, 0, 0, 0.05);
}

.icon-only {
	width: 36px;
	height: 36px;
	border-radius: 50%;
}

.pill-btn {
	padding: 6px 16px;
	height: 36px;
	box-sizing: border-box;

	text {
		font-size: 14px;
		font-weight: 500;
		color: #333;
	}
}

/* 🌟 植物头部卡片 */
.plant-header-card {
	margin: 10px 16px;
	padding: 20px;
	background-color: rgba(255, 255, 255, 0.55);
	/* 纯白背景更干净，或者微透 */
	border-radius: 24px;
	box-shadow: 0 8px 20px rgba(107, 136, 87, 0.08);
}

/* 上半部分：左图右文 */
.card-top {
	display: flex;
	margin-bottom: 20px;
}

.plant-avatar {
	width: 100px;
	height: 100px;
	border-radius: 16px;
	margin-right: 16px;
	background-color: #eee;
	flex-shrink: 0;
	/* 防止图片被挤压 */
}

.placeholder {
	background-color: #E0E0E0;
}

.plant-info {
	flex: 1;
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	/* 上下撑开 */
	padding-top: 2px;
	padding-bottom: 2px;
}

.plant-name {
	font-size: 20px;
	font-weight: bold;
	color: #2F3E25;
	margin-bottom: 4px;
}

.plant-desc {
	font-size: 13px;
	color: #888;
	line-height: 1.4;
	margin-bottom: 8px;
	/* 限制显示两行 */
	display: -webkit-box;
	-webkit-box-orient: vertical;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	overflow: hidden;
}

/* 标签样式 */
.tag-list {
	display: flex;
	flex-wrap: wrap;
	gap: 6px;
}

.tag-item {
	background-color: #566C44;
	border-radius: 100px;
	/* 改大一点确保是正圆角 */

	/* 🌟 改法：使用 Flex 居中，不再单纯依赖 padding */
	display: inline-flex;
	align-items: center;
	justify-content: center;

	/* 🌟 设定固定高度，比 padding 更稳 */
	height: 22px;
	padding: 0 10px;
	/* 只控制左右间距 */

	text {
		font-size: 11px;
		color: #fff;
		line-height: 1;
		/* 消除文字自带行高 */

		/* 🌟 视觉微调：如果觉得还偏下，就给个底部 padding 把它顶上去 */
		padding-bottom: 2px;
	}
}

.add-tag {
	background-color: #ccc;
	padding: 2px 8px;
}

/* 分割线 */
.card-divider {
	height: 1px;
	background-color: #f0f0f0;
	margin-bottom: 16px;
	width: 100%;
}

/* 下半部分：统计数据 */
.stats-row {
	display: flex;
	justify-content: space-between;
	align-items: center;

	.stat-item {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;

		.stat-val {
			font-size: 18px;
			font-weight: bold;
			color: #2F3E25;
			margin-bottom: 4px;
			font-family: 'DIN', sans-serif;
			/* 如果有数字字体更好 */
		}

		.stat-label {
			font-size: 12px;
			color: #888;
		}
	}

	.vertical-line {
		width: 1px;
		height: 24px;
		background-color: #eee;
	}
}

/* ... 下面保持原有的 CSS ... */
.section-container {
	margin: 24px 0;
	padding-left: 16px;
}

.section-header {
	display: flex;
	align-items: center;
	margin-bottom: 16px;
	padding-right: 16px;
}

.section-title {
	font-size: 18px;
	font-weight: bold;
	color: #333;
	margin-right: auto;
}

.care-scroll-view {
	width: 100%;
	white-space: nowrap;
}

.care-list {
	display: flex;
	padding-right: 16px;
}

.care-item {
	display: inline-flex;
	flex-direction: column;
	align-items: center;
	margin-right: 20px;
	width: 70px;
}

.care-icon-box {
	width: 56px;
	height: 56px;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-bottom: 8px;
	box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);
}

.care-name {
	font-size: 12px;
	color: #555;
	text-align: center;
}

/* 快捷入口 */
.quick-links {
	display: flex;
	justify-content: space-between;
	margin: 0 16px 24px 16px;
	gap: 12px;
}

.link-card {
	flex: 1;
	background-color: rgba(255, 255, 255, 0.6);
	border: 1px solid rgba(255, 255, 255, 1);
	border-radius: 16px;
	padding: 16px 12px;
	display: flex;
	align-items: center;
	justify-content: space-between;
}

.link-left {
	display: flex;
	align-items: center;
	gap: 8px;
}

.link-text {
	font-size: 14px;
	font-weight: 500;
	color: #333;
}

/* 日志部分 */
.log-section {
	padding: 0 16px;
}

.header-left {
	display: flex;
	align-items: center;
}

.header-right {
	display: flex;
	align-items: center;
	margin-left: auto;
}

.edit-btn {
	font-size: 14px;
	color: #666;
}

.log-list {
	margin-top: 10px;
}

.log-group {
	margin-bottom: 20px;
}

.log-date-header {
	font-size: 16px;
	font-weight: bold;
	color: #333;
	margin-bottom: 16px;
	display: block;
}

.log-item {
	display: flex;
	position: relative;
	padding-bottom: 20px;
}

.log-time-col {
	width: 50px;
	text-align: right;
	padding-right: 10px;
	display: flex;
	flex-direction: column;
}

.log-time {
	font-size: 16px;
	font-weight: bold;
	color: #333;
}

.log-date-mini {
	font-size: 12px;
	color: #999;
	margin-top: 2px;
}

.log-timeline {
	width: 20px;
	display: flex;
	flex-direction: column;
	align-items: center;
	position: relative;
}

.dot {
	width: 10px;
	height: 10px;
	background-color: #6B8857;
	border-radius: 50%;
	z-index: 1;
	margin-top: 6px;
}

.line {
	flex: 1;
	width: 1px;
	background-color: #ddd;
	margin-top: 4px;
}

.log-content-box {
	flex: 1;
	padding-left: 10px;
	padding-top: 2px;
}

.log-tag {
	display: inline-flex;
	align-items: center;
	background-color: #D6EAF8;
	padding: 4px 10px;
	border-radius: 20px;
	gap: 6px;
	margin-bottom: 6px;
}

.log-content {
	font-size: 13px;
	color: #666;
	line-height: 1.5;
	margin-bottom: 8px;
	word-break: break-all;
}

.log-images {
	display: flex;
	flex-wrap: wrap;
	gap: 6px;
}

.log-img {
	width: 70px;
	height: 70px;
	border-radius: 8px;
}

.delete-icon {
	padding: 10px;
}

.log-text {
	font-size: 13px;
	color: #333;
}
</style>
