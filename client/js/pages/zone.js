var React = require('react'),
    stateTree = require('../stateTree'),
    ChatHeader = require('../components/ChatHeader'),
    ChatWindow = require('../components/ChatWindow'),
    ChatCompose = require('../components/ChatCompose'),
    ChatMap = require('../components/ChatMap'),
    SubscriberList = require('../components/SubscriberList');

var messagesCursor = stateTree.select('messages'),
    subscribersCursor = stateTree.select('subscribers'),
    zoneCursor = stateTree.select('zone');

var ZonePage = React.createClass({

  mixins: [React.addons.PureRenderMixin],

  handleChatEvent: function (chatEvent) {
    switch (chatEvent.type) {
      case "message":
        messagesCursor.push(chatEvent);
        break;
      case "zone":
        stateTree.set('zone', chatEvent);
        for (var i = chatEvent.data.archive.events.length - 1; i >= 0; i--) {
          this.handleChatEvent(chatEvent.data.archive.events[i]);
        }
        break;
      case "join":
      case "online":
      case "offline":
        subscribersCursor.set(chatEvent.data.subscriber.id, chatEvent);
        break;
      case "leave":
        subscribersCursor.unset(chatEvent.data.subscriber.id);
        break;
      default:
    }
    stateTree.commit();
  },

  wsOpened: function () {
    console.log("opened");
  },

  wsClosed: function () {
    console.log("closed")
  },

  wsMessage: function (e) {
    var wsEvent = JSON.parse(e.data);
    return this.handleChatEvent(wsEvent);
  },

  wsError: function (err) {
    console.log("Websocket error", err);
  },

  componentDidMount: function () {
    this.socket = new WebSocket('ws://' + window.location.host + '/s/' + this.props.subscription + '/socket')
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
            <ChatCompose subscription={this.props.subscription} />
          </div>
          <div className="col-md-4 gc-sidebar">
            <div className="row gc-map">
              <div className="col-md-12">
                <ChatMap />
              </div>
            </div>
            <SubscriberList />
          </div>
        </div>
      </div>
    )
  }
})

module.exports = ZonePage
