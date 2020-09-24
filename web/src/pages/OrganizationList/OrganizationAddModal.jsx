import React, { useMemo, useState, useEffect } from "react";
import { Modal, Form, Input, message } from "antd";
import useBROAPI from "shared/hooks";

const layout = {
  labelCol: { span: 8 },
  wrapperCol: { span: 16 },
};

const ruleReuired = [{ required: true }];

function useAddOrganization() {
  const [payload, setPayload] = useState(undefined);
  const args = useMemo(
    () => (payload ? ["/api/v1/organizations", { method: "POST", body: JSON.stringify(payload) }] : [undefined, undefined]),
    [payload]
  );
  const [, status] = useBROAPI(...args);

  return [setPayload, status];
}

function OrganizationAddModal({ onClose, onAdd }) {
  const [add, status] = useAddOrganization();

  useEffect(() => {
    if (status.isSuccess) {
      message.success("New organization added successfully!");
      onClose();
      onAdd();
    } else if (status.isError) {
      message.error("Oops! Failed to add new organization");
    }
  }, [status, onAdd, onClose]);

  return (
    <Modal
      title="Add Organization"
      okButtonProps={{
        htmlType: "submit",
        key: "submit",
        form: "organization-add-form",
        loading: status.isLoading,
      }}
      onCancel={onClose}
      okText="Add Organization"
      width={600}
      visible
    >
      <Form {...layout} id="organization-add-form" onFinish={add}>
        <Form.Item name="name" label="Name" rules={ruleReuired}>
          <Input />
        </Form.Item>
        <Form.Item name="address" label="Address" rules={ruleReuired}>
          <Input.TextArea rows={4} />
        </Form.Item>
        <Form.Item name="zipcode" label="Zip Code" rules={ruleReuired}>
          <Input />
        </Form.Item>
        <Form.Item name="phoneNumber" label="Phone" rules={ruleReuired}>
          <Input />
        </Form.Item>
      </Form>
    </Modal>
  );
}

export default OrganizationAddModal;
