import { fireEvent, render } from "@testing-library/react";
import { expect } from "vitest";

import PopupDemo from "@tdesign-mcp-server/docs-react/popup/Demo";

export default async () => {
  const { container } = render(<PopupDemo />);
  const button = container.querySelector(".t-button");
  if (button) {
    fireEvent.mouseOver(button);
  }
  const popup = document.querySelector(".t-popup");
  expect(popup).toMatchSnapshot();
};
