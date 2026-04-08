import type { App } from 'vue';
import mitt from 'mitt';
import { i18n } from '/@/i18n';
import common, { getUpFileUrl, handleTree, convertToTree, parseTime, selectDictLabel, selectDictColor } from '/@/utils/common';
import { useDict } from '/@/api/system/dict/data';
import { getItems, getOptionValue, isEmpty, setItems } from '/@/api/items';
import config from '/@/config';
import api from '/@/api';
import { directive } from '/@/directives/directive';
import ErrHandle from '/@/utils/errorHandler';
import qpagination from '/@/components/qpagination/index.vue';
import qcontainer from '/@/components/qcontainer/index.vue';
import qcontainersize from '/@/components/qcontainer/autosize.vue'
import qtable from '/@/components/qtable/v1/index.vue';
import qtablev2 from '/@/components/qtable/v2/index.vue';
import qformtable from '/@/components/qtable/form/index.vue';
import qselecttable from '/@/components/qtable/select/index.vue';
import qsearch from '/@/components/qtable/search.vue';
import qpasswordstrength from '/@/components/qpasswordstrength/index.vue';
import qtime from '/@/components/qtime/index.vue';
import qstatistic from '/@/components/qmini/statistic.vue';
import qtrend from '/@/components/qmini/trend.vue';
import qstatusindicator from '/@/components/qmini/statusindicator.vue';
import qtitle from '/@/components/qtitle/index.vue';
import qweather from '/@/components/qweather/index.vue';
import qdialog from '/@/components/qdialog/index.vue';
import qcropper from '/@/components/qcropper/index.vue';
import qnoticebar from '/@/components/noticeBar/index.vue';
import qselect from '/@/components/qselect/index.vue';
import qtreeselect from '/@/components/qselect/tree-select.vue';
import qimageselect from '/@/components/qselect/image-select.vue';
import qinputnumber from '/@/components/qinputnumber/index.vue';
import qlabelvalue from '/@/components/qlabelvalue/index.vue';
import qdate from '/@/components/qdate/index.vue';
import qrate from '/@/components/qrate/index.vue';
import qlinklabel from '/@/components/qlinklabel/index.vue';
import qtag from '/@/components/qtag/index.vue';
import qecharts from '/@/components/qecharts/index.vue';
import qfilepreview from '/@/components/qfilepreview/index.vue';
import qdocxpreview from '/@/components/qfilepreview/component/docxPreview.vue';
import qxlsxpreview from '/@/components/qfilepreview/component/xlsxPreview.vue';
import qpdfpreview from '/@/components/qfilepreview/component/pdfPreview.vue';
import qimgpreview from '/@/components/qfilepreview/component/imgPreview.vue';
import qupload from '/@/components/qupload/index.vue';
import quploadfile from '/@/components/qupload/file.vue';
import quploadmultiple from '/@/components/qupload/multiple.vue';
import qricheditor from '/@/components/qeditor/richEditor.vue';
import qcodeeditor from '/@/components/qeditor/codeEditor.vue';

import * as qicons from '/@/assets/icons';

// @ts-ignore
import VueClipboard from 'vue3-clipboard';

// @ts-ignore
import VueMarkdownEditor from '@kangc/v-md-editor';
// @ts-ignore
import '@kangc/v-md-editor/lib/style/base-editor.css';
// @ts-ignore
import vuepressTheme from '@kangc/v-md-editor/lib/theme/vuepress.js';
import '@kangc/v-md-editor/lib/theme/style/vuepress.css';

// @ts-ignore 代码高亮
import Prism from 'prismjs';

// @ts-ignore 表情插件
import createEmojiPlugin from '@kangc/v-md-editor/lib/plugins/emoji/index';
import '@kangc/v-md-editor/lib/plugins/emoji/emoji.css';

// @ts-ignore 代码行号插件
import createLineNumbertPlugin from '@kangc/v-md-editor/lib/plugins/line-number/index';

// @ts-ignore 高亮代码行插件
import createHighlightLinesPlugin from '@kangc/v-md-editor/lib/plugins/highlight-lines/index';
import '@kangc/v-md-editor/lib/plugins/highlight-lines/highlight-lines.css';

// @ts-ignore 快捷复制插件
import createCopyCodePlugin from '@kangc/v-md-editor/lib/plugins/copy-code/index';
import '@kangc/v-md-editor/lib/plugins/copy-code/copy-code.css';

// @ts-ignore 快捷插入提示(带样式)
import createTipPlugin from '@kangc/v-md-editor/lib/plugins/tip/index';
import '@kangc/v-md-editor/lib/plugins/tip/tip.css';

// @ts-ignore md预览组件
import VMdPreview from '@kangc/v-md-editor/lib/preview';
import '@kangc/v-md-editor/lib/style/preview.css';

// @ts-ignore
import VMdPreviewHtml from '@kangc/v-md-editor/lib/preview-html';
import '@kangc/v-md-editor/lib/style/preview-html.css';

// @ts-ignore 引入你所使用的主题 此处以 github 主题为例
import githubTheme from '@kangc/v-md-editor/lib/theme/github';
import '@kangc/v-md-editor/lib/theme/style/github.css';

// @ts-ignore
import hljs from 'highlight.js';

VMdPreview.use(createLineNumbertPlugin());
VMdPreview.use(createHighlightLinesPlugin());
VMdPreview.use(createCopyCodePlugin());
VMdPreview.use(githubTheme, {
	Hljs: hljs,
});
VMdPreview.use(vuepressTheme, {
	Prism,
})

VueMarkdownEditor.use(createEmojiPlugin());
VueMarkdownEditor.use(createLineNumbertPlugin());
VueMarkdownEditor.use(createHighlightLinesPlugin());
VueMarkdownEditor.use(createCopyCodePlugin());
VueMarkdownEditor.use(createTipPlugin())
VueMarkdownEditor.use(vuepressTheme, {
	Prism,
});

export default {
	install(app:App) {
		app.use(VueClipboard, {autoSetContainer: true, appendToBody: true,});
		app.use(VueMarkdownEditor);
		app.use(VMdPreview);
		app.use(VMdPreviewHtml);

		app.use(VueClipboard, {autoSetContainer: true, appendToBody: true,});

		// 挂载全局对象
		app.config.globalProperties.getUpFileUrl=getUpFileUrl;
		app.config.globalProperties.handleTree=handleTree;
		app.config.globalProperties.convertToTree=convertToTree;
		app.config.globalProperties.useDict=useDict;
		app.config.globalProperties.selectDictLabel=selectDictLabel;
		app.config.globalProperties.selectDictColor=selectDictColor;

		app.config.globalProperties.getItems=getItems;
		app.config.globalProperties.setItems=setItems;
		app.config.globalProperties.getOptionValue=getOptionValue;
		app.config.globalProperties.isEmpty=isEmpty;
		app.config.globalProperties.parseTime=parseTime;

		app.config.globalProperties.$CONFIG = config;
		app.config.globalProperties.$API = api;

		// 注册组件
		app.component(qpagination.name, qpagination);
		app.component(qcontainer.name, qcontainer);
		app.component(qcontainersize.name, qcontainersize);
		app.component(qtable.name, qtable);
		app.component(qtablev2.name, qtablev2);
		app.component(qformtable.name, qformtable);
		app.component(qselecttable.name, qselecttable);
		app.component(qsearch.name, qsearch);
		app.component(qpasswordstrength.name, qpasswordstrength);
		app.component(qtime.name, qtime);
		app.component(qstatistic.name, qstatistic);
		app.component(qtrend.name, qtrend);
		app.component(qstatusindicator.name, qstatusindicator);
		app.component(qtitle.name, qtitle);
		app.component(qweather.name, qweather);
		app.component(qdialog.name, qdialog);
		app.component(qcropper.name, qcropper);
		app.component(qnoticebar.name, qnoticebar);
		app.component(qselect.name, qselect);
		app.component(qtreeselect.name, qtreeselect);
		app.component(qimageselect.name, qimageselect);
		app.component(qinputnumber.name, qinputnumber);
		app.component(qlabelvalue.name, qlabelvalue);
		app.component(qdate.name, qdate);
		app.component(qrate.name, qrate);
		app.component(qlinklabel.name, qlinklabel);
		app.component(qtag.name, qtag);
		app.component(qecharts.name, qecharts);
		app.component(qfilepreview.name, qfilepreview);
		app.component(qdocxpreview.name, qdocxpreview);
		app.component(qxlsxpreview.name, qxlsxpreview);
		app.component(qpdfpreview.name, qpdfpreview);
		app.component(qimgpreview.name, qimgpreview);
		app.component(qupload.name, qupload);
		app.component(quploadfile.name, quploadfile);
		app.component(quploadmultiple.name, quploadmultiple);
		app.component(qricheditor.name, qricheditor);
		app.component(qcodeeditor.name, qcodeeditor);

		const globalProperties={
			mittBus: mitt(),
			i18n
		};

		// 必须合并vue默认的变量，否则有问题
		app.config.globalProperties = Object.assign(
			app.config.globalProperties,
			globalProperties
		);

		// 注册全局指令
		directive(app);

		// 注册全局图标
		common.elSvg(app);

		//注册q-icon自定义图标
		for(let icon in qicons){
			var iconName = icon.replace(/[A-Z]/g, (match)=>'-'+match.toLowerCase());
			// @ts-ignore
			app.component(`q-icon${iconName}`, qicons[icon]);
		}

		// 全局代码错误捕捉
		if (process.env.NODE_ENV !== 'production') {
			app.config.errorHandler = ErrHandle.errorHandler;
			// @ts-ignore
			window.onerror = ErrHandle.errorWin;
			// @ts-ignore 关闭async-validator全局控制台警告
			window.ASYNC_VALIDATOR_NO_WARNING = 1;
		}
	}
}