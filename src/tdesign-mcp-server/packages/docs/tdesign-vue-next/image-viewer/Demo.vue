<template>
  <div>
    <div class="tdesign-demo-image-viewer__base">
      <t-image-viewer :images="[img]">
        <template #trigger="{ open }">
          <div class="tdesign-demo-image-viewer__ui-image">
            <img alt="test" :src="img" class="tdesign-demo-image-viewer__ui-image--img" />
            <div class="tdesign-demo-image-viewer__ui-image--hover" @click="open">
              <span><BrowseIcon size="1.4em" /> 预览</span>
            </div>
          </div>
        </template>
      </t-image-viewer>
    </div>
  </div>
</template>
<script setup>
import { BrowseIcon } from 'tdesign-icons-vue-next';

const img = 'https://tdesign.gtimg.com/demo/demo-image-1.png';
</script>