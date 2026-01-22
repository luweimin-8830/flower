<template>
	<view class="container">
		<loadingPage ref="loading" />
		<!-- 头部 -->
		<navBar />
		<view class="plant-title" :style="{ top: menuButtonInfo.top + 'px' }">添加植物</view>

		<view :style="{ height: topBarHeight + 20 + 'px' }"></view>
		<!-- 双按钮 -->
		<view :style="{ padding: '0 ' + paddingLeft + 'px' }">
			<uni-row :gutter="gutter" :width="nvueWidth">
				<uni-col :span="5">
					<view class="clean-btn" @click="clear">
						<text>取消</text>
					</view>
				</uni-col>
				<uni-col :span="5" :push="14">
					<view class="add-btn" @click="save">
						<text>保存</text>
					</view>
				</uni-col>
			</uni-row>
		</view>
		<!-- 图片上传 -->
		<view class="upload-section">
			<view class="avatar-wrapper">
				<!-- 图片：如果没有src显示背景色，有src显示图片  :src="'cloud://prod-0gr2o3qpe533f1fb.7072-prod-0gr2o3qpe533f1fb-1352691102/020-plant.png'"-->
				<image class="plant-img"
					:src="plant.cover || 'cloud://prod-0gr2o3qpe533f1fb.7072-prod-0gr2o3qpe533f1fb-1352691102/020-plant.png'"
					mode="aspectFill" />
				<!-- 相机图标：绝对定位到右下角 -->
				<view class="camera-badge" @click="uploadImages">
					<uni-icons type="camera-filled" size="18" color="#fff"></uni-icons>
				</view>
			</view>
		</view>
		<view class="form-container">

			<!-- 1. 名称 (胶囊输入框 + 扫码图标) -->
			<view class="form-item input-pill" :class="{ 'focused': nameFocus }">
				<uni-easyinput v-model="plant.name" :inputBorder="false" :clearable="false" placeholder="输入植物名称"
					placeholderStyle="color:#8BA989; font-size:15px;" @focus="nameFocus = true"
					@blur="nameFocus = false"></uni-easyinput>
				<!-- 右侧刷新图标 -->
				<view class="icon-right" @click="rollName">
					<uni-icons type="loop" size="24" color="#566C44"></uni-icons>
				</view>
			</view>

			<!-- 2. 描述 (带标题 + 圆角文本域) -->
			<view class="section-label">
				<uni-icons type="compose" size="18" color="#566C44"></uni-icons>
				<text class="label-text">描述</text>
			</view>

			<view class="form-item textarea-box" :class="{ 'focused': descFocus }">
				<uni-easyinput type="textarea" v-model="plant.desc" :disabled="isCalendarOpen" :inputBorder="false"
					placeholder="记录植物的状态..." placeholderStyle="color:#8BA989; font-size:14px;" autoHeight
					@focus="descFocus = true" @blur="descFocus = false"></uni-easyinput>
			</view>

			<!-- 3. 日期 (居中胶囊按钮) -->
			<view class="date-section">
				<uni-datetime-picker type="date" ref="datePicker" v-model="plant.birthday" :border="false"
					@change="onDateChange" @maskClick="dateClear">
					<!-- 自定义显示样式的插槽 -->
					<view class="date-pill" @click.stop="openCalendar">
						<uni-icons type="calendar-filled" size="18" color="#566C44"
							style="margin-left: 8px;"></uni-icons>
						<text>{{ plant.birthday || '选择日期' }}</text>
					</view>
				</uni-datetime-picker>
			</view>
		</view>
		<!-- 标签 -->
		<view class="section-label" style="margin-top: 40rpx;">
			<uni-icons type="tag" size="18" color="#566C44"></uni-icons>
			<text class="label-text">标签</text>
		</view>

		<view class="tags-wrapper" :style="{ padding: paddingLeft + 'px' }">
			<!-- 循环显示标签：点击切换选中状态 -->
			<view class="tag-item" :class="{ 'active': item.active }" v-for="(item, index) in tags" :key="index"
				@click="toggleTag(index)">
				<text>{{ item.name }}</text>
			</view>

			<!-- 添加按钮 -->
			<view class="tag-add-btn" @click="addTag">
				<uni-icons type="plusempty" size="16" color="#566C44"></uni-icons>
				<text>添加</text>
			</view>
		</view>
	</view>
</template>

<script>
import navBar from '@/components/navBar.vue'
import { callContainer } from '../../utils/request';
import loadingPage from '@/components/loading.vue';
export default {
	components: {
		navBar,
		loadingPage,
	},
	data() {
		return {
			menuButtonInfo: {},
			topBarHeight: 0,
			paddingLeft: 0,
			familyId: 0,
			plant: {
				name: '',
				cover: '',
				coverId: 0,
				desc: '',
				birthday: '',
				tags: [],
			},
			uploadImage: {
				url: "",
				height: "",
				width: 0,
				sha256: 0,
				size: 0,
				isUpload: false,
			},
			nameFocus: false,
			descFocus: false,
			tags: [],
			isCD: false,
			isCalendarOpen: false,
			isSave: true,
		}
	},
	methods: {
		getTimeCode() {
			// 1. 获取秒级时间戳
			const timestamp = Math.floor(new Date().getTime() / 1000);
			// 2. 转为36进制并大写
			return timestamp.toString(36).toUpperCase();
		},
		rollName() {
			if (this.isCD) return
			this.isCD = true
			wx.vibrateShort({ type: "medium" })
			const c = this.getTimeCode()
			this.plant.name = "未知-" + c
			setTimeout(() => { this.isCD = false }, 1000)
		},
		dateClear() {
			this.isCalendarOpen = false;
		},
		openCalendar() {
			this.isCalendarOpen = true;
			this.$refs.datePicker.show();
		},
		onDateChange(e) {
			console.log("date", e)
			this.plant.birthday = e
			this.isCalendarOpen = false;
		},
		async addTag() {
			const tag = await new Promise((resolve, reject) => {
				uni.showModal({
					title: '添加标签',
					editable: true,
					placeholderText: '例如：房间、喜阳、品种...',
					success: resolve,
					fail: reject
				})
			})
			console.log("tag", tag)
			if (tag.confirm) {
				const result = await callContainer("/api/tag/add", {
					name: tag.content,
					familyId: this.familyId
				})
				console.log("add tag", result)
				this.getTagsList()
			}

		},
		// 切换标签选中状态
		toggleTag(index) {
			this.tags[index].active = !this.tags[index].active;
		},
		async getTagsList() {
			try {
				const tagList = await callContainer("/api/tag/", {
					familyId: this.familyId
				})
				console.log("tagList:", tagList)
				const apiTags = tagList?.data || []
				this.tags = [
					...apiTags.map((item) => ({
						name: item.name,
						ID: item.ID,
						active: false
					}))
				]
				console.log("tags", this.tags)
			} catch (error) {
				console.error(error)
			}
		},
		clear() {
			wx.vibrateShort({ type: "medium" })
			wx.navigateBack()
		},
		async save() {
			if (!this.isSave) return;
			this.isSave = false;
			this.$refs.loading.open();

			wx.vibrateShort({ type: "medium" })
			if (!this.plant.name) {
				uni.showToast({
					title: '请输入植物名称',
					icon: 'error',
					mask: true
				})
				return;
			}
			if (this.uploadImage.isUpload) {
				const timestamp = Math.floor(new Date().getTime() / 1000);
				const random = Math.floor(Math.random() * 1000000).toString().padStart(6, '0');
				const name = `${timestamp}${random}.jpg`;

				const uploadFile = await new Promise((resolve, reject) => {
					wx.cloud.uploadFile({
						cloudPath: this.familyId + '/' + name,
						filePath: this.uploadImage.url,
						config: {
							env: 'prod-0gr2o3qpe533f1fb'
						},
						success: resolve,
						fail: reject,
					})
				})
				console.log("upload image", uploadFile)
				const fileID = uploadFile.fileID
				const addImage = await callContainer("/api/image/add", {
					url: fileID,
					width: Number(this.uploadImage.width),
					height: Number(this.uploadImage.height),
					size: Number(this.uploadImage.size),
					hash: this.uploadImage.sha256,
				})
				console.log("save image", addImage)
				this.plant.coverId = addImage.data.id;
			}
			this.plant.tags = this.tags.filter(item => item.active).map(item => ({ id: item.ID })) || []
			if (this.plant.coverId === 0) { this.plant.coverId = 6 }
			console.log("name", this.plant.name)
			console.log("coverId", this.plant.coverId)
			console.log("desc", this.plant.desc)
			console.log("birthday", this.plant.birthday)
			console.log("tags", this.plant.tags)
			try {
				const plant = await callContainer("/api/plant/add", {
					"name": this.plant.name,
					"familyId": this.familyId,
					"coverId": this.plant.coverId,
					"desc": this.plant.desc,
					"birthday": this.plant.birthday,
					"tags": this.plant.tags
				})
				console.log("call container plant add", plant)
				uni.$emit('refreshHomeList', { needRefresh: true });
			} catch (error) {
				console.error(error)
			} finally {
				this.isSave = true;
				this.$refs.loading.close();
				uni.navigateBack()
			}
		},
		async uploadImages() {
			try {
				const result = await new Promise((resolve, reject) => {
					wx.chooseMedia({
						count: 1,
						mediaType: ['image'],
						success: resolve,
						fail: reject
					})
				})
				const file = result.tempFiles[0]
				this.uploadImage.url = file.tempFilePath;
				this.uploadImage.size = file.size;
				const MAX_SIZE = 2 * 1024 * 1024 // 2MB
				if (file.size > MAX_SIZE) {
					wx.showToast({
						title: '不能超过2MB',
						icon: 'error',
						duration: 2000
					})
					return
				}
				this.plant.cover = file.tempFilePath;
				const imageInfo = await new Promise((resolve, reject) => {
					wx.getImageInfo({
						src: file.tempFilePath,
						success: resolve,
						fail: reject
					})
				})
				this.uploadImage.width = imageInfo.width
				this.uploadImage.height = imageInfo.height
				const sha256 = await new Promise((resolve, reject) => {
					wx.getFileSystemManager().getFileInfo({
						filePath: file.tempFilePath,
						digestAlgorithm: "sha256",
						success: resolve,
						fail: reject
					})
				})
				this.uploadImage.sha256 = sha256.digest;
				console.log("upload image", this.uploadImage)
				const isImage = await callContainer("/api/image/check", {
					hash: sha256.digest
				})
				console.log("last image info has image", isImage)
				this.uploadImage.isUpload = isImage.data.uploadRequired
				if (!isImage.data.uploadRequired) {
					this.uploadImage.url = isImage.data.url;
					this.plant.coverId = isImage.data.id;
				}
			} catch (error) {
				console.error(error)
			}
		}
	},
	onLoad(options) {
		if (options) {
			this.familyId = Number(options.familyId);
			console.log("获取家庭ID", this.familyId)
		}
		const menuButtonInfo = wx.getMenuButtonBoundingClientRect()
		this.menuButtonInfo = menuButtonInfo
		const systemInfo = uni.getWindowInfo()
		this.paddingLeft = systemInfo.screenWidth - menuButtonInfo.right
		const app = getApp()
		this.topBarHeight = app.globalData.topBarHeight;
		this.getTagsList()
	}
}
</script>

<style scoped lang="scss">
$focus-border-color: #566C44;
$bg-color: #E8F0E8;
$text-color: #566C44;

.container {
	font-size: 1rem;
}

.plant-title {
	width: 100%;
	color: #333;
	position: fixed;
	font-size: 34rpx;
	font-weight: bold;
	text-align: center;
}

.clean-btn,
.add-btn {
	width: 100%;
	height: 74rpx;
	background: rgba(255, 255, 255, 0.55);
	border-radius: 60rpx;

	display: flex;
	align-items: center;
	justify-content: center;
	transition: all 0.2s cubic-bezier(0.25, 0.8, 0.25, 1);
	border: 1px solid rgba(255, 255, 255, 0.1);

	&:active {
		transform: scale(0.92) translateY(2px);
		/* 点击时下沉 */
	}
}

.upload-section {
	display: flex;
	justify-content: center;
	/* 水平居中 */
	margin-top: 60rpx;
	/* 距离顶部按钮的间距 */
	margin-bottom: 40rpx;
}

.avatar-wrapper {
	width: 240rpx;
	height: 240rpx;
	position: relative;
	/* 作为定位参考点 */
	border-radius: 50rpx;
	background-color: #E8F0E8;
	/* 图片未加载时的占位底色 */
	box-shadow: 0 10rpx 20rpx rgba(0, 0, 0, 0.08);
	/* 柔和阴影 */
}

.plant-img {
	width: 100%;
	height: 100%;
	border-radius: 50rpx;
}

.camera-badge {
	position: absolute;
	bottom: 10rpx;
	right: 10rpx;
	width: 64rpx;
	height: 64rpx;
	background-color: #566C44;
	/* 绿色背景 */
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	border: 4rpx solid #fff;
	/* 白色描边，增加层次感 */
	box-shadow: 0 4rpx 8rpx rgba(0, 0, 0, 0.15);
}

.form-container {
	padding: 0 40rpx;
}

.form-item {
	background-color: $bg-color;
	margin-bottom: 30rpx;

	/* 默认状态：透明边框 (占位，防止抖动) */
	border: 2rpx solid transparent;

	/* 过渡效果 */
	transition: all 0.3s ease;
}

/* 聚焦状态：当添加了 focused 类名时 */
.form-item.focused {
	border-color: $focus-border-color;
	background-color: #f5f5f5;
	/* 聚焦时背景变白，更像输入框，体验更好 */
	box-shadow: 0 4rpx 12rpx rgba(86, 108, 68, 0.15);
	/* 加一点阴影 */
}

/* 1. 名称输入框 */
.input-pill {
	border-radius: 60rpx;
	padding: 10rpx 30rpx;
	display: flex;
	align-items: center;
	height: 100rpx;
}

.icon-right {
	margin-left: 20rpx;
	opacity: 0.8;
}

/* 2. 描述区域 */
.section-label {
	display: flex;
	align-items: center;
	margin-bottom: 16rpx;
	margin-left: 10rpx;

	.label-text {
		font-size: 30rpx;
		color: #333;
		font-weight: 600;
		margin-left: 10rpx;
	}
}

.textarea-box {
	border-radius: 30rpx;
	padding: 20rpx 30rpx;
}

/* 3. 日期选择器 */
.date-section {
	display: flex;
	justify-content: center;
	margin-top: 50rpx;
	margin-bottom: 30rpx;
}

.date-pill {
	background-color: #E8F0E8;
	padding: 16rpx 50rpx;
	border-radius: 50rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	color: #333;
	font-size: 30rpx;
	font-weight: 500;
	/* 日期按钮保持默认无边框或极淡边框 */
	border: 2rpx solid transparent;
}

.tags-wrapper {
	display: flex;
	flex-wrap: wrap;
	gap: 20rpx;
}

/* 单个标签 */
.tag-item {
	/* 默认状态：浅色背景，深色字 */
	background-color: $bg-color;
	padding: 12rpx 32rpx;
	border-radius: 40rpx;
	display: flex;
	align-items: center;
	transition: all 0.2s ease;
	/* 添加过渡动画 */
	border: 2rpx solid transparent;
	/* 预留边框位置防止抖动 */

	text {
		font-size: 26rpx;
		color: $text-color;
		font-weight: 500;
		transition: color 0.2s ease;
	}

	/* 选中状态：深色背景，浅色字 (互换颜色) */
	&.active {
		background-color: $text-color;
		/* #566C44 */
		box-shadow: 0 4rpx 10rpx rgba(86, 108, 68, 0.3);

		text {
			color: #fff;
			/* 或者使用 $bg-color (#E8F0E8) */
		}
	}
}

/* 添加按钮 */
.tag-add-btn {
	padding: 10rpx 26rpx;
	border-radius: 40rpx;
	display: flex;
	align-items: center;
	border: 2rpx dashed $text-color;
	background-color: rgba(255, 255, 255, 0.4);

	text {
		font-size: 26rpx;
		color: $text-color;
		margin-left: 6rpx;
	}

	&:active {
		background-color: rgba(86, 108, 68, 0.1);
	}
}


/* 穿透修改 uni-easyinput */
::v-deep .uni-easyinput__content {
	background-color: transparent !important;
	border: none !important;
}
</style>
