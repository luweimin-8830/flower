<template>
	<view class="container">
		<navBar />

		<!-- 占位符，防止内容被导航栏遮挡 -->
		<view :style="{ height: topBarHeight + 'px' }"></view>

		<scroll-view scroll-y class="main-scroll" :enable-back-to-top="true">

			<!-- 2. 植物基本信息卡片 -->
			<view class="plant-header-card">
				<view class="header-top">
					<image v-if="plant.cover" :src="plant.cover.url" mode="aspectFill" class="plant-avatar"></image>
					<text class="plant-name">{{ plant.name }}</text>
				</view>

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
						<text class="stat-val">{{ 0 }}</text>
						<text class="stat-label">植物日志</text>
					</view>
				</view>
			</view>

			<!-- 3. 日常护理 (横向滚动) -->
			<view class="section-container">
				<view class="section-header">
					<text class="section-title">日常护理</text>
					<uni-icons type="notification" size="20" color="#6B8857"></uni-icons>
				</view>

				<scroll-view scroll-x class="care-scroll-view" :show-scrollbar="false">
					<view class="care-list">
						<view class="care-item" v-for="(item, index) in careActions" :key="index"
							@click="handleCare(item)">
							<view class="care-icon-box" :style="{ backgroundColor: item.color }">
								<!-- 这里用 uni-icons 模拟，实际可用 image -->
								<uni-icons :type="item.icon" size="28" color="#fff" v-if="item.icon"></uni-icons>
								<image v-else :src="item.img" class="care-img"></image>
							</view>
							<text class="care-name">{{ item.name }}</text>
						</view>
					</view>
				</scroll-view>
			</view>

			<!-- 4. 快捷入口 (日历/相册) -->
			<view class="quick-links">
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
						<text class="link-text">植物相册 ({{ plant.photoCount }})</text>
					</view>
					<uni-icons type="right" size="14" color="#999"></uni-icons>
				</view>
			</view>

			<!-- 5. 植物日志 (时间轴) -->
			<view class="log-section">
				<view class="section-header">
					<view class="header-left">
						<uni-icons type="compose" size="20" color="#333"></uni-icons>
						<text class="section-title" style="margin-left: 6px;">植物日志</text>
					</view>
					<view class="header-right">
						<text class="edit-btn">编辑</text>
						<uni-icons type="plusempty" size="24" color="#333" style="margin-left: 15px;"></uni-icons>
					</view>
				</view>

				<!-- 日志列表 -->
				<view class="log-list">
					<view v-for="(group, gIndex) in logList" :key="gIndex" class="log-group">
						<text class="log-date-header">{{ group.dateStr }}</text>

						<view v-for="(log, lIndex) in group.items" :key="lIndex" class="log-item">
							<!-- 左侧时间 -->
							<view class="log-time-col">
								<text class="log-time">{{ log.time }}</text>
								<text class="log-date-mini">{{ log.dateMini }}</text>
							</view>

							<!-- 中间轴线 -->
							<view class="log-timeline">
								<view class="dot"></view>
								<view class="line" v-if="lIndex !== group.items.length - 1"></view>
							</view>

							<!-- 右侧内容气泡 -->
							<view class="log-content-box">
								<view class="log-tag pill">
									<uni-icons type="checkbox-filled" size="16" color="#4A90E2"
										v-if="log.type === 'Watering'"></uni-icons>
									<text class="log-text">{{ log.actionName }}</text>
								</view>
							</view>
						</view>
					</view>
				</view>
			</view>

			<!-- 底部留白 -->
			<view style="height: 100px;"></view>
		</scroll-view>
	</view>
</template>

<script>
import navBar from '@/components/navBar.vue'
import { callContainer } from '../../utils/request';
export default {
	components: {
		navBar,
	},
	data() {
		return {
			topBarHeight: 0, // 默认值，created 中会更新
			plantId: 0,
			plant: {},
			careActions: [
				{ name: "Fertilizing", icon: "", img: "", color: "#DCECC9" }, // 施肥
				{ name: "Pruning", icon: "scissors", color: "#F2D7D5" }, // 修剪
				{ name: "Soil Change", icon: "", img: "", color: "#E8E0D5" }, // 换土
				{ name: "Watering", icon: "checkbox-filled", color: "#D6EAF8" }, // 浇水
				{ name: "修剪", icon: "scissors", color: "#F2D7D5" },
			],
			logList: [
				{
					dateStr: "2026年1月15日 星期四",
					items: [
						{ time: "19:07", dateMini: "01/15", type: "Watering", actionName: "Watering" }
					]
				},
				{
					dateStr: "2026年1月14日 星期三",
					items: [
						{ time: "10:00", dateMini: "01/14", type: "Other", actionName: "记录生长" }
					]
				}
			]
		};
	},
	created() {

	},
	onLoad(options) {
		const app = getApp()
		this.topBarHeight = app.globalData.topBarHeight;
		console.log("传递参数为:", options)
		if (options) {
			this.plantId = options.id
			this.getPlant()
		}
	},
	methods: {
		async getPlant() {
			try {
				const result = await callContainer("/api/plant/",{
					id:Number(this.plantId)
				})
				console.log("call container plant:",result)
				this.plant = result.data
				this.plant.birthday = this.plant.birthday.split('T')[0]
			} catch (error) {
				console.error(error)
			}
		},
		handleCare(item) {
			console.log("点击护理:", item.name);
			uni.showToast({ title: item.name, icon: 'none' });
		},
		goCalendar() {
			console.log("去日历");
		},
		goAlbum() {
			console.log("去相册");
		}
	}
};
</script>

<style lang="scss" scoped>
/* 全局背景色 */
.container {
	min-height: 100vh;
}

/* 2. 植物头部卡片 */
.plant-header-card {
	margin: 10px 16px;
	padding: 20px;
	background-color: rgba(255,255,255,0.55);
	/* 这里的背景色可能需要稍微透一点绿 */
	border-radius: 24px;
	box-shadow: 0 4px 20px rgba(107, 136, 87, 0.05);

	.header-top {
		display: flex;
		align-items: center;
		margin-bottom: 24px;

		.plant-avatar {
			width: 80px;
			height: 80px;
			border-radius: 16px;
			margin-right: 16px;
			background-color: #eee;
		}

		.plant-name {
			font-size: 22px;
			font-weight: bold;
			color: #2c3e50;
		}
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
				color: #333;
				margin-bottom: 4px;
			}

			.stat-label {
				font-size: 12px;
				color: #888;
			}
		}

		.vertical-line {
			width: 1px;
			height: 20px;
			background-color: #eee;
		}
	}
}

/* 3. 日常护理 */
.section-container {
	margin: 24px 0;
	padding-left: 16px;
	/* 标题左对齐 */
}

.section-header {
	display: flex;
	align-items: center;
	margin-bottom: 16px;
	padding-right: 16px;

	.section-title {
		font-size: 18px;
		font-weight: bold;
		color: #333;
		margin-right: auto;
	}
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

	.care-icon-box {
		width: 56px;
		height: 56px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 8px;
		box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);

		/* 默认背景色，会被内联样式覆盖 */
		background-color: #eee;
	}

	.care-name {
		font-size: 12px;
		color: #555;
		text-align: center;
		white-space: normal;
		/* 允许换行 */
		line-height: 1.2;
	}
}

/* 4. 快捷入口 */
.quick-links {
	display: flex;
	justify-content: space-between;
	margin: 0 16px 24px 16px;
	gap: 12px;

	.link-card {
		flex: 1;
		background-color: rgba(255, 255, 255, 0.6);
		border: 1px solid rgba(255, 255, 255, 1);
		border-radius: 16px;
		padding: 16px 12px;
		display: flex;
		align-items: center;
		justify-content: space-between;

		.link-left {
			display: flex;
			align-items: center;
			gap: 8px;

			.link-text {
				font-size: 14px;
				font-weight: 500;
				color: #333;
			}
		}
	}
}

/* 5. 植物日志 */
.log-section {
	padding: 0 16px;

	.header-left {
		display: flex;
		align-items: center;
	}

	.header-right {
		display: flex;
		align-items: center;
		margin-left: auto;

		.edit-btn {
			font-size: 14px;
			color: #666;
		}
	}
}

.log-list {
	margin-top: 10px;
}

.log-group {
	margin-bottom: 20px;

	.log-date-header {
		font-size: 16px;
		font-weight: bold;
		color: #333;
		margin-bottom: 16px;
		display: block;
	}
}

.log-item {
	display: flex;
	position: relative;
	padding-bottom: 20px;

	.log-time-col {
		width: 50px;
		text-align: right;
		padding-right: 10px;
		display: flex;
		flex-direction: column;

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
	}

	.log-timeline {
		width: 20px;
		display: flex;
		flex-direction: column;
		align-items: center;
		position: relative;

		.dot {
			width: 10px;
			height: 10px;
			background-color: #6B8857;
			border-radius: 50%;
			z-index: 1;
			margin-top: 6px;
			/* 对齐时间 */
		}

		.line {
			flex: 1;
			width: 1px;
			background-color: #ddd;
			margin-top: 4px;
		}
	}

	.log-content-box {
		flex: 1;
		padding-left: 10px;
		padding-top: 2px;

		.log-tag {
			display: inline-flex;
			align-items: center;
			background-color: #D6EAF8;
			/* 浅蓝底 */
			padding: 6px 12px;
			border-radius: 20px;
			gap: 6px;

			.log-text {
				font-size: 14px;
				color: #333;
			}
		}
	}
}
</style>