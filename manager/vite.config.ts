import vue from '@vitejs/plugin-vue';
import { resolve } from 'path';
import { defineConfig, loadEnv, ConfigEnv } from 'vite';

const pathResolve = (dir: string): any => {
	return resolve(__dirname, '.', dir);
};

const alias: Record<string, string> = {
	'/@': pathResolve('./src/'),
	'/root': pathResolve('./'),
	'vue-i18n': 'vue-i18n/dist/vue-i18n.cjs.js',
};

const viteConfig = defineConfig((mode: ConfigEnv) => {
	const env = loadEnv(mode.mode, process.cwd());
	// const webpack = require("webpack");
	return {
		plugins: [vue()],
		root: process.cwd(),
		resolve: { alias },
		base: mode.command === 'serve' ? './' : env.VITE_PUBLIC_PATH,
		hmr: true,
		// 引入第三方的配置
		optimizeDeps: {
			include: [
				'element-plus/es/locale/lang/zh-cn', 
				'element-plus/es/locale/lang/en', 
				'element-plus/es/locale/lang/zh-tw'
			],
		},
		server: {
			host: '0.0.0.0',
			port: env.VITE_PORT as unknown as number,
			// 是否开启 https
			https: false,
			// 是否自动在浏览器打开
			open: env.VITE_OPEN,
			// 默认启用并允许任何源
			cors: true,
			proxy: {
				'/api': {
					// 你要跨域访问的网址
					target: env.VITE_API_URL + ":" + env.VITE_API_PORT + "/",
					ws: true,
					// 允许跨域
					changeOrigin: true,
					// rewrite: (path) => path.replace(/^\/api/, ''),
				},
			},
		},
		build: {
			outDir: 'dist',
			// 是否构建source map 文件
			sourcemap: false,
			// chunk 大小警告的限制，默认500KB
			chunkSizeWarningLimit: 1500,
			// 混淆器，terser 构建后文件体积更小，'terser' | 'esbuild'
			minify: 'terser',
			rollupOptions: {
				output: {
					entryFileNames: `assets/[name].${new Date().getTime()}.js`,
					chunkFileNames: `assets/[name].${new Date().getTime()}.js`,
					assetFileNames: `assets/[name].${new Date().getTime()}.[ext]`,
					compact: true,
					manualChunks: {
						vue: ['vue', 'vue-router', 'pinia'],
						echarts: ['echarts'],
					},
				},
			},
			terserOptions: {
				// 生产环境移除console
				compress: {
					drop_console: true,
					drop_debugger: true,
				}
			},
			output: {
				// 去掉注释内容
				comments: true,
			},
		},
		css: {
			preprocessorOptions: {
				css: { charset: false },
				scss: {
					api: 'modern-compiler',
					silenceDeprecations: ['import', 'legacy-js-api'],
				},
			},
		},
		define: {
			__VUE_I18N_LEGACY_API__: JSON.stringify(false),
			__VUE_I18N_FULL_INSTALL__: JSON.stringify(false),
			__INTLIFY_PROD_DEVTOOLS__: JSON.stringify(false),
			__NEXT_VERSION__: JSON.stringify(process.env.npm_package_version),
			__NEXT_NAME__: JSON.stringify(process.env.npm_package_name),
		},
		// configureWebpack: {
		// 	plugins: [
		// 		new webpack.ProvidePlugin({
		// 			process: "process/browser",
		// 			Buffer: ["buffer", "Buffer"],
		// 		}),
		// 	]
		// },
	};
});

export default viteConfig;
