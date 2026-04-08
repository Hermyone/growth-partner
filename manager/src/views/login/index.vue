<template>
  <div class="login-container">
    <div class="login-bg-layer login-bg-image" aria-hidden="true"></div>
    <div class="login-bg-layer login-bg-tint" aria-hidden="true"></div>
    <div class="login-bg-layer login-bg-shapes" aria-hidden="true"></div>
    <div class="login-content-out">
      <div class="login-content">
        <div class="login-content-main">
          <div class="login-icon-group">
            <div class="login-icon-group-title">
              <img :src="logoMini" alt="" />
              <div class="login-icon-group-title-text font25">{{ getThemeConfig.globalViceTitle }}</div>
            </div>
            <p class="login-tagline">陪伴成长 · 温暖记录每一步闪光</p>
            <div class="login-icon-group-title-ver">ver {{ $CONFIG.APP_VER }}</div>
          </div>
          <div v-if="!isScan">
            <Account />
<!--            <el-tabs class="login-content-tabs" v-model="tabsActiveName">-->
<!--              <el-tab-pane :label="$t('message.label.one1')" name="account">-->
<!--                <Account />-->
<!--              </el-tab-pane>-->
<!--              <el-tab-pane :label="$t('message.label.two2')" name="mobile">-->
<!--                <Mobile />-->
<!--              </el-tab-pane>-->
<!--            </el-tabs>-->
          </div>
<!--          <Scan v-if="isScan" />-->
<!--          <div class="login-content-main-sacn" @click="isScan = !isScan">-->
<!--            <i class="iconfont" :class="isScan ? 'icon-shoujidiannao' : 'icon-ico'"></i>-->
<!--            <div class="login-content-main-sacn-delta"></div>-->
<!--          </div>-->
        </div>
      </div>
    </div>
    <div class="login-footer">
      <div class="login-footer-content mt15">
        <div class="login-footer-content-warp">
          <div>Copyright © 2026-2029 otqsoft All Rights Reserved.</div>
          <div class="mt5">{{$t('message.system.copyright')}}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { toRefs, reactive, computed, defineComponent, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useThemeConfig } from '/@/stores/themeConfig';
import logoMini from '/@/assets/logo-mini.svg';
import { NextLoading } from '/@/utils/loading';
import Account from '/@/views/login/component/account.vue';
import Mobile from '/@/views/login/component/mobile.vue';
import Scan from '/@/views/login/component/scan.vue';

export default defineComponent({
	name: 'loginIndex',
	components: { Account, Mobile, Scan },
	setup() {
		const storesThemeConfig = useThemeConfig();
		const { themeConfig } = storeToRefs(storesThemeConfig);
		const state = reactive({
			tabsActiveName: 'account',
			isScan: false,
		});
		// 获取布局配置信息
		const getThemeConfig = computed(() => {
			return themeConfig.value;
		});
		// 页面加载时
		onMounted(() => {
			NextLoading.done();
		});
		return {
			logoMini,
			getThemeConfig,
			...toRefs(state),
		};
	},
});
</script>


<style scoped lang="scss">
/* 中小学生成长主题：阳光橙 #FF9F43 · 天空蓝 #48DBFB · 薄荷绿 #1DD1A1 · 柔和粉 #FF9FF3 */
.login-container {
  --growth-orange: #ff9f43;
  --growth-sky: #48dbfb;
  --growth-mint: #1dd1a1;
  --growth-pink: #ff9ff3;
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;
  background: linear-gradient(165deg, #fff9f2 0%, #e8f9ff 48%, #f5fffb 100%);

  .login-bg-layer {
    position: absolute;
    inset: 0;
    pointer-events: none;
  }

  .login-bg-image {
    background-image: url('/@/assets/login-growth-bg.svg');
    background-size: cover;
    background-position: center;
    opacity: 1;
  }

  .login-bg-tint {
    background: linear-gradient(
      120deg,
      rgba(255, 159, 67, 0.08) 0%,
      rgba(72, 219, 251, 0.12) 45%,
      rgba(29, 209, 161, 0.06) 100%
    );
  }

  .login-bg-shapes {
    background:
      radial-gradient(ellipse 55% 40% at 88% 12%, rgba(255, 159, 67, 0.22) 0%, transparent 65%),
      radial-gradient(ellipse 40% 35% at 8% 75%, rgba(255, 159, 243, 0.14) 0%, transparent 60%),
      radial-gradient(ellipse 50% 45% at 72% 88%, rgba(29, 209, 161, 0.12) 0%, transparent 55%);
  }

  .login-icon-group {
    width: 100%;
    position: relative;
    text-align: center;
    margin-bottom: 8px;

    .login-icon-group-title {
      display: flex;
      justify-content: center;
      align-items: center;
      flex-wrap: wrap;
      gap: 8px;

      img {
        width: 84px;
        height: 54px;
      }

      &-text {
        padding-left: 2px;
        background: linear-gradient(100deg, var(--growth-orange) 0%, #e8892e 40%, var(--growth-sky) 100%);
        -webkit-background-clip: text;
        background-clip: text;
        color: transparent;
        font-weight: 600;
      }
    }

    .login-tagline {
      margin: 10px 0 6px;
      font-size: 14px;
      line-height: 1.5;
      color: rgba(80, 95, 110, 0.88);
      letter-spacing: 0.06em;
    }

    .login-icon-group-title-ver {
      font-size: 13px;
      color: rgba(120, 135, 150, 0.85);
      margin-top: 2px;
    }

    &-icon {
      width: 60%;
      height: 70%;
      position: absolute;
      left: 0;
      bottom: 0;
    }
  }

  .login-content-out {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    position: relative;
    z-index: 1;
  }

  .login-content {
    width: 550px;
    max-width: calc(100vw - 32px);
    padding: 28px 24px 24px;
    margin: auto;
    background: rgba(255, 255, 255, 0.92);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    border: 1px solid rgba(255, 159, 67, 0.28);
    box-shadow:
      0 4px 24px rgba(72, 219, 251, 0.12),
      0 12px 40px rgba(255, 159, 67, 0.08),
      inset 0 1px 0 rgba(255, 255, 255, 0.95);
    border-radius: 20px;
    overflow: hidden;
    position: relative;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 4px;
      background: linear-gradient(
        90deg,
        var(--growth-orange) 0%,
        var(--growth-sky) 34%,
        var(--growth-mint) 68%,
        var(--growth-pink) 100%
      );
      border-radius: 20px 20px 0 0;
    }

    .login-content-main {
      margin: 0 auto;
      width: 80%;
      .login-content-title {
        color: var(--el-text-color-primary);
        font-weight: 500;
        font-size: 22px;
        text-align: center;
        letter-spacing: 4px;
        margin: 15px 0 30px;
        white-space: nowrap;
        z-index: 5;
        position: relative;
        transition: all 0.3s ease;
      }

      :deep(.login-content-form) {
        .el-input__wrapper {
          border-radius: 10px;
          transition: box-shadow 0.2s ease;
        }

        .el-input__wrapper:hover {
          box-shadow: 0 0 0 1px rgba(255, 159, 67, 0.35) inset;
        }

        .el-input__wrapper.is-focus {
          box-shadow: 0 0 0 1px var(--growth-orange) inset, 0 0 0 3px rgba(72, 219, 251, 0.22) !important;
        }

        .login-content-code-img {
          border-color: rgba(72, 219, 251, 0.35) !important;
          border-radius: 8px !important;
        }

        .login-content-submit.el-button--primary {
          border: none;
          background: linear-gradient(95deg, var(--growth-orange) 0%, #ffb35a 48%, var(--growth-sky) 100%);
          color: #fff;
          font-weight: 500;
          box-shadow: 0 8px 22px rgba(255, 159, 67, 0.28);
        }

        .login-content-submit.el-button--primary:hover {
          filter: brightness(1.05);
        }
      }
    }
    .login-content-main-sacn {
      position: absolute;
      top: 2px;
      right: 2px;
      width: 84px;
      height: 84px;
      overflow: hidden;
      cursor: pointer;
      transition: all ease 0.3s;
      color: var(--el-text-color-primary);
      &-delta {
        position: absolute;
        width: 35px;
        height: 95px;
        z-index: 2;
        top: 2px;
        right: 21px;
        background: var(--el-color-white);
        transform: rotate(-45deg);
      }
      &:hover {
        opacity: 1;
        transition: all ease 0.3s;
        color: var(--el-color-primary) !important;
      }
      i {
        width: 47px;
        height: 50px;
        display: inline-block;
        font-size: 48px;
        position: absolute;
        right: 2px;
        top: -1px;
      }
    }
    .login-content-tabs {
      :deep(.el-tabs__item) {
        width: 120px;
      }
    }
  }
  .login-footer {
    position: absolute;
    bottom: 5px;
    width: 100%;
    z-index: 1;

    &-content {
      width: 100%;
      display: flex;

      &-warp {
        margin: auto;
        color: rgba(90, 105, 120, 0.75);
        text-align: center;
        font-size: 13px;
        animation: error-num 1s ease-in-out;
      }
    }
  }
}
</style>
