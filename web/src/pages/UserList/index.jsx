import "./styles.scss";

import { Button, Table } from "antd";
import React, { useState, useMemo } from "react";
import { NSHandler, ListActions, Search } from "shared/components";
import useBROAPI from "shared/hooks";

import UserAddModal from "./UserAddModal";
import { listsearch } from "shared/utils";

const { Column } = Table;

const searchFields = ["id", "name", "username", "team.name"];

function UserList() {
  const [shouldShowUserAddModal, setShouldShowUserAddModal] = useState(false);
  const [users = [], status, refresh] = useBROAPI("/api/v1/users");

  const showUserAddModal = () => setShouldShowUserAddModal(true);
  const closeUserAddModal = () => setShouldShowUserAddModal(false);
  const [searchText, setSearchText] = useState("");
  const searched = useMemo(() => listsearch(users, searchFields, searchText), [users, searchText]);
  return (
    <div className="users-container">
      <ListActions>
        <Search placeholder="Search for anything..." onSearch={setSearchText} style={{ width: 320 }} size="large" />
        <Button type="primary" onClick={showUserAddModal} size="large">
          Add User
        </Button>
      </ListActions>
      <NSHandler status={status}>
        {() => (
          <Table dataSource={searched} rowKey="id">
            <Column title="Userame" dataIndex="username" />
            <Column title="Name" dataIndex="name" />
            <Column title="Role" dataIndex="role" />
          </Table>
        )}
      </NSHandler>
      {shouldShowUserAddModal && <UserAddModal onClose={closeUserAddModal} onAdd={refresh} />}
    </div>
  );
}

export default UserList;
