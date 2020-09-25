import { Form, Input, message, Modal } from "antd";
import React, { useEffect } from "react";
import { useCustomer } from "shared/hooks";

const layout = {
  labelCol: { span: 8 },
  wrapperCol: { span: 16 },
};

const ruleRequired = { required: true };
const ruleJustRequired = [ruleRequired];

const rulePhoneNumberRequired = [ruleRequired, { type: "string", min: 10, max: 12 }, { pattern: "[0-9]+" }];

function CustomerAddModal({ onClose, onAdd }) {
  const [, , , add, status] = useCustomer();

  useEffect(() => {
    if (status.isSuccess) {
      message.success("New customer added successfully!");
      onAdd();
    } else if (status.isError) {
      message.error("Oops! Failed to add new customer");
    }
  }, [status, onAdd, onClose]);

  return (
    <Modal
      title="Add Customer"
      okButtonProps={{
        htmlType: "submit",
        key: "submit",
        form: "customer-add-form",
        loading: status.isLoading,
      }}
      onCancel={onClose}
      okText="Add Customer"
      width={600}
      visible
    >
      <Form {...layout} id="customer-add-form" onFinish={add}>
        <Form.Item name="name" label="Name" rules={ruleJustRequired}>
          <Input />
        </Form.Item>
        <Form.Item name="address" label="Address" rules={ruleJustRequired}>
          <Input.TextArea rows={4} />
        </Form.Item>
        <Form.Item name="zipcode" label="Zipcode" rules={ruleJustRequired}>
          <Input />
        </Form.Item>
        <Form.Item name="phoneNumber" label="Phone Number" rules={rulePhoneNumberRequired}>
          <Input />
        </Form.Item>
      </Form>
    </Modal>
  );
}

export default CustomerAddModal;
