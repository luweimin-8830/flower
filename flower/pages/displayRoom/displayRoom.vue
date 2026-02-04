<template>
	<view class="container">
		<navBar />
		<view class="page-title" :style="{ top: menuButtonInfo.top + 'px' }">陈列室</view>

		<view :style="{ height: topBarHeight + 'px' }"></view>

		<scroll-view scroll-y class="scroll-area" @scrolltolower="loadMore">
			<view v-if="plants.length === 0 && !loading" class="empty-state">
				<uni-icons type="shop" size="64" color="#ccc"></uni-icons>
				<text class="empty-text">这里还没有凋零的植物</text>
			</view>

			<view class="plant-grid">
				<view v-for="item in plants" :key="item.id" class="plant-card" @click="goDetail(item.id)">
					<image class="cover" :src="item.cover && item.cover.url ? item.cover.url : '/static/default.svg'" mode="aspectFill"></image>
					<view class="info">
						<text class="name">{{ item.name }}</text>
						<view class="date-info">
							<uni-icons type="calendar" size="12" color="#999"></uni-icons>
							<text class="date">{{ formatDate(item.deathAt) }} 离去</text>
						</view>
					</view>
				</view>
			</view>

			<view v-if="loading" class="loading-more">
				<text>加载中...</text>
			</view>
			<view v-if="!hasMore && plants.length > 0" class="no-more">
				<text>已经到底啦</text>
			</view>
            <view style="height: 50rpx;"></view>
		</scroll-view>
	</view>
</template>

<script>
import navBar from '@/components/navBar.vue'
import { callContainer } from '@/utils/request'

export default {
	components: {
		navBar
	},
	data() {
		return {
			menuButtonInfo: {},
			topBarHeight: 0,
			plants: [],
			page: 1,
			pageSize: 20,
			total: 0,
			loading: false,
			hasMore: true,
			familyId: 0
		}
	},
	methods: {
		formatDate(dateStr) {
			if (!dateStr) return '';
			return dateStr.split('T')[0];
		},
		async getDeadPlants() {
			if (this.loading || !this.hasMore) return;
			this.loading = true;
			try {
				const res = await callContainer('/api/plant/dead/list', {
					familyId: this.familyId,
					page: this.page,
					pageSize: this.pageSize
				});
				if (res && res.data) {
					const list = res.data.list || [];
					this.plants = [...this.plants, ...list];
					this.total = res.data.total || 0;
					this.hasMore = this.plants.length < this.total;
					this.page++;
				}
			} catch (e) {
				console.error('获取陈列室数据失败:', e);
			} finally {
				this.loading = false;
			}
		},
		loadMore() {
			this.getDeadPlants();
		},
		goDetail(id) {
			uni.navigateTo({
				url: `/pages/plantDetail/plantDetail?id=${id}`
			});
		}
	},
	async onLoad() {
		this.menuButtonInfo = wx.getMenuButtonBoundingClientRect();
		const app = getApp();
		this.topBarHeight = app.globalData.topBarHeight;

		const familyID = await new Promise((resolve, reject) => {
			uni.getStorage({ key: 'familyId', success: resolve, fail: reject })
		});
		this.familyId = familyID.data;
		this.getDeadPlants();
	}
}
</script>

<style scoped lang="scss">
.container {
	min-height: 100vh;
	background-color: var(--bg-color);
}

.page-title {
	width: 100%;
	color: #333;
	position: fixed;
	font-size: 34rpx;
	font-weight: bold;
	text-align: center;
	z-index: 1000;
}

.scroll-area {
	height: calc(100vh - 100rpx);
}

.plant-grid {
	padding: 20rpx;
	display: flex;
	flex-wrap: wrap;
	gap: 20rpx;
}

.plant-card {
	width: calc(50% - 10rpx);
	background-color: #fff;
	border-radius: 20rpx;
	overflow: hidden;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);

	@media (prefers-color-scheme: dark) {
		background-color: #2c2c2c;
	}

	.cover {
		width: 100%;
		height: 300rpx;
		background-color: #f5f5f5;
        filter: grayscale(100%); // 陈列室里的植物显示为黑白，更有氛围感
        opacity: 0.8;
	}

	.info {
		padding: 16rpx;
		display: flex;
		flex-direction: column;
		gap: 8rpx;

		.name {
			font-size: 28rpx;
			font-weight: bold;
			color: #333;
			@media (prefers-color-scheme: dark) {
				color: #eee;
			}
		}

		.date-info {
			display: flex;
			align-items: center;
			gap: 6rpx;

			.date {
				font-size: 22rpx;
				color: #999;
			}
		}
	}
}

.empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding-top: 200rpx;
	gap: 20rpx;

	.empty-text {
		color: #999;
		font-size: 28rpx;
	}
}

.loading-more, .no-more {
	text-align: center;
	padding: 30rpx;
	font-size: 24rpx;
	color: #999;
}
</style>
