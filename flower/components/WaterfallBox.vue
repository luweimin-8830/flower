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
    console.log("WaterfallBox - watch triggered");
    console.log("WaterfallBox - newVal length:", newVal?.length, "oldVal length:", oldVal?.length);

    // 检查是否需要重新渲染
    let needRebuild = false;

    // 1. 如果是新数据或数量差异较大，需要重新构建
    if (!oldVal || newVal.length !== oldVal.length) {
        needRebuild = true;
        console.log("WaterfallBox - 需要重新构建（长度不同）");
    }
    // 2. 如果数量相同，检查每个元素的 ID 是否相同
    else {
        for (let i = 0; i < newVal.length; i++) {
            if (newVal[i][props.idKey] !== oldVal[i][props.idKey]) {
                needRebuild = true;
                console.log(`WaterfallBox - 需要重新构建（ID不同）: ${newVal[i][props.idKey]} !== ${oldVal[i][props.idKey]}`);
                break;
            }
        }
        if (!needRebuild) {
            console.log("WaterfallBox - 数据完全相同，不需要重新构建");
        }
    }

    if (needRebuild) {
        console.log("WaterfallBox - 开始重新构建列");
        columns.value = Array.from({ length: props.cols }, () => []);
        tempQueue.value = [...newVal];
    }

    if (!isRendering.value) {
        nextTick(() => {
            renderNext();
        });
    }
}, { immediate: true });

const renderNext = async () => {
    if (tempQueue.value.length === 0) {
        isRendering.value = false;
        return;
    }

    isRendering.value = true;
    const item = tempQueue.value.shift();

    // 寻找高度最小的列
    let minHeight = Infinity;
    let minColIndex = 0;

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

    // 等待渲染
    setTimeout(() => {
        renderNext();
    }, props.gap);
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
