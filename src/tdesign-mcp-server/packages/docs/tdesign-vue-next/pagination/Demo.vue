<template>
  <t-space direction="vertical" size="large">
    <t-pagination
      v-model="current"
      v-model:pageSize="pageSize"
      :total="100"
      @change="onChange"
      @page-size-change="onPageSizeChange"
      @current-change="onCurrentChange"
    />
    
    <t-pagination
      :total="36"
      :default-current="2"
      :default-page-size="10"
      show-first-and-last-page-btn
      @current-change="onCurrentChange"
      @page-size-change="onPageSizeChange"
      @change="onChange"
    />
  </t-space>
</template>

<script setup>
import { ref } from 'vue';
import { MessagePlugin } from 'tdesign-vue-next';

const current = ref(1);
const pageSize = ref(5);

const onPageSizeChange = (size) => {

  MessagePlugin.success(`pageSize变化为${size}`);
};

const onCurrentChange = (index, pageInfo) => {
  MessagePlugin.success(`转到第${index}页`);

};

const onChange = (pageInfo) => {

};
</script>