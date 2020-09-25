import { Input } from "antd";
import React from "react";
import "./Search.scss";

export default function Search({ onSearch, size, ...rest }) {
  // const [searchText, setSearchText] = useState("");
  // const handleChange = (e) => {
  //   setSearchText(e.target.value);
  // };

  // const handleSubmit = () => onSearch(searchText);
  // const handleClear = () => {
  //   setSearchText("");
  //   onSearch("");
  // };
  return (
    <div className="search-container">
      <Input.Search
        placeholder="Search for anything..."
        onSearch={onSearch}
        // suffix={<CloseCircleOutlined style={{ fontSize: 14, color: "rgba(66, 82, 110, 0.75)" }} />}
        allowClear
      />
    </div>
  );
}
