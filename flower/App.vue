<script>
import { callContainer } from './utils/request';
export default {
	globalData: {
		initPromise: null,
		topBarHeight: 0,
		bottomSafeAreaHeight: 0,
		windowWidth: 0,
		windowHeight: 0,
	},
	onLaunch: async function () {
		console.log('App Launch')
		if(!wx.cloud) {
			console.error('请使用 2.2.3 或以上的基础库以使用云能力')
		}else{
			wx.cloud.init({
				env:"prod-0gr2o3qpe533f1fb",
				traceUser: true,
			})
		}
		//获取手机信息
		const systemInfo = wx.getWindowInfo()
		this.globalData.windowWidth = systemInfo.windowWidth;
		this.globalData.windowHeight = systemInfo.windowHeight;
		const menuButtonInfo = wx.getMenuButtonBoundingClientRect()
		const statusBarHeight = systemInfo.statusBarHeight
		const navBarHeight = (menuButtonInfo.top - statusBarHeight) * 2 + menuButtonInfo.height
		const barHeight = statusBarHeight + navBarHeight
		this.globalData.topBarHeight = barHeight
		const bottomSafeAreaHeight = systemInfo.screenHeight - systemInfo.safeArea.bottom
		this.globalData.bottomSafeAreaHeight = bottomSafeAreaHeight
		// uni.hideTabBar()
		const user = await callContainer("/api/login")
		console.log("callContainer login:", user)
		await new Promise((resolve) => {
			uni.setStorage({ key: "family", data: user.data.family, success: resolve })
		})
		await new Promise((resolve) => {
			uni.setStorage({ key: "userInfo", data: user.data.user, success: resolve })
		})
	},
	onShow: function () {

	},
	onHide: function () {

	}
}
</script>

<style lang="scss">
/*每个页面公共css */
@import '@/uni_modules/uni-scss/index.scss';
/* #ifndef APP-NVUE */
@import '@/static/customicons.css';
/* icon font 引入*/
@import '@/static/iconfont.css';

// 设置整个项目的背景色
:root,
page {
	--bg-color: #C1D0B7;
	--bg-btn-color: #ffffff;
	--card-bg: #f8f8f8;
	/* 卡片背景 */
	--text-color: #333333;
	/* 主文本 */
	--text-sub: #999999;
	/* 副文本 */
	--border-color: #e5e5e5;
	/* 边框 */
	--primary-color: #007aff;
	/* 主题色 */
}

@media (prefers-color-scheme: dark) {

	:root,
	page {
		--bg-color: #202020;
		--bg-btn-color: #151515 --card-bg: #1e1e1e;
		--text-color: #ffffff;
		--text-sub: #777777;
		--border-color: #333333;
		/* 主题色通常在深色模式下需要调亮一点点，或者保持不变 */
		--primary-color: #0a84ff;
	}
}

page {
	background-color: var(--bg-color);
	color: var(--text-color);
	transition: background-color 0.3s, color, 0.3s;
}

/*
#fffff
#FAF2CB
#A7C190
#6B8857
#566C44
*/
</style>
