var React = require('react'),
    stateTree = require('../stateTree'),
    Message = require('./events/Message');

var messagesCursor = stateTree.select('messages'),
    usersCursor = stateTree.select('users'),
    zoneCursor = stateTree.select('zone');

var ChatWindow = React.createClass({

  showMessage: function (e) {
    this.setState({ events: e.data.data.map(function (message) {
      return React.createElement(Message, message)
    })})
  },

  showUser: function (e) {
    console.log("new user");
  },

  showZone: function (e) {
    console.log("new zone");
  },

  getInitialState: function () {
    return { eventElements: [] }
  },

  componentDidMount: function () {
    messagesCursor.on('update', this.showMessage);
    usersCursor.on('update', this.showUser);
    zoneCursor.on('update', this.showZone);
  },

  getInitialState: function () {
    return { events: [] }
  },

  render: function () {
    return (
      <div>
        {this.state.events}
      </div>
    )
  }
});

module.exports = ChatWindow
