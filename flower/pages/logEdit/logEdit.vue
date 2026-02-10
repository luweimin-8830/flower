<template>
	<view class="container">
		<navBar :title="type === 'edit' ? '编辑日志' : '添加日志'" />
		<view :style="{ height: topBarHeight + 'px' }"></view>

		<view class="form-container">
			<!-- 1. 选择类型 -->
			<view class="section">
				<text class="section-title">记录类型</text>
				<scroll-view scroll-x class="care-scroll" :show-scrollbar="false">
					<view class="care-list">
						<view class="care-item" v-for="(item, index) in careOptions" :key="index"
							@click="formData.actionType = (formData.actionType === item.type ? '' : item.type)">
							<view class="care-icon-box"
								:class="{ 'active': formData.actionType === item.type }"
								:style="{ backgroundColor: formData.actionType === item.type ? item.color : '#f5f5f5' }">
								<view class="iconfont" :class="item.icon" 
									style="font-size: 24px;"
									:style="{ color: formData.actionType === item.type ? '#fff' : '#999' }"></view>
							</view>
							<text class="care-name" :class="{ 'active-text': formData.actionType === item.type }">{{ item.name }}</text>
						</view>
					</view>
				</scroll-view>
			</view>

			<!-- 2. 日期选择 -->
			<view class="section">
				<text class="section-title">记录日期</text>
				<picker mode="date" :value="formData.logTime" @change="onDateChange">
					<view class="picker-box">
						<text>{{ formData.logTime }}</text>
						<uni-icons type="calendar" size="18" color="#666"></uni-icons>
					</view>
				</picker>
			</view>

			<!-- 3. 文字内容 -->
			<view class="section">
				<text class="section-title">记录内容</text>
				<textarea class="content-input" v-model="formData.content" placeholder="记录下植物的点滴成长吧..." maxlength="500"></textarea>
			</view>

			<!-- 4. 图片上传 -->
			<view class="section">
				<text class="section-title">记录照片 (最多3张)</text>
				<view class="image-grid">
					<view class="image-item" v-for="(img, index) in uploadImages" :key="index">
						<image :src="img.url" mode="aspectFill" @click="previewImage(index)"></image>
						<view class="delete-btn" @click="removeImage(index)">
							<uni-icons type="closeempty" size="12" color="#fff"></uni-icons>
						</view>
					</view>
					<view class="image-item add-btn" v-if="uploadImages.length < 3" @click="chooseImage">
						<uni-icons type="plusempty" size="30" color="#999"></uni-icons>
					</view>
				</view>
			</view>

			<!-- 保存按钮 -->
			<view class="footer">
				<button class="save-btn" @click="handleSave" :loading="isSaving" :disabled="isSaving">
					保存记录
				</button>
			</view>
		</view>

		<uni-load-more v-if="isSaving" status="loading" iconType="circle"></uni-load-more>
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
			logId: 0,
			type: 'add',
			formData: {
				actionType: '',
				logTime: '',
				content: ''
			},
			careOptions: [],
			uploadImages: [],
			isSaving: false,
			familyId: 0
		};
	},
	async onLoad(options) {
		const app = getApp();
		this.topBarHeight = app.globalData.topBarHeight;
		this.familyId = uni.getStorageSync('familyId');
		
		if (options.plantId) {
			this.plantId = Number(options.plantId);
		}
		
		if (options.id) {
			this.logId = Number(options.id);
			this.type = 'edit';
			this.getLogDetail();
		}
		
		// 设置默认日期
		if (options.date) {
			this.formData.logTime = options.date;
		} else {
			const now = new Date();
			this.formData.logTime = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
		}
		
		this.getCareOptions();
	},
	watch: {
		careOptions: {
			handler(newVal) {
				if (newVal && newVal.length > 0 && !this.formData.actionType && this.type === 'add') {
					const growthAction = newVal.find(item => item.name === '成长记录' || item.type === 'Growth');
					if (growthAction) {
						this.formData.actionType = growthAction.type;
					}
				}
			},
			immediate: false
		}
	},
	methods: {
		async getCareOptions() {
			try {
				const result = await callContainer("/api/care/", { familyId: Number(this.familyId) });
				let options = [];
				// 前端写死第一个：成长记录
				const growthOption = {
					name: '成长记录',
					type: 'record',
					icon: 'plant-zhiwuzhiyuan-duorouzhiwuyuan',
					color: '#A7C190'
				};
				
				if (result.data) {
					// 过滤掉后端可能存在的重复项（兼容老数据）
					options = result.data.filter(item => item.type !== 'record' && item.type !== 'Growth' && item.name !== '成长记录');
				}
				this.careOptions = [growthOption, ...options];
			} catch (e) { console.error(e); }
		},
		async getLogDetail() {
			try {
				const res = await callContainer("/api/plant/log/detail", { id: this.logId });
				if (res.data) {
					const data = res.data;
					this.formData.actionType = data.actionType;
					this.formData.content = data.content;
					this.formData.logTime = data.logTime.split('T')[0];
					if (data.images) {
						this.uploadImages = data.images.map(img => ({
							id: img.id,
							url: img.url,
							isOld: true
						}));
					}
				}
			} catch (e) { console.error(e); }
		},
		onDateChange(e) {
			this.formData.logTime = e.detail.value;
		},
		async chooseImage() {
			try {
				const res = await new Promise((resolve, reject) => {
					wx.chooseMedia({
						count: 3 - this.uploadImages.length,
						mediaType: ['image'],
						success: resolve,
						fail: reject
					})
				});
				
				const files = res.tempFiles;
				uni.showLoading({ title: '处理图片...' });
				
				for (let file of files) {
					// 1. 大小检查
					const MAX_SIZE = 2 * 1024 * 1024; // 2MB
					if (file.size > MAX_SIZE) {
						uni.showToast({ title: '图片不能超过2MB', icon: 'none' });
						continue;
					}

					// 2. 获取图片信息
					const imageInfo = await new Promise((resolve, reject) => {
						wx.getImageInfo({
							src: file.tempFilePath,
							success: resolve,
							fail: reject
						})
					});

					// 3. 计算 SHA256
					const fileInfo = await new Promise((resolve, reject) => {
						wx.getFileSystemManager().getFileInfo({
							filePath: file.tempFilePath,
							digestAlgorithm: "sha256",
							success: resolve,
							fail: reject
						})
					});
					const hash = fileInfo.digest;

					// 4. 秒传检查
					const checkRes = await callContainer("/api/image/check", { hash });
					
					this.uploadImages.push({
						url: checkRes.data.uploadRequired ? file.tempFilePath : checkRes.data.url,
						size: file.size,
						width: imageInfo.width,
						height: imageInfo.height,
						hash: hash,
						isUpload: checkRes.data.uploadRequired,
						id: checkRes.data.id,
						isNew: true
					});
				}
				uni.hideLoading();
			} catch (e) { 
				console.error(e); 
				uni.hideLoading();
			}
		},
		removeImage(index) {
			this.uploadImages.splice(index, 1);
		},
		previewImage(index) {
			uni.previewImage({
				current: index,
				urls: this.uploadImages.map(img => img.url)
			});
		},
		async handleSave() {
			if (!this.formData.actionType) {
				uni.showToast({ title: '请选择记录类型', icon: 'none' });
				return;
			}
			if (this.type === 'add' && !this.plantId) {
				uni.showToast({ title: '参数错误: 缺少植物ID', icon: 'none' });
				return;
			}
			
			this.isSaving = true;
			try {
				// 1. 处理图片
				const imageIds = [];
				for (let img of this.uploadImages) {
					if (img.isOld) {
						imageIds.push(img.id);
						continue;
					}
					
					if (!img.isUpload) {
						// 秒传，直接使用 ID
						imageIds.push(img.id);
						continue;
					}

					// 需要上传
					const timestamp = Math.floor(Date.now() / 1000);
					const random = Math.floor(Math.random() * 1000000).toString().padStart(6, '0');
					const cloudPath = `${this.familyId}/logs/${timestamp}${random}.jpg`;
					
					const uploadRes = await wx.cloud.uploadFile({
						cloudPath,
						filePath: img.url,
						config: {
							env: 'prod-0gr2o3qpe533f1fb'
						}
					});
					
					// 保存图片记录到数据库
					const saveImgRes = await callContainer("/api/image/add", {
						url: uploadRes.fileID,
						width: img.width,
						height: img.height,
						size: img.size,
						hash: img.hash
					});
					imageIds.push(saveImgRes.data.id);
				}
				
				// 2. 保存日志
				if (this.type === 'add') {
					await callContainer("/api/plant/log/add", {
						plantIds: [this.plantId],
						actionType: this.formData.actionType,
						content: this.formData.content,
						logTime: this.formData.logTime,
						imageIds: imageIds
					});
				} else {
					await callContainer("/api/plant/log/update", {
						id: this.logId,
						actionType: this.formData.actionType,
						content: this.formData.content,
						logTime: this.formData.logTime,
						imageIds: imageIds
					});
				}
				
				uni.showToast({ title: '保存成功', icon: 'success' });
				uni.$emit('refreshHomeList');
				uni.$emit('refreshPlantDetail');
				setTimeout(() => uni.navigateBack(), 1500);
			} catch (e) {
				console.error(e);
				uni.showToast({ title: e.message || '保存失败', icon: 'none' });
			} finally {
				this.isSaving = false;
			}
		}
	}
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #C1D0B7;
	padding-bottom: 50px;
}

.form-container {
	padding: 20px 16px;
}

.section {
	margin-bottom: 24px;
	background-color: rgba(255, 255, 255, 0.55);
	padding: 16px;
	border-radius: 16px;
}

.section-title {
	font-size: 14px;
	color: #4A6139;
	font-weight: bold;
	margin-bottom: 12px;
	display: block;
}

.care-scroll {
	width: 100%;
	height: 150rpx;
}

.care-list {
	display: inline-flex;
	padding: 12px 60rpx 12px 10px;
	height: 100%;
}

.care-item {
	display: flex;
	flex-direction: column;
	align-items: center;
	margin-right: 12px;
	width: 110rpx;
	flex-shrink: 0;
	
	.care-icon-box {
		width: 90rpx;
		height: 90rpx;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 6px;
		border: 1px solid transparent;
		transition: all 0.2s;
		
		&.active {
			transform: scale(1.1);
			box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
		}
	}
	
	.care-name {
		font-size: 11px;
		color: #999;
		text-align: center;
		width: 100%;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
		
		&.active-text {
			color: #4A6139;
			font-weight: bold;
		}
	}
}

.picker-box {
	height: 44px;
	background-color: rgba(255,255,255,0.6);
	border-radius: 8px;
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 0 12px;
	font-size: 14px;
	color: #333;
}

.content-input {
	width: 100%;
	height: 100px;
	background-color: rgba(255,255,255,0.6);
	border-radius: 8px;
	padding: 12px;
	font-size: 14px;
	box-sizing: border-box;
}

.image-grid {
	display: flex;
	flex-wrap: wrap;
	gap: 10px;
}

.image-item {
	width: 80px;
	height: 80px;
	border-radius: 8px;
	position: relative;
	overflow: hidden;
	
	image {
		width: 100%;
		height: 100%;
	}
	
	&.add-btn {
		background-color: rgba(255,255,255,0.6);
		display: flex;
		align-items: center;
		justify-content: center;
		border: 1px dashed #ccc;
	}
	
	.delete-btn {
		position: absolute;
		right: 4px;
		top: 4px;
		width: 18px;
		height: 18px;
		background-color: rgba(0, 0, 0, 0.5);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
	}
}

.footer {
	margin-top: 40px;
}

.save-btn {
	background-color: #4A6139;
	color: #fff;
	height: 48px;
	border-radius: 24px;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 16px;
	font-weight: bold;
	border: none;
	
	&:disabled {
		opacity: 0.7;
	}
}
</style>
