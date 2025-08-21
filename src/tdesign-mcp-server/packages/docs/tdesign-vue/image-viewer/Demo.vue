<template>
  <div>
    <div class="tdesign-demo-image-viewer__base">
      <t-image-viewer v-model="visible" :images="[img]" :closeOnEscKeydown="false">
        <template #trigger="{ open }">
          <div class="tdesign-demo-image-viewer__ui-image" @click="open">
            <img alt="test" :src="img" class="tdesign-demo-image-viewer__ui-image--img" />
            <div class="tdesign-demo-image-viewer__ui-image--hover">
              <span><browse-icon size="1.4em" /> 预览</span>
            </div>
          </div>
        </template>
      </t-image-viewer>
    </div>
  </div>
</template>
<script>
import { BrowseIcon } from 'tdesign-icons-vue';

export default {
  components: {
    BrowseIcon,
  },
  data() {
    return {
      img: 'https://tdesign.gtimg.com/demo/demo-image-1.png',
      visible: false,
    };
  },
};
</script>