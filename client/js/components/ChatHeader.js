var React = require('react'),
    stateTree = require('../stateTree');

var zoneCursor = stateTree.select('zone');

var ChatHeader = React.createClass({

  updateHeader: function () {
    this.forceUpdate();
  },

  componentDidMount: function () {
    zoneCursor.on('update', this.updateHeader);
  },

  render: function () {
    return (
      <div className="row gc-header">
        <div className="col-md-12">
          <h4>Location: {zoneCursor.get().id}</h4>
        </div>
      </div>
    )
  }
})

module.exports = ChatHeader
