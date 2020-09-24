import React, { useMemo, useState, useEffect } from "react";
import useAPI from "shared/hooks";
import { Modal, Form, Input, message } from "antd";

const layout = {
  labelCol: { span: 8 },
  wrapperCol: { span: 16 },
};

const ruleRequired = { required: true };
const ruleJustRequired = [ruleRequired];

function useAddCategory() {
  const [payload, setPayload] = useState(undefined);
  const args = useMemo(
    () => (payload ? ["/api/v1/categories", { method: "POST", body: JSON.stringify(payload) }] : [undefined, undefined]),
    [payload]
  );
  const [, status] = useAPI(...args);

  return [setPayload, status];
}

function CategoryAddModal({ onClose, onAdd }) {
  const [add, status] = useAddCategory();

  useEffect(() => {
    if (status.isSuccess) {
      message.success("New category added successfully!");
      onAdd();
    } else if (status.isError) {
      message.error("Oops! Failed to add new category");
    }
  }, [status, onAdd, onClose]);

  return (
    <Modal
      title="Add Category"
      okButtonProps={{
        htmlType: "submit",
        key: "submit",
        form: "category-add-form",
        loading: status.isLoading,
      }}
      onCancel={onClose}
      okText="Add Category"
      width={600}
      visible
    >
      <Form {...layout} id="category-add-form" onFinish={add}>
        <Form.Item name="name" label="Name" rules={ruleJustRequired}>
          <Input />
        </Form.Item>
      </Form>
    </Modal>
  );
}

export default CategoryAddModal;
