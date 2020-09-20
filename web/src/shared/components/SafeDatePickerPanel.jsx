import React, { useState } from "react";
import { Popconfirm } from "antd";
import { DatePickerPanel } from "shared/components";

export default function SafeDatePickerPanel({ onChange, defaultValue, okText, cancelText, confirmTitleRender, ...rest }) {
  const [[newDate, confirmed], setNewDate] = useState([undefined, true]);
  const handleChange = (d) => setNewDate([d, false]);
  const handleConfirm = () => {
    onChange(newDate);
    setNewDate(([d]) => [d, true]);
  };
  const handleCancel = () => setNewDate([undefined, true]);

  return (
    <Popconfirm
      title={newDate ? confirmTitleRender(newDate) : ""}
      placement="left"
      okText={okText}
      cancelText={cancelText}
      onCancel={handleCancel}
      onConfirm={handleConfirm}
      visible={!confirmed}
    >
      <DatePickerPanel value={newDate || defaultValue} defaultValue={defaultValue} onChange={handleChange} {...rest} />
    </Popconfirm>
  );
}
