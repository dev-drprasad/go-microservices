import React from "react";
import { Spin } from "antd";

function StyleWrapper({ children, style }) {
  return (
    <div
      style={{
        display: "flex",
        height: "100%",
        justifyContent: "center",
        alignItems: "center",
        color: "inherit",
        padding: 16,
        ...style,
      }}
    >
      {children}
    </div>
  );
}

export default function NSHandler({ status, children, style, noDataMessage }) {
  if (status.isInit) return null;

  if (status.isSuccess) {
    if (!status.hasData && noDataMessage) return <StyleWrapper style={style}>{noDataMessage}</StyleWrapper>;
    return children();
  }

  if (status.isLoading || status.isError) {
    return (
      <StyleWrapper style={style}>
        {status.isLoading && <Spin />}
        {status.isError && "Oops! Something went wrong."}
        {status.isSuccess && !status.hasData && noDataMessage}
      </StyleWrapper>
    );
  }

  return null;
}
