import { EyeOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, Card, Form, Input, InputNumber, Select, Upload } from "antd";
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

const ruleRequired = { required: true };
const ruleGreaterThan0 = [ruleRequired, { type: "number", min: 1 }];
const ruleJustRequired = [ruleRequired];

function getBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = (error) => reject(error);
  });
}

const round = (n, fractions = 2) => Math.round(n.toFixed(fractions) * 10 ** fractions) / 10 ** fractions;

const renderSellPrice = (getFieldValue, formatter) => {
  const cost = getFieldValue("cost");
  const sellPrice = getFieldValue("sellPrice");
  return (
    <div style={{ fontSize: "1rem", textAlign: "center" }}>
      {formatter(round(sellPrice - cost))}
      <span style={{ fontSize: "0.6rem" }}>({round((sellPrice - cost) / 100)}%)</span>
    </div>
  );
};

function ProductForm({ id, initialValues, uploadHeaders, onFinish }) {
  const [form] = Form.useForm();

  const [currency] = useContext(CurrencyContext);
  const [locale] = useContext(LocaleContext);

  const [imagePreviewUrl, setImagePreviewUrl] = useState();
  const [shouldShowBrandAddModal, setShouldShowBrandAddModal] = useState(false);
  const [shouldShowCategoryAddModal, setShouldShowCategoryAddModal] = useState(false);
  const [uploadedImages, setUploadedImages] = useState(initialValues.imageUrls.size);
  const [brands, brandsStatus, refreshBrands] = useBrands();
  const [categories, categoriesStatus, refreshCategories] = useCategories();

  const handleSubmit = (payload) => {
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
    const imageURLs = form.getFieldValue("imageUrls");
    if (file.status === "done" && file.response) {
      // this event will fire multiple times
      // So, made form.imageUrls as `Set` to avoid duplication
      file.response.forEach(imageURLs.add, imageURLs);
    } else if (file.status === "removed") {
      imageURLs.delete(file.uid);
    }
    setUploadedImages(imageURLs.size);
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
        <div className="flex">
          <Form.Item
            label="Price"
            style={{ width: "40%" }}
            name="sellPrice"
            rules={ruleJustRequired}
            labelCol={{ span: 10 }}
            wrapperCol={{ span: 14 }}
          >
            <InputNumber style={{ width: "100%" }} />
          </Form.Item>
          <Form.Item
            label="Cost"
            style={{ width: "35%" }}
            name="cost"
            rules={ruleJustRequired}
            labelCol={{ span: 8 }}
            wrapperCol={{ span: 16 }}
          >
            <InputNumber style={{ width: "100%" }} />
          </Form.Item>
          <Form.Item style={{ width: "25%" }} wrapperCol={{ span: 24 }} dependencies={["cost", "sellPrice"]}>
            {({ getFieldValue }) => renderSellPrice(getFieldValue, (v) => formatcurrency(locale, currency, v))}
          </Form.Item>
        </div>
        <div className="flex">
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
        </div>
        <Form.Item name="stock" label="Stock" rules={ruleGreaterThan0}>
          <InputNumber min={0} />
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
          {uploadedImages < 6 && (
            <div>
              <PlusOutlined />
              <div style={{ marginLeft: 8 }}>Upload</div>
            </div>
          )}
        </Upload>
      </Card>
    </div>
  );
}

export default ProductForm;
