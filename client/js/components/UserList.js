var React = require('react'),
    stateTree = require('../stateTree'),
    User = require('./User');


var UserList = React.createClass({
  getDefaultProps: function() {
    return {
      users: []
    }
  },

  render: function () {
    var users = this.props.users.map(function (user) {
      return (
        <User user={user} key={user.id} />
      )
    });

    return (
	    <div className="row gc-users">
	      <div className="col-md-12 gc-user-container">
          {users}
	      </div>
	    </div>
    )
  }
})

module.exports = UserList
