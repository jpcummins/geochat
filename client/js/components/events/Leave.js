var React = require('react'),
    stateTree = require('../../stateTree');

var usersCursor = stateTree.select('users');

var Leave = React.createClass({
  render: function () {
    return (
      <div className="row gc-message">
        <div className="col-md-offset-1 col-md-10">
          {this.props.data.user.name} left
        </div>
      </div>
    )
  }
})

module.exports = Leave
