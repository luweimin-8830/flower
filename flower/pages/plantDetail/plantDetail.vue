<template>
	<view class="container">
		<navBar />

		<!-- 🌟 新增：右上角操作按钮 (悬浮在导航栏之上) -->
		<!-- topBarHeight 是状态栏高度，我们需要往下一点点 -->
		<view class="top-actions" :style="{ top: (topBarHeight + 8) + 'px' }">
			<view class="action-btn icon-only" @click="handleNotification">
				<uni-icons type="notification" size="20" color="#333"></uni-icons>
			</view>
			<button class="action-btn icon-only share-btn" open-type="share" @tap="handleShareClick">
				<uni-icons type="upload" size="20" color="#333"></uni-icons>
			</button>
			<view class="action-btn pill-btn" @click="goEdit">
				<text>编辑</text>
			</view>
		</view>

		<!-- 占位符 -->
		<view :style="{ height: topBarHeight + 44 + 'px' }"></view>

		<scroll-view scroll-y class="main-scroll" :enable-back-to-top="true">
			<!-- 提醒时间选择弹窗 -->
			<uni-calendar
				ref="remindCalendar"
				:insert="false"
				:clear-date="true"
				:startDate="todayDate"
				@confirm="onRemindConfirm"
			/>

			<!-- 提醒事项输入弹窗 -->
			<uni-popup ref="remindPopup" type="center">
				<view class="remind-modal">
					<view class="modal-header">
						<text class="modal-title">预约提醒</text>
					</view>
					<view class="modal-body">
						<view class="input-group">
							<text class="label">日期</text>
							<text class="value">{{ remindData.date }}</text>
						</view>
						<view class="input-group">
							<text class="label">时间</text>
							<picker mode="time" :value="remindData.time" @change="e => remindData.time = e.detail.value">
								<view class="time-picker-value">{{ remindData.time }}</view>
							</picker>
						</view>
						<view class="input-group vertical">
							<text class="label">事项内容</text>
							<textarea 
								class="remind-textarea" 
								v-model="remindData.content" 
								placeholder="输入提醒内容,如:该浇水了"
								maxlength="20"
							/>
						</view>
					</view>
					<view class="modal-footer">
						<button class="btn cancel" @click="$refs.remindPopup.close()">取消</button>
						<button class="btn confirm" @click="saveRemind">保存预约</button>
					</view>
				</view>
			</uni-popup>

			<!-- 🌟 修改：植物基本信息卡片 -->
			<view class="plant-header-card">
				<!-- 上半部分：左图右文 -->
				<view class="card-top">
					<!-- 左侧图片 -->
					<image class="plant-avatar" :src="(plant.cover && plant.cover.url) ? plant.cover.url : '/static/default.svg'" mode="aspectFill" @click="previewCover"></image>

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
								<view class="iconfont" :class="item.icon" style="font-size: 28px; color: #fff;" v-if="item.icon"></view>
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
									<view class="iconfont" :class="log.icon || 'plant-jiaoshui1'" style="font-size: 16px;" :style="{ color: log.color || '#4A90E2' }"></view>
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
	onShareAppMessage(res) {
		const id = String(this.plantId);
		const name = this.plant.name || '我的植物';
		const imageUrl = (this.plant.cover && this.plant.cover.url) ? this.plant.cover.url : '';
		
		return {
			title: '分享我的' + name,
			path: '/pages/shareDetail/shareDetail?id=' + id,
			imageUrl: imageUrl
		};
	},
	onShareTimeline() {
		const id = String(this.plantId);
		const name = this.plant.name || '我的植物';
		const imageUrl = (this.plant.cover && this.plant.cover.url) ? this.plant.cover.url : '';
		
		return {
			title: '看看我养的' + name,
			query: 'id=' + id,
			imageUrl: imageUrl
		};
	},
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
			isManageMode: false,
			todayDate: '',
			remindData: {
				date: '',
				time: '08:00',
				content: ''
			}
		};
	},
	async onLoad(options) {
		const app = getApp()
		const menuButtonInfo = uni.getMenuButtonBoundingClientRect();
		this.topBarHeight = app.globalData.topBarHeight || menuButtonInfo.top;

		const now = new Date();
		this.todayDate = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;

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
		
		// #ifdef MP-WEIXIN
		uni.showShareMenu({
			withShareTicket: true,
			menus: ['shareAppMessage', 'shareTimeline']
		});
		// #endif
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
				let action = this.careActions.find(a => a.type === log.actionType);
				
				// 处理特殊类型 'record' (成长记录)
				if (!action && (log.actionType === 'record' || log.actionType === 'Growth')) {
					action = {
						name: '成长记录',
						icon: 'plant-zhiwuzhiyuan-duorouzhiwuyuan',
						color: '#A7C190'
					};
				}
				
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
		previewCover() {
			if (this.plant.cover && this.plant.cover.url) {
				uni.previewImage({
					urls: [this.plant.cover.url]
				});
			}
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
		handleShareClick() {
			console.log('用户点击了分享按钮');
		},
		handleNotification() {
			this.requestSubscribe();
		},
		requestSubscribe() {
			const templateId = 'jtmMCRDxoFP3AEDRAi0yNGo9PXI_3BGyb7bcqdlSJk4';
			wx.requestSubscribeMessage({
				tmplIds: [templateId],
				success: (res) => {
					if (res[templateId] === 'accept') {
						uni.showToast({
							title: '订阅成功',
							icon: 'success'
						});
						// 订阅成功后弹出时间选择器
						setTimeout(() => {
							this.$refs.remindCalendar.open();
						}, 500);
					} else if (res[templateId] === 'reject') {
						uni.showToast({
							title: '已取消订阅',
							icon: 'none'
						});
					}
				},
				fail: (err) => {
					console.error('订阅消息失败：', err);
					if (err.errCode === 20004) {
						uni.showModal({
							title: '提示',
							content: '请在小程序设置中开启订阅消息通知',
							confirmText: '去设置',
							success: (modalRes) => {
								if (modalRes.confirm) {
									wx.openSetting();
								}
							}
						});
					}
				}
			});
		},
		onRemindConfirm(e) {
			this.remindData.date = e.fulldate;
			this.remindData.content = ''; // 重置内容
			
			// 优先从配置中获取默认提醒时间
			const userInfo = uni.getStorageSync('userInfo');
			if (userInfo && userInfo.remindTime) {
				this.remindData.time = userInfo.remindTime;
			} else {
				this.remindData.time = '08:00';
			}
			
			this.$refs.remindPopup.open();
		},
		async saveRemind() {
			if (!this.remindData.content.trim()) {
				return uni.showToast({ title: '请输入提醒内容', icon: 'none' });
			}
			
			try {
				uni.showLoading({ title: '保存预约...' });
				
				const fullRemindTime = `${this.remindData.date} ${this.remindData.time}`;
				const familyId = uni.getStorageSync('familyId');
				
				await callContainer('/api/remind/add', {
					familyId: Number(familyId),
					plantId: Number(this.plantId),
					remindTime: fullRemindTime,
					content: this.remindData.content,
					actionType: 'remind'
				});
				
				// 同时记录一条日志
				await callContainer('/api/plant/log/add', {
					familyId: Number(familyId),
					plantId: Number(this.plantId),
					actionType: 'remind',
					content: `【预约提醒】${this.remindData.content}`,
					logTime: `${this.remindData.date}T${this.remindData.time}:00Z`,
					images: []
				});
				
				uni.showToast({
					title: '预约成功',
					icon: 'success'
				});
				this.$refs.remindPopup.close();
				this.getLogs(); // 刷新详情页日志列表
			} catch (e) {
				console.error(e);
				uni.showToast({ title: '预约失败', icon: 'none' });
			} finally {
				uni.hideLoading();
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

.share-btn {
	padding: 0;
	margin: 0;
	line-height: 1;
	background-color: rgba(255, 255, 255, 0.55);
	&::after {
		border: none;
	}
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
	height: 160rpx; // 确保高度足够展示 120rpx 的内容
}

.care-list {
	display: inline-flex;
	padding-right: 60rpx;
	height: 100%;
}

.care-item {
	display: inline-flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	margin-right: 20px;
	width: 120rpx;
	flex-shrink: 0;
}

.care-icon-box {
	width: 100rpx;
	height: 100rpx;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-bottom: 8px;
	box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);
}

.care-name {
	font-size: 11px;
	color: #555;
	text-align: center;
	width: 100%;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	line-height: 1.2;
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

.remind-modal {
	width: 600rpx;
	background-color: #fff;
	border-radius: 20rpx;
	padding: 30rpx;
	
	.modal-header {
		text-align: center;
		margin-bottom: 30rpx;
		.modal-title {
			font-size: 18px;
			font-weight: bold;
			color: #333;
		}
	}
	
	.input-group {
		display: flex;
		align-items: center;
		padding: 20rpx 0;
		border-bottom: 1px solid #f5f5f5;
		
		&.vertical {
			flex-direction: column;
			align-items: flex-start;
			border-bottom: none;
		}
		
		.label {
			width: 140rpx;
			font-size: 15px;
			color: #666;
		}
		
		.value, .time-picker-value {
			font-size: 15px;
			color: #333;
			font-weight: 500;
		}
		
		.remind-textarea {
			width: 100%;
			height: 160rpx;
			background-color: #f9f9f9;
			border-radius: 10rpx;
			padding: 20rpx;
			margin-top: 16rpx;
			font-size: 14px;
		}
	}
	
	.modal-footer {
		display: flex;
		gap: 20rpx;
		margin-top: 40rpx;
		
		.btn {
			flex: 1;
			height: 80rpx;
			line-height: 80rpx;
			border-radius: 40rpx;
			font-size: 15px;
			
			&::after { border: none; }
			
			&.cancel {
				background-color: #f5f5f5;
				color: #666;
			}
			
			&.confirm {
				background-color: #4A6139;
				color: #fff;
			}
		}
	}
}
</style>
