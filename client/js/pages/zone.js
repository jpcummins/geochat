var React = require('react'),
    stateTree = require('../stateTree'),
    ChatHeader = require('../components/ChatHeader'),
    ChatWindow = require('../components/ChatWindow'),
    ChatCompose = require('../components/ChatCompose'),
    ChatMap = require('../components/ChatMap'),
    UserList = require('../components/UserList');

var eventsCursor = stateTree.select('visibleEvents'),
    zoneCursor = stateTree.select('zone'),
    usersCursor = stateTree.select('users');

var ZonePage = React.createClass({

  mixins: [React.addons.PureRenderMixin],

  getInitialState: function () {
    return { zone: {}, users: {} }
  },

  handleChatEvent: function (chatEvent) {
    switch (chatEvent.type) {
      case "message":
        eventsCursor.push(chatEvent);
        break;
      case "join":
        eventsCursor.push(chatEvent)
        usersCursor.set(chatEvent.data.user.id, chatEvent.data.user)

        if (chatEvent.data.zone) {
          this.setState({
            zone: chatEvent.data.zone,
            users: chatEvent.data.zone.users
          });
          usersCursor.set(chatEvent.data.zone.users);
        }
        break;
      default:
    }
    stateTree.commit();
  },

  wsOpened: function (e) {
    console.log("opened", e);
  },

  wsClosed: function (e) {
    console.log("closed", e);
  },

  wsMessage: function (e) {
    var wsEvent = JSON.parse(e.data);
    return this.handleChatEvent(wsEvent);
  },

  wsError: function (err) {
    console.log("Websocket error", err);
  },

  componentDidMount: function () {
    this.socket = new WebSocket('ws://' + window.location.host + '/chat/socket')
    this.socket.onopen = this.wsOpened;
    this.socket.onclose = this.wsClosed;
    this.socket.onmessage = this.wsMessage;
    this.socket.onerror = this.wsError;
  },

  componentWillUnmount: function () {
    this.socket.close();
  },

  render: function () {
    return (
      <div className="container-fluid">
        <ChatHeader zone={this.state.zone} />
        <div className="row gc-content">
          <div className="col-md-8">
            <ChatWindow />
            <ChatCompose />
          </div>
          <div className="col-md-4 gc-sidebar">
            <div className="row gc-map">
              <div className="col-md-12">
                <ChatMap zone={this.state.zone} />
              </div>
            </div>
            <UserList users={this.state.users} />
          </div>
        </div>
      </div>
    )
  }
})

module.exports = ZonePage
