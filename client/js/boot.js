var React = require('react'),
    Zone = require('./pages/zone');

var pages = {
  'zone': Zone
}

var pageElement = React.createElement(pages[initData.page], initData);
React.render(pageElement, document.getElementById(initData.page));