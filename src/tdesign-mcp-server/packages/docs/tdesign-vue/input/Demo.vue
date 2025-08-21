<template>
  <t-space direction="vertical">
    <t-input @focus="onFocus" />
    <t-input v-model="input" placeholder="请输入内容（有默认值）" @enter="onEnter" @change="onChange" />
    <t-input v-model="count" label="价格：" type="number" suffix="元" @change="onNumberChange" @focus="onNumberFocus" />

    <t-input placeholder="Search Input">
      <template #suffixIcon>
        <search-icon />
      </template>
    </t-input>
  </t-space>
</template>

<script>
import { SearchIcon } from 'tdesign-icons-vue';

export default {
  components: { SearchIcon },

  data() {
    return {
      input: '有默认值',
      count: undefined,
    };
  },

  methods: {
    onFocus(val) {

    },
    onEnter() {

    },
    onChange(val) {

    },
    onNumberChange(val) {

    },
    onNumberFocus(val) {

    },
  },
};
</script>