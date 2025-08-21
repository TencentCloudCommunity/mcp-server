### Avatar Props

名称 | 类型 | 默认值 | 说明 | 必传
-- | -- | -- | -- | --
alt | String | - | 头像替换文本，仅当图片加载失败时有效 | N
content | String / Slot / Function | - | 子元素内容。TS 类型：`string \| TNode`。 | N
default | String / Slot / Function | - | 子元素内容，同 content。TS 类型：`string \| TNode`。 | N
hideOnLoadFailed | Boolean | false | 加载失败时隐藏图片 | N
icon | Slot / Function | - | 图标。TS 类型：`TNode`。 | N
image | String | - | 图片地址 | N
imageProps | Object | - | 透传至 Image 组件。TS 类型：`ImageProps`，[Image API Documents](./image?tab=api)。 | N
shape | String | circle | 形状。可选项：circle/round。TS 类型：`ShapeEnum ` `type ShapeEnum = 'circle' \| 'round'`。 | N
size | String | - | 尺寸，示例值：small/medium/large/24px/38px 等。优先级高于 AvatarGroup.size 。Avatar 单独存在时，默认值为 medium。如果父组件存在 AvatarGroup，默认值便由 AvatarGroup.size 决定 | N
onError | Function |  | TS 类型：`(context: { e: Event }) => void`<br/>图片加载失败时触发 | N

### Avatar Events

名称 | 参数 | 描述
-- | -- | --
error | `(context: { e: Event })` | 图片加载失败时触发

### AvatarGroup Props

名称 | 类型 | 默认值 | 说明 | 必传
-- | -- | -- | -- | --
cascading | String | 'right-up' | 图片之间的层叠关系，可选值：左侧图片在上和右侧图片在上。可选项：left-up/right-up。TS 类型：`CascadingValue` `type CascadingValue = 'left-up' \| 'right-up'`。 | N
collapseAvatar | String / Slot / Function | - | 头像数量超出时，会出现一个头像折叠元素。该元素内容可自定义。默认为 `+N`。示例：`+5`，`...`, `更多`。TS 类型：`string \| TNode`。 | N
max | Number | - | 能够同时显示的最多头像数量 | N
popupProps | Object | - | 头像右上角提示信息。TS 类型：`PopupProps`，[Popup API Documents](./popup?tab=api)。 | N
size | String | medium | 尺寸，示例值：small/medium/large/24px/38px 等。优先级低于 Avatar.size | N