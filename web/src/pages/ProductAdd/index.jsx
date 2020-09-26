import React, { useState, useMemo, useEffect, useContext } from "react";
import { Button, message } from "antd";

import { useProduct } from "shared/hooks";

import { AuthContext, CurrencyContext } from "shared/contexts";
import ProductForm from "pages/ProductForm";
import "./styles.scss";

const formID = "product-add-form";

function ProductAdd({ navigate }) {
  const currency = useContext(CurrencyContext);
  const [user] = useContext(AuthContext);

  const [, , , add, status] = useProduct();

  useEffect(() => {
    if (status.isSuccess) {
      message.success("New product added successfully!");
      navigate("/products");
    } else if (status.isError) {
      message.error("Oops! Failed to add new product");
    }
  }, [status, navigate]);

  const initialValues = {
    cost: 0,
    sellPrice: 0,
    imageUrls: new Set([]),
  };

  return (
    <div className="product-add-container">
      <ProductForm
        id={formID}
        onFinish={add}
        uploadHeaders={{ Authorization: `Bearer ${user.token}` }}
        currency={currency}
        initialValues={initialValues}
      />
      <Button className="right-align" type="primary" form={formID} htmlType="submit" loading={status.isLoading}>
        Add Product
      </Button>
    </div>
  );
}

export default ProductAdd;
