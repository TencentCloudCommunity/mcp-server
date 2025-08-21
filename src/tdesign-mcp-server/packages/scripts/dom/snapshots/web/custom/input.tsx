import { render } from "@testing-library/react";
import { expect } from "vitest";

import { LockOnIcon } from "tdesign-icons-react";
import { Input } from "tdesign-react";

export default async () => {
  const { container } = render(
    <>
      <p>
        基础输入框
        <Input />
      </p>
      <p>
        密码输入框
        <Input
          prefixIcon={<LockOnIcon />}
          type="password"
        />
      </p>
      <p>
        错误状态输入框
        <Input
          tips="校验存在严重问题文本提示"
          status="error"
        />
      </p>
    </>
  );

  expect(container).toMatchSnapshot();
};
