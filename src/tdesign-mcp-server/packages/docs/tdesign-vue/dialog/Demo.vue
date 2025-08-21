<template>
  <div>
    <t-button theme="primary" @click="visible = true">基础确认对话框</t-button>

    <t-dialog
      header="对话框标题"
      body="对话框内容"
      :visible.sync="visible"
      @confirm="onConfirm"
      width="40%"
      :confirmOnEnter="true"
      :onConfirm="onConfirmAnother"
      :onCancel="onCancel"
      :onEscKeydown="onKeydownEsc"
      :onCloseBtnClick="onClickCloseBtn"
      :onOverlayClick="onClickOverlay"
      :onClose="close"
    >
      <t-space direction="vertical">
        <t-button theme="primary" @click="visible1 = true">弹窗二</t-button>
        <t-pagination v-model="current" :total="30" />
      </t-space>
    </t-dialog>
    <t-dialog header="对话框标题二" body="对话框内容二" :visible.sync="visible1" />
  </div>
</template>
<script>
export default {
  data() {
    return {
      visible: false,
      visible1: false,
      current: 1,
    };
  },
  methods: {
    onConfirm(context) {

      this.visible = false;
    },
    onConfirmAnother(context) {

    },
    close(context) {

    },
    onCancel(context) {

    },
    onKeydownEsc(context) {

    },
    onClickCloseBtn(context) {

    },
    onClickOverlay(context) {

    },
  },
};
</script>