import React, { useState } from 'react';
import { Dialog, Button } from 'tdesign-react';
import type { DialogProps } from 'tdesign-react';

export default function BasicUsage() {
  const [visible, setVisible] = useState(false);
  const handleClick = () => {
    setVisible(true);
  };
  const onConfirm: DialogProps['onConfirm'] = (context) => {

    setVisible(false);
  };
  const onCancel: DialogProps['onCancel'] = (context) => {

  };
  const onClickCloseBtn: DialogProps['onCloseBtnClick'] = (context) => {

  };
  const onKeydownEsc: DialogProps['onEscKeydown'] = (context) => {

  };
  const onClickOverlay: DialogProps['onOverlayClick'] = (context) => {

  };
  const handleClose: DialogProps['onClose'] = (context) => {

    setVisible(false);
  };
  return (
    <>
      <Button theme="primary" onClick={handleClick}>
        Open Modal
      </Button>
      <Dialog
        header="Basic Modal"
        visible={visible}
        confirmOnEnter
        onClose={handleClose}
        onConfirm={onConfirm}
        onCancel={onCancel}
        onEscKeydown={onKeydownEsc}
        onCloseBtnClick={onClickCloseBtn}
        onOverlayClick={onClickOverlay}
      >
        <p>This is a dialog</p>
      </Dialog>
    </>
  );
}