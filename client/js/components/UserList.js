var React = require('react')

var UserList = React.createClass({
  render: function() {
    var createUser = function(user, index) {
      return <li key={index}>{user}</li>;
    };
    return <ul>{this.props.users.map(createUser)}</ul>
  }
});


module.exports = UserList