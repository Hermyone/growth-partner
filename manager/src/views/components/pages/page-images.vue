<template>
    <el-image-viewer v-if="image.visible"
                     :url-list="[image.url]"
                     hide-on-click-modal
                     @close="image.visible = false">
    </el-image-viewer>
</template>

<script lang="ts">
	import { defineComponent, reactive, toRefs } from 'vue';

	export default defineComponent({
		name: 'q-page-image-viewer',
		setup() {
			const state = reactive({
                image: {
                    visible: false,
                    url: ''
                },
			});

            const openViewer = (url: string) => {
                state.image.url = url;
                state.image.visible = true;
            };

			return {
                openViewer,
				...toRefs(state),
			}
		},
	});

</script>

<style scoped lang="scss">
</style>