import React, { useMemo, useState, useEffect } from "react";
import { Modal, Form, Input, message, Select } from "antd";
import useBROAPI, { useBranches } from "shared/hooks";

const layout = {
  labelCol: { span: 8 },
  wrapperCol: { span: 16 },
};

const roleOptions = [
  { label: "admin", value: "admin" },
  { label: "staff", value: "staff" },
];

const ruleConfirmPassword = ({ getFieldValue }) => ({
  validator(_, value) {
    if (!value || getFieldValue("password") === value) {
      return Promise.resolve();
    }
    return Promise.reject("The two passwords you entered do not match!");
  },
});

const ruleRequired = { required: true };
const ruleJustRequired = [ruleRequired];

function useAddUser() {
  const [user, setUser] = useState(undefined);
  const args = useMemo(() => (user ? ["/api/v1/users", { method: "POST", body: JSON.stringify(user) }] : [undefined, undefined]), [
    user,
  ]);
  const [, status] = useBROAPI(...args);

  return [setUser, status];
}

function UserAddModal({ onClose, onAdd }) {
  const [branches, branchesStatus] = useBranches();
  const [add, status] = useAddUser();

  useEffect(() => {
    if (status.isSuccess) {
      message.success("New user added successfully!");
      onClose();
      onAdd();
    } else if (status.isError) {
      message.error("Oops! Failed to add new user");
    }
  }, [status]);

  const branchOptions = branches.map(({ id, name }) => ({ value: id, label: name }));

  return (
    <Modal
      title="Add User"
      // onOk={addInsurer}
      okButtonProps={{ htmlType: "submit", key: "submit", form: "user-add-form", loading: status.isLoading }}
      onCancel={onClose}
      okText="Add User"
      width={600}
      visible
    >
      <Form {...layout} id="user-add-form" onFinish={add}>
        <Form.Item name="name" label="Name" rules={ruleJustRequired}>
          <Input />
        </Form.Item>
        <Form.Item name="username" label="Username" rules={ruleJustRequired}>
          <Input />
        </Form.Item>
        <Form.Item name="password" label="Password" rules={ruleJustRequired}>
          <Input.Password />
        </Form.Item>
        <Form.Item name="confirmPassword" label="Re-type password" rules={[...ruleJustRequired, ruleConfirmPassword]}>
          <Input.Password />
        </Form.Item>
        <Form.Item name="branchId" label="Branch" rules={ruleJustRequired}>
          <Select options={branchOptions} loading={branchesStatus.isLoading} />
        </Form.Item>
        <Form.Item name="role" label="Role" rules={ruleJustRequired}>
          <Select options={roleOptions} />
        </Form.Item>
      </Form>
    </Modal>
  );
}

export default UserAddModal;
