export default function includes(input, searchText) {
  if (!input) return false;

  const trimmed = searchText.trim()
  if (!trimmed) return true;

  return input.toString().toLowerCase().includes(trimmed.toString());
}
