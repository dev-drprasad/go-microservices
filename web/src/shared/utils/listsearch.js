import includes from './includes';
import safeget from './safeget';

export default function listsearch(list, fields, searchText) {
  return list.filter(
      (patient) => fields.some(
          (f) => includes(safeget(patient, f.split('.'), ''), searchText)))
}
