var React = require('react'),
    UserList = require('../components/UserList'),
    archive = require('../events/archive'),
    join = require('../events/join'),
    leave = require('../events/leave'),
    zone = require ('../events/zone');

var Users = React.createClass({
  getInitialState: function() {
    return {users: ['foo', 'bar', 'baz']};
  },
  render: function() {
    return (
      <UserList users={this.state.users} />
    );
  }
});

var init = function (options) {

  React.render(<Users />, $('.gc-user-container')[0]);

  var socket = new WebSocket('ws://' + window.location.host + '/s/' + initData.subscriptionId + '/socket')
  socket.onmessage = function(event) {
    var eventData = JSON.parse(event.data)
    switch (eventData.type) {
      case "archive": archive(eventData.data); break;
      case "join": join(eventData.data); break;
      case "leave": leave(eventData.data); break;
      case "zone": zone(eventData.data); break;
    }
  }

  socket.onclose = function(event) {
    console.log("disconnected")
  }

}

module.exports = init