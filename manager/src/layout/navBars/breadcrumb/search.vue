<template>
	<div class="search-input">
		<el-autocomplete
				v-model="menuQuery"
				:fetch-suggestions="menuSearch"
				:placeholder="$t('message.user.searchPlaceholder')"
				ref="layoutMenuAutocompleteRef"
				@select="onHandleSelect"
				@blur="onSearchBlur"
		>
			<template #prefix>
				<el-icon class="el-input__icon">
					<ele-Search />
				</el-icon>
			</template>
			<template #default="{ item }">
				<div>
					<SvgIcon :name="item.meta.icon" class="mr5" />
					<span v-html="item.meta.title_new"></span>
				</div>
			</template>

		</el-autocomplete>
	</div>
</template>

<script lang="ts">
import { reactive, toRefs, defineComponent, ref, nextTick, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { storeToRefs } from 'pinia';
import { useTagsViewRoutes } from '/@/stores/tagsViewRoutes';
import { RouteItem } from '/@/types/global';

export default defineComponent({
	name: 'layoutBreadcrumbSearch',
	setup(props,{ emit }) {
		const storesTagsViewRoutes = useTagsViewRoutes();
		const { tagsViewRoutes } = storeToRefs(storesTagsViewRoutes);
		const layoutMenuAutocompleteRef = ref();
		const { t } = useI18n();
		const router = useRouter();
		const state = reactive({
			menuQuery: '',
			tagsViewList: [],
		});

		onMounted(() => {
			state.menuQuery = '';
			initTageView();
			nextTick(() => {
				setTimeout(() => {
					layoutMenuAutocompleteRef.value.focus();
				});
			});
		});
		// 搜索弹窗关闭
		const closeSearch = () => {
			emit('close', false)
		};
		// 菜单搜索数据过滤
		const menuSearch = (queryString: string, cb: Function) => {
			let results = queryString ? state.tagsViewList.filter(createFilter(queryString)) : state.tagsViewList;
			results.forEach((item:any) => {
				let title = t(item.meta.title)
				if(state.menuQuery != '') title = title.replaceAll(state.menuQuery, '<em>'+state.menuQuery+'</em>')
				item.meta.title_new = title
			})
			cb(results);
		};
		// 菜单搜索过滤
		const createFilter: any = (queryString: string) => {
			return (restaurant: RouteItem) => {
				return (
					restaurant.path.toLowerCase().indexOf(queryString.toLowerCase()) > -1 ||
					restaurant.meta!.title!.toLowerCase().indexOf(queryString.toLowerCase()) > -1 ||
					t(restaurant.meta!.title!).indexOf(queryString.toLowerCase()) > -1
				);
			};
		};
		// 初始化菜单数据
		const initTageView = () => {
			if (state.tagsViewList.length > 0) return false;
			tagsViewRoutes.value.map((v: any) => {
				// 不加载隐藏的菜单, 不加载菜单
				if (!v.meta.isHide && (!v.children || (v.children && (v.children.length == 0)))) {
					state.tagsViewList.push({ ...v });
				}
			});
		};
		// 当前菜单选中时
		const onHandleSelect = (item: RouteItem) => {
			let { path, redirect } = item;
			if (item.meta?.isLink && !item.meta?.isIframe) window.open(item.meta?.isLink);
			else if (redirect) router.push(redirect);
			else router.push(path);
			closeSearch();
		};
		// input 失去焦点时
		const onSearchBlur = () => {
			closeSearch();
		};

		return {
			layoutMenuAutocompleteRef,
			closeSearch,
			menuSearch,
			onHandleSelect,
			onSearchBlur,
			...toRefs(state),
		};
	},
});
</script>

<style scoped lang="scss">
	.search-input {
		width: 300px;
		animation: expandFromRight .3s forwards;
		:deep(.el-autocomplete) {
			width: 100%;
		}
	}
</style>
