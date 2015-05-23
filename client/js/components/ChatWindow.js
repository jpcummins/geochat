var React = require('react'),
    stateTree = require('../stateTree'),
    Message = require('./events/Message'),
    Zone = require('./events/Zone'),
    Subscription = require('./events/Subscription');

var visibleEvents = stateTree.select('visibleEvents');

var eventClasses = {
  "message": Message,
  "zone": Zone,
  "online": Subscription,
  "offline": Subscription,
  "join": Subscription,
  "leave": Subscription
}

var ChatWindow = React.createClass({

  showEvent: function (e) {
    var event = e.data.data[e.data.data.length - 1];
    event.key = event.id;

    // todo: optimize - use a hashtable instead.
    for (var i = 0; i < this.state.events.length; i++) {
      var key = this.state.events[i].key;
      if (key == event.key) {
        return;
      }
    }
    
    var element = React.createElement(eventClasses[event.type], event);
    this.setState({ events: this.state.events.concat(element) });
  },

  getInitialState: function () {
    return { eventElements: [] }
  },

  componentDidUpdate: function () {
    var chatWindow = $(this.refs.chatWindow.getDOMNode());
    chatWindow.scrollTop(chatWindow[0].scrollHeight - chatWindow.height());
  },

  componentDidMount: function () {
    visibleEvents.on('update', this.showEvent);
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
