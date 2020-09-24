import React, { useContext, useEffect } from "react";
import ProductForm from "pages/ProductForm";
import NotFound from "pages/NotFound";
import { useProduct } from "shared/hooks";
import { NSHandler } from "shared/components";
import { Button, message } from "antd";
import { AuthContext, CurrencyContext } from "shared/contexts";
import "./styles.scss";

const formID = "product-edit-form";

const denormalize = (product) => ({
  ...product,
  imageUrls: new Set(product.imageUrls),
  priceCalcValue: ((product.sellPrice - product.cost) / product.cost) * 100,
  priceCalcMode: "%",
});

function ProductEdit({ navigate, id: idStr }) {
  const idParsed = Number(idStr);
  const id = !Number.isNaN(idParsed) && idParsed > 0 ? idParsed : undefined;

  const currency = useContext(CurrencyContext);
  const [user] = useContext(AuthContext);

  const [product, status, , update, updateStatus] = useProduct(id);

  useEffect(() => {
    if (updateStatus.isSuccess) {
      message.success("Product updated successfully!");
      navigate("/products");
    } else if (updateStatus.isError) {
      message.error("Oops! Failed to update product");
    }
  }, [updateStatus, navigate]);

  if (!id) return <NotFound />;

  return (
    <NSHandler status={status}>
      {() => (
        <div className="product-edit-container">
          <ProductForm
            id={formID}
            onFinish={update}
            uploadHeaders={{ Authorization: `Bearer ${user.token}` }}
            currency={currency}
            initialValues={denormalize(product)}
          />
          <Button className="right-align" type="primary" form={formID} htmlType="submit" loading={status.isLoading}>
            Update Product
          </Button>
        </div>
      )}
    </NSHandler>
  );
}

export default ProductEdit;
