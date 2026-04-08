<template>
  <div v-if="showChatWindow" class="chat-window">
    <div class="chat-window-header">
      <div style="display: flex; gap: 5px;">
        <el-icon size="20"><ele-ChatDotRound /></el-icon>
        <span>{{ $t('message.user.chat') }}</span>
      </div>
      <el-icon class="close-icon" @click="showChatWindow = false">
        <ele-Close />
      </el-icon>
    </div>
    <iframe :src="getChatUrl()" allow="camera; microphone; display-capture; clipboard-write; clipboard-read" frameborder="0" class="chat-iframe"></iframe>
  </div>
</template>

<script lang="ts">
	import { reactive, toRefs, defineComponent, ref, getCurrentInstance, onMounted } from 'vue';
  import {Session} from "/@/utils/storage";
  import {getConfig} from "/@/api/system/config";

export default defineComponent({
	setup(props, { emit }) {
    const showChatWindow = ref(false);
    const chatBaseUrl = ref('');
		const state = reactive({
			loading: false,
      userid: "",
		});

    const getChatUrl = () => {
      let token = encodeURIComponent(Session.get('token'))
      let url = `${chatBaseUrl.value}#token=${token}`
      if(state.userid){
        url = url + `#id=${state.userid}`
      }
      url = url + `&_t=${Date.now()}`
      return url
    }

    const openDialog = (to_user: string) => {
      state.userid = to_user;
      showChatWindow.value = true;
    }

		onMounted(() => {
      getConfig(14).then(resp => {
        let data = resp.data.data || {};
        chatBaseUrl.value = data.configValue;
      })
		});

		return {
      getChatUrl,
      openDialog,
      showChatWindow,
			...toRefs(state),
		};
	},
});
</script>

<style scoped lang="scss">
.chat-window {
  position: fixed;
  bottom: -35px;
  right: -65px;
  width: 1000px;
  height: 650px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15), 0 0 0 1px rgba(0, 0, 0, 0.05);
  z-index: 9999;
  display: flex;
  flex-direction: column;
  transform: scale(0.85);

  &-header {
    padding: 10px 15px;
    border-bottom: 1px solid #eee;
    display: flex;
    justify-content: space-between;
    align-items: center;

    .close-icon {
      cursor: pointer;
      &:hover {
        color: var(--el-color-primary);
      }
    }
  }
}
.chat-iframe {
  width: 100%;
  height: calc(100% - 40px);
  transform: scale(1.0);
  transform-origin: 0 0;
  border: none;
}
</style>
