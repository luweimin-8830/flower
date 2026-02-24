<template>
	<view class="calendar-wrapper">
		<!-- 顶部固定区域 -->
		<view class="fixed-header-group">
			<navBar v-if="showNav" />
			<view
				:style="{ height: (weatherInfo.city ? (menuButtonTop + 70) : (menuButtonTop + menuButtonHeight + 10)) + 'px' }"
				style="position: relative; transition: height 0.3s;">
				<!-- 天气小组件（长条长方形布局） -->
				<view class="weather-widget-card" v-if="currentDayWeather" :style="{ top: menuButtonTop + 'px' }"
					@click="handleWeatherClick">
					<view class="widget-left">
						<image class="temp-icon" :src="getWeatherIconfont(currentDayWeather.weather)" mode="aspectFit"></image>
					</view>
					<view class="widget-center">
						<view class="city-info">
							<text class="city-name">{{ currentDayWeather.city }}</text>
							<uni-icons type="paperplane-filled" size="10" color="#fff"></uni-icons>
						</view>
						<view class="status-info">
							<image :src="getWeatherIconfont(currentDayWeather.weather)" mode="aspectFit"
								style="width: 14px; height: 14px; margin-right: 4px;"></image>
							<text class="weather-status">{{ currentDayWeather.weather || '多云' }}</text>
						</view>
					</view>
					<view class="widget-right">
						<text class="week-tag">{{ currentDayWeather.week }}</text>
						<text class="temp-range">{{ currentDayWeather.nighttemp }}°~{{ currentDayWeather.daytemp }}°</text>
					</view>
				</view>

				<view class="weather-empty-mini" v-else
					:style="{ top: menuButtonTop + 'px', height: menuButtonHeight + 'px' }" @click="handleWeatherClick">
					<uni-icons type="location-filled" size="14" color="#6B8857"></uni-icons>
					<text class="placeholder-btn-text">获取天气</text>
				</view>
			</view>
		</view>

		<scroll-view scroll-y class="scroll-body">
			<view class="log-calendar">
				<uni-calendar ref="calendar" :insert="true" :lunar="false" :selected="selectedDates" color="#566C44"
					@change="onDateChange" />

				<!-- 提醒时间选择弹窗 -->
				<uni-calendar ref="remindCalendar" :insert="false" :clear-date="true" :startDate="todayDate"
					@confirm="onRemindConfirm" />

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
								<picker mode="time" :value="remindData.time"
									@change="e => remindData.time = e.detail.value">
									<view class="time-picker-value">{{ remindData.time }}</view>
								</picker>
							</view>
							<view class="input-group vertical">
								<text class="label">事项内容</text>
								<textarea class="remind-textarea" v-model="remindData.content"
									placeholder="输入提醒内容,如:该浇水了" maxlength="20" />
							</view>
						</view>
						<view class="modal-footer">
							<button class="btn cancel" @click="$refs.remindPopup.close()">取消</button>
							<button class="btn confirm" @click="saveRemind">保存预约</button>
						</view>
					</view>
				</uni-popup>

				<view class="day-detail">
					<view class="detail-header">
						<text class="detail-title">{{ displayDate || selectedDate }} 记录</text>
						<uni-icons type="notification" size="20" color="#6B8857" @click="requestSubscribe"></uni-icons>
					</view>

					<view class="log-list">
						<view v-if="currentLogs.length === 0" class="empty-box">
							<text class="empty-text">当日无养护记录</text>
						</view>
						<view v-for="log in currentLogs" :key="log.id" class="log-item" @click="goEdit(log)">
							<view class="log-icon" :style="{ backgroundColor: (log.color || '#4A6139') + '22' }">
								<view class="iconfont" :class="log.icon || 'plant-jiaoshui1'" style="font-size: 20px;"
									:style="{ color: log.color || '#4A6139' }"></view>
							</view>
							<view class="log-info">
								<view class="log-top">
									<text class="log-name">{{ log.displayName }}</text>
									<text class="log-time">{{ formatTime(log.logTime) }}</text>
								</view>
								<text class="log-content" v-if="log.content">{{ log.content }}</text>
								<view class="log-images" v-if="log.images && log.images.length > 0">
									<image v-for="(img, idx) in log.images" :key="idx" :src="img.url" mode="aspectFill"
										class="thumb" @click.stop="previewImages(log.images, idx)"></image>
								</view>
							</view>
							<!-- 🌟 所有记录都可以删除 -->
							<view class="delete-batch-btn" @click.stop="handleDeleteBatch(log)">
								<uni-icons type="trash" size="18" color="#dd524d"></uni-icons>
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
			menuButtonTop: 0,
			menuButtonHeight: 32,
			selectedDate: '',
			displayDate: '', // 用于显示日志的日期（可能与日历选择不同）
			weatherInfo: {
				city: '',
				temp: '',
				icon: '',
				adcode: '',
				weather: '',
				casts: []
			},
			allLogs: [],
			careActions: [],
			familyId: 0,
			selectedDates: [], // uni-calendar markers
			todayDate: '',
			remindData: {
				date: '',
				time: '08:00',
				content: ''
			},
			allLogImages: [] // 所有日志图片列表
		};
	},
	computed: {
		currentLogs() {
			// 使用 displayDate 来决定显示哪天的日志
			const targetDate = this.displayDate || this.selectedDate;
			if (!targetDate) return [];
			
			const filtered = this.allLogs.filter(log => {
				const logDate = log.logTime.split('T')[0];
				return logDate === targetDate;
			});

			// 收集当日所有图片
			const todayImages = [];
			filtered.forEach(log => {
				if (log.images && log.images.length > 0) {
					log.images.forEach(img => {
						todayImages.push(img.url);
					});
				}
			});
			this.allLogImages = todayImages;

			// 按时间、动作类型、内容进行分组合并
			const groups = {};
			filtered.forEach(log => {
				// 精确到分钟的 key，用于判定是否是同一次批量操作
				const timeKey = log.logTime.substring(0, 16);
				const key = `${timeKey}_${log.actionType}_${log.content}`;

				if (!groups[key]) {
					groups[key] = {
						...log,
						count: 1,
						plants: [log.name],
						logIds: [log.id] // 收集同批次的所有日志 ID
					};
				} else {
					groups[key].count++;
					if (log.name) groups[key].plants.push(log.name);
					groups[key].logIds.push(log.id);
				}
			});

			return Object.values(groups).map(g => {
				return {
					...g,
					displayName: g.count > 1 ? `${g.actionName || g.name} (共 ${g.count} 棵)` : g.name
				};
			}).sort((a, b) => new Date(b.logTime) - new Date(a.logTime));
		},
		currentDayWeather() {
			if (!this.weatherInfo.city) return null;

			// 如果没有选中日期，默认显示今日
			const targetDate = this.selectedDate || this.todayDate;
			const weekMap = { '1': '周一', '2': '周二', '3': '周三', '4': '周四', '5': '周五', '6': '周六', '7': '周日' };

			// 在 casts 中查找匹配日期的预报
			const forecast = (this.weatherInfo.casts || []).find(c => c.date === targetDate);

			if (forecast) {
				return {
					city: this.weatherInfo.city,
					weather: forecast.dayweather,
					iconfont: this.getWeatherIconfont(forecast.dayweather),
					daytemp: forecast.daytemp,
					nighttemp: forecast.nighttemp,
					week: weekMap[forecast.week] || ''
				};
			}

			// 找不到或者是今天，显示今天的天气 (casts[0] 通常是今日)
			const todayForecast = (this.weatherInfo.casts && this.weatherInfo.casts.length > 0) ? this.weatherInfo.casts[0] : null;

			let weekStr = '';
			if (todayForecast) {
				weekStr = weekMap[todayForecast.week];
			} else {
				const dateObj = new Date(this.todayDate.replace(/-/g, '/'));
				const weeks = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'];
				weekStr = weeks[dateObj.getDay()];
			}

			return {
				city: this.weatherInfo.city,
				weather: this.weatherInfo.weather,
				iconfont: this.getWeatherIconfont(this.weatherInfo.weather),
				daytemp: todayForecast ? todayForecast.daytemp : this.weatherInfo.temp,
				nighttemp: todayForecast ? todayForecast.nighttemp : '',
				week: weekStr
			};
		}
	},
	mounted() {
		this.initData();
		uni.$on('refreshLogCalendar', () => {
			this.getLogs();
		});
	},
	destroyed() {
		uni.$off('refreshLogCalendar');
	},
	methods: {
		async initData() {
			const systemInfo = uni.getSystemInfoSync();
			this.statusBarHeight = systemInfo.statusBarHeight || 44;

			const menuButtonInfo = uni.getMenuButtonBoundingClientRect();
			this.menuButtonTop = menuButtonInfo.top;
			this.menuButtonHeight = menuButtonInfo.height;

			const app = getApp();
			if (app && app.globalData) {
				this.topBarHeight = app.globalData.topBarHeight || this.statusBarHeight;
			} else {
				this.topBarHeight = this.statusBarHeight;
			}
			this.familyId = uni.getStorageSync('familyId');
			const now = new Date();
			this.todayDate = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
			this.selectedDate = this.todayDate;
			this.displayDate = this.todayDate; // 初始化显示日期为今天

			// 从用户信息中获取默认提醒时间
			const userInfo = uni.getStorageSync('userInfo');
			if (userInfo && userInfo.remindTime) {
				this.remindData.time = userInfo.remindTime;
			}

			// 优先从缓存加载天气数据
			const cache = uni.getStorageSync('weatherCache');
			if (cache && cache.date === this.todayDate) {
				this.weatherInfo = cache;
			}

			// 必须先获取养护配置，因为日志列表依赖它来匹配图标
			await this.getCareActions();
			await this.getLogs();
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
						// 1. 尝试从养护配置中匹配
						let action = this.careActions.find(a => a.type === log.actionType);

						// 2. 特殊处理成长记录 (record)
						const isRecord = log.actionType === 'record' || log.actionType === 'Growth';

						// 3. 构建最终显示的日志对象
						const finalLog = {
							...log,
							name: log.name || (action ? action.name : (isRecord ? '成长记录' : log.actionType)),
							// 如果是成长记录，强制使用指定图标，否则使用 log.icon 或 action.icon
							icon: isRecord ? 'plant-zhiwuzhiyuan-duorouzhiwuyuan' : (log.icon || (action ? action.icon : 'plant-jiaoshui1')),
							color: isRecord ? '#A7C190' : (log.color || (action ? action.color : '#4A6139'))
						};

						return finalLog;
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
			this.displayDate = e.fulldate; // 同步更新显示日期
		},
		formatTime(timeStr) {
			if (!timeStr) return '';
			const date = new Date(timeStr);
			return `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`;
		},
		goEdit(log) {
			// 如果是批量操作记录，不允许编辑
			if (log.count > 1) {
				uni.showToast({ title: '批量记录无法编辑', icon: 'none' });
				return;
			}
			uni.navigateTo({
				url: `/pages/logEdit/logEdit?id=${log.id}&plantId=${log.plantId}`
			});
		},
		async handleDeleteBatch(log) {
			const title = log.count > 1 
				? `确定要删除这 ${log.count} 条批量操作记录吗？\n涉及植物：${log.plants.join('、')}`
				: `确定要删除"${log.displayName}"的这条记录吗？`;
			
			uni.showModal({
				title: '提示',
				content: title,
				confirmColor: '#dd524d',
				success: async (res) => {
					if (res.confirm) {
						uni.showLoading({ title: '正在删除...' });
						try {
							// 删除所有相关日志（单条时 logIds 只有一个元素）
							await Promise.all(
								log.logIds.map(id => callContainer("/api/plant/log/delete", { id }))
							);
							
							uni.showToast({ title: '已删除', icon: 'success' });
							
							// 刷新日志列表
							await this.getLogs();
							
							// 通知首页和详情页刷新
							uni.$emit('refreshHomeList');
							uni.$emit('refreshPlantDetail');
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
		previewImages(images, index) {
			// 找到当前点击图片在所有图片列表中的位置
			const currentUrl = images[index].url;
			const currentIndex = this.allLogImages.indexOf(currentUrl);
			
			uni.previewImage({
				current: currentIndex >= 0 ? currentIndex : 0,
				urls: this.allLogImages,
				success: () => {
					// 预览完成后，将日志显示切换回今天
					this.displayDate = this.todayDate;
				}
			});
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

				await callContainer('/api/remind/add', {
					familyId: Number(this.familyId),
					plantId: 0,
					remindTime: fullRemindTime,
					content: this.remindData.content,
					actionType: 'remind'
				});

				// 同时记录一条日志（可选，或者由后端自动生成）暂时不考虑
				// await callContainer('/api/plant/log/add', {
				// 	familyId: Number(this.familyId),
				// 	plantId: 0,
				// 	actionType: 'remind',
				// 	content: `【预约提醒】${this.remindData.content}`,
				// 	logTime: `${this.remindData.date}T${this.remindData.time}:00Z`,
				// 	images: []
				// });

				uni.showToast({
					title: '预约成功',
					icon: 'success'
				});
				this.$refs.remindPopup.close();
				this.getLogs(); // 刷新列表
			} catch (e) {
				console.error(e);
				uni.showToast({ title: '预约失败', icon: 'none' });
			} finally {
				uni.hideLoading();
			}
		},
		handleWeatherClick() {
			this.fetchLocationWeather();
		},
		/**
		 * 映射天气到 iconfont 图标名称
		 * 预留图标名称空位，待 iconfont 库更新后填入
		 */
		getWeatherIconfont(weather) {
			if (!weather) return '../static/weather/sun.svg';
			const map = {
				'雷': '../static/weather/lei.svg',
				'大雪': '../static/weather/daxue.svg',
				'中雪': '../static/weather/zhongxue.svg',
				'小雪': '../static/weather/xiaoxue.svg',
				'大雨': '../static/weather/dayu.svg',
				'中雨': '../static/weather/zhongyu.svg',
				'小雨': '../static/weather/xiaoyu.svg',
				'阴': '../static/weather/duoyun.svg',
				'多云': '../static/weather/duoyun.svg',
				'晴': '../static/weather/sun.svg',
				'雨': '../static/weather/xiaoyu.svg',
				'雪': '../static/weather/xiaoxue.svg',
			};
			for (let key in map) {
				if (weather.includes(key)) return map[key];
			}
			return '../static/weather/sun.svg';
		},
		// 获取模糊位置并请求后端换取天气信息
		async fetchLocationWeather() {
			// #ifdef MP-WEIXIN
			try {
				uni.showLoading({ title: '定位中...' });
				const locRes = await new Promise((resolve, reject) => {
					wx.getFuzzyLocation({
						type: 'gcj02',
						success: resolve,
						fail: reject
					});
				});
				console.log('模糊位置获取成功:', locRes);

				// 检查本地缓存
				const cache = uni.getStorageSync('weatherCache');
				if (cache && cache.date === this.todayDate) {
					const lonDiff = Math.abs(locRes.longitude - cache.longitude);
					const latDiff = Math.abs(locRes.latitude - cache.latitude);
					// 如果位置变动在 3 位小数以内，认为位置未变，直接用缓存
					if (lonDiff < 0.001 && latDiff < 0.001) {
						console.log('位置未变，使用本地天气缓存');
						this.weatherInfo = cache;
						uni.hideLoading();
						return;
					}
				}

				// 请求后端接口获取天气信息
				const res = await callContainer('/api/weather', {
					longitude: locRes.longitude,
					latitude: locRes.latitude
				});

				if (res && res.data) {
					this.weatherInfo = {
						city: res.data.city,
						temp: res.data.temp,
						icon: res.data.icon,
						adcode: res.data.adcode,
						weather: res.data.weather,
						casts: res.data.casts || []
					};
					// 存入本地缓存，包含位置信息
					uni.setStorageSync('weatherCache', {
						...this.weatherInfo,
						longitude: locRes.longitude,
						latitude: locRes.latitude,
						date: this.todayDate
					});
				} else {
					throw new Error('获取天气数据失败');
				}

				uni.hideLoading();
			} catch (error) {
				uni.hideLoading();
				console.error('获取位置或天气失败:', error);
				uni.showToast({ title: '获取失败，请重试', icon: 'none' });
			}
			// #endif
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
	background-color: var(--bg-color);
}

// 固定头部
.fixed-header-group {
	position: sticky;
	top: 0;
	z-index: 998;
	background-color: var(--bg-color);
}

.scroll-body {
	flex: 1;
	height: 0;
}

.log-calendar {
	margin: 12px;
	background-color: var(--bg-btn-color);
	border-radius: 16px;
	overflow: hidden;
	position: relative;
}

// 匹配长条形的天气卡片样式
.weather-widget-card {
	position: absolute;
	left: 15px;
	width: 420rpx;
	height: 100rpx;
	background: linear-gradient(135deg, #A7C190 0%, #566C44 100%);
	border-radius: 20rpx;
	padding: 0 20rpx;
	display: flex;
	align-items: center;
	color: #fff;
	box-shadow: 0 6rpx 16rpx rgba(122, 153, 173, 0.2);
	z-index: 100;

	@media (prefers-color-scheme: dark) {
		background: linear-gradient(135deg, #6B8857 0%, #A7C190 100%);
	}

	.widget-left {
		margin-right: 16rpx;

		.temp-icon {
			width: 64rpx;
			height: 64rpx;
		}
	}

	.widget-center {
		flex: 1;
		display: flex;
		flex-direction: column;
		justify-content: center;
		gap: 4rpx;

		.city-info {
			display: flex;
			align-items: center;
			gap: 6rpx;

			.city-name {
				font-size: 26rpx;
				font-weight: bold;
			}
		}

		.status-info {
			display: flex;
			align-items: center;
			gap: 6rpx;

			.weather-status {
				font-size: 22rpx;
				opacity: 0.9;
			}
		}
	}

	.widget-right {
		display: flex;
		flex-direction: column;
		align-items: flex-end;
		justify-content: center;
		gap: 6rpx;

		.week-tag {
			font-size: 20rpx;
			opacity: 0.8;
			font-weight: 500;
		}

		.temp-range {
			font-size: 22rpx;
			font-weight: 500;
		}
	}
}

.weather-empty-mini {
	position: absolute;
	left: 15px;
	z-index: 10;
	padding: 0 16rpx;
	background-color: var(--bg-btn-color);
	border-radius: 30rpx;
	display: flex;
	align-items: center;
	gap: 6rpx;
	border: 1px solid var(--border-color);
	box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);

	.placeholder-btn-text {
		font-size: 22rpx;
		color: var(--text-color);
		font-weight: 500;
	}
}

.weather-placeholder {
	display: none;
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
	color: var(--text-color);
}

.add-btn {
	display: flex;
	align-items: center;
	gap: 4px;
	background-color: var(--bg-btn-color);
	padding: 4px 10px;
	border-radius: 20px;
	border: 1px solid var(--border-color);

	text {
		font-size: 13px;
		color: var(--text-color);
	}
}

.log-list {
	.empty-box {
		padding: 30px 0;
		text-align: center;

		.empty-text {
			font-size: 14px;
			color: var(--text-sub);
		}
	}
}

.log-item {
	display: flex;
	padding: 12px;
	background-color: var(--bg-btn-color);
	border-radius: 12px;
	margin-bottom: 10px;
	box-shadow: 0 2px 6px rgba(0, 0, 0, 0.02);

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
				color: var(--text-color);
			}

			.log-time {
				font-size: 12px;
				color: var(--text-sub);
			}
		}

		.log-content {
			font-size: 13px;
			color: var(--text-color);
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

	.delete-batch-btn {
		padding: 0 12rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
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

	// 暗黑模式日历适配
	.uni-calendar-item__weeks-box-text {
		color: var(--text-color) !important;
	}
}

.remind-modal {
	width: 600rpx;
	background-color: var(--card-bg);
	border-radius: 20rpx;
	padding: 30rpx;

	.modal-header {
		text-align: center;
		margin-bottom: 30rpx;

		.modal-title {
			font-size: 18px;
			font-weight: bold;
			color: var(--text-color);
		}
	}

	.input-group {
		display: flex;
		align-items: center;
		padding: 20rpx 0;
		border-bottom: 1px solid var(--border-color);

		&.vertical {
			flex-direction: column;
			align-items: flex-start;
			border-bottom: none;
		}

		.label {
			width: 140rpx;
			font-size: 15px;
			color: var(--text-sub);
		}

		.value,
		.time-picker-value {
			font-size: 15px;
			color: var(--text-color);
			font-weight: 500;
		}

		.remind-textarea {
			width: 90%;
			height: 160rpx;
			background-color: var(--bg-btn-color);
			border-radius: 10rpx;
			padding: 20rpx;
			margin-top: 16rpx;
			font-size: 14px;
			color: var(--text-color);
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

			&::after {
				border: none;
			}

			&.cancel {
				background-color: var(--border-color);
				color: var(--text-sub);
			}

			&.confirm {
				background-color: #4A6139;
				color: #fff;
			}
		}
	}
}
</style>
