<template>
	<view class="container">
		<navBar :isHome="true" />

		<!-- 占位符 - 减小高度让卡片上移 -->
		<view :style="{ height: topBarHeight + 10 + 'px' }"></view>

		<scroll-view scroll-y class="main-scroll" :enable-back-to-top="true">

			<!-- 植物基本信息卡片 -->
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

			<!-- 植物日志 (只读模式) -->
			<view class="log-section">
				<view class="section-header">
					<view class="header-left">
						<uni-icons type="compose" size="20" color="#333"></uni-icons>
						<text class="section-title" style="margin-left: 6px;">生长日志</text>
					</view>
				</view>
				<!-- 日志列表 -->
				<view class="log-list">
					<view v-for="(group, gIndex) in logList" :key="gIndex" class="log-group">
						<text class="log-date-header">{{ group.dateStr }}</text>
						<view v-for="(log, lIndex) in group.items" :key="lIndex" class="log-item">
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
						</view>
					</view>
					<view v-if="logList.length === 0" class="empty-logs">
						<text>暂无日志记录</text>
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
			totalLogCount: 0
		};
	},
	async onLoad(options) {
		const app = getApp()
		const menuButtonInfo = uni.getMenuButtonBoundingClientRect();
		this.topBarHeight = app.globalData.topBarHeight || menuButtonInfo.top;

		if (options && options.id) {
			this.plantId = options.id
			await this.getPlant();
			// 获取养护项定义用于日志显示
			await this.getCareActions();
			this.getLogs();
		}
	},
	methods: {
		async getCareActions() {
			try {
				// 获取养护项定义，分享页面根据日志中的 actionType 匹配，即使没有 familyId 也可以通过日志本身信息或默认映射
				const result = await callContainer("/api/care/", { familyId: this.plant.familyId });
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
		}
	}
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
}

.plant-header-card {
	margin: 10px 16px;
	padding: 20px;
	background-color: rgba(255, 255, 255, 0.55);
	border-radius: 24px;
	box-shadow: 0 8px 20px rgba(107, 136, 87, 0.08);
}

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
	background-color: #E0E0E0;
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

.card-divider {
	height: 1px;
	background-color: #f0f0f0;
	margin-bottom: 16px;
	width: 100%;
}

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

.log-section {
	padding: 0 16px;
}

.section-header {
	display: flex;
	align-items: center;
	margin-bottom: 16px;
}

.section-title {
	font-size: 18px;
	font-weight: bold;
	color: #333;
	@media (prefers-color-scheme: dark) {
		color: #f5f5f5;
	}
}

.header-left {
	display: flex;
	align-items: center;
	@media (prefers-color-scheme: dark) {
		::v-deep .uni-icons {
			color: #f5f5f5 !important;
		}
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
	background-color: #ddd;
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
		// 保持原色背景但降低亮度或保持透明度
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
	gap: 6px;
}

.log-img {
	width: 70px;
	height: 70px;
	border-radius: 8px;
}

.log-text {
	font-size: 13px;
	color: #333;
	@media (prefers-color-scheme: dark) {
		color: #f5f5f5;
	}
}

.empty-logs {
	text-align: center;
	padding: 40px 0;
	color: #999;
	font-size: 14px;
	@media (prefers-color-scheme: dark) {
		color: rgba(245, 245, 245, 0.4);
	}
}
</style>
