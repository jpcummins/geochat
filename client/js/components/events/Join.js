var React = require('react'),
    stateTree = require('../../stateTree');

var usersCursor = stateTree.select('users');

var Join = React.createClass({
  render: function () {
    return (
      <div className="row gc-message">
        <div className="col-md-offset-1 col-md-10">
          {this.props.data.user.name} joined
        </div>
      </div>
    )
  }
})

module.exports = Join
