<script>
export default {
	globalData: {
		topBarHeight: 0,
		bottomSafeAreaHeight: 0,
		windowWidth: 0,
		windowHeight: 0,
	},
	onLaunch: function () {
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

:root,
page {
	--page-bg: #F1F7F4;
	--bg-color: var(--page-bg);
	--card-bg: #FCFFFD;
	--bg-btn-color: var(--card-bg);
	--surface-bg: #E8F1ED;
	--border-color: #D6E5DE;
	--primary-color: #72A990;
	--primary-strong: #4F8B70;
	--primary-soft: #DCEDE4;
	--text-color: #203029;
	--text-sub: #6C7A73;
	--button-text: #FFFFFF;
	--shadow-color: rgba(49, 96, 73, 0.10);
}

page {
	background-color: var(--page-bg);
	color: var(--text-color);
	color-scheme: light;
	transition: background-color 0.3s, color 0.3s;
}

view, scroll-view, input, textarea, button {
	box-sizing: border-box;
}

input, textarea {
	color: var(--text-color) !important;
}

.theme-page {
	min-height: 100vh;
	background: var(--page-bg);
	color: var(--text-color);
}

.theme-card {
	background: var(--card-bg);
	border: 1px solid var(--border-color);
	box-shadow: 0 8rpx 24rpx var(--shadow-color);
}

/*
#fffff
#FAF2CB
#A7C190
#6B8857
#566C44
*/
</style>
