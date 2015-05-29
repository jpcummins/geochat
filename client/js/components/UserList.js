var React = require('react'),
    stateTree = require('../stateTree'),
    User = require('./User');

var usersCursor = stateTree.select('users');

var UserList = React.createClass({

  getInitialState: function() {
    return { users: [] };
  },

  handleUserUpdate: function (e) {
    var users = $.map(usersCursor.get(), function(e) { return e; });
    this.setState({ users: users });
  },

  componentDidMount: function () {
    usersCursor.on('update', this.handleUserUpdate);
  },

  render: function () {
    var users = this.state.users.map(function (user) {
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
