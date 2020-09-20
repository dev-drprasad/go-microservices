export { default as debounce } from "./debounce";
export { default as includes } from "./includes";
export { default as listsearch } from "./listsearch";
export { default as NS } from "./NS";
export { default as safeget } from "./safeget";
export { default as sorters } from "./sorters";

export const getInitialsFromName = (str) =>
  str
    .toUpperCase()
    .split(" ")
    .slice(0, 2)
    .map((w) => w[0])
    .join("");
