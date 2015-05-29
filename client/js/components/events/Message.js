var React = require('react'),
    stateTree = require('../../stateTree');

var subscribersCursor = stateTree.select('subscribers');

var Message = React.createClass({
  render: function () {
    return (
      <div className="row gc-message">
        <div className="col-md-1 gc-name">
          {subscribersCursor.get(this.props.data.subscription).name}
        </div>
        <div className="col-md-10">
          {this.props.data.text}
        </div>
      </div>
    )
  }
})

module.exports = Message
