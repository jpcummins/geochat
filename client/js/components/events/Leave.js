var React = require('react'),
    stateTree = require('../../stateTree');

var usersCursor = stateTree.select('users');

var Leave = React.createClass({
  render: function () {
    return (
      <div className="row gc-message">
        <div className="col-md-1 gc-name">

        </div>
          <div className="col-md-10">
          {this.props.data.user.firstName} left
        </div>
      </div>
    )
  }
})

module.exports = Leave
