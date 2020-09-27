import {message} from 'antd';
import {useCallback, useContext, useEffect, useMemo, useState} from 'react';
// import userService from "src/shared/services/user.service";
import NS from 'shared/utils/NS';
import {AuthContext} from './contexts';

export const cr = '\n';
export const tab = '\t';

function getType(data) {
  if (data === null) return 'Null';
  if (data === undefined) return 'Undefined';
  if (typeof data === 'string') return 'String';
  if (typeof data === 'number' && !Number.isNaN(data)) return 'Number';
  if (Number.isNaN(data)) return 'NaN';
  if (typeof data === 'boolean') return 'Boolean';
  if (data instanceof Array)
    return 'Array';  // always should be before `Object`
  if (data instanceof Object) return 'Object';

  return '';
}

// https://stackoverflow.com/a/2117523
function uuidv4() {
  return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11)
      .replace(
          /[018]/g,
          (c) =>
              (c ^
               (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (c / 4))))
                  .toString(16));
}

const defaultFetchOptions = {
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json',
  },
};

export function useFetch(url, opts) {
  const [[rId, fresh], setParams] = useState(() => ['', false]);
  const [response, setResponse] = useState([undefined, new NS('INIT')]);

  const refresh = useCallback(() => setParams([uuidv4(), true]), []);

  useEffect(() => {
    if (url && opts) {
      setParams([uuidv4(), false]);
    }
  }, [url, opts]);

  useEffect(() => {
    if (!url || !rId) return;

    const abortctrl = new AbortController();

    setResponse(([, s]) => [undefined, s.clone('LOADING')]);
    const startTime = performance.now();

    // recursive merge might be better solution
    const finalopts = {
      ...defaultFetchOptions,
      ...opts,
      headers: {
        ...defaultFetchOptions.headers,
        ...opts?.headers,
        ...(fresh && {'X-Clear-Cache': true}),
        'X-Request-ID': rId,
      },
    };

    if (finalopts.headers['Content-Type'] === null) {
      delete finalopts.headers['Content-Type'];
    }
    fetch(url, finalopts)
        .then(async (res) => {
          if (abortctrl.signal.aborted) return;

          const responseTime = performance.now() - startTime;
          const cached = !!res.headers.get('X-Browser-Cache');
          let body;

          try {
            body = await res.json();
          } catch (e) {
            const message = 'Invalid JSON response from API';
            console.error(`${cr}API Error:${cr}${tab}URL: ${url}${cr}${
                tab}Msg: ${message}${cr}${tab}Code: ${res.status}`);
            setResponse(
                ([, s]) =>
                    [undefined,
                     s.clone(
                         'ERROR', '', res.status, responseTime, rId, cached)]);
            return;
          }

          if (res.status >= 400) {
            const errorType = body.error || '';
            const isInternalError =
                !errorType || errorType === 'Internal Server Error';
            const message = !isInternalError ? body.message || '' : '';
            setResponse(
                ([, s]) =>
                    [undefined,
                     s.clone(
                         'ERROR', message, res.status, responseTime, rId,
                         cached, false)]);
            return;
          }

          const dataType = getType(body);
          const hasData = dataType !== 'Null' &&
              (dataType === 'Array' ? body.length > 0 : true);

          setResponse(
              ([, s]) =>
                  [body,
                   s.clone(
                       'SUCCESS', '', res.status, responseTime, rId, cached,
                       hasData)]);
        })
        .catch((err) => {
          if (abortctrl.signal.aborted) return;
          const responseTime = performance.now() - startTime;
          console.error(`${cr}API Error:${cr}${tab}URL: ${url}${cr}${tab}Msg: ${
              err.message}${cr}${tab}Code: 0`);
          setResponse(
              ([,
                s]) => [undefined, s.clone('ERROR', '', 0, responseTime, rId)]);
        });

    return () => abortctrl.abort();
  }, [fresh, opts, rId, url]);

  return [response[0], response[1], refresh];
}

const API_BASE_URL =
    process.env.REACT_APP_API_BASE_URL || window.location.origin;
export default function useAPI(urlpath, extraOptions) {
  const [user, , logout] = useContext(AuthContext);
  const url = urlpath && new URL(urlpath, API_BASE_URL).toString();

  const options = useMemo(
      () => ({
        ...extraOptions,
        headers: {
          Authorization: `Bearer ${user?.token || ''}`,
          ...extraOptions?.headers,
        },
      }),
      [
        extraOptions, user
      ]);  // Don't add `user` . Adding will cause re-render when
  // token change (ex: /login api)

  const [data, status, refresh] = useFetch(url, options);

  useEffect(() => {
    if (status.statusCode === 401) {
      logout();
    }
    if (status.isError && !status.errorCaught && status.statusCode >= 400) {
      message.error('Oops! Something went wrong.', 3);
    }
  }, [status, logout]);

  return [data, status, refresh];
}

export const mergeStatuses = (...statuses) => {
  const hasData = statuses.some((s) => s.hasData);
  const statusWithError = statuses.find((s) => s.isError);
  const message = statusWithError ? statusWithError.message : '';
  const statusCode = statusWithError ? statusWithError.statusCode : 0;

  const status = statuses.every((s) => s.isSuccess) ?
      'SUCCESS'  // [SUCCESS, SUCCESS, SUCCESS]
      :
      statusWithError ?
      'ERROR'  // [SUCCESS, ERROR] or [LOADING, ERROR]
      :
      statuses.find((s) => s.isLoading) ?
      'LOADING'  // [SUCCESS, LOADING] or [LOADING, INIT] or [LOADING, LOADING]
      :
      statuses.find((s) => s.isSuccess) ? 'LOADING'  // [SUCCESS, INIT]
                                          :
                                          'INIT';  // [INIT, INIT]

  return new NS(status, message, statusCode, 0, '', false, hasData);
};

export function useInsurerIdAndNames() {
  const queryParams = new URLSearchParams([
    ['fields', ['id', 'name']],
    ['sortBy', 'name'],
  ]);
  const [insurers = [], status] = useAPI(`/api/v1/insurers?${queryParams}`);
  return [insurers, status];
}
export function useBrands() {
  const queryParams = new URLSearchParams([
    ['fields', ['id', 'name']],
    ['sortBy', 'name'],
  ]);
  const [patients = [], status, refresh] =
      useAPI(`/api/v1/brands?${queryParams}`);
  return [patients, status, refresh];
}
export function useCategories() {
  const queryParams = new URLSearchParams([
    ['fields', ['id', 'name']],
    ['sortBy', 'name'],
  ]);
  const [physicians = [], status] = useAPI(`/api/v1/categories?${queryParams}`);
  return [physicians, status];
}

export function useBranches() {
  const [branches = [], status] = useAPI(`/api/v1/branches`);
  return [branches, status];
}

export function useOrganizations() {
  const [organizations = [], status] = useAPI('/api/v1/organizations');
  return [organizations, status];
}

export function useOrderStatuses() {
  const [statuses = [], status] = useAPI('/api/v1/orderstatuses');
  return [statuses, status];
}

export function useEquiments() {
  const [equiments = [], status] = useAPI('/api/v1/equipments');
  return [equiments, status];
}

export function useSalesUsers() {
  const [users = [], status] = useAPI(`/api/v1/sales-users`);
  return [users, status];
}

export function useProduct(id) {
  const [payload, setPayload] = useState(undefined);
  const [product, status, refresh] =
      useAPI(id ? `/api/v1/products/${id}` : undefined);

  const saveOpts = useMemo(() => {
    if (!payload) return [undefined];
    const [urlpath, method] =
        id ? [`/api/v1/products/${id}`, 'PUT'] : ['/api/v1/products', 'POST'];
    return [urlpath, {method, body: JSON.stringify(payload)}];
  }, [payload, id]);
  const [, saveStatus] = useAPI(...saveOpts);

  return [product, status, refresh, setPayload, saveStatus];
}

export function useCustomer(id) {
  const [payload, setPayload] = useState(undefined);
  const [product, status, refresh] =
      useAPI(id ? `/api/v1/customers/${id}` : undefined);

  const saveOpts = useMemo(() => {
    if (!payload) return [undefined];
    const [urlpath, method] =
        id ? [`/api/v1/customers/${id}`, 'PUT'] : ['/api/v1/customers', 'POST'];
    return [urlpath, {method, body: JSON.stringify(payload)}];
  }, [payload, id]);
  const [, saveStatus] = useAPI(...saveOpts);

  return [product, status, refresh, setPayload, saveStatus];
}

export function useOrder(id) {
  const [payload, setPayload] = useState(undefined);
  const [order, status, refresh] =
      useAPI(id ? `/api/v1/orders/${id}` : undefined);

  const saveOpts = useMemo(() => {
    console.log('p :>> ', payload, id);
    if (!payload) return [undefined];
    const [urlpath, method] =
        id ? [`/api/v1/orders/${id}`, 'PUT'] : ['/api/v1/orders', 'POST'];
    return [urlpath, {method, body: JSON.stringify(payload)}];
  }, [payload, id]);
  const [, saveStatus] = useAPI(...saveOpts);

  return [order, status, refresh, setPayload, saveStatus];
}

export function useCustomers() {
  const [customers = [], status, refresh] = useAPI('/api/v1/customers');
  return [customers, status, refresh];
}
export function useProducts() {
  const [products = [], status, refresh] = useAPI('/api/v1/products');
  return [products, status, refresh];
}
export function useOrders() {
  const [orders = [], status, refresh] = useAPI('/api/v1/orders');
  return [orders, status, refresh];
}
