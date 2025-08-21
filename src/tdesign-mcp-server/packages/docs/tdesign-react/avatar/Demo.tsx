import React from 'react';
import { Avatar, Space } from 'tdesign-react';
import { UserIcon } from 'tdesign-icons-react';

export default function BasicAvatar() {
  return (
    <Space>
      <Avatar icon={<UserIcon />} />
      <Avatar
        image="https://tdesign.gtimg.com/site/avatar.jpg"
        hideOnLoadFailed={false}
      />
      <Avatar>W</Avatar>
    </Space>
  );
}