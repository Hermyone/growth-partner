<!--
 * @Descripttion: 自定义表格搜索条件
 * @version: 1.0
 * @Author: OTQ
 * @Date: 2023-06-08
 * @LastEditors:
 * @LastEditTime:
-->

<template>
	<div class="table-search-container" v-if="search.length > 0">
		<el-form ref="tableSearchRef" :model="form" size="default" label-width="100px" class="table-form">
			<el-row>
				<el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="mb10" v-for="(val, key) in search" :key="key" v-show="key <= showFormNum || isToggle">
					<template v-if="val.type !== ''">
						<el-form-item
								:label="val.label"
								:prop="val.prop"
								:rules="[{ required: val.required, message: `${val.label}不能为空`, trigger: val.type === 'input' ? 'blur' : 'change' }]"
								class="w100"
						>
							<el-input
									v-model="form[val.prop]"
									:placeholder="val.placeholder"
									clearable
									v-if="val.type === 'input'"
									:style="{width: val.width ? val.width : `100%`}"
									@change="val.change?onSearch(tableSearchRef):null"
									@clear="val.clear?onReset(tableSearchRef):null"
							/>
							<el-date-picker
									v-model="form[val.prop]"
									:type="val.kind !=='' ? val.kind : 'date'"
									:placeholder="val.placeholder"
									v-else-if="val.type === 'date'"
									:value-format="val.format"
									:style="{width: val.width ? val.width : `100%`}"
							/>
							<el-date-picker
									v-model="form[val.prop]"
									type="daterange"
									unlink-panels
									:placeholder="val.placeholder"
									v-else-if="val.type === 'daterange'"
									range-separator="-"
									:value-format="val.format"
									:shortcuts="shortcuts"
									start-placeholder="开始日期"
									end-placeholder="结束日期"
									:style="{width: val.width ? val.width : `100%`}"
							/>
							<el-select v-model="form[val.prop]" :placeholder="val.placeholder" v-else-if="val.type === 'select'" :style="{width: val.width ? val.width : `100%`}" clearable>
								<el-option v-for="item in val.options" :key="item.value" :label="item.label" :value="item.value"> </el-option>
							</el-select>
						</el-form-item>
					</template>
				</el-col>
        <slot name="expandWhere"></slot>
				<el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="mb10">
					<el-form-item class="table-form-btn" :label-width="search.length <= showFormNum ? '10px' : '90px'">
						<template #label v-if="search.length > showFormNum+1">
							<div class="table-form-btn-toggle ml10" @click="isToggle = !isToggle">
								<el-button type="primary" link size="small">
									{{ wordSearch }}
									<el-icon v-show="isToggle"><ele-ArrowUp /></el-icon>
									<el-icon v-show="!isToggle"><ele-ArrowDown /></el-icon>
								</el-button>
							</div>
						</template>
						<div style="min-width: 200px" >
							<el-button size="default" type="primary" round class="ml10" icon="ele-Search" @click="onSearch(tableSearchRef)">查询</el-button>
							<el-button size="default" round @click="onReset(tableSearchRef)" icon="ele-Refresh" >重置</el-button>
						</div>
					</el-form-item>
				</el-col>
				<div :style="{'min-width': `${buttonMinWidth}px`, 'padding-left': '25px' }" class="mb10">
					<slot></slot>
				</div>
			</el-row>
		</el-form>
	</div>
</template>

<script lang="ts">
	import { reactive, toRefs, ref, defineComponent, computed, onMounted } from 'vue';
	import type { FormInstance } from 'element-plus';

	// 定义父组件传过来的值
	const props = {
		// 搜索表单
		search: {
			type: Array<TableSearchType>,
			default: () => [],
		},
		showFormNum: {
			type: Number,
			default: 1,
		},
		buttonMinWidth: {
			type: Number,
			default: () => 200,
		}
	};

	export default defineComponent({
		name: 'q-search',
		props: props,
		setup(props, { emit }) {
			// 定义变量内容
			const tableSearchRef = ref<FormInstance>();
			const state = reactive({
				form: [] as any,
				isToggle: false,

				shortcuts: [
					{
						text: '最近一周',
						value: () => {
							const end = new Date()
							const start = new Date()
							start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
							return [start, end]
						},
					},
					{
						text: '最近一月',
						value: () => {
							const end = new Date()
							const start = new Date()
							start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
							return [start, end]
						},
					},
					{
						text: '最近3个月',
						value: () => {
							const end = new Date()
							const start = new Date()
							start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
							return [start, end]
						},
					},
				]
			});

			const wordSearch = computed(() => {
				if (!state.isToggle) {
					return '展开搜索';
				} else {
					return '收起搜索';
				}
			});

			// 查询
			const onSearch = (formEl: FormInstance | undefined) => {
				if (!formEl) return;
				formEl.validate((valid: boolean) => {
					if (valid) {
						emit('search', state.form);
					} else {
						return false;
					}
				});
			};
			// 重置
			const onReset = (formEl: FormInstance | undefined) => {
				if (!formEl) return;
				formEl.resetFields();
				emit('search', state.form, true);
			};
			// 初始化 form 字段，取自父组件 search.prop
			const initFormField = () => {
				if (props.search.length <= 0) return false;
				props.search.forEach((v: any) => (state.form[v.prop] = ''));
			};
			// 根据ID设置默认值
			const setDefaultValue = (prop: string, val: any) => {
				props.search.forEach((v: any) => {
					if(prop === v.prop)
						state.form[v.prop] = val;
				});
			};

			// 页面加载时
			onMounted(() => {
				initFormField();
			});

			return {
				tableSearchRef,
				wordSearch,
				onSearch,
				onReset,
				setDefaultValue,
				...toRefs(state),
			}
		}
	});
</script>

<style scoped lang="scss">
	.table-search-container {
		display: flex;
		.table-form {
			flex: 1;
			padding-left: 0 !important;
			.table-form-btn-toggle {
				white-space: nowrap;
				user-select: none;
				display: flex;
				align-items: center;
				height:32px;
				color: var(--el-color-primary);
			}
		}
	}
</style>
