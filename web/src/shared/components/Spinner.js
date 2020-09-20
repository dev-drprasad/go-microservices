import React from "react";
import { Spin } from "antd";
import Icon from "@ant-design/icons";

const Spinner = () => (
  <div className="spinner">
    <Spin indicator={<Icon type="loading" />} />
  </div>
);
export default Spinner;
