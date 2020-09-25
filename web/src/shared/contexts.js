import { createContext } from "react";
import { NS } from "shared/utils";

export const AuthContext = createContext([{ isAuthenticated: true }, new NS("INIT"), () => {}, () => {}]);

export const CurrencyContext = createContext();

export const LocaleContext = createContext();
