import Icon from "@ant-design/icons";
import { Link } from "@reach/router";
import { Button, Card, Carousel, Radio, Table } from "antd";
import React, { memo, useContext, useEffect, useMemo, useRef, useState } from "react";
import { ListActions, NSHandler, Search } from "shared/components";
import { CurrencyContext, LocaleContext } from "shared/contexts";
import useAPI from "shared/hooks";
import { ReactComponent as InStockIcon } from "shared/icons/instock.svg";
import { formatcurrency, listsearch } from "shared/utils";
import "./styles.scss";

const { Column } = Table;

function getPatientNameFromOrder(order) {
  return order.orderedBy.firstName + " " + order.orderedBy.lastName;
}

function patientNameAnchored(_, order) {
  const patientName = getPatientNameFromOrder(order);
  return <Link to={`/patients/${order.orderedBy.accountId}`}>{patientName}</Link>;
}

function physicianNameAnchored(_, order) {
  return <Link to={`/physicians/${order.prescribedBy.id}`}>{order.prescribedBy.name}</Link>;
}
function orderIdAnchored(id) {
  return <Link to={`/orders/${id}`}>{id}</Link>;
}

const strictEqCheck = (v1) => (v2) => v1 === v2;

const [InStock, LowStock, OutOfStock] = [1, 2, 0];
const filterFns = {
  stock: (value) => {
    if (value === InStock) {
      return (v) => v > 0;
    } else if (value === LowStock) {
      return (v) => v <= 5;
    } else if (value === OutOfStock) {
      return (v) => v === 0;
    }
  },
};

const searchFields = ["id", "name"];

const StockFilter = memo(
  function ({ onChange }) {
    const didMount = useRef(false);
    const [value, setValue] = useState();
    const handleChange = (e) => {
      setValue(e.target.value);
    };
    const handleClick = (e) => {
      if (e.target.checked) {
        setValue(undefined);
      }
    };

    useEffect(() => {
      if (!didMount.current) {
        didMount.current = true;
        return;
      }
      onChange(value);
    }, [onChange, value]);

    return (
      <Radio.Group onChange={handleChange} value={value}>
        <Radio.Button value={InStock} onClick={handleClick}>
          In Stock
        </Radio.Button>
        <Radio.Button value={LowStock} onClick={handleClick}>
          Low Stock
        </Radio.Button>
        <Radio.Button value={OutOfStock} onClick={handleClick}>
          Out of Stock
        </Radio.Button>
      </Radio.Group>
    );
  },
  () => true
);

function ProductList({ navigate }) {
  const [{ searchText, filters }, setSearchText] = useState({ searchText: "", filters: {} });
  const [currency] = useContext(CurrencyContext);
  const [locale] = useContext(LocaleContext);

  const [products = [], status] = useAPI("/api/v1/products");

  const handleStockFilterChange = (value) => {
    setSearchText((s) => ({ ...s, filters: { ...s.filters, stock: value } }));
  };

  const handleSearch = (searchText) => {
    setSearchText((s) => ({ ...s, searchText }));
  };

  const onRow = (product) => {
    return {
      onClick: (e) => {
        if (e.target.tagName !== "A") {
          navigate(`/products/${product.id}/edit`);
        }
      },
    };
  };

  const handleCardClick = (id) => () => {
    navigate(`/products/${id}/edit`);
  };

  const filtered = useMemo(
    () =>
      products.filter((o) =>
        Object.entries(filters).every(([fieldName, filterValue]) => {
          const fn = (filterFns[fieldName] || strictEqCheck)(filterValue);
          return fn ? fn(o[fieldName]) : true;
        })
      ),
    [products, filters]
  );

  const searched = useMemo(() => listsearch(filtered, searchFields, searchText), [searchText, filtered]);

  return (
    <div className="products-container">
      <ListActions>
        <Search placeholder="Search for anything..." onSearch={handleSearch} style={{ width: 320 }} />
        <StockFilter onChange={handleStockFilterChange} />
        <Button type="primary" onClick={() => navigate("/products/new")} size="large">
          Add Product
        </Button>
      </ListActions>
      <NSHandler status={status}>
        {() => (
          <div className="product-list">
            {searched.map((product) => (
              <Card
                size="small"
                key={product.id}
                className="product-card"
                cover={
                  <Carousel draggable={false} lazyLoad="ondemand">
                    {product.imageUrls.map((url) => (
                      <img key={url} alt={product.name} src={`/static/${url}`} />
                    ))}
                  </Carousel>
                }
              >
                <Card.Meta
                  title={
                    <a href={`/products/${product.id}/edit`} title={product.name}>
                      {product.name}
                    </a>
                  }
                  description=""
                />
                <div style={{ display: "flex", justifyContent: "space-between" }}>
                  <div className="stock">
                    <Icon component={InStockIcon} />
                    <span>{product.stock}</span>
                  </div>
                  <div className="price">{formatcurrency(locale, currency, product.sellPrice)}</div>
                </div>
              </Card>
            ))}
          </div>
        )}
      </NSHandler>
    </div>
  );
}

export default ProductList;
