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
    return { zone: {}, users: {}, events: [] }
  },

  handleChatEvent: function (chatEvent) {
    switch (chatEvent.type) {
      case "message":
        chatEvent.data.user = usersCursor.get(chatEvent.data.user_id)
        eventsCursor.push(chatEvent)
        break;
      case "zone":
        eventsCursor.push(chatEvent)
        zoneCursor.set(chatEvent.data.zone)
        usersCursor.set(chatEvent.data.zone.users)
        break;
      case "join":
        usersCursor.set(chatEvent.data.user.id, chatEvent.data.user)
        eventsCursor.push(chatEvent)
        break;
      case "leave":
        chatEvent.data.user = usersCursor.get(chatEvent.data.user_id)
        eventsCursor.push(chatEvent)
        usersCursor.unset(chatEvent.data.user.id)
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
        <ChatHeader />
        <div className="row gc-content">
          <div className="col-md-8">
            <ChatWindow />
            <ChatCompose />
          </div>
          <div className="col-md-4 gc-sidebar">
            <div className="row gc-map">
              <div className="col-md-12">
                <ChatMap />
              </div>
            </div>
            <UserList />
          </div>
        </div>
      </div>
    )
  }
})

module.exports = ZonePage
