import React from 'react';
import { Pagination } from 'tdesign-react';

import type { PaginationProps } from 'tdesign-react';

export default function PaginationExample() {
  const onChange: PaginationProps['onChange'] = (pageInfo) => {

  };

  const onPageSizeChange: PaginationProps['onPageSizeChange'] = (size) => {

  };
  const onCurrentChange: PaginationProps['onCurrentChange'] = (index, pageInfo) => {

  };

  return (
    <Pagination
      total={100}
      defaultPageSize={5}
      onChange={onChange}
      onCurrentChange={onCurrentChange}
      onPageSizeChange={onPageSizeChange}
    />
  );
}