import { EyeOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, Card, Form, Input, InputNumber, Radio, Select, Upload } from "antd";
import BrandAddModal from "pages/BrandList/BrandAddModal";
import CategoryAddModal from "pages/CategoryList/CategoryAddModal";
import React, { useContext, useState } from "react";
import { CurrencyContext, LocaleContext } from "shared/contexts";
import { useBrands, useCategories } from "shared/hooks";
import { formatcurrency } from "shared/utils";
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

function ProductForm({ id, initialValues, uploadHeaders, onFinish }) {
  const [form] = Form.useForm();

  const [currency] = useContext(CurrencyContext);
  const [locale] = useContext(LocaleContext);

  const [imagePreviewUrl, setImagePreviewUrl] = useState();
  const [shouldShowBrandAddModal, setShouldShowBrandAddModal] = useState(false);
  const [shouldShowCategoryAddModal, setShouldShowCategoryAddModal] = useState(false);
  const [brands, brandsStatus, refreshBrands] = useBrands();
  const [categories, categoriesStatus, refreshCategories] = useCategories();

  const handleSubmit = (payload) => {
    payload.sellPrice = getSellPrice(payload.cost, payload.priceCalcMode, payload.priceCalcValue);
    payload.imageUrls = [...form.getFieldValue("imageUrls")];
    console.log("payload :>> ", payload);
    onFinish(payload);
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

  const handleImageChange = ({ file }) => {
    if (file.status === "done" && file.response) {
      // this event will fire multiple times
      // So, made form.imageUrls as `Set` to avoid duplication
      const imageURLs = form.getFieldValue("imageUrls");
      file.response.forEach(imageURLs.add, imageURLs);
    } else if (file.status === "removed") {
      form.getFieldValue("imageUrls").delete(file.uid);
    }
  };

  const brandOptions = brands.map(({ id, name }) => ({
    value: id,
    label: name,
  }));
  const categoryOptions = categories.map(({ id, name }) => ({
    value: id,
    label: name,
  }));

  return (
    <div className="product-form-container">
      <Form form={form} {...layout} className="product-form" id={id} onFinish={handleSubmit} initialValues={initialValues}>
        <Form.Item name="name" label="Name" rules={ruleJustRequired}>
          <Input />
        </Form.Item>
        <Form.Item label="Price">
          <Input.Group compact>
            <Form.Item name="cost" rules={ruleJustRequired}>
              <InputNumber style={{ width: 200 }} />
            </Form.Item>
            <Seperator>+</Seperator>
            <Form.Item name="priceCalcValue" noStyle>
              <InputNumber />
            </Form.Item>
            <Form.Item name="priceCalcMode" noStyle>
              <Radio.Group options={priceCalcMethods} optionType="button" buttonStyle="solid" />
            </Form.Item>
            <Form.Item dependencies={["cost", "priceCalcValue", "priceCalcMode"]}>
              {({ getFieldValue }) => renderSellPrice(getFieldValue, (v) => formatcurrency(locale, currency, v))}
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
        <Form.Item></Form.Item>
        <Upload
          multiple
          action="/api/v1/products/images/upload"
          headers={uploadHeaders}
          accept=".jpg, .jpeg"
          listType="picture-card"
          name="files"
          onPreview={handleImagePreview}
          onChange={handleImageChange}
          defaultFileList={[...initialValues.imageUrls].map((u) => ({ url: `/static/${u}`, uid: u }))}
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

export default ProductForm;
