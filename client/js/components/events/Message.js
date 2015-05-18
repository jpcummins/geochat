var React = require('react');

var Message = React.createClass({
  render: function () {
    return (
      <div key={this.props.timestamp} className="row gc-message">
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