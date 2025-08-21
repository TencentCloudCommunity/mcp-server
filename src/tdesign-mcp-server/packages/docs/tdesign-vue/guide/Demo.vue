<template>
  <t-row justify="center">
    <t-button @click="handleClick">新手引导</t-button>
    <t-drawer :visible.sync="visible" header="演示新手引导" size="60%" :show-overlay="false" destroy-on-close>
      <template #footer>
        <t-button @click="visible = false"> 关闭抽屉 </t-button>
      </template>
      <div class="guide-container">
        <div class="main-title">
          <div class="title-major">Guide 用户引导</div>
          <div class="title-sub">按钮用于开启一个闭环的操作任务，如“删除”对象、“购买”商品等。</div>
        </div>
        <div class="field label-field">
          <div class="label">Label</div>
          <t-input placeholder="请输入内容" />
        </div>
        <div class="field">
          <div class="label">Label</div>
          <t-input placeholder="请输入内容" />
        </div>
        <t-row class="action">
          <t-button>确定</t-button>
          <t-button theme="default" variant="base">取消</t-button>
        </t-row>
      </div>

      <t-guide
        :current.sync="current"
        :steps="steps"
        :highlightPadding="32"
        @change="handleChange"
        @prev-step-click="handlePrevStepClick"
        @next-step-click="handleNextStepClick"
        @finish="handleFinish"
        @skip="handleSkip"
      />
    </t-drawer>
  </t-row>
</template>

<script>
export default {
  data() {
    return {
      visible: false,
      current: -1,
      steps: [
        {
          element: '.main-title',
          title: '新手引导标题',
          body: '新手引导的说明文案',
          placement: 'bottom-right',
        },
        {
          element: '.label-field',
          title: '新手引导标题',
          body: '新手引导的说明文案',
          placement: 'bottom',
        },
        {
          element: '.action',
          title: '新手引导标题',
          body: '新手引导的说明文案',
          placement: 'right',
        },
      ],
    };
  },
  methods: {
    handleClick() {
      this.visible = true;
      setTimeout(() => {
        this.current = 0;
      }, 1000);
    },
    handleChange(current, { e, total }) {

    },
    handlePrevStepClick({
      e, prev, current, total,
    }) {

    },
    handleNextStepClick({
      e, next, current, total,
    }) {

    },
    handleFinish({ e, current, total }) {
      this.visible = false;

    },
    handleSkip({ e, current, total }) {
      this.visible = false;

    },
  },
};
</script>