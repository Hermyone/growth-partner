<template>
	<div ref="backRef" class="pdf-preview" :style="{ height: `${clientHeight}px` }">
		<iframe :src="url" frameborder="0" style="width: 100%; height: 100%;"></iframe>
	</div>
</template>
   
<script lang="ts">
	import { defineComponent, ref, onMounted, onUnmounted, reactive, toRefs } from 'vue';

	export default defineComponent({
		name: 'q-pdf-preview',
		setup(props) {
			const url = ref("");
			const state = reactive({
				clientHeight: 100,
			});

			const open = (address: string) => {
				url.value = address;
			};

			const onLayoutResize = () => {
				state.clientHeight = document.body.clientHeight;
			};

			onMounted(() => {
				onLayoutResize();
			})

			onUnmounted(() => {
				window.removeEventListener('resize', onLayoutResize);
			})

			return {
				url,
				open,
				...toRefs(state),
			};
		},
	});
</script>

<style scoped lang="scss">
	.pdf-preview {
		padding: 0!important;
		width: 100%;
	}
</style>