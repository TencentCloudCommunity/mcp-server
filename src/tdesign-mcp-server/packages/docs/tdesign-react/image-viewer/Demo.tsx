import React from 'react';
import { BrowseIcon } from 'tdesign-icons-react';
import { Image, ImageViewer, Space } from 'tdesign-react';

import type { ImageViewerProps } from 'tdesign-react';

const img = 'https://tdesign.gtimg.com/demo/demo-image-1.png';

export default function BasicImageViewer() {
  const trigger: ImageViewerProps['trigger'] = ({ open }) => {
    const mask = (
      <div
        onClick={open}
      >
        <span>
          <BrowseIcon size="16px" name={'browse'} /> 预览
        </span>
      </div>
    );

    return (
      <Image
        alt={'test'}
        src={img}
        overlayContent={mask}
        overlayTrigger="hover"
        fit="contain"
      />
    );
  };

  return (
    <Space breakLine size={16}>
      <ImageViewer trigger={trigger} images={[img]} zIndex={10000} />
    </Space>
  );
}