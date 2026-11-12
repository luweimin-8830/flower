<template>
	<view class="uni-calendar-item__weeks-box" :class="{
		'uni-calendar-item--disable':weeks.disable,
		'uni-calendar-item--isDay':calendar.fullDate === weeks.fullDate && weeks.isDay,
		'uni-calendar-item--checked':(calendar.fullDate === weeks.fullDate && !weeks.isDay) ,
		'uni-calendar-item--before-checked':weeks.beforeMultiple,
		'uni-calendar-item--multiple': weeks.multiple,
		'uni-calendar-item--after-checked':weeks.afterMultiple,
		}"
	 @click="choiceDate(weeks)">
		<view class="uni-calendar-item__weeks-box-item">
			<text v-if="selected&&weeks.extraInfo" class="uni-calendar-item__weeks-box-circle"></text>
			<text class="uni-calendar-item__weeks-box-text" :class="{
				'uni-calendar-item--isDay-text': weeks.isDay,
				'uni-calendar-item--isDay':calendar.fullDate === weeks.fullDate && weeks.isDay,
				'uni-calendar-item--checked':calendar.fullDate === weeks.fullDate && !weeks.isDay,
				'uni-calendar-item--before-checked':weeks.beforeMultiple,
				'uni-calendar-item--multiple': weeks.multiple,
				'uni-calendar-item--after-checked':weeks.afterMultiple,
				'uni-calendar-item--disable':weeks.disable,
				}">{{weeks.date}}</text>
			<text v-if="!lunar&&!weeks.extraInfo && weeks.isDay" class="uni-calendar-item__weeks-lunar-text" :class="{
				'uni-calendar-item--isDay-text':weeks.isDay,
				'uni-calendar-item--isDay':calendar.fullDate === weeks.fullDate && weeks.isDay,
				'uni-calendar-item--checked':calendar.fullDate === weeks.fullDate && !weeks.isDay,
				'uni-calendar-item--before-checked':weeks.beforeMultiple,
				'uni-calendar-item--multiple': weeks.multiple,
				'uni-calendar-item--after-checked':weeks.afterMultiple,
				}">{{todayText}}</text>
			<text v-if="lunar&&!weeks.extraInfo" class="uni-calendar-item__weeks-lunar-text" :class="{
				'uni-calendar-item--isDay-text':weeks.isDay,
				'uni-calendar-item--isDay':calendar.fullDate === weeks.fullDate && weeks.isDay,
				'uni-calendar-item--checked':calendar.fullDate === weeks.fullDate && !weeks.isDay,
				'uni-calendar-item--before-checked':weeks.beforeMultiple,
				'uni-calendar-item--multiple': weeks.multiple,
				'uni-calendar-item--after-checked':weeks.afterMultiple,
				'uni-calendar-item--disable':weeks.disable,
				}">{{weeks.isDay ? todayText : (weeks.lunar.IDayCn === '初一'?weeks.lunar.IMonthCn:weeks.lunar.IDayCn)}}</text>
			<text v-if="weeks.extraInfo&&weeks.extraInfo.info" class="uni-calendar-item__weeks-lunar-text" :class="{
				'uni-calendar-item--extra':weeks.extraInfo.info,
				'uni-calendar-item--isDay-text':weeks.isDay,
				'uni-calendar-item--isDay':calendar.fullDate === weeks.fullDate && weeks.isDay,
				'uni-calendar-item--checked':calendar.fullDate === weeks.fullDate && !weeks.isDay,
				'uni-calendar-item--before-checked':weeks.beforeMultiple,
				'uni-calendar-item--multiple': weeks.multiple,
				'uni-calendar-item--after-checked':weeks.afterMultiple,
				'uni-calendar-item--disable':weeks.disable,
				}">{{weeks.extraInfo.info}}</text>
		</view>
	</view>
</template>

<script>
	import { initVueI18n } from '@dcloudio/uni-i18n'
	import i18nMessages from './i18n/index.js'
	const {	t	} = initVueI18n(i18nMessages)

	export default {
		emits:['change'],
		props: {
			weeks: {
				type: Object,
				default () {
					return {}
				}
			},
			calendar: {
				type: Object,
				default: () => {
					return {}
				}
			},
			selected: {
				type: Array,
				default: () => {
					return []
				}
			},
			lunar: {
				type: Boolean,
				default: false
			}
		},
		computed: {
			todayText() {
				return t("uni-calender.today")
			},
		},
		methods: {
			choiceDate(weeks) {
				this.$emit('change', weeks)
			}
		}
	}
</script>

<style lang="scss" scoped>
	$uni-font-size-base:14px;
	$uni-text-color:#333;
	$uni-font-size-sm:12px;
	$uni-color-error: #e43d33;
	$uni-opacity-disabled: 0.3;
	$uni-text-color-disable:#c0c0c0;
	$uni-primary: #566C44 !default;
	.uni-calendar-item__weeks-box {
		flex: 1;
		/* #ifndef APP-NVUE */
		display: flex;
		/* #endif */
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}

	.uni-calendar-item__weeks-box-text {
		font-size: $uni-font-size-base;
		color: $uni-text-color;

	}

	.uni-calendar-item__weeks-lunar-text {
		font-size: $uni-font-size-sm;
		color: $uni-text-color;

	}

	.uni-calendar-item__weeks-box-item {
		position: relative;
		/* #ifndef APP-NVUE */
		display: flex;
		/* #endif */
		flex-direction: column;
		justify-content: center;
		align-items: center;
		width: 100rpx;
		height: 100rpx;
	}

	.uni-calendar-item__weeks-box-circle {
		position: absolute;
		top: 5px;
		right: 5px;
		width: 8px;
		height: 8px;
		border-radius: 8px;
		background-color: $uni-color-error;

	}

	.uni-calendar-item--disable {
		background-color: rgba(249, 249, 249, $uni-opacity-disabled);
		color: $uni-text-color-disable;

	}

	.uni-calendar-item--isDay-text {
		color: #566C44;
	}

	.uni-calendar-item--isDay {
		background-color: #566C44;
		opacity: 0.8;
		color: #fff;

	}

	.uni-calendar-item--extra {
		color: $uni-color-error;
		opacity: 0.8;
	}

	.uni-calendar-item--checked {
		background-color: #566C44;
		color: #fff;
		opacity: 0.8;

	}

	.uni-calendar-item--multiple {
		background-color: #566C44;
		color: #fff;
		opacity: 0.8;

	}
	.uni-calendar-item--before-checked {
		background-color: #ff5a5f;
		color: #fff;
	}
	.uni-calendar-item--after-checked {
		background-color: #ff5a5f;
		color: #fff;
	}

	/* 单选日期使用固定正方形内层绘制，确保选中背景始终为圆形而不是方块。 */
	.uni-calendar-item__weeks-box.uni-calendar-item--isDay,
	.uni-calendar-item__weeks-box.uni-calendar-item--checked {
		background-color: transparent;
		opacity: 1;

		> .uni-calendar-item__weeks-box-item {
			width: 76rpx;
			height: 76rpx;
			border-radius: 50%;
			background-color: var(--primary-color, #566C44);
			overflow: hidden;
		}

		.uni-calendar-item--isDay,
		.uni-calendar-item--checked {
			background-color: transparent;
			color: #fff;
			opacity: 1;
		}
	}

	/* 日期文字跟随当前主题；微信小程序组件隔离下需在组件内部声明。 */
	.uni-calendar-item__weeks-box-text,
	.uni-calendar-item__weeks-lunar-text {
		color: var(--text-color, #333) !important;
	}

	.uni-calendar-item--disable {
		background-color: transparent !important;
		color: var(--text-sub, #c0c0c0) !important;
		opacity: .36;
	}

	.uni-calendar-item--isDay-text {
		color: var(--primary-color, #566c44) !important;
	}

	.uni-calendar-item__weeks-box.uni-calendar-item--isDay,
	.uni-calendar-item__weeks-box.uni-calendar-item--checked {
		> .uni-calendar-item__weeks-box-item {
			background-color: var(--primary-color, #566c44) !important;
		}

		.uni-calendar-item__weeks-box-text,
		.uni-calendar-item__weeks-lunar-text,
		.uni-calendar-item--isDay,
		.uni-calendar-item--checked {
			color: var(--button-text, #fff) !important;
		}
	}
</style>
