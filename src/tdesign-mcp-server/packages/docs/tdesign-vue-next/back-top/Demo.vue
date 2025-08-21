<template>
  <div>
    <t-list class="baseList-container">
      <t-list-item v-for="index in 15" :key="index">列表内容{{ index }}</t-list-item>
    </t-list>
    <t-back-top
      container=".baseList-container"
      :visible-height="0"
      :offset="['24px', '80px']"
    ></t-back-top>
  </div>
</template>