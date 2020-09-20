import React, { useCallback, useState } from "react";

import { Router } from "@reach/router";
import Dashboard from "Dashboard";
import Login from "pages/Login";
import UserList from "pages/UserList";
import NotFound from "pages/NotFound";
import { AuthContext } from "shared/contexts";
import { ProtectedRoute } from "shared/components";
import "./App.less";

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

  const login = useCallback((user) => {
    localStorage.setItem(LS_USER_KEY, JSON.stringify(user));
    setUser(user);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem(LS_USER_KEY);
    setUser(undefined);
  }, []);

  return (
    <AuthContext.Provider value={[user, login]}>
      <Router id="router">
        <Login path="login" />
        <ProtectedRoute user={user} component={Dashboard} logout={logout} path="/">
          <UserList path="users" />
          <NotFound default />
        </ProtectedRoute>
        <NotFound default />
      </Router>
    </AuthContext.Provider>
  );
}

export default App;
