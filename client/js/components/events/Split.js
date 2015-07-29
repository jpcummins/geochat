var React = require('react'),
    stateTree = require('../../stateTree');

var usersCursor = stateTree.select('users');

var Split = React.createClass({
  render: function () {
    return (
      <div className="row gc-message announcement">
        <div className="col-md-offset-1 col-md-10">
          '{this.props.data.previousZone.id}' was split. New zone: '{this.props.data.zone.id}' ({Object.keys(this.props.data.zone.users).length} users)
        </div>
      </div>
    )
  }
})

module.exports = Split
