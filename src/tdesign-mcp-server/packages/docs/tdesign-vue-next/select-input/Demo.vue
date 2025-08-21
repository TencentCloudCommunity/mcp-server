<template>
  <t-select-input
    :value="selectValue"
    :popup-visible="popupVisible"
    :popup-props="{ overlayInnerStyle: { padding: '6px' } }"
    placeholder="请输入任意关键词"
    allow-input
    clearable
    @input-change="onInputChange"
    @popup-visible-change="onPopupVisibleChange"
  >
    <template #panel>
      <ul class="tdesign-demo__select-input-ul-autocomplete">
        <li v-for="item in options" :key="item" @click="() => onOptionClick(item)">
          {{ item }}
        </li>
      </ul>
    </template>
    <template #suffixIcon><search-icon /></template>
  </t-select-input>
</template>
<script setup>
import { ref } from 'vue';
import { SearchIcon } from 'tdesign-icons-vue-next';

const OPTIONS = ['Student A', 'Student B', 'Student C', 'Student D', 'Student E', 'Student F'];

const popupVisible = ref(false);
const selectValue = ref();
const options = ref(OPTIONS);

const onOptionClick = (item) => {
  selectValue.value = item;
  popupVisible.value = false;
};

const onInputChange = (keyword) => {
  selectValue.value = keyword;
  options.value = new Array(5).fill(null).map((t, index) => `${keyword} Student ${index}`);
};

const onPopupVisibleChange = (val) => {
  popupVisible.value = val;
};
</script>