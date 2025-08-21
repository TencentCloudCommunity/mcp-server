import React, { useState } from 'react';
import { TagInput, Space } from 'tdesign-react';
import type { TagInputProps, TagInputValue } from 'tdesign-react';

export default function TagInputBaseExample() {
  const [tags1, setTags1] = useState<TagInputValue>(['Vue', 'React', 'Angular']);
  const [tags2] = useState(['Vue', 'React']);
  const [tags3] = useState(['Vue', 'React']);

  const onTagInputEnter: TagInputProps['onEnter'] = (val, context) => {

  };

  const onChange: TagInputProps['onChange'] = (val, context) => {

    setTags1(val);
  };

  const onPaste: TagInputProps['onPaste'] = (context) => {

  };

  return (
    <Space direction="vertical">
      <TagInput
        value={tags1}
        onChange={onChange}
        clearable
        onPaste={onPaste}
        onEnter={onTagInputEnter}
        placeholder="请输入"
      />
      <TagInput value={tags2} label="Controlled: " placeholder="请输入" clearable />
      <TagInput defaultValue={tags3} label="UnControlled: " placeholder="请输入" clearable />
    </Space>
  );
}