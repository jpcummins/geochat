var React = require('react'),
    stateTree = require('../stateTree'),
    Message = require('./events/Message'),
    Zone = require('./events/Zone');

var messagesCursor = stateTree.select('messages'),
    usersCursor = stateTree.select('users'),
    zoneCursor = stateTree.select('zone');

var ChatWindow = React.createClass({

  showMessage: function (e) {
    var event = e.data.data[e.data.data.length - 1];
    var message = React.createElement(Message, event, event.timestamp);
    this.setState({ events: this.state.events.concat(message) });
  },

  showUser: function (e) {
    console.log("new user");
  },

  showZone: function (e) {
    var event = e.data.data;
    var zone = React.createElement(Zone, event, event.timestamp);
    this.setState({ events: this.state.events.concat(zone) });
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
