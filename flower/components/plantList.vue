<template>
    <!-- 最外层包裹一个全屏容器 -->
    <view class="plantlist-container">

        <!-- 顶部固定区域 -->
        <view class="fixed-header-group">
            <!-- 顶部占位 -->
            <view :style="{ height: topBarHeight + 'px' }"></view>
        </view>

        <!-- 中间独立滚动区域 -->
        <view class="content-container">
            <!-- 加载中状态 -->
            <view v-if="loading" class="loading-wrapper">
                <text class="loading-text">加载中...</text>
            </view>

            <!-- 空状态 -->
            <view v-else-if="indexedOptions.length === 0" class="empty-wrapper">
                <image src="/static/icon/c2m.svg" class="empty-icon" mode="aspectFit"></image>
                <text class="empty-text">暂无植物</text>
            </view>

            <!-- 索引列表 -->
            <view v-else class="indexed-list-wrapper">
                <scroll-view scroll-y class="plant-scroll-view" :scroll-into-view="scrollIntoView">
                    <view v-for="(group, index) in indexedOptions" :key="index" :id="'group-' + group.key" class="plant-group">
                        <view class="group-title">{{ group.key }}</view>
                        <view v-for="item in group.data" :key="item.id" class="plant-item" @click="onPlantItemClick(item)">
                            <image 
                                v-if="item.data.cover && item.data.cover.url"
                                class="plant-image" 
                                :src="item.data.cover.url" 
                                mode="aspectFill"
                            ></image>
                            <view v-else class="plant-image placeholder">
                                <uni-icons type="image" size="30" color="#ccc"></uni-icons>
                            </view>
                            <view class="plant-info">
                                <text class="plant-name">{{ item.name }}</text>
                                <text v-if="item.data.desc" class="plant-desc">{{ item.data.desc }}</text>
                            </view>
                        </view>
                    </view>
                </scroll-view>
                
                <!-- 右侧索引条 -->
                <view class="index-bar">
                    <view 
                        v-for="(group, index) in indexedOptions" 
                        :key="index" 
                        class="index-item"
                        @click="scrollToGroup(group.key)"
                    >
                        {{ group.key }}
                    </view>
                </view>
            </view>
        </view>

    </view>
</template>

<script>
import { callContainer } from '../utils/request.js';

export default {
    name: 'plantList',
    emits: [],
    data() {
        return {
            loading: true,
            plantsList: [],
            indexedOptions: [],
            scrollIntoView: '',
            currentFamilyId: null,
            
            // 系统信息
            statusBarHeight: 0,
            topBarHeight: 0
        };
    },

    mounted() {
        this.initPageInfo();
        this.loadPlantsList();
        
        // 监听家庭切换事件
        uni.$on('familyChanged', this.handleFamilyChanged);
    },
    
    beforeUnmount() {
        // 移除监听
        uni.$off('familyChanged', this.handleFamilyChanged);
    },

    // 页面显示时检查家庭是否变化
    onShow() {
        this.checkFamilyChange();
    },

    methods: {
        // 初始化页面信息
        initPageInfo() {
            try {
                const systemInfo = uni.getSystemInfoSync();
                this.statusBarHeight = systemInfo.statusBarHeight || 44;
                const app = getApp();
                if (app && app.globalData) {
                    this.topBarHeight = app.globalData.topBarHeight || this.statusBarHeight;
                } else {
                    this.topBarHeight = this.statusBarHeight;
                }
            } catch (error) {
                console.error("获取系统信息失败:", error);
                this.statusBarHeight = 44;
                this.topBarHeight = 44;
            }
        },
        
        // 检查家庭是否变化
        async checkFamilyChange() {
            try {
                const familyIdResult = await new Promise((resolve, reject) => {
                    uni.getStorage({ key: 'familyId', success: resolve, fail: reject })
                });
                
                const newFamilyId = familyIdResult?.data;
                
                // 如果家庭ID变化了，重新加载
                if (newFamilyId && newFamilyId !== this.currentFamilyId) {
                    console.log('检测到家庭切换，重新加载植物列表');
                    this.currentFamilyId = newFamilyId;
                    await this.loadPlantsList();
                }
            } catch (error) {
                console.error('检查家庭变化失败:', error);
            }
        },
        
        // 处理家庭切换事件
        async handleFamilyChanged(newFamilyId) {
            console.log('收到家庭切换事件:', newFamilyId);
            this.currentFamilyId = newFamilyId;
            this.plantsList = [];
            this.indexedOptions = [];
            await this.loadPlantsList();
        },

        // 加载植物列表
        async loadPlantsList() {
            try {
                this.loading = true;

                // 获取当前家庭ID
                const familyIdResult = await new Promise((resolve, reject) => {
                    uni.getStorage({ key: 'familyId', success: resolve, fail: reject })
                });
                
                const familyId = familyIdResult?.data;
                
                if (!familyId) {
                    uni.showToast({
                        title: '请先选择家庭',
                        icon: 'none'
                    });
                    this.loading = false;
                    return;
                }
                
                // 更新当前家庭ID
                this.currentFamilyId = familyId;

                const result = await callContainer('/api/plant/list', {
                    familyId: familyId
                });

                console.log('植物列表数据:', result);
                this.plantsList = result?.data || [];
                this.processIndexedData();

            } catch (error) {
                console.error('加载植物列表失败:', error);
                uni.showToast({
                    title: '加载失败，请重试',
                    icon: 'none'
                });
            } finally {
                this.loading = false;
            }
        },

        // 处理数据为索引列表格式
        processIndexedData() {
            if (!this.plantsList || this.plantsList.length === 0) {
                this.indexedOptions = [];
                return;
            }

            const groupMap = {};

            this.plantsList.forEach(plant => {
                if (!plant.name) return;

                let firstChar = this.getFirstLetter(plant.name);

                if (!groupMap[firstChar]) {
                    groupMap[firstChar] = [];
                }

                // uni-indexed-list 默认显示 name 字段
                groupMap[firstChar].push({
                    name: plant.name,
                    id: plant.ID,
                    data: plant
                });
            });

            this.indexedOptions = Object.keys(groupMap)
                .sort()
                .map(key => ({
                    key: key.toUpperCase(),
                    data: groupMap[key]
                }));

            console.log('处理后的索引数据:', this.indexedOptions);
        },

        // 获取字符串首字母
        getFirstLetter(str) {
            if (!str) return '#';

            const firstChar = str.charAt(0).toUpperCase();

            // 如果已经是英文字母，直接返回
            if (/^[A-Z]$/.test(firstChar)) {
                return firstChar;
            }

            // 如果是数字，返回 #
            if (/^[0-9]$/.test(firstChar)) {
                return '#';
            }

            // 使用 Unicode 范围判断并转换中文到拼音首字母
            const code = str.charCodeAt(0);
            
            // 常用汉字 Unicode 范围判断
            if (code >= 19968 && code <= 40869) {
                // 使用更完整的拼音首字母映射
                return this.getChinesePinyin(str.charAt(0));
            }

            return '#';
        },
        
        // 汉字转拼音首字母（字典映射 + 边界对比，彻底解决环境兼容性问题）
        getChinesePinyin(char) {
            if (!char) return '#';

            // 1. 特殊覆盖与多音字纠正
            const overrides = {
                '梦': 'M', '重': 'C', '长': 'C', '行': 'X', '厦': 'X', '地': 'D', '重': 'Z'
            };
            if (overrides[char]) return overrides[char];

            // 2. 核心常用字字典 (解决 localeCompare 在部分环境失效的问题)
            // 包含您提到的 W、J、H、A 等易错区域
            const pinyinDict = {
                'A': '阿啊啊爱安按暗案奥',
                'B': '八巴吧把白百柏摆败班般板版办半帮绑榜棒包炮薄饱宝报抱',
                'C': '擦猜才材财采彩菜参餐藏操草册侧测层叉插查茶察差拆柴产单产阐',
                'D': '搭答打大呆代带待袋贷单担胆诞但弹淡蛋当挡党荡刀导岛到盗道稻',
                'E': '鹅额俄恶饿儿而尔耳二',
                'F': '发乏伐罚阀法帆翻凡烦反返范贩方防坊房放非飞肥匪费废分芬粉',
                'G': '旮该改盖概干甘杆肝赶敢干冈刚岗纲港杠高糕搞稿告戈歌各哥个给',
                'H': '哈孩海害含寒喊汉汗行号毫豪好喝河荷核合和何合河盒贺赫黑痕很',
                'J': '讥击机积基绩激及吉级即极急疾集及给籍几己挤济季记计既继寄加夹佳家架价甲假钾价尖坚间肩艰兼监渐践鉴键箭江将讲奖降匠酱交郊娇浇骄胶焦角脚教阶皆接揭街节杰结捷截竭姐解介戒届界巾斤今金津筋仅',
                'K': '咔开刊看康慷扛抗考烤靠科壳可渴克刻客肯恳坑空孔恐口扣哭苦库裤',
                'L': '啦喇拉落垃圾蓝篮栏拦兰婪澜览懒烂滥郎廊狼浪捞劳牢老落涝乐雷累',
                'M': '嘛妈麻马玛码买卖脉满慢忙芒茫毛矛貌貌没枚眉梅媒煤美每门们萌盟猛梦孟眯米秘密',
                'N': '拿哪那纳娜奶耐男南难囊挠脑闹内嫩能尼泥你拟逆年念娘鸟尿捏您宁',
                'O': '喔喔哦噢欧偶沤',
                'P': '妑趴啪怕爬帕排牌派判攀盘判盼旁庞抛炮跑步陪培赔佩配喷盆蓬棚蓬',
                'Q': '七妻戚期欺齐其奇骑旗棋起启岂气弃汽契砌器恰洽千牵迁签前钱潜浅',
                'R': '然燃染嚷让饶惹热人仁忍认任扔仍日绒荣容溶熔融柔肉如茹儒乳辱入',
                'S': '仨洒萨赛三伞散桑嗓丧扫嫂色涩森僧杀沙纱傻啥晒删山闪陕扇善伤商',
                'T': '他塔踏台抬太态泰贪摊滩谈坛昙谈坦炭探汤糖躺趟涛掏逃陶陶讨套特',
                'W': '哇歪外弯湾完玩顽晚碗万汪王网往忘望危威微为维围委伟伪尾纬未味谓喂慰温文纹稳问翁嗡我沃卧握乌污',
                'X': '夕汐西吸希析悉稀溪习席袭洗喜戏系细隙虾峡侠狭下夏仙先鲜闲',
                'Y': '丫牙雅亚讶烟延严言岩沿炎研盐颜掩眼演厌宴艳验样羊阳养仰样腰',
                'Z': '匝杂灾栽载宰再在咱攒暂赞脏葬遭糟早枣造燥躁则择泽责贼怎增赠扎'
            };

            for (let key in pinyinDict) {
                if (pinyinDict[key].indexOf(char) !== -1) return key;
            }

            // 3. 兜底策略：localeCompare (用于字典未涵盖的生僻字)
            try {
                const letters = "ABCDEFGHJKLMNOPQRSTWXYZ".split("");
                const boundary = "阿八嚓哒妸发旮哈讥咔垃妈拿喔妑七然仨他哇夕丫匝".split("");
                for (let i = 0; i < boundary.length; i++) {
                    if (char.localeCompare(boundary[i], 'zh') < 0) {
                        return i === 0 ? 'A' : letters[i - 1];
                    }
                }
            } catch (e) {
                console.error('拼音转换失败:', e);
            }
            
            return 'Z';
        },

        // 点击植物项
        onPlantClick(e) {
            console.log('点击植物:', e);
            const plantData = e.item?.data;

            if (plantData && plantData.ID) {
                uni.navigateTo({
                    url: `/pages/plantDetail/plantDetail?id=${plantData.ID}`
                });
            }
        },
        
        // 自定义列表项点击
        onPlantItemClick(item) {
            console.log('点击植物项:', item);
            if (item && item.data && item.data.ID) {
                uni.navigateTo({
                    url: `/pages/plantDetail/plantDetail?id=${item.data.ID}`
                });
            }
        },
        
        // 滚动到指定分组
        scrollToGroup(key) {
            this.scrollIntoView = 'group-' + key;
            setTimeout(() => {
                this.scrollIntoView = '';
            }, 300);
        }
    }
};
</script>

<style lang="scss" scoped>
.plantlist-container {
    width: 100%;
    height: 100vh;
    display: flex;
    flex-direction: column;
    padding-bottom: calc(100rpx + env(safe-area-inset-bottom));
    box-sizing: border-box;
    background-color: #C1D0B7;
}

// 固定头部
.fixed-header-group {
    position: sticky;
    top: 0;
    z-index: 998;
    background-color: #C1D0B7;
}

// 内容区域
.content-container {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

.loading-wrapper {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 100rpx 0;
}

.loading-text {
    font-size: 14px;
    color: #999;
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
    filter: grayscale(100%);
}

.empty-text {
    font-size: 14px;
    color: #999;
}

.indexed-list-wrapper {
    flex: 1;
    overflow: hidden;
    position: relative;
}

// 滚动视图
.plant-scroll-view {
    height: 100%;
    width: 100%;
}

// 分组容器
.plant-group {
    margin-bottom: 20rpx;
}

// 分组标题
.group-title {
    padding: 16rpx 32rpx;
    font-size: 14px;
    font-weight: 600;
    color: #6B8857;
    background: rgba(107, 136, 87, 0.1);
    position: sticky;
    top: 0;
    z-index: 10;
}

// 自定义植物项样式
.plant-item {
    padding: 24rpx 32rpx;
    background: rgba(255, 255, 255, 0.55);
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);
    transition: background 0.2s;
    display: flex;
    align-items: center;
    gap: 24rpx;
    
    &:active {
        background: rgba(107, 136, 87, 0.1);
    }
}

.plant-image {
    width: 100rpx;
    height: 100rpx;
    border-radius: 12rpx;
    flex-shrink: 0;
    background: #f5f5f5;
    
    &.placeholder {
        display: flex;
        align-items: center;
        justify-content: center;
        background: #f5f5f5;
    }
}

.plant-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8rpx;
    overflow: hidden;
}

.plant-name {
    font-size: 15px;
    color: #333;
    font-weight: 500;
}

.plant-desc {
    font-size: 12px;
    color: #999;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

// 右侧索引条
.index-bar {
    position: absolute;
    right: 10rpx;
    top: 50%;
    transform: translateY(-50%);
    display: flex;
    flex-direction: column;
    align-items: center;
    z-index: 100;
}

.index-item {
    padding: 4rpx 8rpx;
    font-size: 12px;
    color: #6B8857;
    font-weight: 500;
    line-height: 1.2;
    
    &:active {
        color: #fff;
        background: #6B8857;
        border-radius: 50%;
    }
}

// 索引列表样式优化
::v-deep .uni-indexed-list {
    height: 100%;
}

::v-deep .uni-indexed-list__scroll {
    height: 100%;
}

::v-deep .uni-indexed-list__item {
    padding: 12px 20px;
    background: rgba(255, 255, 255, 0.5);
    border-radius: 8px;
    margin: 8px 16px;
    transition: all 0.2s;

    &:active {
        background: rgba(255, 255, 255, 0.8);
        transform: scale(0.98);
    }
}

::v-deep .uni-indexed-list__item-content {
    font-size: 15px;
    color: #333;
}

::v-deep .uni-indexed-list__menu {
    right: 5px;
}
</style>
