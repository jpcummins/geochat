var React = require('react'),
    stateTree = require('../stateTree'),
    Subscriber = require('./Subscriber');

var subscribersCursor = stateTree.select('subscribers');

var SubscriberList = React.createClass({

  getInitialState: function() {
    return { subscribers: [] };
  },

  handleSubscriptionUpdate: function (e) {
    var subscribers = $.map(subscribersCursor.get(), function(e) { return e; });
    this.setState({ subscribers: subscribers });
  },

  componentDidMount: function () {
    subscribersCursor.on('update', this.handleSubscriptionUpdate);
  },

  render: function () {
    var subscribers = this.state.subscribers.map(function (subscriber) {
      return (
        <Subscriber subscriber={subscriber.data.subscriber} key={subscriber.id} />
      )
    });

    return (
	    <div className="row gc-users">
	      <div className="col-md-12 gc-user-container">
          {subscribers}
	      </div>
	    </div>
    )
  }
})

module.exports = SubscriberList
