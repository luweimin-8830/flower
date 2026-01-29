<template>
	<view class="calendar-wrapper">
		<!-- 顶部固定区域 -->
		<view class="fixed-header-group">
			<navBar v-if="showNav" />
			<view class="header-title-container" :style="{ paddingTop: statusBarHeight + 'px' }">
				<text class="page-title">养护历程</text>
			</view>
		</view>
		
		<scroll-view scroll-y class="scroll-body">
			<view class="log-calendar">
				<uni-calendar 
					ref="calendar"
					:insert="true" 
					:lunar="false" 
					:selected="selectedDates" 
					color="#566C44"
					@change="onDateChange"
				/>
				
				<view class="day-detail">
					<view class="detail-header">
						<text class="detail-title">{{ selectedDate }} 记录</text>
					</view>
					
					<view class="log-list">
						<view v-if="currentLogs.length === 0" class="empty-box">
							<text class="empty-text">当日无养护记录</text>
						</view>
						<view v-for="log in currentLogs" :key="log.ID" class="log-item" @click="goEdit(log)">
							<view class="log-icon" :style="{ backgroundColor: log.color + '22' }">
								<uni-icons :type="log.icon || 'checkbox-filled'" size="20" :color="log.color || '#4A6139'"></uni-icons>
							</view>
							<view class="log-info">
								<view class="log-top">
									<text class="log-name">{{ log.actionName }}</text>
									<text class="log-time">{{ formatTime(log.logTime) }}</text>
								</view>
								<text class="log-content" v-if="log.content">{{ log.content }}</text>
								<view class="log-images" v-if="log.images && log.images.length > 0">
									<image v-for="(img, idx) in log.images" :key="idx" :src="img.url" mode="aspectFill" class="thumb"></image>
								</view>
							</view>
						</view>
					</view>
				</view>
			</view>
			<view class="safe-area-bottom"></view>
		</scroll-view>
	</view>
</template>

<script>
import { callContainer } from '../utils/request';
import navBar from './navBar.vue';

export default {
	name: 'logCalendar',
	components: { navBar },
	props: {
		showNav: {
			type: Boolean,
			default: false
		}
	},
	data() {
		return {
			statusBarHeight: 0,
			topBarHeight: 0,
			selectedDate: '',
			allLogs: [],
			careActions: [],
			familyId: 0,
			selectedDates: [] // uni-calendar markers
		};
	},
	computed: {
		currentLogs() {
			if (!this.selectedDate) return [];
			return this.allLogs.filter(log => {
				const logDate = log.logTime.split('T')[0];
				return logDate === this.selectedDate;
			});
		}
	},
	mounted() {
		this.initData();
	},
	methods: {
		async initData() {
			const systemInfo = uni.getSystemInfoSync();
			this.statusBarHeight = systemInfo.statusBarHeight || 44;
			const app = getApp();
			if (app && app.globalData) {
				this.topBarHeight = app.globalData.topBarHeight;
			}
			this.familyId = uni.getStorageSync('familyId');
			const now = new Date();
			this.selectedDate = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
			
			await Promise.all([
				this.getCareActions(),
				this.getLogs()
			]);
		},
		async getCareActions() {
			try {
				const res = await callContainer("/api/care/", { familyId: Number(this.familyId) });
				if (res.data) {
					this.careActions = res.data;
				}
			} catch (e) { console.error(e); }
		},
		async getLogs() {
			try {
				const params = {
					familyId: Number(this.familyId)
				};
				
				const res = await callContainer("/api/plant/log/list", params);
				if (res.data) {
					this.allLogs = res.data.map(log => {
						const action = this.careActions.find(a => a.type === log.actionType);
						return {
							...log,
							actionName: action ? action.name : log.actionType,
							icon: action ? action.icon : '',
							color: action ? action.color : ''
						};
					});
					
					// 生成日历标记
					const dateMap = {};
					this.allLogs.forEach(log => {
						const date = log.logTime.split('T')[0];
						if (!dateMap[date]) {
							dateMap[date] = {
								date: date,
								info: ''
							};
						}
					});
					this.selectedDates = Object.values(dateMap);
					
				}
			} catch (e) { console.error(e); }
		},
		onDateChange(e) {
			this.selectedDate = e.fulldate;
		},
		formatTime(timeStr) {
			if (!timeStr) return '';
			const date = new Date(timeStr);
			return `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`;
		},
		goEdit(log) {
			uni.navigateTo({
				url: `/pages/logEdit/logEdit?id=${log.ID}&plantId=${log.plantId}`
			});
		}
	}
};
</script>

<style lang="scss" scoped>
.calendar-wrapper {
	width: 100%;
	height: 100vh;
	display: flex;
	flex-direction: column;
	padding-bottom: calc(100rpx + env(safe-area-inset-bottom));
	box-sizing: border-box;
}

// 固定头部
.fixed-header-group {
    position: sticky;
    top: 0;
    z-index: 998;
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
}

.header-title-container {
    padding: 16px 20px 12px;
    text-align: center;
}

.page-title {
    font-size: 18px;
    font-weight: 600;
    color: #333;
}

.scroll-body {
	flex: 1;
	height: 0;
}

.log-calendar {
	margin: 12px;
	background-color: rgba(255,255,255,0.55);
	border-radius: 16px;
	overflow: hidden;
}

.safe-area-bottom {
	height: env(safe-area-inset-bottom);
	padding-bottom: 20px;
}

.day-detail {
	padding: 16px;
}

.detail-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 15px;
}

.detail-title {
	font-size: 16px;
	font-weight: bold;
	color: #333;
}

.add-btn {
	display: flex;
	align-items: center;
	gap: 4px;
	background-color: #fff;
	padding: 4px 10px;
	border-radius: 20px;
	border: 1px solid #eee;
	text {
		font-size: 13px;
		color: #4A6139;
	}
}

.log-list {
	.empty-box {
		padding: 30px 0;
		text-align: center;
		.empty-text {
			font-size: 14px;
			color: #999;
		}
	}
}

.log-item {
	display: flex;
	padding: 12px;
	background-color: rgba(255,255,255,0.6);
	border-radius: 12px;
	margin-bottom: 10px;
	box-shadow: 0 2px 6px rgba(0,0,0,0.02);
	
	.log-icon {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-right: 12px;
		flex-shrink: 0;
	}
	
	.log-info {
		flex: 1;
		.log-top {
			display: flex;
			justify-content: space-between;
			align-items: center;
			margin-bottom: 4px;
			.log-name {
				font-size: 15px;
				font-weight: 500;
				color: #333;
			}
			.log-time {
				font-size: 12px;
				color: #999;
			}
		}
		
		.log-content {
			font-size: 13px;
			color: #666;
			line-height: 1.4;
			margin-bottom: 8px;
			display: block;
		}
		
		.log-images {
			display: flex;
			flex-wrap: wrap;
			gap: 6px;
			.thumb {
				width: 60px;
				height: 60px;
				border-radius: 4px;
			}
		}
	}
}

// 隐藏日历文字，确保显示标记点
::v-deep .uni-calendar-item__weeks-box-item {
	.uni-calendar-item__weeks-box-info {
		display: none !important;
	}
	// 标记点保持红色
	.uni-calendar-item__weeks-box-circle {
		display: block !important;
		background-color: #ff4d4f !important;
	}
}
</style>
