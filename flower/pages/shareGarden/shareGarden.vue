<template>
    <view class="share-page theme-page" :style="themeStyle">
        <navBar title="分享长图" />
        <view :style="{ height: topBarHeight + 'px' }"></view>

        <view class="preview-layout">
            <scroll-view scroll-y class="preview-scroll">
                <view class="poster-preview">
                    <view class="poster-header">
                        <view class="poster-eyebrow"><view class="eyebrow-line"></view><text>MY GARDEN</text><view class="eyebrow-line"></view></view>
                        <text class="poster-title">{{ shareTitle }}</text>
                        <text class="poster-meta">{{ displayTagName }}收藏 · {{ plants.length }}株植物</text>
                        <view class="poster-leaf-mark"><view class="leaf leaf-left"></view><view class="leaf-dot"></view><view class="leaf leaf-right"></view></view>
                    </view>
                    <view v-if="plants.length" class="poster-columns">
                        <view v-for="(column, columnIndex) in previewColumns" :key="columnIndex" class="poster-column">
                            <view v-for="plant in column" :key="plant.id" class="poster-plant">
                                <image :src="plant.cover?.url || '/static/default.svg'" mode="widthFix"></image>
                                <view class="poster-plant-info"><text class="poster-plant-name">{{ plant.name || '未命名植物' }}</text><view class="plant-mark"></view></view>
                            </view>
                        </view>
                    </view>
                    <view v-else class="empty-state">
                        <image src="/static/icon/c2m.svg" mode="aspectFit"></image>
                        <text>这个标签下还没有植物</text>
                    </view>
                    <view v-if="plants.length" class="poster-footer"><view class="footer-line"></view><text>每一株植物，都有自己的生长故事</text><view class="footer-line"></view></view>
                </view>
            </scroll-view>

            <view class="preview-actions">
                <button class="outline-button" :loading="generating" @click="savePoster">保存到相册</button>
                <button class="share-button" open-type="share">分享给朋友</button>
            </view>
        </view>

        <canvas canvas-id="gardenPoster" class="poster-canvas"
            :style="{ width: canvasWidth + 'px', height: canvasHeight + 'px' }"></canvas>
    </view>
</template>

<script>
import navBar from '@/components/navBar.vue'
import { callContainer } from '../../utils/request.js'

export default {
    components: { navBar },
    data() {
        return {
            topBarHeight: 0,
            familyId: 0,
            tagId: 0,
            tagName: '全部',
            familyName: '我的花园',
            plants: [],
            posterPath: '',
            generating: false,
            canvasWidth: 690,
            canvasHeight: 1200
        }
    },
    computed: {
        displayTagName() {
            return this.tagName && this.tagName !== '全部' ? this.tagName : '花园'
        },
        shareTitle() {
            return this.tagName && this.tagName !== '全部'
                ? `${this.familyName}的${this.displayTagName}`
                : this.familyName
        },
        previewColumns() {
            const columns = [[], []]
            this.plants.forEach((plant, index) => columns[index % 2].push(plant))
            return columns
        }
    },
    async onLoad(options) {
        this.topBarHeight = getApp().globalData.topBarHeight || 44
        this.familyId = Number(options.familyId || uni.getStorageSync('familyId') || 0)
        this.tagId = Number(options.tagId || 0)
        this.tagName = options.tagName ? decodeURIComponent(options.tagName) : '全部'

        const families = uni.getStorageSync('family') || []
        const currentFamily = families.find(item => String(item.id) === String(this.familyId))
        this.familyName = currentFamily?.name || '我的花园'
        await this.loadPlants()
    },
    methods: {
        async loadPlants() {
            try {
                const result = await callContainer('/api/plant/list', {
                    familyId: this.familyId,
                    page: 1,
                    pageSize: 100,
                    tagId: this.tagId,
                    keyword: ''
                })
                this.plants = result?.data?.list || result?.data || []
            } catch (error) {
                console.error('加载分享植物失败:', error)
                uni.showToast({ title: '植物加载失败', icon: 'none' })
            }
        },
        getLocalImage(src) {
            return new Promise(resolve => uni.getImageInfo({
                src,
                success: res => resolve({ path: res.path, width: res.width, height: res.height }),
                fail: () => resolve(null)
            }))
        },
        drawRoundedRect(ctx, x, y, width, height, radius, color) {
            this.roundedRectPath(ctx, x, y, width, height, radius)
            ctx.setFillStyle(color)
            ctx.fill()
        },
        roundedRectPath(ctx, x, y, width, height, radius) {
            const r = Math.min(radius, width / 2, height / 2)
            ctx.beginPath()
            ctx.moveTo(x + r, y)
            ctx.lineTo(x + width - r, y)
            ctx.arcTo(x + width, y, x + width, y + r, r)
            ctx.lineTo(x + width, y + height - r)
            ctx.arcTo(x + width, y + height, x + width - r, y + height, r)
            ctx.lineTo(x + r, y + height)
            ctx.arcTo(x, y + height, x, y + height - r, r)
            ctx.lineTo(x, y + r)
            ctx.arcTo(x, y, x + r, y, r)
            ctx.closePath()
        },
        async createPoster() {
            if (!this.plants.length) throw new Error('当前标签下没有植物')
            this.generating = true

            try {
                const images = await Promise.all(this.plants.map(plant =>
                    this.getLocalImage(plant.cover?.url || '/static/default.svg')
                ))
                const cardWidth = 303
                const gap = 18
                const side = 33
                const titleHeight = 190
                const nameHeight = 66
                const columnY = [titleHeight, titleHeight]
                const cards = []

                this.plants.forEach((plant, index) => {
                    const image = images[index]
                    const ratio = image?.width && image?.height ? image.height / image.width : 1
                    const imageHeight = Math.max(150, Math.min(455, Math.round(cardWidth * ratio)))
                    const column = columnY[0] <= columnY[1] ? 0 : 1
                    const x = side + column * (cardWidth + gap)
                    const y = columnY[column]
                    cards.push({ plant, image, x, y, imageHeight })
                    columnY[column] += imageHeight + nameHeight + gap
                })

                this.canvasHeight = Math.max(520, Math.max(...columnY) + 78)
                await this.$nextTick()

                const colors = this.activeTheme.vars
                const ctx = uni.createCanvasContext('gardenPoster', this)
                ctx.setFillStyle(colors['--surface-bg'])
                ctx.fillRect(0, 0, this.canvasWidth, this.canvasHeight)
                ctx.setFillStyle(colors['--page-bg'])
                ctx.fillRect(0, 0, this.canvasWidth, this.canvasHeight)
                ctx.setTextAlign('center')
                ctx.setFillStyle(colors['--primary-strong'])
                ctx.setFontSize(16)
                ctx.fillText('—  M Y  G A R D E N  —', this.canvasWidth / 2, 55)
                ctx.setFillStyle(colors['--text-color'])
                ctx.setFontSize(38)
                ctx.fillText(this.shareTitle, this.canvasWidth / 2, 103)
                ctx.setFillStyle(colors['--text-sub'])
                ctx.setFontSize(19)
                ctx.fillText(`${this.displayTagName}收藏 · ${this.plants.length}株植物`, this.canvasWidth / 2, 137)
                ctx.setStrokeStyle(colors['--border-color'])
                ctx.setLineWidth(2)
                ctx.beginPath()
                ctx.moveTo(250, 162)
                ctx.lineTo(322, 162)
                ctx.moveTo(368, 162)
                ctx.lineTo(440, 162)
                ctx.stroke()
                ctx.setFillStyle(colors['--primary-color'])
                ctx.beginPath()
                ctx.arc(345, 162, 6, 0, Math.PI * 2)
                ctx.fill()

                cards.forEach(card => {
                    const cardHeight = card.imageHeight + nameHeight
                    this.drawRoundedRect(ctx, card.x, card.y + 5, cardWidth, cardHeight, 20, colors['--border-color'])
                    this.drawRoundedRect(ctx, card.x, card.y, cardWidth, cardHeight, 20, colors['--card-bg'])
                    if (card.image) {
                        ctx.save()
                        this.roundedRectPath(ctx, card.x, card.y, cardWidth, card.imageHeight + 20, 20)
                        ctx.clip()
                        ctx.drawImage(card.image.path, card.x, card.y, cardWidth, card.imageHeight)
                        ctx.restore()
                    }
                    ctx.setTextAlign('left')
                    ctx.setFillStyle(colors['--text-color'])
                    ctx.setFontSize(23)
                    const name = String(card.plant.name || '未命名植物')
                    ctx.fillText(name.length > 12 ? `${name.slice(0, 12)}…` : name, card.x + 16, card.y + card.imageHeight + 40)
                    ctx.setFillStyle(colors['--primary-color'])
                    ctx.beginPath()
                    ctx.arc(card.x + cardWidth - 20, card.y + card.imageHeight + 33, 4, 0, Math.PI * 2)
                    ctx.fill()
                })

                ctx.setTextAlign('center')
                ctx.setFillStyle(colors['--text-sub'])
                ctx.setFontSize(16)
                ctx.fillText('每一株植物，都有自己的生长故事', this.canvasWidth / 2, this.canvasHeight - 43)

                await new Promise(resolve => ctx.draw(false, resolve))
                this.posterPath = await new Promise((resolve, reject) => uni.canvasToTempFilePath({
                    canvasId: 'gardenPoster',
                    width: this.canvasWidth,
                    height: this.canvasHeight,
                    destWidth: this.canvasWidth,
                    destHeight: this.canvasHeight,
                    fileType: 'jpg',
                    quality: .92,
                    success: res => resolve(res.tempFilePath),
                    fail: reject
                }, this))
                return this.posterPath
            } finally {
                this.generating = false
            }
        },
        async savePoster() {
            if (!this.plants.length) return uni.showToast({ title: '当前标签下没有植物', icon: 'none' })
            try {
                const path = await this.createPoster()
                await new Promise((resolve, reject) => uni.saveImageToPhotosAlbum({
                    filePath: path,
                    success: resolve,
                    fail: reject
                }))
                uni.showToast({ title: '已保存到相册', icon: 'success' })
            } catch (error) {
                console.error('保存分享长图失败:', error)
                uni.showToast({ title: '保存失败，请检查相册权限', icon: 'none' })
            }
        }
    },
    onShareAppMessage() {
        return {
            title: this.shareTitle,
            path: `/pages/sharePlantList/sharePlantList?familyId=${this.familyId}&tagId=${this.tagId}&tagName=${encodeURIComponent(this.tagName)}`,
            imageUrl: this.posterPath || this.plants[0]?.cover?.url || '/static/images/share-garden-decorated.jpg'
        }
    }
}
</script>

<style scoped lang="scss">
.share-page { height: 100vh; overflow: hidden; color: var(--text-color); }
.preview-layout { display: flex; flex-direction: column; height: calc(100vh - 70px); }
.preview-scroll { flex: 1; height: 0; padding: 0 0 18px; background: var(--page-bg); }
.poster-preview { width: 100%; min-height: 300px; box-sizing: border-box; padding: 0 16px 18px; background: var(--page-bg); }
.poster-header { display: flex; flex-direction: column; align-items: center; padding: 25px 8px 21px; }
.poster-eyebrow { display: flex; align-items: center; justify-content: center; gap: 9px; color: var(--primary-strong); font-size: 9px; font-weight: 600; letter-spacing: 2px; }
.eyebrow-line { width: 31px; height: 1px; background: var(--border-color); }
.poster-title { display: block; max-width: 100%; padding-top: 9px; overflow: hidden; text-align: center; color: var(--text-color); font-size: 23px; line-height: 1.35; font-weight: 700; white-space: nowrap; text-overflow: ellipsis; }
.poster-meta { margin-top: 5px; color: var(--text-sub); font-size: 11px; }
.poster-leaf-mark { display: flex; align-items: center; gap: 3px; margin-top: 13px; }
.leaf { width: 8px; height: 5px; background: var(--primary-soft); }
.leaf-left { border-radius: 8px 1px 8px 1px; transform: rotate(18deg); }
.leaf-right { border-radius: 1px 8px 1px 8px; transform: rotate(-18deg); }
.leaf-dot { width: 5px; height: 5px; border-radius: 50%; background: var(--primary-color); }
.poster-columns { display: flex; align-items: flex-start; gap: 9px; }
.poster-column { flex: 1; width: 0; }
.poster-plant { overflow: hidden; margin-bottom: 9px; border: 1px solid var(--border-color); border-radius: 14px; background: var(--card-bg); box-shadow: 0 4px 12px var(--shadow-color); }
.poster-plant image { display: block; width: 100%; background: var(--surface-bg); }
.poster-plant-info { display: flex; align-items: center; gap: 6px; min-height: 41px; padding: 8px 9px; }
.poster-plant-name { display: block; flex: 1; min-width: 0; overflow: hidden; font-size: 12px; line-height: 1.3; font-weight: 500; white-space: nowrap; text-overflow: ellipsis; }
.plant-mark { flex: 0 0 auto; width: 5px; height: 5px; border-radius: 50%; background: var(--primary-color); }
.poster-footer { display: flex; align-items: center; justify-content: center; gap: 9px; padding: 15px 3px 2px; color: var(--text-sub); font-size: 9px; }
.footer-line { flex: 1; height: 1px; background: var(--border-color); }
.empty-state { display: flex; flex-direction: column; align-items: center; gap: 10px; padding: 70px 20px; color: var(--text-sub); font-size: 13px; }
.empty-state image { width: 76px; height: 76px; opacity: .7; }
.preview-actions { display: flex; gap: 10px; padding: 11px 16px calc(11px + env(safe-area-inset-bottom)); border-top: 1px solid var(--border-color); background: var(--card-bg); box-shadow: 0 -5px 18px var(--shadow-color); }
.outline-button, .share-button { flex: 1; height: 44px; margin: 0; line-height: 44px; border-radius: 999px; font-size: 14px; font-weight: 500; }
.outline-button { color: var(--primary-strong); background: var(--card-bg); border: 1px solid var(--primary-color); }
.share-button { color: var(--button-text); background: var(--primary-color); }
.outline-button::after, .share-button::after { border: 0; }
.poster-canvas { position: fixed; left: -10000px; top: -10000px; pointer-events: none; }
</style>
