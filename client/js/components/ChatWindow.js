var React = require('react'),
    stateTree = require('../stateTree'),
    Message = require('./events/Message'),
    Zone = require('./events/Zone'),
    Join = require('./events/Join'),
    Leave = require('./events/Leave'),
    Split = require('./events/Split'),
    Merge = require('./events/Merge'),
    Announcement = require('./events/Announcement');

var eventClasses = {
  "message": Message,
  "zone": Zone,
  "join": Join,
  "leave": Leave,
  "split": Split,
  "merge": Merge,
  "announcement": Announcement,
}

var eventsCursor = stateTree.select('visibleEvents');

var ChatWindow = React.createClass({

  events: [],

  showEvent: function (e) {
    var that = this;
    var newEvents = e.data.data.filter(function(i) {return e.data.previousData.indexOf(i) < 0;});

    newEvents.forEach(function (event, index) {
      event.key = event.id
      element = React.createElement(eventClasses[event.type], event)
      that.events = that.events.concat(element);
    })

    if (this.events.length > 200) {
      this.events = this.events.splice(200)
    }

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
