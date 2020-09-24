import React, { useState, useMemo, useContext } from "react";
import { Button, Table, DatePicker } from "antd";
import { Link } from "@reach/router";
import useBROAPI from "shared/hooks";
import { NSHandler, Search, ListActions } from "shared/components";
import "./styles.scss";
import { formatcurrency, listsearch } from "shared/utils";
import { CurrencyContext } from "shared/contexts";

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
const searchFields = ["id", "serviceDate", "status", "orderedBy.firstName", "orderedBy.lastName", "prescribedBy.name"];
function ProductList({ navigate }) {
  const [{ searchText, filters }, setSearchText] = useState({ searchText: "", filters: [] });
  const currency = useContext(CurrencyContext);

  const [products = [], status] = useBROAPI("/api/v1/products");
  const handleDateSelect = (_, dateStr) => {
    setSearchText({
      searchText,
      filters: dateStr
        ? [
            ["serviceDate", dateStr],
            ["status", "scheduled"],
          ]
        : [],
    });
  };
  const handleSearch = (searchText) => {
    setSearchText({ filters, searchText });
  };

  const onRow = (product) => {
    return {
      onClick: (e) => {
        if (e.target.tagName !== "A") {
          navigate(`/products/${product.id}`);
        }
      },
    };
  };

  const filtered = useMemo(() => products.filter((o) => filters.every(([f, v]) => o[f] === v)), [products, filters]);

  const searched = useMemo(() => listsearch(filtered, searchFields, searchText), [searchText, filtered]);

  return (
    <div className="products-container">
      <ListActions>
        <Search placeholder="Search for anything..." onSearch={handleSearch} style={{ width: 320 }} size="large" />
        <DatePicker
          placeholder="Delivery Scheduled On"
          style={{ marginLeft: 8, width: 220 }}
          onChange={handleDateSelect}
          size="large"
        />

        <Button type="primary" onClick={() => navigate("/products/new")} size="large">
          Add Product
        </Button>
      </ListActions>
      <NSHandler status={status}>
        {() => (
          <Table className="row-clickable" dataSource={searched} rowKey="id" onRow={onRow}>
            <Column title="Name" dataIndex="name" />
            <Column title="Price" dataIndex="sellPrice" render={(p) => formatcurrency(currency, p)} />
            <Column title="Stock" dataIndex="stock" />
          </Table>
        )}
      </NSHandler>
    </div>
  );
}

export default ProductList;
