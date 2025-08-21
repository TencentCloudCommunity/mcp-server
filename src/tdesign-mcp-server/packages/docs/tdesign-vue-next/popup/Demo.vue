<template>
  <t-space size="large">
    
    <t-popup content="可以在浮层内容中描述禁用原因">
      <div><t-button variant="outline" disabled>禁用状态</t-button></div>
    </t-popup>

    <t-popup>
      <t-button variant="outline">触发元素</t-button>
      
    <t-popup :content="renderPopupContent">
      <t-button variant="outline">浮层内容</t-button>
    </t-popup>
  </t-space>
</template>

<script setup lang="jsx">
function renderPopupContent() {
  return <div>浮层内容是指浮层里面显示的全部内容，这里的浮层内容使用渲染函数定义</div>;
}
</script>