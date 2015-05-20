var React = require('react');

var Subscriber = React.createClass({
  render: function () {
    return (
	    <div>
        {this.props.subscriber.user.name} { this.props.subscriber.is_online ? '(online)' : '(offline)' }
	    </div>
    )
  }
})

module.exports = Subscriber
