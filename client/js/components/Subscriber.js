var React = require('react');

var Subscriber = React.createClass({
  render: function () {
    return (
	    <div className={ this.props.subscriber.is_online ? 'online' : 'offline' }>
        {this.props.subscriber.user.name}
	    </div>
    )
  }
})

module.exports = Subscriber
