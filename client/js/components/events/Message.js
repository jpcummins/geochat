var React = require('react'),
    stateTree = require('../../stateTree');

var usersCursor = stateTree.select('users');

var Message = React.createClass({
  render: function () {
    return (
      <div className="row gc-message">
        <div className="col-md-1 gc-name">
          {this.props.data.user.name}
        </div>
        <div className="col-md-10">
          {this.props.data.text}
        </div>
      </div>
    )
  }
})

module.exports = Message
