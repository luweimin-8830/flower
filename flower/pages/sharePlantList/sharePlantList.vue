<template>
	<view class="home-container theme-page" :style="themeStyle">
		<!-- 顶部固定区域 -->
		<view class="fixed-header-group">
			<navBar :isHome="true" />
			<view :style="{ height: topBarHeight + 'px' }"></view>

			<!-- 搜索框 (保持与 home.vue 一致) -->
			<view class="search-container">
				<view class="search-box-wrapper">
					<uni-search-bar @confirm="searchPlant" placeholder="搜索分享列表中的植物" radius="20" :focus="false"
						v-model="searchValue" bgColor="rgba(255,255,255,0.5)" clearButton="auto"
						cancelButton="none">
					</uni-search-bar>
				</view>
			</view>
		</view>

		<!-- 滚动列表 -->
		<scroll-view scroll-y class="content-scroll-view" @scrolltolower="onScrollToLower">
			<!-- 空状态 -->
			<view v-if="plantsList.length === 0 && !isLoading" class="empty-wrapper">
				<image src="/static/icon/c2m.svg" class="empty-icon" mode="aspectFit"></image>
				<text class="empty-text">这里还没有植物记录哦</text>
			</view>

			<!-- 瀑布流列表 -->
			<view v-else class="waterfall-wrapper">
				<WaterfallBox :list="plantsList" idKey="id" cols="2">
					<template #item="{ item }">
						<view class="plant-card" @click="gotoDetail(item)">
							<view class="image-wrapper"
								:style="{ paddingBottom: getCoverRatio(item.cover) + '%' }">
								<image :src="getCoverUrl(item.cover)" mode="aspectFill" class="plant-image show" lazy-load
									@load="onImgLoad(item)" @error="onImgError(item)"></image>
							</view>
							<view class="plant-info">
								<text class="plant-name">{{ item.name }}</text>
							</view>
						</view>
					</template>
				</WaterfallBox>
				<view class="share-list-bottom-space"></view>
			</view>
		</scroll-view>
	</view>
</template>

<script>
import { callContainer } from '../../utils/request';
import WaterfallBox from '../../components/WaterfallBox.vue';
import navBar from '../../components/navBar.vue';

export default {
	components: {
		WaterfallBox,
		navBar
	},
	data() {
		return {
			topBarHeight: 0,
			familyId: 0,
			tagId: 0,
			tagName: '',
			searchValue: '',
			searchTimer: null,
			plantsList: [],
			loadedImagesMap: {},
			page: 1,
			pageSize: 20,
			total: 0,
			isLoading: false,
			isNoMore: false
		}
	},
	watch: {
		searchValue() {
			if (this.searchTimer) clearTimeout(this.searchTimer);
			this.searchTimer = setTimeout(() => {
				this.getPlantsList();
			}, 500);
		}
	},
	async onLoad(options) {
		const app = getApp();
		this.topBarHeight = app.globalData.topBarHeight;
		
		if (options.familyId) {
			this.familyId = Number(options.familyId);
			this.tagId = options.tagId ? Number(options.tagId) : 0;
			this.tagName = options.tagName || '';
			
			await this.getPlantsList();
		}
	},
	methods: {
		async getPlantsList(isLoadMore = false) {
			if (this.isLoading || (isLoadMore && this.isNoMore)) return;

			if (!isLoadMore) {
				this.page = 1;
				this.isNoMore = false;
			}

			this.isLoading = true;
			try {
				const res = await callContainer("/api/plant/list", {
					familyId: this.familyId,
					page: this.page,
					pageSize: this.pageSize,
					tagId: this.tagId,
					keyword: this.searchValue
				});

				let rawData = [];
				if (res?.data?.list) {
					rawData = res.data.list;
					this.total = res.data.total;
				} else if (Array.isArray(res?.data)) {
					rawData = res.data;
					this.total = rawData.length;
				}

				if (isLoadMore) {
					this.plantsList = [...this.plantsList, ...rawData];
				} else {
					this.plantsList = rawData;
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
		onImgLoad(item) {
			if (item.id) this.loadedImagesMap[item.id] = true;
		},
		onImgError(item) {
			if (item.id) this.loadedImagesMap[item.id] = false;
		},
		getCoverUrl(cover) {
			return cover && cover.url ? cover.url : '/static/default.svg';
		},
		getCoverRatio(cover) {
			const width = Number(cover && cover.width);
			const height = Number(cover && cover.height);
			if (!width || !height) return 100;
			return Math.min(150, Math.max(75, height / width * 100));
		},
		searchPlant() {
			this.getPlantsList();
		},
		gotoDetail(item) {
			uni.navigateTo({
				url: `/pages/shareDetail/shareDetail?id=${item.id}&from=shareList`
			});
		}
	}
}
</script>

<style scoped lang="scss">
.home-container {
	height: 100vh;
	display: flex;
	flex-direction: column;
	overflow: hidden;
	background-color: var(--bg-color);
}

.fixed-header-group {
	flex-shrink: 0;
	z-index: 10;
	background-color: var(--bg-color);
	padding-bottom: 10px;
}

.search-container {
	padding: 0 16px 10px;
}

.search-box-wrapper {
	width: 100%;
}

::v-deep .uni-searchbar {
	padding: 0 !important;
}

.content-scroll-view {
	flex: 1;
	height: 0;
	overflow: hidden;
}

.waterfall-wrapper {
	padding: 10px;
}

.share-list-bottom-space { height: calc(30px + env(safe-area-inset-bottom)); }

.plant-card {
	background-color: rgba(255, 255, 255, 0.5);
	border-radius: 12px;
	overflow: hidden;
	margin-bottom: 10px;
	box-shadow: 0 4px 12px rgba(0,0,0,0.03);
}

.image-wrapper {
	position: relative;
	width: 100%;
	height: 0;
	background-color: rgba(255, 255, 255, 0.6);
}

.plant-image {
	position: absolute;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	opacity: 0;
	transition: opacity 0.4s ease-in-out;
}

.plant-image.show {
	opacity: 1;
}

.plant-info {
	padding: 10px;
	text-align: center;
}

.plant-name {
	font-size: 14px;
	color: #333;
	font-weight: bold;
}

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
}

.empty-text {
	font-size: 14px;
	color: #999;
}
</style>
