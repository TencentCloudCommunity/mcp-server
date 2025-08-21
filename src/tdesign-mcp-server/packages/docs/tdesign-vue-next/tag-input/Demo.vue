<template>
  <t-space direction="vertical">
    <t-tag-input v-model="tags1" clearable @paste="onPaste" @enter="onTagInputEnter" @focus="onFocus" @blur="onBlur" />

    <t-tag-input :value="tags2" label="Controlled: " clearable @change="onChange" />

    <t-tag-input :default-value="tags3" label="UnControlled: " clearable />
  </t-space>
</template>

<script setup>
import { ref } from 'vue';

const tags1 = ref(['Vue', 'React']);
const tags2 = ref(['Vue', 'React']);
const tags3 = ref(['Vue', 'React']);

const onTagInputEnter = (val, context) => {

};
const onChange = (val, context) => {

  tags2.value = val;
};
const onPaste = (context) => {

};

const onFocus = (val, ctx) => {

};

const onBlur = (val, ctx) => {

};
</script>