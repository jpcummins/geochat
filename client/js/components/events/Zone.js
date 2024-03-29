var React = require('react'),
    stateTree = require('../../stateTree');

var usersCursor = stateTree.select('users');

var Zone = React.createClass({
  render: function () {
    return (
      <div className="row gc-message">
        <div className="col-md-offset-1 col-md-10">
          Joined zone: "{this.props.data.zone.id}" with {Object.keys(this.props.data.zone.users).length - 1} other users.
        </div>
      </div>
    )
  }
})

module.exports = Zone
