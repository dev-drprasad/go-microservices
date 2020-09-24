import "./styles.scss";

import { Button, Table } from "antd";
import React, { useMemo, useState } from "react";
import { ListActions, NSHandler, Search } from "shared/components";
import useAPI from "shared/hooks";
import { listsearch } from "shared/utils";

import OrganizationAddModal from "./OrganizationAddModal";

const { Column } = Table;

const searchFields = ["id", "name", "address", "phoneNumber", "zipcode"];

function OrganizationList() {
  const [shouldShowOrganizationAddModal, setShouldShowOrganizationAddModal] = useState(false);
  const [physicians = [], status, refresh] = useAPI("/api/v1/organizations");

  const showOrganizationAddModal = () => setShouldShowOrganizationAddModal(true);
  const closeOrganizationAddModal = () => setShouldShowOrganizationAddModal(false);

  const [searchText, setSearchText] = useState("");
  const searched = useMemo(() => listsearch(physicians, searchFields, searchText), [physicians, searchText]);

  return (
    <div className="organizations-container">
      <ListActions>
        <Search placeholder="Search for anything..." onSearch={setSearchText} style={{ width: 320 }} size="large" />
        <Button type="primary" onClick={showOrganizationAddModal} size="large">
          Add Organization
        </Button>
      </ListActions>
      <NSHandler status={status}>
        {() => (
          <Table dataSource={searched} rowKey="id">
            <Column title="Name" dataIndex="name" />
            <Column title="Address" dataIndex="address" />
            <Column title="Zip Code" dataIndex="zipcode" />
            <Column title="Phone Number" dataIndex="phoneNumber" />
          </Table>
        )}
      </NSHandler>
      {shouldShowOrganizationAddModal && <OrganizationAddModal onClose={closeOrganizationAddModal} onAdd={refresh} />}
    </div>
  );
}

export default OrganizationList;
