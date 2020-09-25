import { Router } from "@reach/router";
import { ConfigProvider } from "antd";
import enUS from "antd/es/locale/en_US";
import hiIN from "antd/es/locale/hi_IN";
import ptBR from "antd/es/locale/pt_BR";
import Dashboard from "Dashboard";
import BranchList from "pages/BranchList";
import CustomerList from "pages/CustomerList";
import Home from "pages/Home";
import Login from "pages/Login";
import NotFound from "pages/NotFound";
import OrganizationList from "pages/OrganizationList";
import ProductAdd from "pages/ProductAdd";
import ProductEdit from "pages/ProductEdit";
import ProductList from "pages/ProductList";
import UserList from "pages/UserList";
import React, { useCallback, useState } from "react";
import { ProtectedRoute } from "shared/components";
import { AuthContext, CurrencyContext, LocaleContext } from "shared/contexts";
import "./App.less";

const localeFileMap = { "en-US": enUS, "hi-IN": hiIN, "pt-BR": ptBR };
const localeOpts = [
  { label: "🇺🇸 en-US", value: "en-US" },
  { label: "🇮🇳 hi-IN", value: "hi-IN" },
  { label: "🇧🇷 pt-BR", value: "pt-BR" },
];

const currencyOpts = [
  { label: "₹ INR", value: "INR" },
  { label: "$ USD", value: "USD" },
  { label: "R$ BRL", value: "BRL" },
];
// const currencies = Object.keys(currencyOpts);

const LS_USER_KEY = "user";

function getUserFromStorage() {
  let user;
  try {
    user = JSON.parse(localStorage.getItem(LS_USER_KEY)) || undefined;
  } catch (err) {
    console.err(err);
  }
  return user;
}

function App() {
  const [user, setUser] = useState(getUserFromStorage);
  const [locale, setLocale] = useState([localeOpts[0].value]);
  const [currency, setCurrency] = useState([currencyOpts[0].value]);

  const login = useCallback((user) => {
    localStorage.setItem(LS_USER_KEY, JSON.stringify(user));
    setUser(user);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem(LS_USER_KEY);
    setUser(undefined);
  }, []);

  return (
    <ConfigProvider locale={localeFileMap[locale]}>
      <LocaleContext.Provider value={[locale, localeOpts, setLocale]}>
        <AuthContext.Provider value={[user, login, logout]}>
          <CurrencyContext.Provider value={[currency, currencyOpts, setCurrency]}>
            <Router id="router">
              <Login path="login" />
              <ProtectedRoute user={user} component={Dashboard} path="/">
                <Home path="/" />
                <CustomerList path="customers" />
                <UserList path="users" />
                <OrganizationList path="organizations" />
                <BranchList path="branches" />
                <ProductList path="products" />
                <ProductAdd path="products/new" />
                <ProductEdit path="products/:id/edit" />
                <NotFound default />
              </ProtectedRoute>
              <NotFound default />
            </Router>
          </CurrencyContext.Provider>
        </AuthContext.Provider>
      </LocaleContext.Provider>
    </ConfigProvider>
  );
}

export default App;
