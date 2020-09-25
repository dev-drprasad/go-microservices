import { Button, Table } from "antd";
import React, { useMemo, useState } from "react";
import { ListActions, NSHandler, Search } from "shared/components";
import { useCustomers } from "shared/hooks";
import { listsearch } from "shared/utils";
import CustomerAddModal from "./CustomerAddModal";
import "./styles.scss";

const { Column } = Table;

const searchFields = ["id", "name", "address", "zipcode", "phoneNumber"];

function CustomerList({ navigate }) {
  const [shouldShowCustomerAddModal, setShouldShowCustomerAddModal] = useState(false);
  const [{ searchText, filters }, setSearchText] = useState({ searchText: "", filters: [] });

  const [customers, status, refresh] = useCustomers();

  const showCustomerAddModal = () => setShouldShowCustomerAddModal(true);
  const closeCustomerAddModal = () => setShouldShowCustomerAddModal(false);

  const handleCustomerAdd = () => {
    closeCustomerAddModal();
    refresh();
  };

  const handleSearch = (searchText) => {
    setSearchText({ filters, searchText });
  };

  const onRow = (product) => {
    return {
      onClick: (e) => {
        if (e.target.tagName !== "A") {
          navigate(`/customers/${product.id}/edit`);
        }
      },
    };
  };

  const searched = useMemo(() => listsearch(customers, searchFields, searchText), [searchText, customers]);

  return (
    <div className="customers-container">
      <ListActions>
        <Search placeholder="Search for anything..." onSearch={handleSearch} style={{ width: 320 }} size="large" />
        <Button type="primary" onClick={showCustomerAddModal} size="large">
          Add Customer
        </Button>
      </ListActions>
      <NSHandler status={status}>
        {() => (
          <Table className="row-clickable" dataSource={searched} rowKey="id" onRow={onRow}>
            <Column title="Name" dataIndex="name" />
            <Column title="Address" dataIndex="address" />
            <Column title="Zipcode" dataIndex="zipcode" />
            <Column title="Phone Number" dataIndex="phoneNumber" />
          </Table>
        )}
      </NSHandler>
      {shouldShowCustomerAddModal && <CustomerAddModal onClose={closeCustomerAddModal} onAdd={handleCustomerAdd} />}
    </div>
  );
}

export default CustomerList;
