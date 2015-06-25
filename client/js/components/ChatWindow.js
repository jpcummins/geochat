var React = require('react'),
    stateTree = require('../stateTree'),
    Message = require('./events/Message'),
    Zone = require('./events/Zone'),
    Join = require('./events/Join'),
    Leave = require('./events/Leave');

var eventClasses = {
  "message": Message,
  "zone": Zone,
  "join": Join,
  "leave": Leave,
}

var eventsCursor = stateTree.select('visibleEvents');

var ChatWindow = React.createClass({

  events: [],

  showEvent: function (e) {
    var that = this;
    var newEvents = e.data.data.filter(function(i) {return e.data.previousData.indexOf(i) < 0;});

    newEvents.forEach(function (event) {
      event.key = event.id
      element = React.createElement(eventClasses[event.type], event)
      that.events = that.events.concat(element);
    })

    this.forceUpdate();
  },

  componentDidMount: function () {
    eventsCursor.on('update', this.showEvent);
  },

  componentDidUpdate: function () {
    var chatWindow = $(this.refs.chatWindow.getDOMNode());
    chatWindow.scrollTop(chatWindow[0].scrollHeight - chatWindow.height());
  },

  render: function () {
    return (
      <div className="row gc-chat-window">
        <div className="col-md-12" ref="chatWindow">
          {this.events}
        </div>
      </div>
    )
  }
});

module.exports = ChatWindow
