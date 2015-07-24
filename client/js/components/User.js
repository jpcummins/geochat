var React = require('react');

var User = React.createClass({
  render: function () {
    return (
	    <div className="online" >
        {this.props.user.name}
	    </div>
    )
  }
})

module.exports = User
