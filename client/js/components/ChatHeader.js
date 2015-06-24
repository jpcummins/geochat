var React = require('react'),
    stateTree = require('../stateTree');

var zoneCursor = stateTree.select('zone');

var ChatHeader = React.createClass({

  componentDidMount: function () {
    zoneCursor.on('update', this.showZone);
  },

  render: function () {
    return (
      <div className="row gc-header">
        <div className="col-md-12">
          <h4>Location: {this.props.zone.id}</h4>
        </div>
      </div>
    )
  }
})

module.exports = ChatHeader
