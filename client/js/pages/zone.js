var React = require('react'),
    ChatHeader = require('../components/ChatHeader'),
    ChatWindow = require('../components/ChatWindow'),
    ChatCompose = require('../components/ChatCompose'),
    ChatMap = require('../components/ChatMap'),
    ChatUsers = require('../components/ChatUsers');

var ZonePage = React.createClass({

  getInitialState: function () {
    return { events: [], zone: {} }
  },

  wsOpened: function () {
    console.log("opened");
  },

  wsClosed: function () {
    console.log("closed")
  },

  wsMessage: function (e) {
    var eventData = JSON.parse(e.data);

    if (eventData.type == "zone") {
      this.state.zone = eventData;
    }

    this.setState({ events: this.state.events.concat(eventData) });
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
                <ChatWindow events={this.state.events} />
              </div>
            </div>
            <ChatCompose subscription={this.props.subscription} />
          </div>
          <div className="col-md-4 gc-sidebar">
            <div className="row gc-map">
              <div className="col-md-12">
                <ChatMap zone={this.state.zone} />
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