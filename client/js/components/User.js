var React = require('react');

var User = React.createClass({
  render: function () {
    return (
	    <div className="user" >
        <img src={this.props.user.fbPictureUrl} />
        {this.props.user.name}
	    </div>
    )
  }
})

module.exports = User
