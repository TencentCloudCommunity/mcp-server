### Empty Props

名称 | 类型 | 默认值 | 描述 | 必传
-- | -- | -- | -- | --
action | Slot / Function | - | 操作按钮。TS 类型：`TNode`。 | N
description | String / Slot / Function | - | 描述文字。TS 类型：`string \| TNode`。 | N
image | String / Slot / Function | - | 组件图片，可以完全自定义内容。值类型为字符串时，表示图片地址；值类型为对象时，则表示透传全部属性到图片组件。TS 类型：`string \| ImageProps \| TNode `，[Image API Documents](./image?tab=api)。。 | N
imageStyle | Object | - | 透传图片样式表。TS 类型：`Styles`。 | N
size | String | medium | 空状态的尺寸，默认为 `medium`。可选项：small/medium/large。TS 类型：`SizeEnum`。 | N
title | String / Slot / Function | - | 错误标题。TS 类型：`string \| TNode`。 | N
type | String | empty | 组件类型，如：空数据/成功/失败/网络错误/建设中。可选项：empty/success/fail/network-error/maintenance | N