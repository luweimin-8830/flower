<template>
    <!-- 家庭选择框 -->
    <view>
        <uni-data-select v-model="value" :localdata="range" @change="change"></uni-data-select>
    </view>
    <view>

    </view>
</template>

<script>
export default {
    /**
     * 组件名称，也就是开发者使用的标签
     */
    name: 'home',
    /**
     * 组件涉及的事件声明，只有声明过的事件，才能被正常发送
     */
    emits: [],
    /**
     * 属性声明，组件的使用者会传递这些属性值到组件
     */
    props: {
        navText: {
            type: String,
            default: "",
        },
    },
    /**
     * 组件内部变量声明
     */
    data() {
        return {
            
        }
    },
    /**
     * 属性变化监听器实现
     */
    watch: {

    },
    /**
     * 规则：如果没有配置expose，则methods中的方法均对外暴露，如果配置了expose，则以expose的配置为准向外暴露
     * ['publicMethod'] 含义为：只有 `publicMethod` 在实例上可用
     * 
     * 注意：如果在data中声明了一个变量，此时组件配置了 expose字段，但未在expose字段中包含此变量。会导致该变量被标记为`private`：仅能在组件内使用，不能在组件外访问
     */
    //expose: [''],
    methods: {
        
        /**
      * 内部使用的组件方法
      */
        //privateMethod() {}
    },
    /**
     * [可选实现] 组件被创建，组件第一个生命周期，
     * 在内存中被占用的时候被调用，开发者可以在这里执行一些需要提前执行的初始化逻辑
     */
    created() {
        
    },
}
</script>

<style scoped lang="scss">
    .nav-back-button{
        position: fixed;
        z-index: 999;
        opacity:0.55;
        background-color: #fff;
        // backdrop-filter: blur(10px);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        border: 1px solid rgba(0,0,0,0.08);
        box-shadow: 0 2px 8px rgba(0,0,0,0.15);
        transition: transform 0.1s,background-color 0.2s;
    }
    .nav-back-icon{
        font-size: 20px;
        // color: var(--text-color,#333);
        color:#333
    
    }
    .nav-back-button-hover {
        transform: scale(0.92);
        background-color: #f5f5f5; 
        opacity:0.65;
        box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
    }
    @media (prefers-color-scheme: dark) {
        .nav-back-button {
            background-color: rgba(30, 30, 30, 0.6);
            border: 1px solid rgba(255, 255, 255, 0.1);
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
        }
        .nav-back-button-hover {
            background-color: rgba(50, 50, 50, 0.8);
        }
    }
</style>