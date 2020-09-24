import React, { useMemo, useState, useEffect } from "react";
import useBROAPI from "shared/hooks";
import { Modal, Form, Input, message } from "antd";

const layout = {
  labelCol: { span: 8 },
  wrapperCol: { span: 16 },
};

const ruleRequired = { required: true };
const ruleJustRequired = [ruleRequired];

function useAddBrand() {
  const [payload, setPayload] = useState(undefined);
  const args = useMemo(
    () => (payload ? ["/api/v1/brands", { method: "POST", body: JSON.stringify(payload) }] : [undefined, undefined]),
    [payload]
  );
  const [, status] = useBROAPI(...args);

  return [setPayload, status];
}

function BrandAddModal({ onClose, onAdd }) {
  const [add, status] = useAddBrand();

  useEffect(() => {
    if (status.isSuccess) {
      message.success("New brand added successfully!");
      onAdd();
    } else if (status.isError) {
      message.error("Oops! Failed to add new brand");
    }
  }, [status, onAdd, onClose]);

  return (
    <Modal
      title="Add Brand"
      okButtonProps={{
        htmlType: "submit",
        key: "submit",
        form: "brand-add-form",
        loading: status.isLoading,
      }}
      onCancel={onClose}
      okText="Add Brand"
      width={600}
      visible
    >
      <Form {...layout} id="brand-add-form" onFinish={add}>
        <Form.Item name="name" label="Name" rules={ruleJustRequired}>
          <Input />
        </Form.Item>
      </Form>
    </Modal>
  );
}

export default BrandAddModal;
