import "./styles.scss";

import { Button, Table } from "antd";
import React, { useMemo, useState } from "react";
import { ListActions, NSHandler, Search } from "../CategoryList/node_modules/shared/components";
import useAPI from "../CategoryList/node_modules/shared/hooks";
import { listsearch } from "../CategoryList/node_modules/shared/utils";

import BrandAddModal from "./BrandAddModal";

const { Column } = Table;

const searchFields = ["id", "name", "address", "phoneNumber", "zipcode"];

function BrandList() {
  const [shouldShowBranchAddModal, setShouldShowBranchAddModal] = useState(false);
  const [physicians = [], status, refresh] = useAPI("/api/v1/brands");

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
      {shouldShowBranchAddModal && <BrandAddModal onClose={closeBranchAddModal} onAdd={refresh} />}
    </div>
  );
}

export default BrandList;
