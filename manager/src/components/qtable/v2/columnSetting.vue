<template>
	<div v-if="usercolumn.length>0" class="setting-column" v-loading="isSave">
		<div class="setting-column__title">
			<slot></slot>
		</div>
		<div class="setting-column__title">
			<span class="f_move">
				<el-tooltip content="拖动进行排序" placement="top-start">
					<SvgIcon name="fa fa-question-circle-o" :size="14" color="#909399" />
				</el-tooltip>
			</span>
			<span class="f_show">显示</span>
			<span class="f_name">名称</span>
			<span class="f_width">宽度</span>
			<span class="f_sortable">排序</span>
			<span class="f_align">对齐</span>
			<span class="f_fixed">固定</span>
		</div>
		<div class="setting-column__list" ref="list">
			<ul>
				<li v-for="item in usercolumn" :key="item.key">
					<span class="f_move">
						<svg class="move" style="cursor: move;" viewBox="64 64 896 896" focusable="false" data-icon="holder" width="1em" height="1em" fill="currentColor" aria-hidden="true">
							<path d="M300 276.5a56 56 0 1056-97 56 56 0 00-56 97zm0 284a56 56 0 1056-97 56 56 0 00-56 97zM640 228a56 56 0 10112 0 56 56 0 00-112 0zm0 284a56 56 0 10112 0 56 56 0 00-112 0zM300 844.5a56 56 0 1056-97 56 56 0 00-56 97zM640 796a56 56 0 10112 0 56 56 0 00-112 0z"></path>
						</svg>
					</span>

					<span class="f_show">
						<el-checkbox v-model="item.show" size="default" />
					</span>
					<span class="f_name" :title="item.key">{{ item.title }}</span>
					<span class="f_width">
						<el-input v-model="item.width" placeholder="auto" size="small"></el-input>
					</span>
					<span class="f_sortable">
						<el-switch size="small" v-model="item.sortable"></el-switch>
					</span>
					<span class="f_align">
						<el-radio-group v-model="item.align">
							<el-radio-button class="radio-mini-button" v-for="(item, index) in alignoptions" size="small" :value="item.value">{{ item.label }}</el-radio-button>
						</el-radio-group>
					</span>
					<span class="f_fixed">
						<el-radio-group v-model="item.fixed">
							<el-radio-button class="radio-mini-button" v-for="(item, index) in fixedoptions" size="small" :value="item.value">{{ item.label }}</el-radio-button>
						</el-radio-group>
					</span>
				</li>
			</ul>
		</div>
		<div class="setting-column__bottom" v-if="optionButton">
			<el-button @click="backDefaul" size="small" :disabled="isSave">重置</el-button>
			<el-button @click="save" size="small" type="primary">保存</el-button>
		</div>
	</div>
	<el-empty v-else description="暂无可配置的列" :image-size="80"></el-empty>
</template>

<script>
	import Sortable from 'sortablejs'

	export default {
		components: {
			Sortable
		},
		props: {
			column: { type: Array, default: () => {} },
			optionButton: { type: Boolean, default: true },
		},
		data() {
			return {
				isSave: false,
				usercolumn: this.column, // JSON.parse(JSON.stringify(this.column||[])),
				alignoptions: [
					{ label: '左', value: 'left' },
					{ label: '中', value: 'center' },
					{ label: '右', value: 'right' },
				],
				fixedoptions: [
					{ label: '否', value: false },
					{ label: '左', value: 'left' },
					{ label: '右', value: 'right' },
				],
			}
		},
		watch:{
			usercolumn: {
				handler(){
					this.$emit('userChange', this.usercolumn)
				},
				deep: true
			}
		},
		mounted() {
			this.usercolumn.length>0 && this.rowDrop()
		},
		methods: {
			rowDrop(){
				const _this = this
				const tbody = this.$refs.list.querySelector('ul')
				Sortable.create(tbody, {
					handle: ".move",
					animation: 300,
					ghostClass: "ghost",
					onEnd({ newIndex, oldIndex }) {
						const tableData = _this.usercolumn
						const currRow = tableData.splice(oldIndex, 1)[0]
						tableData.splice(newIndex, 0, currRow)
					}
				})
			},
			backDefaul(){
				this.$emit('back', this.usercolumn)
			},
			save(){
				this.$emit('save', this.usercolumn)
			},
		}
	}
</script>

<style scoped lang="scss">
	.setting-column {}

	.setting-column__title span {display: inline-block;font-weight: bold;color: #909399;font-size: 12px;}
	.setting-column__title {border-bottom: 1px solid var(--next-border-color-light);padding-bottom:5px; margin-bottom: 5px;}
	.setting-column__title span.f_move {width: 16px;}
	.setting-column__title span.f_show {width: 32px;}
	.setting-column__title span.f_name {width: 120px;}
	.setting-column__title span.f_width {width: 40px;margin-right:15px;}
	.setting-column__title span.f_sortable {width: 50px;}
	.setting-column__title span.f_align {width: 100px;}
	.setting-column__title span.f_fixed {width: 80px;}

	.setting-column__list {max-height:320px;overflow: auto;}
	.setting-column__list li {list-style: none;margin:2px 0;display: flex;align-items: center;}
	.setting-column__list li>span {display: inline-block;font-size: 12px;}
	.setting-column__list li span.f_move {width: 16px;margin-right:5px;}
	.setting-column__list li span.f_show {width: 32px;}
	.setting-column__list li span.f_name {width: 110px;white-space: nowrap;text-overflow: ellipsis;overflow: hidden;cursor:default;}
	.setting-column__list li span.f_width {width: 50px;margin-right:10px;}
	.setting-column__list li span.f_sortable {width: 50px;}
	.setting-column__list li span.f_align {width: 85px;}
	.setting-column__list li span.f_fixed {width: 85px; margin-left: 15px}
	.setting-column__list li.ghost {opacity: 0.3;}

	.setting-column__bottom {border-top: 1px solid var(--next-border-color-light);padding-top:15px;text-align: right; margin-top: 5px; margin-right: 25px}
	
	.dark .setting-column__title {border-color: var(--el-border-color-light);}
	.dark .setting-column__bottom {border-color: var(--el-border-color-light);}

	.radio-mini-button {
		:deep(.el-radio-button__inner) {
			width: 26px;
			height: 21px;
			font-size: 12px;
			padding: 4px 2px;
		}
	}
</style>
