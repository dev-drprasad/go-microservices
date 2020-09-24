import React, { useContext } from "react";
import { LogoutOutlined, UserAddOutlined, ShoppingCartOutlined, SettingOutlined, GiftOutlined, HomeOutlined } from "@ant-design/icons";
import { Link } from "@reach/router";
import { Layout, Menu } from "antd";
import "./Dashboard.scss";
import { AuthContext } from "shared/contexts";

const { Content, Sider } = Layout;
const { SubMenu } = Menu;

const routes = [
  { path: "organizations", label: "Organizations", allowed: ["superadmin"] },
  { path: "branches", label: "Branches", allowed: ["superadmin", "admin"] },
  { path: "users", label: "Users", allowed: ["superadmin", "admin"] },
];

function Dashboard({ children, logout, location }) {
  const [user] = useContext(AuthContext);
  const defaultSelectedKey = location.pathname === "/" ? "home" : location.pathname.slice(1);
  return (
    <Layout
      style={{
        height: "100%",
      }}
    >
      <Sider>
        <div className="logo" />
        <Menu theme="dark" defaultSelectedKeys={[defaultSelectedKey]} defaultOpenKeys={["administration"]} mode="inline">
          <Menu.Item key="home" icon={<HomeOutlined style={{ fontSize: 20 }} />}>
            <Link to="/">Home</Link>
          </Menu.Item>
          <Menu.Item key="products" icon={<UserAddOutlined style={{ fontSize: 20 }} />}>
            <Link to="products">Products</Link>
          </Menu.Item>
          <Menu.Item key="orders" icon={<ShoppingCartOutlined style={{ fontSize: 20 }} />}>
            <Link to="orders">Orders</Link>
          </Menu.Item>
          <Menu.Item key="delivery" icon={<GiftOutlined style={{ fontSize: 20 }} />}>
            <Link to="delivery">Delivery Tech</Link>
          </Menu.Item>
          <SubMenu key="administration" icon={<SettingOutlined style={{ fontSize: 18 }} />} title="Administration">
            {routes
              .filter(({ allowed }) => allowed.includes(user.role))
              .map((i) => (
                <Menu.Item key={i.path}>
                  <Link to={i.path}>{i.label}</Link>
                </Menu.Item>
              ))}
          </SubMenu>
          <Menu.Item key="logout" icon={<LogoutOutlined style={{ fontSize: 18 }} />} onClick={logout}>
            Logout
          </Menu.Item>
        </Menu>
      </Sider>
      <Layout className="site-layout">
        <Content style={{ padding: "32px 16px 16px 16px" }} id="main">
          {children}
        </Content>
      </Layout>
    </Layout>
  );
}

export default Dashboard;
