export const addTdPrefix = (suffix: string) => `tdesign-${suffix}`;

export const PLATFORMS = {
  web: "web",
  mobile: "mobile"
} as const;

export type Platform = (typeof PLATFORMS)[keyof typeof PLATFORMS];

export const WEB_FRAMEWORKS = ["react", "vue", "vue-next"] as const;

// TODO：未来再引入移动端，现阶段先测试 PC 端
// export const MOBILE_FRAMEWORKS = ["mobile-react", "mobile-vue", "miniprogram"] as const;
export const MOBILE_FRAMEWORKS = [] as const;

export const FRAMEWORKS = [...WEB_FRAMEWORKS, ...MOBILE_FRAMEWORKS] as const;

export type Framework = (typeof FRAMEWORKS)[number];

/* 组件库是否为大仓结构 */
const MONOREPO: Framework[] = ["react", "vue-next"] as const;
export const isMonorepo = (fw: Framework) => MONOREPO.includes(fw);

export const getSuffixByFramework = (fw: Framework) => (fw.includes("react") ? "tsx" : "vue");

export const getPlatformByFramework = (fw: Framework) => {
  return WEB_FRAMEWORKS.includes(fw) ? PLATFORMS.web : PLATFORMS.mobile;
};
