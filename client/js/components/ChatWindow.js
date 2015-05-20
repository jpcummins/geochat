var React = require('react'),
    stateTree = require('../stateTree'),
    Message = require('./events/Message'),
    Zone = require('./events/Zone'),
    Subscription = require('./events/Subscription');

var messagesCursor = stateTree.select('messages'),
    subscribersCursor = stateTree.select('subscribers'),
    zoneCursor = stateTree.select('zone');

var ChatWindow = React.createClass({

  showMessage: function (e) {
    var event = e.data.data[e.data.data.length - 1];
    event.key = event.id;
    var message = React.createElement(Message, event);
    this.setState({ events: this.state.events.concat(message) });
  },

  updateSubscription: function (e) {
    var subscriptionEvent = e.data.data[Object.keys(e.data.data)[0]]; // gross!
    subscriptionEvent.key = subscriptionEvent.id;
    var subscription = React.createElement(Subscription, subscriptionEvent);
    this.setState({ events: this.state.events.concat(subscription) });
  },

  showZone: function (e) {
    var event = e.data.data;
    event.key = event.id;
    var zone = React.createElement(Zone, event);
    this.setState({ events: this.state.events.concat(zone) });
  },

  getInitialState: function () {
    return { eventElements: [] }
  },

  componentDidUpdate: function () {
    var chatWindow = $(this.refs.chatWindow.getDOMNode());
    chatWindow.scrollTop(chatWindow[0].scrollHeight - chatWindow.height());
  },

  componentDidMount: function () {
    messagesCursor.on('update', this.showMessage);
    subscribersCursor.on('update', this.updateSubscription);
    zoneCursor.on('update', this.showZone);
  },

  getInitialState: function () {
    return { events: [] }
  },

  render: function () {
    return (
      <div className="row gc-chat-window">
        <div className="col-md-12" ref="chatWindow">
          {this.state.events}
        </div>
      </div>
    )
  }
});

module.exports = ChatWindow
