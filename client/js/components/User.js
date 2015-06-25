var React = require('react');

var User = React.createClass({
  render: function () {
    return (
	    <div className={ this.props.user.is_online ? 'online' : 'offline' }>
        {this.props.user.name}
	    </div>
    )
  }
})

module.exports = User
