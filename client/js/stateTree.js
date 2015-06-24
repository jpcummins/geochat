var ReactAddons = require('react/addons'),
    Baobab = require('baobab');

var stateTree = new Baobab({
  visibleEvents: [],
  zone: {},
  users: {},
}, {
  mixins: [ReactAddons.PureRenderMixin],
  shiftReferences: true
});

module.exports = stateTree;
