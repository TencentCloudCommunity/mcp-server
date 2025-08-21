<template>
  <div :id="`${path}#default`" class="anchor-demo">
    <t-anchor :affix-props="{ offsetTop: 150 }">
      <t-anchor-item href="#基础锚点" title="基础锚点" />
      <t-anchor-item href="#多级锚点" title="多级锚点" />
      <t-anchor-item href="#指定容器锚点" title="指定容器锚点" />
      <t-anchor-item href="#特定交互锚点" title="特定交互锚点" />
      <t-anchor-item href="#不同尺寸的锚点" title="不同尺寸的锚点" />
    </t-anchor>
  </div>
</template>
<script setup>
import { computed, getCurrentInstance } from 'vue';
import { get } from 'lodash-es';

const { appContext } = getCurrentInstance() || {};
const path = computed(() => get(appContext, '$route.path', ''));
</script>