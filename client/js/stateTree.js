var ReactAddons = require('react/addons'),
    Baobab = require('baobab');

var stateTree = new Baobab({
  visibleEvents: [],
  users: {},
  zone: {}
}, {
  mixins: [ReactAddons.PureRenderMixin],
  shiftReferences: true
});

module.exports = stateTree;
