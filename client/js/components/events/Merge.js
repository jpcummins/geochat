var React = require('react'),
    stateTree = require('../../stateTree');

var usersCursor = stateTree.select('users');

var Merge = React.createClass({
  render: function () {
    return (
      <div className="row gc-message">
        <div className="col-md-offset-1 col-md-10">
          Zone merged
        </div>
      </div>
    )
  }
})

module.exports = Merge
