<template>
	<view class="container">
		<navBar title="新手手册" />
		<view :style="{ height: topBarHeight + 'px' }"></view>

		<scroll-view scroll-y class="content-scroll">
			<view class="manual-content">
				<!-- 1. 欢迎介绍 -->
				<view class="manual-section">
					<view class="section-header">
						<uni-icons type="heart-filled" size="24" color="#566C44"></uni-icons>
						<text class="section-title">欢迎来到我的花园</text>
					</view>
					<view class="section-body">
						<text class="p">这是一个为您和您的家人设计的植物养护日志应用。在这里，您可以记录每一株植物的成长点滴，与家人共同分享绿色的喜悦。</text>
					</view>
				</view>

				<!-- 2. 植物管理 -->
				<view class="manual-section">
					<view class="section-header">
						<view class="iconfont plant-zhiwuzhiyuan-duorouzhiwuyuan" style="font-size: 24px; color: #566C44; margin-right: 8px;"></view>
						<text class="section-title">植物管理</text>
					</view>
					<view class="section-body">
						<view class="step-item">
							<text class="step-num">1</text>
							<text class="step-text">点击首页右上角的“+”按钮，可以添加新植物。您可以上传照片、设置植物名称、到家日期，并为它打上标签。</text>
						</view>
						<view class="step-item">
							<text class="step-num">2</text>
							<text class="step-text">在植物详情页，点击右上角的“编辑”可以修改植物信息，点击“凋零”可以将已枯萎的植物移入陈列室。</text>
						</view>
					</view>
				</view>

				<!-- 3. 日常养护 -->
				<view class="manual-section">
					<view class="section-header">
						<view class="iconfont plant-kongzhonghuayuan" style="font-size: 24px; color: #566C44; margin-right: 8px;"></view>
						<text class="section-title">日常养护</text>
					</view>
					<view class="section-body">
						<view class="step-item">
							<text class="step-num">1</text>
							<text class="step-text">在植物详情页，您可以通过“快捷记录”栏快速记录今天的浇水、施肥等操作。</text>
						</view>
						<view class="step-item">
							<text class="step-num">2</text>
							<text class="step-text">如果您想记录更多细节（如文字描述、多张照片），可以点击“植物日志”旁的“+”号进行详细记录。</text>
						</view>
						<view class="step-item">
							<text class="step-num">3</text>
							<text class="step-text">在首页长按任意植物卡片，可以进入批量操作模式，一次性为多株植物记录相同的养护操作（如集体浇水）。</text>
						</view>
					</view>
				</view>

				<!-- 4. 家庭共享 -->
				<view class="manual-section">
					<view class="section-header">
						<uni-icons type="staff" size="24" color="#566C44"></uni-icons>
						<text class="section-title">家庭共享</text>
					</view>
					<view class="section-body">
						<text class="p">进入“我的 → 家庭管理”，您可以生成邀请链接发送给家人。家人加入后，大家可以共同查看和记录同一个花园里的植物状态。</text>
					</view>
				</view>

				<!-- 5. 数据分析 -->
				<view class="manual-section">
					<view class="section-header">
						<view class="iconfont plant-tubiao-bingtu" style="font-size: 24px; color: #566C44; margin-right: 8px;"></view>
						<text class="section-title">图表与分析</text>
					</view>
					<view class="section-body">
						<text class="p">通过“我的 → 图表分析”，您可以直观地看到花园里植物的分类比例，以及每个月的养护操作统计，帮助您更好地掌握养护节奏。</text>
					</view>
				</view>

				<!-- 6. 提醒设置 -->
				<view class="manual-section">
					<view class="section-header">
						<uni-icons type="notification" size="24" color="#566C44"></uni-icons>
						<text class="section-title">养护提醒</text>
					</view>
					<view class="section-body">
						<text class="p">在“我的 → 提醒时间”中，您可以设置一个固定的提醒时间。当您在植物详情页预约了提醒后，系统会在指定时间通过订阅消息通知您。</text>
					</view>
				</view>

				<view class="manual-footer">
					<text>祝您的植物茁壮成长！🌿</text>
				</view>
				<view style="height: 50px;"></view>
			</view>
		</scroll-view>
	</view>
</template>

<script>
import navBar from '@/components/navBar.vue'

export default {
	components: { navBar },
	data() {
		return {
			topBarHeight: 0
		};
	},
	onLoad() {
		const app = getApp();
		this.topBarHeight = app.globalData.topBarHeight;
	}
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
}

.content-scroll {
	height: calc(100vh - 100rpx);
}

.manual-content {
	padding: 20px 16px;
}

.manual-section {
	background-color: rgba(255,255,255,0.7);
	border-radius: 16px;
	padding: 20px;
	margin-bottom: 20px;
	box-shadow: 0 4px 12px rgba(86, 108, 68, 0.05);
}

.section-header {
	display: flex;
	align-items: center;
	margin-bottom: 16px;
	border-bottom: 1px solid #f0f0f0;
	padding-bottom: 12px;
}

.section-title {
	font-size: 18px;
	font-weight: bold;
	color: #2F3E25;
}

.section-body {
	.p {
		font-size: 15px;
		color: #666;
		line-height: 1.6;
		display: block;
	}
}

.step-item {
	display: flex;
	margin-bottom: 12px;
	
	&:last-child {
		margin-bottom: 0;
	}
}

.step-num {
	width: 20px;
	height: 20px;
	background-color: #566C44;
	color: #fff;
	font-size: 12px;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-right: 10px;
	flex-shrink: 0;
	margin-top: 2px;
}

.step-text {
	font-size: 14px;
	color: #666;
	line-height: 1.5;
	flex: 1;
}

.manual-footer {
	text-align: center;
	padding: 20px 0;
	text {
		font-size: 14px;
		color: #566C44;
		font-weight: 500;
		@media (prefers-color-scheme: dark) {
            color: #f5f5f5;
        }
	}
}
</style>