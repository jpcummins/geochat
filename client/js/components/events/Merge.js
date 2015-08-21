var React = require('react'),
    stateTree = require('../../stateTree');

var usersCursor = stateTree.select('users');

var Merge = React.createClass({
  render: function () {
    return (
      <div className="row gc-message announcement">
        <div className="col-md-offset-1 col-md-10">
          {this.props.data.leftID} and {this.props.data.rightID} merged to {this.props.data.zone.id} ({Object.keys(this.props.data.zone.users).length} users)
        </div>
      </div>
    )
  }
})

module.exports = Merge
