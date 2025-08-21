type Rule = [RegExp, string];
type RuleMap = Record<string, Rule>;

const MD_LINK_RULE: Rule = [/\[([^\]]+)\]\((https?:\/\/[^\)]+)\)/g, ""];
const ATTR_INLINE_STYLE_RULE: Rule = [/\s*:?style="([^"]*)"/g, ""];

const COMMON_RULE_MAP: RuleMap = {
  extraLine: [/^\s*\n\s*\n/gm, "\n"],
  blankHeadTail: [/^(\s*\n)+|(\s*\n)+$/g, ""]
};

export const API_RULE_MAP: RuleMap = {
  base: [/:: BASE_DOC ::[\s\S]*?^## API.*$/m, ""],
  mdLink: MD_LINK_RULE,
  ...COMMON_RULE_MAP
};

export const DEMO_RULE_MAP: RuleMap = {
  meta: [/---\n([\s\S]*?)\n---/, ""],
  style: [/\s<style\b[^>]*>[\s\S]*?<\/style>/gi, ""],
  jsxInlineStyle: [/\s*style={{([^}]*)*}}/g, ""],
  attrInlineStyle: ATTR_INLINE_STYLE_RULE,
  selfClosingDiv: [/^\s*<div\s*(\/>|>\s*<\/div>)[\n]?/gm, ""],
  comment: [/<!--[\s\S]*?\n[\s\S]*?-->/g, ""],
  log: [/^\s*console\.log\([\s\S]*?\);\s*$/gm, ""],
  ...COMMON_RULE_MAP
};

export const DOM_RULE_MAP: RuleMap = {
  outerDiv: [/^<div>\s*([\s\S]*?)\s*<\/div>$/gm, "$1"],
  attrInlineStyle: ATTR_INLINE_STYLE_RULE,
  link: [/https?[^"]*/g, ""],
  linkTag: [/<link\s[^>]*\/?>/g, ""],
  svgPathData: [/<path([^>]*)\sd="[^"]*"([^>]*)>/g, '<path$1 d=""$2>'],
  allNewLines: [/\n/g, ""],
  space: [/\s+/g, " "],
  spaceBetweenTag: [/\s*(<|>)\s*/g, "$1"]
};

export const CHANGELOG_RULE_MAP: RuleMap = {
  mdLink: MD_LINK_RULE,
  user: [/(?<!`)\B@[\w-]+(?!`)/g, ""],
  emptyBrackets: [/\(\s*\)/g, ""]
};

export const cleanText = (text: string, cleaningRules: RuleMap) => {
  return Object.values(cleaningRules).reduce((acc, [regex, replacement]) => {
    return acc.replace(regex, replacement);
  }, text);
};
