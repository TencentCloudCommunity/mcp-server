<template>
  <t-space>
    <t-button theme="primary" @click="onClick">基础确认对话框</t-button>
    <t-dialog
      v-model:visible="visible"
      header="对话框标题"
      width="40%"
      :confirm-on-enter="true"
      :on-cancel="onCancel"
      :on-esc-keydown="onEscKeydown"
      :on-close-btn-click="onCloseBtnClick"
      :on-overlay-click="onOverlayClick"
      :on-close="close"
      :on-confirm="onConfirmAnother"
    >
      <t-space direction="vertical">
        <div>
          <p>这是弹框内容</p>
          <p>This is Dialog Content</p>
        </div>
        <t-pagination v-model="current" v-model:pageSize="pageSize" :total="30" />
      </t-space>
    </t-dialog>
  </t-space>
</template>
<script setup>
import { ref } from 'vue';

const visible = ref(false);
const current = ref(1);
const pageSize = ref(10);
const onClick = (context) => {

  visible.value = true;
};
const onConfirmAnother = (context) => {

  visible.value = false;
};
const close = (context) => {

};
const onCancel = (context) => {

};
const onEscKeydown = (context) => {

};
const onCloseBtnClick = (context) => {

};
const onOverlayClick = (context) => {

};
</script>