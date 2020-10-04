import { Button, Form, Input, InputNumber, message, Modal, Select } from "antd";
import { List } from "antd";
import { DeleteOutlined } from "@ant-design/icons";
import React, { useContext, useEffect, useState } from "react";
import { useCustomers, useOrder, useProducts } from "shared/hooks";
import { formatcurrency, listsearch } from "shared/utils";
import { CurrencyContext, LocaleContext } from "shared/contexts";

const layout = {
  labelCol: { span: 8 },
  wrapperCol: { span: 16 },
};

const ruleRequired = { required: true };
const ruleJustRequired = [ruleRequired];

const initialValues = {
  products: {},
};

const computeOrderTotal = (products, selectedProducts) => {
  return products.reduce((acc, p) => acc + (selectedProducts[p.id] || 0) * p.sellPrice, 0);
};

const searchFields = ["name", "id"];

function OrderAddModal({ onClose, onAdd }) {
  const [locale] = useContext(LocaleContext);
  const [currency] = useContext(CurrencyContext);

  const [, , , add, status] = useOrder();
  const [products, productsStatus] = useProducts();
  const [customers, customersStatus] = useCustomers();

  const [form] = Form.useForm();
  const [searchText, setSearchText] = useState("");

  const handleFinish = (payload) => {
    console.log("payload :>> ", payload);
    add({
      customerId: payload.customerId,
      products: Object.entries(form.getFieldValue("products")).map(([id, qty]) => ({ productId: Number(id), quantity: qty })),
    });
  };

  const handleProductAdd = (productId) => () => {
    const selectedProducts = form.getFieldValue("products");
    form.setFieldsValue({ products: { ...selectedProducts, [productId]: (selectedProducts[productId] || 0) + 1 } });
  };

  const handleQuantityChange = (id) => (quantity) => {
    const selectedProducts = form.getFieldValue("products");
    form.setFieldsValue({ products: { ...selectedProducts, [id]: quantity } });
  };

  const handleDelete = (id) => () => {
    handleQuantityChange(id)(0);
  };

  const handleSearch = (searchText) => {
    setSearchText(searchText);
  };

  useEffect(() => {
    if (status.isSuccess) {
      message.success("New order placed successfully!");
      onAdd();
    } else if (status.isError) {
      message.error("Oops! Failed to place new order");
    }
  }, [status, onAdd]);

  const customerOpts = customers.map(({ id, name }) => ({ value: id, label: name }));
  const searched = listsearch(products, searchFields, searchText);
  console.log("products :>> ", products);
  return (
    <Modal
      className="order-place-form-modal"
      title="Place Order"
      okButtonProps={{
        htmlType: "submit",
        key: "submit",
        form: "order-place-form",
        loading: status.isLoading,
      }}
      onCancel={onClose}
      okText="Place Order"
      width={1000}
      visible
    >
      <div className="order-place-form-container">
        <div style={{ display: "grid", gridTemplateRows: "minmax(min-content, max-content) auto", minHeight: 0 }}>
          <Input.Search onSearch={setSearchText} allowClear />
          <List
            className="product-list-form"
            itemLayout="horizontal"
            dataSource={searched}
            loading={productsStatus.isLoading}
            renderItem={(product) => (
              <List.Item
                actions={[
                  <Button onClick={handleProductAdd(product.id)} size="small">
                    + ADD
                  </Button>,
                ]}
              >
                <List.Item.Meta avatar={<img src={`/static/${product.imageUrls[0]}`} />} title={product.name}></List.Item.Meta>
              </List.Item>
            )}
          ></List>
        </div>
        <Form {...layout} id="order-place-form" form={form} onFinish={handleFinish} initialValues={initialValues}>
          <Form.Item className="products-form-item" wrapperCol={{ span: 24 }} rules={ruleJustRequired} shouldUpdate>
            {({ getFieldValue }) => (
              <List
                className="product-list-form"
                itemLayout="horizontal"
                dataSource={Object.entries(getFieldValue("products")).filter(([, q]) => !!q)}
                renderItem={([id, quantity]) => (
                  <List.Item>
                    <List.Item.Meta
                      avatar={<img src={`/static/${products.find((p) => p.id === Number(id))?.imageUrls[0]}`} />}
                      title={products.find((p) => p.id === Number(id))?.name}
                    ></List.Item.Meta>
                    <div className="sel-product">
                      <div className="price">
                        {formatcurrency(locale, currency, products.find((p) => p.id === Number(id))?.sellPrice)}
                      </div>
                      <InputNumber value={quantity} onChange={handleQuantityChange(id)} size="small" />
                      <Button icon={<DeleteOutlined />} onClick={handleDelete(id)} size="small" />
                    </div>
                  </List.Item>
                )}
              ></List>
            )}
          </Form.Item>
          <Form.Item labelAlign="left" name="customerId" label="Customer" rules={ruleJustRequired}>
            <Select options={customerOpts} loading={customersStatus.isLoading} showSearch />
          </Form.Item>
          <Form.Item wrapperCol={{ span: 24 }} shouldUpdate>
            {({ getFieldValue }) => (
              <div style={{ textAlign: "right", fontSize: "1.2em" }}>
                Order Total: {formatcurrency(locale, currency, computeOrderTotal(products, getFieldValue("products")))}
              </div>
            )}
          </Form.Item>
        </Form>
      </div>
    </Modal>
  );
}

export default OrderAddModal;
