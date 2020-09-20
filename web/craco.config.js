const CracoLessPlugin = require('craco-less');

const bodyBackground = '#f0f2f5';

module.exports = {
  plugins: [
    {
      plugin: CracoLessPlugin,
      options: {
        lessLoaderOptions: {
          lessOptions: {
            modifyVars: {
              '@primary-color': '#00cc83',
              '@text-color': 'rgb(23, 43, 77)',
              '@heading-color': 'rgba(23, 43, 77, 0.85)',
              '@body-background': bodyBackground,
              '@background-color-light': bodyBackground,
              '@card-head-background': bodyBackground,
              '@table-row-hover-bg': '@primary-2'
            },
            javascriptEnabled: true,
          },
        },
      },
    },
  ],
};
