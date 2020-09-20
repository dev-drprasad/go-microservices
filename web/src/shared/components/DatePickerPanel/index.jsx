import React from "react";
import { PickerPanel } from "rc-picker";
import locale from "rc-picker/lib/locale/en_US";
import momentGenerateConfig from "rc-picker/lib/generate/moment";

function DatePickerPanel(props) {
  return <PickerPanel locale={locale} prefixCls="ant-picker" generateConfig={momentGenerateConfig} {...props} />;
}

export default DatePickerPanel;
