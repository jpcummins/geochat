var ReactAddons = require('react/addons'),
    Baobab = require('baobab');

var stateTree = new Baobab({
  visibleEvents: [],
  subscribers: {},
  zone: {}
}, {
  mixins: [ReactAddons.PureRenderMixin],
  shiftReferences: true
});

module.exports = stateTree;
