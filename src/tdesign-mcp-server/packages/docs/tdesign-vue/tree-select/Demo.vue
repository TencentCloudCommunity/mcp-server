<template>
  <t-space direction="vertical">
    
    <t-tree-select
      ref="single-tree-select"
      :data="options"
      v-model="value2"
      label="单选："
      :treeProps="{ expandAll: true, valueMode: 'onlyLeaf' }"
      clearable
      placeholder="请选择（显示选中节点的全部路径）"
      @blur="onBlurTrigger"
      @change="onChange"
    >
      <template #valueDisplay="{ value }">
        <div v-if="value.label">{{ getNodePath(value) }}</div>
      </template>
    </t-tree-select>

  </t-space>
</template>

<script>
export default {
  data() {
    return {
      value1: '',
      value2: '',
      // popupVisible: false,
      options: [
        {
          label: '广东省',
          value: 'guangdong',
          children: [
            {
              label: '广州市',
              value: 'guangzhou',
            },
            {
              label: '深圳市',
              value: 'shenzhen',
              children: [
                {
                  label: '南山区海天二路33号腾讯滨海大厦',
                  value: 'Tencent',
                },
              ],
            },
          ],
        },
        {
          label: '江苏省',
          value: 'jiangsu',
          disabled: true,
          children: [
            {
              label: '南京市',
              value: 'nanjing',
            },
            {
              label: '苏州市',
              value: 'suzhou',
            },
          ],
        },
      ],
    };
  },
  methods: {
    onBlurTrigger(context) {

    },
    onChange(value, context) {

    },
    // onVisibleChange(v, c) {
    //   if (c.trigger || c.node?.label !== '广州市') {
    //     this.popupVisible = v;
    //   }
    // },
    getNodePath(nodeInfo) {
      if (!nodeInfo.value || !this.$refs['single-tree-select']) return;
      const path = this.$refs['single-tree-select'].$refs.treeRef?.getPath(nodeInfo.value);
      return path.map((node) => node.data.label).join('/');
    },
  },
};
</script>