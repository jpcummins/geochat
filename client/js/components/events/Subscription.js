var React = require('react');

var Subscription = React.createClass({
  render: function () {
    return (
      <div className="row gc-message">
        <div className="col-md-1 gc-name">
          {this.props.data.subscriber.user.name}
        </div>
        <div className="col-md-10">
          {this.props.type}
        </div>
      </div>
    )
  }
})

module.exports = Subscription