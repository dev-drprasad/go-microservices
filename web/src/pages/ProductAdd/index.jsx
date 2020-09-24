import React, { useState, useMemo, useEffect, useContext } from "react";
import { Form, Select, Button, message, Input, InputNumber, Radio, Upload, Card } from "antd";
import BrandAddModal from "pages/BrandList/BrandAddModal";
import CategoryAddModal from "pages/CategoryList/CategoryAddModal";
import { PlusOutlined, EyeOutlined } from "@ant-design/icons";
import useBROAPI, { useCategories, useBrands } from "shared/hooks";
import { formatcurrency } from "shared/utils";
import { AuthContext, CurrencyContext } from "shared/contexts";
import "./styles.scss";

const layout = {
  labelCol: { span: 4 },
  wrapperCol: { span: 20 },
};

const priceCalcMethods = [
  { label: "%", value: "%" },
  { label: "#", value: "#" },
];

const ruleRequired = { required: true };
const ruleJustRequired = [ruleRequired];

const getSellPrice = (cost, priceCalcMode, priceCalcValue) => {
  switch (priceCalcMode) {
    case "%":
      return cost + (cost * priceCalcValue) / 100;
    default:
      return cost + priceCalcValue;
  }
};

function getBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = (error) => reject(error);
  });
}

const Seperator = ({ children }) => <span style={{ padding: "0 16px", fontSize: 16 }}>{children}</span>;

const renderSellPrice = (getFieldValue, formatter) => {
  const sellPrice = getSellPrice(getFieldValue("cost"), getFieldValue("priceCalcMode"), getFieldValue("priceCalcValue"));
  return <Seperator> = {formatter(sellPrice)}</Seperator>;
};

function useProductAdd() {
  const [order, setOrder] = useState(undefined);
  const args = useMemo(() => (order ? ["/api/v1/products", { method: "POST", body: JSON.stringify(order) }] : [undefined, undefined]), [
    order,
  ]);
  const [, status] = useBROAPI(...args);

  return [setOrder, status];
}

function ProductAdd({ navigate }) {
  const currency = useContext(CurrencyContext);
  const [user] = useContext(AuthContext);
  const [form] = Form.useForm();

  const [imagePreviewUrl, setImagePreviewUrl] = useState();
  const [shouldShowBrandAddModal, setShouldShowBrandAddModal] = useState(false);
  const [shouldShowCategoryAddModal, setShouldShowCategoryAddModal] = useState(false);
  const [brands, brandsStatus, refreshBrands] = useBrands();
  const [categories, categoriesStatus, refreshCategories] = useCategories();
  const [add, status] = useProductAdd();

  const addProduct = (payload) => {
    payload.sellPrice = getSellPrice(payload.cost, payload.priceCalcMode, payload.priceCalcValue);
    console.log("payload :>> ", payload);
    add(payload);
  };

  const showBrandAddModal = () => setShouldShowBrandAddModal(true);
  const closeBrandAddModal = () => setShouldShowBrandAddModal(false);

  const handleBrandAdd = () => {
    closeBrandAddModal();
    refreshBrands();
  };

  const showCategoryAddModal = () => setShouldShowCategoryAddModal(true);
  const closeCategoryAddModal = () => setShouldShowCategoryAddModal(false);

  const handleCategoryAdd = () => {
    closeCategoryAddModal();
    refreshCategories();
  };

  const handleImagePreview = async (file) => {
    const imageUrl = file.url || file.preview;
    if (imageUrl) {
      setImagePreviewUrl(file.url);
    } else {
      const imageUrl = await getBase64(file.originFileObj);
      setImagePreviewUrl(imageUrl);
    }
  };

  useEffect(() => {
    if (status.isSuccess) {
      message.success("New product added successfully!");
      navigate("/products");
    } else if (status.isError) {
      message.error("Oops! Failed to add new product");
    }
  }, [status, navigate]);

  const brandOptions = brands.map(({ id, name }) => ({
    value: id,
    label: name,
  }));
  const categoryOptions = categories.map(({ id, name }) => ({
    value: id,
    label: name,
  }));

  return (
    <div className="product-add-container">
      <Form form={form} {...layout} id="product-create-form" onFinish={addProduct}>
        <Form.Item name="name" label="Name" rules={ruleJustRequired}>
          <Input />
        </Form.Item>
        <Form.Item label="Price">
          <Input.Group compact>
            <Form.Item name="cost" initialValue={0} rules={ruleJustRequired}>
              <InputNumber style={{ width: 200 }} />
            </Form.Item>
            <Seperator>+</Seperator>
            <Form.Item name="priceCalcValue" initialValue={10} noStyle>
              <InputNumber />
            </Form.Item>
            <Form.Item name="priceCalcMode" initialValue="%" noStyle>
              <Radio.Group options={priceCalcMethods} optionType="button" buttonStyle="solid" />
            </Form.Item>
            <Form.Item dependencies={["cost", "priceCalcValue", "priceCalcMode"]}>
              {({ getFieldValue }) => renderSellPrice(getFieldValue, (v) => formatcurrency(currency, v))}
            </Form.Item>
          </Input.Group>
        </Form.Item>
        <Form.Item className="inline" wrapperCol={{ span: 24 }}>
          <Form.Item
            name="brandId"
            label="Brand"
            style={{ width: "50%" }}
            labelCol={{ span: 8 }}
            wrapperCol={{ span: 16 }}
            rules={ruleJustRequired}
          >
            <Select
              options={brandOptions}
              loading={brandsStatus.isLoading}
              dropdownRender={(menu) => (
                <>
                  {menu}
                  <hr />
                  <Button type="link" icon={<PlusOutlined />} onClick={showBrandAddModal} block>
                    Add new brand
                  </Button>
                </>
              )}
              showSearch
            />
          </Form.Item>
          <Form.Item
            name="categoryId"
            label="Category"
            style={{ width: "50%" }}
            rules={ruleJustRequired}
            labelCol={{ span: 8 }}
            wrapperCol={{ span: 16 }}
          >
            <Select
              options={categoryOptions}
              loading={categoriesStatus.isLoading}
              dropdownRender={(menu) => (
                <>
                  {menu}
                  <hr />
                  <Button type="link" icon={<PlusOutlined />} onClick={showCategoryAddModal} block>
                    Add new category
                  </Button>
                </>
              )}
              showSearch
            />
          </Form.Item>
        </Form.Item>
        <Form.Item name="stock" label="Stock" rules={ruleJustRequired}>
          <InputNumber min={1} />
        </Form.Item>
        <Button className="right-align" type="primary" htmlType="submit" loading={status.isLoading}>
          Add Product
        </Button>
        {shouldShowBrandAddModal && <BrandAddModal onClose={closeBrandAddModal} onAdd={handleBrandAdd} />}
        {shouldShowCategoryAddModal && <CategoryAddModal onClose={closeCategoryAddModal} onAdd={handleCategoryAdd} />}
      </Form>
      <Card title="Product Images" className="image-upload-container" size="small">
        <div className="preview-image-container">
          {imagePreviewUrl ? (
            <img src={imagePreviewUrl} alt="preview" />
          ) : (
            <span>
              click on <EyeOutlined /> in thumbline to preview
            </span>
          )}
        </div>
        <Upload
          multiple
          action="/api/v1/products/images/upload"
          headers={{ Authorization: `Bearer ${user.token}` }}
          accept=".jpg, .jpeg"
          listType="picture-card"
          name="files"
          onPreview={handleImagePreview}
        >
          <div>
            <PlusOutlined />
            <div style={{ marginLeft: 8 }}>Upload</div>
          </div>
        </Upload>
      </Card>
    </div>
  );
}

export default ProductAdd;
