<template>
	<el-row :gutter="15">
		<el-col :lg="24">
			<el-card shadow="never" class="aboutTop">
				<div class="aboutTop-info">
					<el-avatar :size="150" :src="logo"></el-avatar>
					<h2>{{ config.APP_NAME }}</h2>
					<p>ver: {{ config.APP_VER }}</p>
				</div>
			</el-card>
<!--			<el-card shadow="never" header="dependencies">-->
<!--				<el-descriptions size="small" border :column="3">-->
<!--					<el-descriptions-item v-for="(value, key) in packageJson.dependencies" :key="key" :label="key">{{value}}</el-descriptions-item>-->
<!--				</el-descriptions>-->
<!--			</el-card>-->
<!--			<el-card shadow="never" header="devDependencies">-->
<!--				<el-descriptions size="small" border :column="3">-->
<!--					<el-descriptions-item v-for="(value, key) in packageJson.devDependencies" :key="key" :label="key">{{value}}</el-descriptions-item>-->
<!--				</el-descriptions>-->
<!--			</el-card>-->
		</el-col>
	</el-row>
</template>

<script lang="ts">
	import { defineComponent, reactive, toRefs } from 'vue';
	import logo from '/@/assets/logo-mini.svg';
	import packageJson from '/root/package.json';
	import config from '/@/config';

	export default defineComponent({
		name: "personal-about",
		setup(){
			const state = reactive({
				logo,
				packageJson,
			});

			return {
				config,
				...toRefs(state),
			};
		},
	});

</script>

<style scoped lang="scss">
	.aboutTop {
		border:0;
		//background: linear-gradient(to bottom, var(--el-color-primary), var(--el-color-primary-light-2));color: #fff;
		padding: 50px 0;
		&-info {
			text-align: center;

			img {
				width: 100px;
			}
			h2 {
				font-size: 26px;
				margin-top: 15px;
			}
			p {
				font-size: 16px;
				margin-top: 10px;
			}
		}
	}
	.el-card {
		margin-bottom: 15px;
	}
</style>

