<template>
	<view class="container">
		<view>
			
		</view>
		<view class="flex_c_c page">
			<!-- Tabbar页面 -->
			<view class="flex1 custom-tabbar-page" v-if="list['a'].is"
				:style="{ display: onKey === 'a' ? '' : 'none' }">
				<scroll-view class="wh_100" scroll-y enable-back-to-top @scrolltolower="tabbarPageScrollLower">
					<!-- 页面a -->
				</scroll-view>
			</view>
			<view class="flex1 custom-tabbar-page" v-if="list['b'].is"
				:style="{ display: onKey === 'b' ? '' : 'none' }">
				<scroll-view class="wh_100" scroll-y enable-back-to-top @scrolltolower="tabbarPageScrollLower">
					<!-- 页面b -->
				</scroll-view>
			</view>
			<view class="flex1 custom-tabbar-page" v-if="list['c'].is"
				:style="{ display: onKey === 'c' ? '' : 'none' }">
				<scroll-view class="wh_100" scroll-y enable-back-to-top @scrolltolower="tabbarPageScrollLower">
					<!-- 页面c -->
					<home />
				</scroll-view>
			</view>
			<view class="flex1 custom-tabbar-page" v-if="list['d'].is"
				:style="{ display: onKey === 'd' ? '' : 'none' }">
				<scroll-view class="wh_100" scroll-y enable-back-to-top @scrolltolower="tabbarPageScrollLower">
					<!-- 页面d -->
				</scroll-view>
			</view>
			<view class="flex1 custom-tabbar-page" v-if="list['e'].is"
				:style="{ display: onKey === 'e' ? '' : 'none' }">
				<scroll-view class="wh_100" scroll-y enable-back-to-top @scrolltolower="tabbarPageScrollLower">
					<!-- 页面e -->
				</scroll-view>
			</view>

			<view class="pan-tabbar">
				<!-- 液面 -->
				<view class="liquid-box">
					<view class="flex_c_c liquid" :style="{ ...liquidStyle }">
						<view class="liquid-str01" :class="{ liquid_str01: animation01 }">
							<image class="wh_100 image-"
								src="data:image/jpeg;base64,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"
								mode="scaleToFill"></image>
							<!-- <image class="wh_100 image-"
								src="data:image/png;base64,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"
								mode="scaleToFill"></image> -->
						</view>
						<view class="flex1 liquid-str02"></view>
					</view>
				</view>
				<!-- 球 -->
				<view class="ball" :class="{ ball_: animation02 }" :style="{ ...ballStyle }">
				</view>
				<view class="flex_r fj_b pan-tabbar-list">
					<view class="pan-tabbar-item" :id="key" :data-key="key" :class="{ pan_tabbar_item: key === onKey }"
						v-for="(item, key) in list" :key='key' @click="onTabbar(item, key)">
						<view class="pan-tabbar-item-icon pan-tabbar-item-icon-off">
							<image class="wh_100" :src="item.iconOff" mode="aspectFill"></image>
						</view>
						<view class="pan-tabbar-item-icon pan-tabbar-item-icon-on">
							<image class="wh_100" :src="item.iconOn" mode="aspectFill"></image>
						</view>
					</view>
				</view>
				<!-- 底部安全区 -->
				<view class="safe-area-inset-bottom"></view>
			</view>

		</view>
	</view>
</template>

<script>
import a1 from "@/static/icon/a1.png"
import b1 from "@/static/icon/b1.png"
import c1 from "@/static/icon/c1.png"
import d1 from "@/static/icon/d1.png"
import e1 from "@/static/icon/e1.png"

import a2 from "@/static/icon/a2.png"
import b2 from "@/static/icon/b2.png"
import c2 from "@/static/icon/c2.png"
import d2 from "@/static/icon/d2.png"
import e2 from "@/static/icon/e2.png"
// 节流函数
let timer;
let flag;
const {
	windowWidth
} = uni.getSystemInfoSync();
import navBar from '@/components/navBar.vue';
import home from '@/components/home.vue';
export default {
	components: {
		navBar,
		home,
	},
	data() {
		return {
			animation01: false,
			animation02: false,
			ballStyleLeft: 0,
			liquidStyleLeft: 0,
			onKey: 'c',
			list: {
				a: {
					is: false,
					title: '01',
					iconOff: a1,
					iconOn: a2
				},
				b: {
					is: false,
					title: '02',
					iconOff: b1,
					iconOn: b2
				},
				c: {
					is: true,
					title: '03',
					iconOff: c1,
					iconOn: c2
				},
				d: {
					is: false,
					title: '04',
					iconOff: d1,
					iconOn: d2
				},
				e: {
					is: false,
					title: '05',
					iconOff: e1,
					iconOn: e2
				}
			}
		}
	},
	computed: {
		ballStyle() {
			let style = {
				left: `${this.ballStyleLeft}px`
			}
			return style;
		},
		liquidStyle() {
			let style = {
				left: `${this.liquidStyleLeft}px`
			}
			return style;
		},
	},
	methods: {
		tabbarPageScrollLower() { },
		init(e) {
			console.log('Tabbar init elements:', e);
			if (!e || e.length === 0) return;
			const keys = Object.keys(this.list)
			e.forEach((item, index) => {
				const key = keys[index]
				// 确保获取到了 id 且 list 中有这个 key
				if (key && this.list[key]) {
					// 修正位置
					item.left = item.left + 20;
					// 关键：将样式信息绑定到 list 数据中
					this.list[key]['style'] = item;
				}
			})
			this.onKey = 'c'
			const item = this.list[this.onKey];
			// 只有当 item 和 item.style 都存在时才执行
			if (item && item.style) {
				this.onTabbar(item, this.onKey);
			} else {
				// 如果走到这里，说明 rects 数量和 keys 数量可能不一致，或者渲染未完成
				console.warn("初始化匹配失败，尝试使用默认值");
				// 兜底逻辑：如果获取失败，手动给一个大概的位置，防止报错
				const defaultWidth = windowWidth / 5;
				this.ballStyleLeft = (defaultWidth * 2) + (defaultWidth / 2) - 22; // 假设 c 在中间
				this.liquidStyleLeft = (defaultWidth * 2) + (defaultWidth / 2) - (windowWidth / 2);
			}
		},
		onTabbar(item, key) {

			this.throttle(() => {
				if (!item || !item.style) {
					console.error('onTabbar error: item or item.style is undefined');
					return;
				}
				this.switchTabbarPage(key);
				this.onKey = key
				let left = item.style.left
				this.ballStyleLeft = left - 22;
				this.animation01 = true;
				this.animation02 = true;
				setTimeout(() => {
					this.liquidStyleLeft = left - (windowWidth / 2);
					this.animation01 = false;
				}, 300)
				setTimeout(() => {
					this.animation02 = false;
				}, 610)
			}, 610)
		},
		switchTabbarPage(key) {
			if (this.onKey === key) return
			const selectPageFlag = this.list[key]['is']
			if (selectPageFlag === undefined) {
				return
			}
			if (selectPageFlag === false) {
				this.list[key]['is'] = true
			}
		},
		// 获取元素位置
		getDemRefAll(e) {
			const {
				selector,
				success = () => { }
			} = e;
			return new Promise((r) => {
				let view = uni.createSelectorQuery().in(this).selectAll(selector);
				view.boundingClientRect((ref) => {
					success(ref);
					r(ref)
				}).exec();
			})
		},
		// 节流函数
		throttle(func, wait = 500, immediate = true) {
			if (immediate) {
				if (!flag) {
					flag = true
					// 如果是立即执行，则在wait毫秒内开始时执行
					typeof func === 'function' && func()
					timer = setTimeout(() => {
						flag = false
					}, wait)
				}
			} else if (!flag) {
				flag = true
				// 如果是非立即执行，则在wait毫秒内的结束处执行
				timer = setTimeout(() => {
					flag = false
					typeof func === 'function' && func()
				}, wait)
			}
		}
	},
	onLoad() {
		this.$nextTick(() => {
			this.getDemRefAll({
				selector: '.pan-tabbar-item',
				success: (e) => {
					this.init(e)
				}
			})
		})
	}
}
</script>

<style lang="scss">
.page {
	position: relative;
	width: 100%;
	height: 100vh;
	// background-color: #C1D0B7;

	.custom-tabbar-page {
		width: 100%;
		height: 0;
	}
}

.pan-tabbar {
	position: absolute;
	bottom: 0;
	left: 0;
	z-index: 2;
	width: 100vw;
	min-height: 100rpx;
	opacity: 0.55;
	// background-color: #0A3323;
}

.safe-area-inset-bottom {
	opacity: 0.55;
	background-color: #fff;
	padding-bottom: calc(constant(safe-area-inset-bottom));
	padding-bottom: calc(env(safe-area-inset-bottom));
	// margin-top: -20rpx
}

.liquid-box {
	position: absolute;
	left: 0;
	top: 0;
	z-index: 0;
	width: 100vw;
	height: 100rpx;
	overflow: hidden;
	opacity: 0.55;
	// background-color: #0A3323;
}

.liquid {
	position: absolute;
	width: 100%;
	height: 100%;
	top: 0;

	.liquid-str01 {
		position: relative;
		box-sizing: border-box;
		width: 100%;
		height: 100%;
		transition: height .3s cubic-bezier(0.333, 1.01, 0.32, 1.275), width 0.3s cubic-bezier(0.133, 1.01, 0.32, 1.275);

		.image- {
			display: block;
			position: absolute;
			left: 0;
			right: 0;
			bottom: 0;
			top: 0;
			height: calc(100% + 20rpx);
			transform-origin: bottom center;
			z-index: 1;
		}
	}

	.liquid-str01::after {
		position: absolute;
		left: -100vw;
		top: 0;
		content: '';
		width: 100vw;
		height: 100%;
		background-color: #fff;
	}

	.liquid-str01::before {
		position: absolute;
		right: -100vw;
		top: 0;
		content: '';
		width: 100vw;
		height: 100rpx;
		background-color: #fff;
	}

	.liquid_str01 {
		width: calc(100% + 100%);
		height: 0rpx;

	}

	.liquid-str02 {
		width: 100%;
		height: 0rpx;
		background-color: #fff;
	}

	.liquid-str02::after {
		position: absolute;
		left: -100vw;
		top: 0;
		content: '';
		width: 100vw;
		height: 100rpx;
		background-color: #fff;
	}

	.liquid-str02::before {
		position: absolute;
		right: -100vw;
		top: 0;
		content: '';
		width: 100vw;
		height: 100rpx;
		background-color: #fff;
	}

}


.ball {
	position: absolute;
	z-index: 3;
	top: -16px;
	width: 44px;
	height: 44px;
	border-radius: 50%;
	background-image: linear-gradient(220deg, #BC3823, #BC3823, rgba(105, 55, 248, 0.3));
	// background-image: linear-gradient(220deg, #566C44, #566C44, rgba(255, 255, 255, 0.8));
	// background-color: #8F3430;
	// background-color: red;
	// opacity: 0.85;
	// border: 1px solid rgba(0,0,0,0.08);
	// box-shadow: 0 2px 8px rgba(0,0,0,0.15);
	transition: all .6s;
}

.ball_ {
	animation: moveUpDownBall .6s ease-in-out;
}





.pan-tabbar-list {
	box-sizing: border-box;
	padding: 0 20px;
	position: relative;
	z-index: 3;
	width: 100%;
	height: 100rpx;


	.pan-tabbar-item {
		position: relative;
		width: 40px;
		height: 100%;

		.pan-tabbar-item-icon {
			position: absolute;
			top: 10px;
			left: 5px;
			width: 30px;
			height: 30px;
			transition: all .6s ease-in-out;
		}

		.pan-tabbar-item-icon-off {
			opacity: 1;
		}

		.pan-tabbar-item-icon-on {
			opacity: 0;
		}
	}

	.pan_tabbar_item {
		.pan-tabbar-item-icon {
			top: -10px;
		}

		.pan-tabbar-item-icon-off {
			opacity: 0;
		}

		.pan-tabbar-item-icon-on {
			opacity: 1;
		}
	}
}


/* 定义上下移动的动画 */
@keyframes moveUpDownBall {
	0% {
		transform: translateY(0);
	}

	40% {
		transform: translateY(-60px);
	}

	80% {
		transform: translateY(6px);
	}

	100% {
		transform: translateY(0);

	}
}

// 以下是通用css
.flex_r {
	display: flex;
	flex-direction: row;
}

.flex_c {
	display: flex;
	flex-direction: column;
}

.fj_b {
	justify-content: space-between;
}

.fj_a {
	justify-content: space-around;
}

.fj_c {
	justify-content: center;
}

.fa_c {
	align-items: center;
}

.flex_c_c {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
}

.flex_r_c {
	display: flex;
	flex-direction: row;
	align-items: center;
	justify-content: center;
}

.flex1 {
	flex: 1;
}

.wh_100 {
	width: 100%;
	height: 100%;
}

.z_index2 {
	position: relative;
	z-index: 2;
}

.hover_class {
	position: relative;
	overflow: hidden;
}

.hover_class::after {
	position: absolute;
	z-index: 2;
	content: '';
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background-color: rgba(0, 0, 0, 0.1);
}

.rotate_45 {
	transform: rotate(45deg);
}

.rotate_90 {
	transform: rotate(90deg);
}

.rotate_180 {
	transform: rotate(180deg);
}
</style>
