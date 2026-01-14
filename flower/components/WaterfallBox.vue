<template>
    <view class="waterfall-container">
        <!-- 
      核心修改：
      使用嵌套循环。外层循环列(columns)，内层循环数据(col)。
      这样 <slot> 标签在代码里只出现了一次，避开了微信小程序的编译 Bug。
    -->
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
    if (newVal.length > (oldVal?.length || 0)) {
        const newItems = newVal.slice(oldVal?.length || 0);
        tempQueue.value.push(...newItems);
    } else {
        // 重置：清空每一列
        columns.value = Array.from({ length: props.cols }, () => []);
        tempQueue.value = [...newVal];
    }

    if (!isRendering.value) {
        nextTick(() => {
            renderNext();
        });
    }
}, {  immediate: true });

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
    // 动态计算宽度：(100% - 间隙) / 列数
    // 这里简化处理，假设是双列，留一点间隙
    width: 48%;
    flex: 0 0 48%;
    max-width: 48%;
    // max-height: 100px;
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
