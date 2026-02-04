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
import { ref, watch, nextTick, getCurrentInstance } from 'vue';

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

// columns 是一个二维数组，例如 [[item1, item3], [item2, item4]]
const columns = ref(Array.from({ length: props.cols }, () => []));
const tempQueue = ref([]);
const isRendering = ref(false);

const instance = getCurrentInstance();

watch(() => props.list, (newVal, oldVal) => {
    if (!newVal || newVal.length === 0) {
        columns.value = Array.from({ length: props.cols }, () => []);
        tempQueue.value = [];
        return;
    }

    // 🚀 优化：判断是“追加”还是“重置”
    // 如果新列表的前半部分和旧列表一致，说明是加载更多
    const isAppend = oldVal && oldVal.length > 0 && newVal.length > oldVal.length && 
                   newVal.slice(0, oldVal.length).every((item, i) => item[props.idKey] === oldVal[i][props.idKey]);

    if (isAppend) {
        console.log("WaterfallBox - 检测到追加数据，增量渲染");
        const newItems = newVal.slice(oldVal.length);
        tempQueue.value.push(...newItems);
    } else {
        console.log("WaterfallBox - 检测到重置数据，全量重绘");
        columns.value = Array.from({ length: props.cols }, () => []);
        tempQueue.value = [...newVal];
    }

    if (!isRendering.value) {
        nextTick(() => {
            renderNext();
        });
    }
}, { immediate: true, deep: false });

const renderNext = async () => {
    if (tempQueue.value.length === 0) {
        isRendering.value = false;
        return;
    }

    isRendering.value = true;
    
    // 🚀 优化：一次处理 1 个，但大幅减小延迟
    const item = tempQueue.value.shift();

    // 寻找高度最小的列
    let minHeight = Infinity;
    let minColIndex = 0;

    try {
        // 遍历所有列获取高度
        for (let i = 0; i < props.cols; i++) {
            const height = await getContainerHeight(`#waterfall-column-${i}`);
            if (height < minHeight) {
                minHeight = height;
                minColIndex = i;
            }
        }

        // 将数据加入最短的那一列
        columns.value[minColIndex].push(item);

        // 🚀 使用 nextTick 代替 setTimeout，显著提升渲染速度
        // nextTick 保证 DOM 已更新，此时下一次 getContainerHeight 能拿到正确高度
        await nextTick();
        renderNext();
    } catch (e) {
        console.error("WaterfallBox render error:", e);
        isRendering.value = false;
    }
};

const getContainerHeight = (selector) => {
    return new Promise((resolve) => {
        if (!instance) {
            resolve(0);
            return;
        }
        const query = uni.createSelectorQuery().in(instance.proxy || instance);
        query.select(selector).boundingClientRect((res) => {
            resolve(res ? res.height : 0);
        }).exec();
    });
};
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
