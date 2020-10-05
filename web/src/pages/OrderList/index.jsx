import { Button, Table, Tag } from "antd";
import React, { useContext, useMemo, useState } from "react";
import { ListActions, NSHandler, Search } from "shared/components";
import { CurrencyContext, LocaleContext } from "shared/contexts";
import { useOrders } from "shared/hooks";
import { formatcurrency, listsearch } from "shared/utils";
import OrderAddModal from "./OrderAddModal";
import "./styles.scss";

const { Column } = Table;

const statusTagColor = {
  new: "success",
  preparation: "processing",
  ready: "purple",
  delivered: "default",
};

const searchFields = ["id", "status", "customer.name"];

const getOrderTotal = (o) => o.products.reduce((acc, { unitPrice, quantity }) => acc + unitPrice * quantity, 0);

const OrderTotal = ({ order, currency, locale }) => {
  return formatcurrency(locale, currency, getOrderTotal(order));
};

function OrderList({ navigate }) {
  const [currency] = useContext(CurrencyContext);
  const [locale] = useContext(LocaleContext);

  const [shouldShowOrderAddModal, setShouldShowOrderAddModal] = useState(false);
  const [{ searchText, filters }, setSearchText] = useState({ searchText: "", filters: [] });

  const [orders, status, refresh] = useOrders();

  const showOrderAddModal = () => setShouldShowOrderAddModal(true);
  const closeOrderAddModal = () => setShouldShowOrderAddModal(false);

  const handleCustomerAdd = () => {
    closeOrderAddModal();
    refresh();
  };

  const handleSearch = (searchText) => {
    setSearchText({ filters, searchText });
  };

  const onRow = (product) => {
    return {
      onClick: (e) => {
        if (e.target.tagName !== "A") {
          navigate(`/orders/${product.id}/edit`);
        }
      },
    };
  };

  const searched = useMemo(() => listsearch(orders, searchFields, searchText), [searchText, orders]);

  return (
    <div className="customers-container">
      <ListActions>
        <Search placeholder="Search for anything..." onSearch={handleSearch} style={{ width: 320 }} size="large" />
        <Button type="primary" onClick={showOrderAddModal} size="large">
          Place Order
        </Button>
      </ListActions>
      <NSHandler status={status}>
        {() => (
          <Table className="row-clickable" dataSource={searched} rowKey="id" onRow={onRow}>
            <Column title="ID" dataIndex="id" />
            <Column title="Status" dataIndex="status" render={(s) => <Tag color={statusTagColor[s]}>{s.toUpperCase()}</Tag>} />
            <Column title="Total" dataIndex="total" render={(_, o) => <OrderTotal locale={locale} currency={currency} order={o} />} />
            <Column title="Customer" dataIndex={["customer", "name"]} />
          </Table>
        )}
      </NSHandler>
      {shouldShowOrderAddModal && <OrderAddModal onClose={closeOrderAddModal} onAdd={handleCustomerAdd} />}
    </div>
  );
}

export default OrderList;
