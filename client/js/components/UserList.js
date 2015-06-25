var React = require('react'),
    stateTree = require('../stateTree'),
    User = require('./User');

var usersCursor = stateTree.select('users')

var UserList = React.createClass({

  users: [],

  updateUsers: function(e) {
    this.users = [];
    for (var id in e.data.data) {
      this.users.push(e.data.data[id]);
    }
    this.forceUpdate()
  },

  componentDidMount: function () {
    usersCursor.on('update', this.updateUsers);
  },

  render: function () {

    var users = this.users.map(function (user) {
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
