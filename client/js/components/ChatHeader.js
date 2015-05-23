var React = require('react'),
    stateTree = require('../stateTree');

var zoneCursor = stateTree.select('zone');

var ChatHeader = React.createClass({

  getInitialState: function () {
    return { location: '' };
  },

  componentDidMount: function () {
    zoneCursor.on('update', this.showZone);
  },

  showZone: function (e) {
    this.setState({ location: e.data.data.data.id });
  },

  render: function () {
    return (
      <div className="row gc-header">
        <div className="col-md-12">
          <h4>Location: {this.state.location}</h4>
        </div>
      </div>
    )
  }
})

module.exports = ChatHeader
