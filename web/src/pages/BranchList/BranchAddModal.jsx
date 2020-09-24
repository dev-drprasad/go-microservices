import React, { useMemo, useState, useEffect } from "react";
import { Modal, Form, Input, message, Select } from "antd";
import useAPI, { useOrganizations } from "shared/hooks";

const layout = {
  labelCol: { span: 8 },
  wrapperCol: { span: 16 },
};

const ruleRequired = { required: true };
const ruleJustRequired = [ruleRequired];

function useAddBranch() {
  const [payload, setPayload] = useState(undefined);
  const args = useMemo(
    () => (payload ? ["/api/v1/branches", { method: "POST", body: JSON.stringify(payload) }] : [undefined, undefined]),
    [payload]
  );
  const [, status] = useAPI(...args);

  return [setPayload, status];
}

function BranchAddModal({ onClose, onAdd }) {
  const [organizations, organizationsStatus] = useOrganizations();
  const [add, status] = useAddBranch();

  useEffect(() => {
    if (status.isSuccess) {
      message.success("New branch added successfully!");
      onClose();
      onAdd();
    } else if (status.isError) {
      message.error("Oops! Failed to add new branch");
    }
  }, [status, onAdd, onClose]);

  const organizationOptions = organizations.map(({ id, name }) => ({ value: id, label: name }));

  return (
    <Modal
      title="Add Branch"
      okButtonProps={{
        htmlType: "submit",
        key: "submit",
        form: "branch-add-form",
        loading: status.isLoading,
      }}
      onCancel={onClose}
      okText="Add Branch"
      width={600}
      visible
    >
      <Form {...layout} id="branch-add-form" onFinish={add}>
        <Form.Item name="name" label="Name" rules={ruleJustRequired}>
          <Input />
        </Form.Item>
        <Form.Item name="organizationId" label="Organization" rules={ruleJustRequired}>
          <Select options={organizationOptions} loading={organizationsStatus.isLoading} />
        </Form.Item>
        <Form.Item name="address" label="Address" rules={ruleJustRequired}>
          <Input.TextArea rows={4} />
        </Form.Item>
        <Form.Item name="zipcode" label="Zip Code" rules={ruleJustRequired}>
          <Input />
        </Form.Item>
        <Form.Item name="phoneNumber" label="Phone" rules={ruleJustRequired}>
          <Input />
        </Form.Item>
      </Form>
    </Modal>
  );
}

export default BranchAddModal;
