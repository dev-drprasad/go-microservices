import { Form, message, Modal, Select } from "antd";
import React, { useEffect } from "react";
import { useCustomers, useOrder, useProducts } from "shared/hooks";

const layout = {
  labelCol: { span: 8 },
  wrapperCol: { span: 16 },
};

const ruleRequired = { required: true };
const ruleJustRequired = [ruleRequired];

function OrderAddModal({ onClose, onAdd }) {
  const [, , , add, status] = useOrder();
  const [products, productsStatus] = useProducts();
  const [customers, customersStatus] = useCustomers();

  const handleFinish = (payload) => {
    console.log("payload :>> ", payload);
    add({ customerId: payload.customerId, products: [{ productId: payload.productId, quantity: 1 }] });
  };

  useEffect(() => {
    if (status.isSuccess) {
      message.success("New order placed successfully!");
      onAdd();
    } else if (status.isError) {
      message.error("Oops! Failed to place new order");
    }
  }, [status, onAdd]);

  const productOpts = products.map(({ id, name }) => ({ value: id, label: name }));
  const customerOpts = customers.map(({ id, name }) => ({ value: id, label: name }));

  return (
    <Modal
      title="Place Order"
      okButtonProps={{
        htmlType: "submit",
        key: "submit",
        form: "order-place-form",
        loading: status.isLoading,
      }}
      onCancel={onClose}
      okText="Place Order"
      width={600}
      visible
    >
      <Form {...layout} id="order-place-form" onFinish={handleFinish}>
        <Form.Item name="productId" label="Product" rules={ruleJustRequired}>
          <Select options={productOpts} loading={productsStatus.isLoading} />
        </Form.Item>
        <Form.Item name="customerId" label="Customer" rules={ruleJustRequired}>
          <Select options={customerOpts} loading={customersStatus.isLoading} />
        </Form.Item>
      </Form>
    </Modal>
  );
}

export default OrderAddModal;
