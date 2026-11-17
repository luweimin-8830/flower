<template>
    <view class="lazy-plant-image">
        <image v-for="picture in pictures" :key="picture.key" :src="picture.src"
            :data-request="picture.key" mode="aspectFill"
            class="picture" :class="{ ready: status === 'loaded' }"
            @load="onLoad" @error="onError" />
        <view v-if="status !== 'loaded'" class="placeholder">
            <text v-if="status === 'error'" class="retry" @click.stop="retry">加载失败，点击重试</text>
            <text v-else>图片加载中</text>
        </view>
    </view>
</template>

<script>
import { enqueueImage } from '../utils/imageLoadQueue';
import { thumbnailWidth, thumbnailSource } from '../utils/imageThumbnail';

export default {
    props: {
        src: { type: String, default: '/static/default.svg' },
        originalWidth: Number,
        originalHeight: Number
    },
    data() {
        return { source: '', status: 'idle', generation: 0 };
    },
    computed: {
        pictures() {
            return this.source ? [{ key: `${this.generation}:${this.source}`, src: this.source }] : [];
        }
    },
    watch: {
        src() {
            this.reset();
            if (this._nearby) this.request();
        }
    },
    mounted() {
        this._alive = true;
        this.$nextTick(() => {
            if (!this._alive) return;
            // Keep observing queued images so fast scrolling removes obsolete work.
            this._observer = uni.createIntersectionObserver(this, { thresholds: [0] });
            this._observer.relativeToViewport({ top: 160, bottom: 160 })
                .observe('.lazy-plant-image', result => {
                    this._box = result.boundingClientRect;
                    this._nearby = result.intersectionRatio > 0;
                    if (this._nearby) this.request();
                    else if (this.status === 'queued') this.reset();
                });
        });
    },
    beforeUnmount() {
        this._alive = false;
        if (this._observer) this._observer.disconnect();
        this.reset();
    },
    methods: {
        release() {
            clearTimeout(this._timer);
            const release = this._release;
            this._release = null;
            if (release) release();
        },
        reset() {
            this.generation++;
            this._originalFallback = false;
            this.source = '';
            this.status = 'idle';
            this.release();
        },
        request() {
            if (!this._alive || this.status !== 'idle') return;
            this.status = 'queued';
            this._release = enqueueImage(async () => {
                this.status = 'loading';
                const generation = this.generation;
                const original = this.src || '/static/default.svg';
                this._timer = setTimeout(() => this.fail(), 15000);
                try {
                    const info = uni.getWindowInfo ? uni.getWindowInfo() : uni.getSystemInfoSync();
                    const width = thumbnailWidth(this._box, info.pixelRatio, this.originalWidth, this.originalHeight);
                    const url = await thumbnailSource(original, width, async fileID => {
                        const result = await wx.cloud.getTempFileURL({ fileList: [fileID] });
                        const file = result.fileList && result.fileList[0];
                        if (!file || !file.tempFileURL || file.status) throw new Error('Image URL unavailable');
                        return file.tempFileURL;
                    });
                    if (this._alive && generation === this.generation && this.status === 'loading') this.source = url;
                } catch (error) {
                    if (this._alive && generation === this.generation && this.status === 'loading') {
                        this._originalFallback = true;
                        this.source = original;
                    }
                }
            });
        },
        isCurrentEvent(event) {
            return event?.currentTarget?.dataset?.request === `${this.generation}:${this.source}`;
        },
        onLoad(event) {
            if (this.status !== 'loading' || !this.isCurrentEvent(event)) return;
            this.status = 'loaded';
            this.release();
        },
        onError(event) {
            if (this.status !== 'loading' || !this.isCurrentEvent(event)) return;
            // A bucket without image processing must still display its original.
            if (!this._originalFallback && this.source !== this.src) {
                this._originalFallback = true;
                this.source = this.src || '/static/default.svg';
                clearTimeout(this._timer);
                this._timer = setTimeout(() => this.fail(), 15000);
                return;
            }
            this.fail();
        },
        fail() {
            if (this.status !== 'loading') return;
            this.source = '';
            this.status = 'error';
            this.release();
        },
        retry() {
            this.reset();
            this.request();
        }
    }
};
</script>

<style scoped>
.lazy-plant-image { position: relative; width: 100%; height: 100%; overflow: hidden; background: #edf0e9; }
.picture, .placeholder { position: absolute; top: 0; left: 0; width: 100%; height: 100%; }
.picture { opacity: 0; transition: opacity .2s; }
.picture.ready { opacity: 1; }
.placeholder { display: flex; align-items: center; justify-content: center; color: #8a9484; font-size: 11px; }
.retry { padding: 12px 4px; text-align: center; }
</style>
