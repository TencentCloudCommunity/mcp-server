import { fireEvent, render } from "@testing-library/react";
import { expect } from "vitest";

import TooltipDemo from "@tdesign-mcp-server/docs-react/tooltip/Demo";

export default async () => {
  const { container } = render(<TooltipDemo />);
  const button = container.querySelector(".t-button");
  if (button) {
    fireEvent.mouseOver(button);
  }
  const popup = document.querySelector(".t-tooltip");
  expect(popup).toMatchSnapshot();
};
