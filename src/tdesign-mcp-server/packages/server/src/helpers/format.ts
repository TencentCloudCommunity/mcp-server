export const convert2PascalCase = (name: string) => {
  // 已经是帕斯卡命名的情况，直接返回
  if (/^[A-Z][a-zA-Z]*$/.test(name)) {
    return name;
  }

  return name
    .split("-")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join("");
};
