var React = require('react'),
    stateTree = require('../stateTree'),
    ChatHeader = require('../components/ChatHeader'),
    ChatWindow = require('../components/ChatWindow'),
    ChatCompose = require('../components/ChatCompose'),
    ChatMap = require('../components/ChatMap'),
    ChatUsers = require('../components/ChatUsers');

var messagesCursor = stateTree.select('messages'),
    usersCursor = stateTree.select('users'),
    zoneCursor = stateTree.select('zone');

var ZonePage = React.createClass({

  mixins: [React.addons.PureRenderMixin],

  wsOpened: function () {
    console.log("opened");
  },

  wsClosed: function () {
    console.log("closed")
  },

  wsMessage: function (e) {
    var event = JSON.parse(e.data);

    switch (event.type) {
      case "message":
        messagesCursor.push(event);
        break;
      default:
    }
    //
    // if (eventData.type == "zone") {
    //   this.setProps({ zone: eventData })
    // }
    //
    // this.setState({ events: this.state.events.concat(eventData) });
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
            <div className="row gc-chat-window">
              <div className="col-md-12" id="gc-message-area">
                <ChatWindow />
              </div>
            </div>
            <ChatCompose subscription={this.props.subscription} />
          </div>
          <div className="col-md-4 gc-sidebar">
            <div className="row gc-map">
              <div className="col-md-12">
                <ChatMap />
              </div>
            </div>
            <ChatUsers />
          </div>
        </div>
      </div>
    )
  }
})

module.exports = ZonePage




  // React.render(<Users />, $('.gc-user-container')[0]);

  // var socket = new WebSocket('ws://' + window.location.host + '/s/' + initData.subscriptionId + '/socket')
  // socket.onmessage = function(event) {
  //   var eventData = JSON.parse(event.data)
  //   switch (eventData.type) {
  //     case "archive": archive(eventData.data); break;
  //     case "join": join(eventData.data); break;
  //     case "leave": leave(eventData.data); break;
  //     case "zone": zone(eventData.data); break;
  //   }
  // }

  // socket.onclose = function(event) {
  //   console.log("disconnected")
  // }
