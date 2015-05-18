var React = require('react'),
    Message = require('./events/Message');

var eventClasses = {
  "message": Message
}

var ChatWindow = React.createClass({

  getInitialState: function () {
    return { events: [] }
  },

  render: function () {
    var events = this.props.events.map(function (event) {
      if (eventClasses[event.type]) {
        return React.createElement(eventClasses[event.type], event);
      }
    });

    return (
      <div>
        {events}
      </div>
    )
  }
});

module.exports = ChatWindow