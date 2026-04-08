<template>
	<div class="img-preview">
		<el-image style="width: 100px; height: 100px; border-radius: 5px" :src="getDefaultImage()" :preview-src-list="imageList" title="点击查看大图预览" />
	</div>
</template>
   
<script lang="ts">
import { defineComponent } from 'vue';

	const props = {
		imageList: { type: Object, default: () => [] },
	};

	export default defineComponent({
		name: 'q-img-preview',
		props: props,
		setup(props) {

			const getDefaultImage = () => {
				if(props.imageList.length > 0) return props.imageList[0]
				else return '';
			};

			return {
				getDefaultImage,
			};
		},
	});
</script>