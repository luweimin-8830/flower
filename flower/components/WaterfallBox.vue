<template>
    <view class="waterfall-container">
        <view v-for="(col, colIndex) in columns" :key="colIndex" class="waterfall-column"
            :id="'waterfall-column-' + colIndex">
            <view v-for="(item, index) in col" :key="item[idKey] || index" class="waterfall-item">
                <!-- 插槽只写在这里一次 -->
                <slot name="item" :item="item" :index="index"></slot>
            </view>
        </view>
    </view>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
    list: {
        type: Array,
        required: true,
        default: () => []
    },
    idKey: {
        type: String,
        default: 'id'
    },
    gap: {
        type: Number,
        default: 50
    },
    cols: {
        type: Number,
        default: 2
    }
});

// 微信小程序生产包中，逐项查询节点高度可能在组件首次渲染时拿不到节点，
// 导致队列停止且整页空白。按顺序均匀分列无需查询 DOM，分页追加也会立即更新。
const columns = computed(() => {
    const count = Math.max(1, Number(props.cols) || 1)
    const result = Array.from({ length: count }, () => [])
    ;(props.list || []).forEach((item, index) => {
        result[index % count].push(item)
    })
    return result
})
</script>

<style scoped lang="scss">
.waterfall-container {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    width: 100%;
    padding: 10px;
    box-sizing: border-box;
}

.waterfall-column {
    display: flex;
    flex-direction: column;
    // 简化计算：直接使用 50% 宽度，确保每列均匀分布
    width: 49%;  // 留一点间隙
    flex: 0 0 49%;
    max-width: 49%;
    margin: 0 0.5%;  // 添加列之间的间距
}

@keyframes fadeInUp {
    0% {
        opacity: 0;
        transform: translateY(40px);
        /* 初始位置向下偏移 */
    }

    100% {
        opacity: 1;
        transform: translateY(0);
    }
}

.waterfall-item {
    width: 100%;
    margin-bottom: 10px;
    // background-color: rgba(240,240,240,0.3);
    border-radius: 8px;
    overflow: hidden;
    box-sizing: border-box;
    animation: fadeInUp 0.3s ease-out forwards;
    opacity: 0; 
}
</style>
