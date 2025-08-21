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
<script>
import { SearchIcon } from 'tdesign-icons-vue';

export default {
  components: {
    SearchIcon,
  },
  data() {
    return {
      selectValue: '',
      popupVisible: false,
      options: ['Student A', 'Student B', 'Student C', 'Student D', 'Student E', 'Student F'],
    };
  },
  methods: {
    onOptionClick(item) {
      this.selectValue = item;
      this.popupVisible = false;
    },
    onInputChange(keyword) {
      this.selectValue = keyword;
      this.options = new Array(5).fill(null).map((t, index) => `${keyword} Student ${index}`);
    },
    onPopupVisibleChange(val) {
      this.popupVisible = val;
    },
  },
};
</script>