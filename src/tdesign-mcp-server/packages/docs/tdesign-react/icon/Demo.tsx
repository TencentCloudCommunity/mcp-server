import React from 'react';
import { Icon } from 'tdesign-icons-react';
import type { IconProps } from 'tdesign-icons-react';

const CustomUrlIcon = ({ name, size, style }: Partial<Pick<IconProps, 'name' | 'size' | 'style'>>) => (
  <Icon name={name} size={size} style={style} url="https://tdesign.gtimg.com/icon/default-demo/index.js" />
);

function EnhancedIconExample() {
  return (
    <div>
      <CustomUrlIcon name="cps-icon-home-sheep" />
      <CustomUrlIcon name="cps-icon-home-sheep" size="medium" />
      <CustomUrlIcon name="cps-icon-home-sheep" size="large" />
      <CustomUrlIcon name="cps-icon-home-sheep" size="25px" />
      <br />
      <br />
      <CustomUrlIcon name="cps-icon-home-sheep" />
      <CustomUrlIcon name="cps-icon-home-sheep" />
      <CustomUrlIcon name="cps-icon-home-sheep" />
      <CustomUrlIcon name="t-icon-home" />
    </div>
  );
}

export default EnhancedIconExample;