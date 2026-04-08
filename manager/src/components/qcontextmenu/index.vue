<!--
 * @Descripttion: 自定义上下文菜单
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-04-24
 * @LastEditors:
 * @LastEditTime:
-->

<template>
    <div class="contextmenu-style">
    </div>
</template>

<script lang="ts">
	import { defineComponent } from 'vue';

	const props = {
        title: { type: String, default: "" },
	};

	export default defineComponent({
		name: 'q-contextmenu',
		props: props,

		setup(props) {
        },
	});
</script>

<style scoped lang="scss">
    .contextmenu-style {

    }
</style>
