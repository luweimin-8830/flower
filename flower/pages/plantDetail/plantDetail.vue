<template>
	<view class="container">
		<navBar :isHome="!isOwner" />

		<!-- 🌟 新增：右上角操作按钮 (悬浮在导航栏之上) -->
		<!-- topBarHeight 是状态栏高度，我们需要往下一点点 -->
		<view class="top-actions" :style="{ top: (topBarHeight + 8) + 'px' }">
			<template v-if="isOwner">
				<template v-if="!plant.deathAt">
					<view class="action-btn icon-only" @click="handleNotification">
						<uni-icons type="notification" size="20" color="var(--text-color)"></uni-icons>
					</view>
					<view class="action-btn icon-only" @click="handleShareTimeline">
						<uni-icons type="pyq" size="20" color="var(--text-color)"></uni-icons>
					</view>
					<button class="action-btn icon-only share-btn" open-type="share" @tap="handleShareClick">
						<uni-icons type="upload" size="20" color="var(--text-color)"></uni-icons>
					</button>
					<view class="action-btn pill-btn" @click="goEdit">
						<text>编辑</text>
					</view>
					<view class="action-btn pill-btn wither-btn" @click="handleDeath">
						<text>凋零</text>
					</view>
				</template>
				<template v-else>
					<view class="action-btn pill-btn revive-btn" @click="handleRevive">
						<text>复活</text>
					</view>
					<view class="action-btn pill-btn dead-badge">
						<text>已凋零</text>
					</view>
				</template>
			</template>
			<template v-else>
				<!-- 非主人模式，仅显示分享给朋友/朋友圈按钮 -->
				<view class="action-btn icon-only" @click="handleShareTimeline">
					<uni-icons type="pyq" size="20" color="var(--text-color)"></uni-icons>
				</view>
				<button class="action-btn icon-only share-btn" open-type="share" @tap="handleShareClick">
					<uni-icons type="upload" size="20" color="var(--text-color)"></uni-icons>
				</button>
			</template>
		</view>

		<!-- 占位符 -->
		<view :style="{ height: topBarHeight + 44 + 'px' }"></view>

		<view class="main-content">
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

			<!-- 🌟 新增：植物凋零弹窗 -->
			<uni-popup ref="witherPopup" type="center">
				<view class="remind-modal">
					<view class="modal-header">
						<text class="modal-title">植物凋零</text>
					</view>
					<view class="modal-body">
						<view class="hint-text">凋零的植物会移入陈列室中</view>
						<view class="input-group" style="margin-top: 20rpx;">
							<text class="label">凋零日期</text>
							<picker mode="date" :value="witherDate" :start="plant.birthday" :end="todayDate" @change="e => witherDate = e.detail.value">
								<view class="time-picker-value">{{ witherDate }}</view>
							</picker>
						</view>
					</view>
					<view class="modal-footer">
						<button class="btn cancel" @click="$refs.witherPopup.close()">取消</button>
						<button class="btn confirm wither-confirm" @click="confirmWither">确认</button>
					</view>
				</view>
			</uni-popup>

			<!-- 🌟 修改：植物基本信息卡片 -->
			<view class="plant-header-card">
				<!-- 上半部分：左图右文 -->
				<view class="card-top">
					<!-- 左侧图片 -->
					<image class="plant-avatar" :class="{ 'grayscale': plant.deathAt }" :src="(plant.cover && plant.cover.url) ? plant.cover.url : '/static/default.svg'" mode="aspectFill" @click="previewCover"></image>

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
						<text class="stat-val">{{ plant.deathAt ? plant.deathAt : plant.birthday }}</text>
						<text class="stat-label">{{ plant.deathAt ? '离去日期' : '到家日期' }}</text>
					</view>
					<view class="vertical-line"></view>
					<view class="stat-item">
						<text class="stat-val">{{ plant.deathAt ? companionDays : plant.daysCount }}</text>
						<text class="stat-label">{{ plant.deathAt ? '陪伴天数' : '养护天数' }}</text>
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
			<view class="section-container" v-if="isOwner && !plant.deathAt">
				<!-- 代码省略，保持原样 -->
				<view class="section-header">
					<text class="section-title">日常护理</text>
				</view>
				<scroll-view scroll-x class="care-scroll-view" :show-scrollbar="false">
					<view class="care-list">
						<view class="care-item" v-for="(item, index) in careActions" :key="index"
							@click="handleCare(item)">
							<view class="care-icon-box" :style="{ backgroundColor: item.color }">
								<view class="iconfont" :class="item.icon" style="font-size: 28px; color: #333;" v-if="item.icon"></view>
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
						<uni-icons type="compose" size="20" color="var(--primary-color)"></uni-icons>
						<text class="section-title" style="margin-left: 6px;">植物日志</text>
					</view>
					<view class="header-right" v-if="isOwner && !plant.deathAt">
						<text class="edit-btn" @click="isManageMode = !isManageMode">{{ isManageMode ? '完成' : '编辑' }}</text>
						<uni-icons type="plusempty" size="24" color="var(--primary-color)" style="margin-left: 15px;" @click="goAddLog"></uni-icons>
					</view>
				</view>
				<!-- 日志列表 -->
				<scroll-view scroll-y class="log-scroll-inner" :enable-back-to-top="true">
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
					<view style="height: 100px;"></view>
				</scroll-view>
			</view>
		</view>
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
			query: `id=${id}&isShare=1`,
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
			witherDate: '',
			isOwner: true,
			remindData: {
				date: '',
				time: '08:00',
				content: ''
			},
			allLogImages: [] // 所有日志中的图片 URL 列表
		};
	},
	computed: {
		companionDays() {
			if (!this.plant.deathAt || !this.plant.birthday) return 0;
			const start = new Date(this.plant.birthday.replace(/-/g, '/'));
			const end = new Date(this.plant.deathAt.replace(/-/g, '/'));
			const diff = end.getTime() - start.getTime();
			return Math.floor(diff / (1000 * 60 * 60 * 24)) + 1; // 包含起始和结束当天
		}
	},
	async onLoad(options) {
		const app = getApp()
		const menuButtonInfo = uni.getMenuButtonBoundingClientRect();
		this.topBarHeight = app.globalData.topBarHeight || menuButtonInfo.top;

		const now = new Date();
		this.todayDate = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;

		if (options && options.id) {
			this.plantId = options.id
			
			// 如果是从朋友圈分享进入，强制设为非主人模式（只读）
			// 注意：options 中的值通常是字符串
			if (options && (options.isShare === '1' || options.isShare === 1 || options.isShare === 'true')) {
				this.isOwner = false;
			} else {
				this.isOwner = true;
			}

			await Promise.all([
				this.getPlant(),
				this.getCareActions()
			]);
			// 确保有了 careActions 后再获取并格式化日志
			this.getLogs();
		}

		uni.$on('refreshPlantDetail', async (data) => {
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
		uni.$off('refreshPlantDetail');
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
				this.plant.deathAt = this.plant.deathAt ? this.plant.deathAt.split('T')[0] : '';

				// 移除 openId 强校验，仅保留 isShare 标识判定
				const userInfo = uni.getStorageSync('userInfo');
				const userOpenId = userInfo ? (userInfo.openId || userInfo.openid || userInfo.openID) : null;
				const plantOpenId = this.plant.openId || this.plant.openid || this.plant.openID;
				
				console.log('Permission Info:', {
					userOpenId,
					plantOpenId,
					isOwner: this.isOwner
				});
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
			const allImages = []; // 收集所有图片
			
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
				
				// 收集当前日志的图片
				if (log.images && log.images.length > 0) {
					log.images.forEach(img => {
						allImages.push(img.url);
					});
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
			this.allLogImages = allImages; // 保存所有图片 URL
		},
		goEdit() {
			uni.navigateTo({ url: `/pages/plantEdit/plantEdit?id=${this.plant.id}&type=edit` })
		},
		handleDeath() {
			wx.vibrateShort({ type: "light" });
			this.witherDate = this.todayDate;
			this.$refs.witherPopup.open();
		},
		async handleRevive() {
			uni.showModal({
				title: '提示',
				content: '确定要复活这棵植物吗？它将重新回到花园首页。',
				success: async (res) => {
					if (res.confirm) {
						uni.showLoading({ title: '正在处理...' });
						try {
							await callContainer('/api/plant/update', {
								id: Number(this.plantId),
								deathAt: 'null' // 触发后端清除死亡日期逻辑
							});
							
							uni.hideLoading();
							uni.showToast({ title: '已重获新生', icon: 'success' });
							
							// 触发首页刷新
							uni.$emit('refreshHomeList');
							
							// 刷新当前页面数据
							this.getPlant();
						} catch (error) {
							uni.hideLoading();
							uni.showToast({ title: error.message || '操作失败', icon: 'none' });
						}
					}
				}
			});
		},
		async confirmWither() {
			uni.showLoading({ title: '正在处理...' });
			try {
				await callContainer('/api/plant/update', {
					id: Number(this.plantId),
					deathAt: this.witherDate
				});
				
				uni.hideLoading();
				uni.showToast({ title: '已移入陈列室', icon: 'success' });
				this.$refs.witherPopup.close();
				
				// 触发首页刷新
				uni.$emit('refreshHomeList');
				
				// 延迟返回首页
				setTimeout(() => {
					uni.reLaunch({ url: '/pages/index/index' });
				}, 1500);
			} catch (error) {
				uni.hideLoading();
				uni.showToast({ title: error.message || '操作失败', icon: 'none' });
			}
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
			// 找到当前点击图片在所有图片列表中的位置
			const currentUrl = images[index].url;
			const currentIndex = this.allLogImages.indexOf(currentUrl);
			
			uni.previewImage({
				current: currentIndex >= 0 ? currentIndex : 0,
				urls: this.allLogImages // 使用所有日志的图片列表
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
		handleShareTimeline() {
			uni.showToast({
				title: '点击右上角“...”菜单分享至朋友圈',
				icon: 'none'
			});
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
	height: 100vh;
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

.main-content {
	flex: 1;
	display: flex;
	flex-direction: column;
	overflow: hidden;
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
	background-color: rgba(255, 255, 255, 0.5);
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
	background-color: rgba(255, 255, 255, 0.5);
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
		@media (prefers-color-scheme: dark) {
			color: #FAF2CB;
		}
	}
}

.wither-btn {
	background-color: rgba(0, 0, 0, 0.05);
	border: 1px solid rgba(0, 0, 0, 0.03);

	text {
		color: #999 !important;
	}

	@media (prefers-color-scheme: dark) {
		background-color: rgba(255, 255, 255, 0.08);
		border-color: rgba(255, 255, 255, 0.05);

		text {
			color: rgba(245, 245, 245, 0.4) !important;
		}
	}
}

.dead-badge {
	background-color: rgba(139, 69, 19, 0.1);
	border: 1px solid rgba(139, 69, 19, 0.2);
	text {
		color: #8B4513 !important;
	}
}

.revive-btn {
	background-color: var(--primary-color);
	border: 1px solid rgba(0, 0, 0, 0.05);
	text {
		color: #fff !important;
	}
}

.grayscale {
	filter: grayscale(100%);
	opacity: 0.8;
}

/* 🌟 植物头部卡片 */
.plant-header-card {
	margin: 10px 16px;
	padding: 20px 15px;
	background-color: rgba(255, 255, 255, 0.55);
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

}

.placeholder {
	background-color: #eee;
}

.plant-info {
	flex: 1;
	display: flex;
	flex-direction: column;
	justify-content: space-between;
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
	display: inline-flex;
	align-items: center;
	justify-content: center;
	height: 22px;
	padding: 0 10px;

	text {
		font-size: 11px;
		color: #fff;
		line-height: 1;
		padding-bottom: 2px;
	}
}

/* 分割线 */
.card-divider {
	height: 1px;
	background-color: #f0f0f0;
	margin-bottom: 16px;
	width: 100%;

	@media (prefers-color-scheme: dark) {
		background-color: rgba(255, 255, 255, 0.1);
	}
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
			font-size: 15px;
			font-weight: bold;
			color: #2F3E25;
			margin-bottom: 4px;
			font-family: 'DIN', sans-serif;
			white-space: nowrap;
		}

		.stat-label {
			font-size: 11px;
			color: #888;
			white-space: nowrap;
		}
	}

	.vertical-line {
		width: 1px;
		height: 24px;
		background-color: #eee;
	}
}

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

	@media (prefers-color-scheme: dark) {
		color: #f5f5f5;
	}
}

.care-scroll-view {
	width: 100%;
	white-space: nowrap;
	height: 160rpx;
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

	@media (prefers-color-scheme: dark) {
		box-shadow: none;
	}
}

.care-name {
	font-size: 11px;
	color: #666;
	text-align: center;
	width: 100%;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	line-height: 1.2;

	@media (prefers-color-scheme: dark) {
		color: rgba(245, 245, 245, 0.6);
	}
}

/* 日志部分 */
.log-section {
	flex: 1;
	display: flex;
	flex-direction: column;
	padding: 0 16px;
	overflow: hidden;
}

.log-scroll-inner {
	flex: 1;
	height: 0;
}

.header-left {
	display: flex;
	align-items: center;

	@media (prefers-color-scheme: dark) {
		::v-deep .uni-icons {
			color: #FAF2CB !important;
		}
	}
}

.header-right {
	display: flex;
	align-items: center;
	margin-left: auto;

	@media (prefers-color-scheme: dark) {
		::v-deep .uni-icons {
			color: #FAF2CB !important;
		}
	}
}

.edit-btn {
	font-size: 14px;
	color: #888;

	@media (prefers-color-scheme: dark) {
		color: rgba(245, 245, 245, 0.6);
	}
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

	@media (prefers-color-scheme: dark) {
		color: #FAF2CB;
	}
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

	@media (prefers-color-scheme: dark) {
		color: #f5f5f5;
	}
}

.log-date-mini {
	font-size: 12px;
	color: #999;
	margin-top: 2px;

	@media (prefers-color-scheme: dark) {
		color: rgba(245, 245, 245, 0.4);
	}
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

	@media (prefers-color-scheme: dark) {
		background-color: #FAF2CB;
	}
}

.line {
	flex: 1;
	width: 1px;
	background-color: #eee;
	margin-top: 4px;

	@media (prefers-color-scheme: dark) {
		background-color: rgba(255, 255, 255, 0.1);
	}
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

	@media (prefers-color-scheme: dark) {
		filter: brightness(0.9);
	}
}

.log-content {
	font-size: 13px;
	color: #666;
	line-height: 1.5;
	margin-bottom: 8px;
	word-break: break-all;

	@media (prefers-color-scheme: dark) {
		color: rgba(245, 245, 245, 0.7);
	}
}

.log-images {
	display: flex;
	flex-wrap: wrap;
	gap: 8px;
}

.log-img {
	width: 80px;
	height: 80px;
	border-radius: 10px;
	overflow: hidden;
}

.delete-icon {
	padding: 10px;
}

.log-text {
	font-size: 13px;
	color: #333;

	@media (prefers-color-scheme: dark) {
		color: #f5f5f5;
	}
}

.remind-modal {
	width: 600rpx;
	background-color: #fff;
	border-radius: 20rpx;
	padding: 30rpx;

	@media (prefers-color-scheme: dark) {
		background-color: #1a1a1a;
	}
	
	.modal-header {
		text-align: center;
		margin-bottom: 30rpx;
		.modal-title {
			font-size: 18px;
			font-weight: bold;
			color: #333;
			@media (prefers-color-scheme: dark) {
				color: #f5f5f5;
			}
		}
	}
	
	.input-group {
		display: flex;
		align-items: center;
		padding: 20rpx 0;
		border-bottom: 1px solid #eee;

		@media (prefers-color-scheme: dark) {
			border-bottom-color: rgba(255, 255, 255, 0.1);
		}
		
		&.vertical {
			flex-direction: column;
			align-items: flex-start;
			border-bottom: none;
		}
		
		.label {
			width: 140rpx;
			font-size: 15px;
			color: #666;
			@media (prefers-color-scheme: dark) {
				color: rgba(245, 245, 245, 0.6);
			}
		}
		
		.value, .time-picker-value {
			font-size: 15px;
			color: #333;
			font-weight: 500;
			@media (prefers-color-scheme: dark) {
				color: #f5f5f5;
			}
		}
		
		.remind-textarea {
			width: 100%;
			height: 160rpx;
			background-color: #f9f9f9;
			border-radius: 10rpx;
			padding: 20rpx;
			margin-top: 16rpx;
			font-size: 14px;
			color: #333;

			@media (prefers-color-scheme: dark) {
				background-color: rgba(255, 255, 255, 0.05);
				color: #f5f5f5;
			}
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
				color: #888;
				@media (prefers-color-scheme: dark) {
					background-color: rgba(255, 255, 255, 0.1);
					color: rgba(245, 245, 245, 0.6);
				}
			}
			
			&.confirm {
				background-color: #566C44;
				color: #fff;
				@media (prefers-color-scheme: dark) {
					background-color: #FAF2CB;
					color: #0A3323;
				}
			}

			&.wither-confirm {
				background-color: #8B4513; 
				color: #fff;
				@media (prefers-color-scheme: dark) {
					background-color: #4A3728;
					color: rgba(245, 245, 245, 0.8);
				}
			}
		}
	}
}

.hint-text {
	font-size: 14px;
	color: #888;
	text-align: center;
	line-height: 1.6;
}
</style>
