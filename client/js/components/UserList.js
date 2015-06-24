var React = require('react'),
    stateTree = require('../stateTree'),
    User = require('./User');


var UserList = React.createClass({

  render: function () {

    var users = []
    for (var id in this.props.users) {
      var user = this.props.users[id];
      users = users.concat(<User user={user} key={user.id} />);
    }

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
