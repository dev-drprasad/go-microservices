import "./styles.scss";

import { Button, Table } from "antd";
import React, { useMemo, useState } from "react";
import { ListActions, NSHandler, Search } from "shared/components";
import useBROAPI from "shared/hooks";
import { listsearch } from "shared/utils";

import BranchAddModal from "./BranchAddModal";

const { Column } = Table;

const searchFields = ["id", "name", "address", "phoneNumber", "zipcode"];

function BranchList() {
  const [shouldShowBranchAddModal, setShouldShowBranchAddModal] = useState(false);
  const [physicians = [], status, refresh] = useBROAPI("/api/v1/branches");

  const showBranchAddModal = () => setShouldShowBranchAddModal(true);
  const closeBranchAddModal = () => setShouldShowBranchAddModal(false);

  const [searchText, setSearchText] = useState("");
  const searched = useMemo(() => listsearch(physicians, searchFields, searchText), [physicians, searchText]);

  return (
    <div className="branches-container">
      <ListActions>
        <Search placeholder="Search for anything..." onSearch={setSearchText} style={{ width: 320 }} size="large" />
        <Button type="primary" onClick={showBranchAddModal} size="large">
          Add Branch
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
      {shouldShowBranchAddModal && <BranchAddModal onClose={closeBranchAddModal} onAdd={refresh} />}
    </div>
  );
}

export default BranchList;
