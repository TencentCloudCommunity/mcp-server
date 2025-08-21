import React, { useState, useEffect } from 'react';
import { SelectInput } from 'tdesign-react';
import { SearchIcon } from 'tdesign-icons-react';

const classStyles = `
`;

const OPTIONS = ['Student A', 'Student B', 'Student C', 'Student D', 'Student E', 'Student F'];

export default function SelectInputAutocomplete() {
  const [popupVisible, setPopupVisible] = useState(false);
  const [selectValue, setSelectValue] = useState('');
  const [options, setOptions] = useState(OPTIONS);

  const onOptionClick = (item: string) => {
    setSelectValue(item);
    setPopupVisible(false);
  };

  const onInputChange = (keyword: string) => {
    setSelectValue(keyword);
    const options = new Array(5).fill(null).map((t, index) => `${keyword} Student ${index}`);
    setOptions(options);
  };

  const onPopupVisibleChange = (val: boolean) => {
    setPopupVisible(val);
  };

  useEffect(() => {
    // 添加示例代码所需样式
    document.head.insertAdjacentHTML('beforeend', classStyles);
  }, []);

  // 如果需要输入框宽度自适应，可以使用 autoWidth
  return (
    <div>
      <SelectInput
        value={selectValue}
        popupVisible={popupVisible}
        placeholder="请输入任意关键词"
        allowInput
        clearable
        onInputChange={onInputChange}
        onPopupVisibleChange={onPopupVisibleChange}
        panel={
          <ul className="tdesign-demo__select-input-ul-autocomplete">
            {options.map((item) => (
              <li key={item} onClick={() => onOptionClick(item)}>
                {item}
              </li>
            ))}
          </ul>
        }
        suffixIcon={<SearchIcon />}
      />
    </div>
  );
}