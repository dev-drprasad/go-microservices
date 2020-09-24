export default (currency, value) => {
  return new Intl
      .NumberFormat(
          'en-US', {style: 'currency', currency, maximumFractionDigits: 2})
      .format(value)
}
