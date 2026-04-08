<!--
 * @Descripttion: 自适应屏幕高度容器
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-08-15
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <div class="back-container" :style="{ height: `${clientHeight}px` }">
        <slot></slot>
    </div>
</template>

<script lang="ts">
	import { defineComponent, watch, toRefs, onBeforeMount, onUnmounted, onMounted, onActivated, reactive } from 'vue';

	const props = {
        offset: {
            type: Number,
            default: 0,
        },
        offsetClass: {
            type: String,
            default: '',
        }
	};

	export default defineComponent({
		name: 'q-container-size',
		props: props,
		setup(props) {
            const state = reactive({
                clientHeight: 100,
            });

            const onLayoutResize = () => {
            	var h : any = 0;
            	if(props.offsetClass !== ''){
            		h = document.querySelector('.' + props.offsetClass)?.clientHeight;
            		h = h + 30;
                }
                state.clientHeight = document.body.clientHeight - props.offset - h;
            };

            // 页面加载前
            onBeforeMount(() => {
                onLayoutResize();
                window.addEventListener('resize', onLayoutResize);
            });
            // 页面卸载时
            onUnmounted(() => {
                window.removeEventListener('resize', onLayoutResize);
            });

            onMounted(() => {
                onLayoutResize();
            });

            onActivated(() => {
                onLayoutResize();
            });

            watch(
                () => props.offset,
                () => {
                    onLayoutResize();
                },
            );

            return {
                onLayoutResize,
                ...toRefs(state)
            }
		},
	});
</script>

<style scoped lang="scss">
    .back-container {
        padding: 0!important;
    }
</style>