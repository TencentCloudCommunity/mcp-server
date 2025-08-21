import { fireEvent, render } from "@testing-library/react";

import DropdownDemo from "@tdesign-mcp-server/docs-react/dropdown/Demo";
import { expect } from "vitest";

export default async () => {
  const { container } = render(<DropdownDemo />);
  const button = container.querySelector(".t-button");
  if (button) {
    fireEvent.mouseOver(button);
  }
  const dropdown = document.querySelector(".t-dropdown");
  expect(dropdown).toMatchSnapshot();
};
