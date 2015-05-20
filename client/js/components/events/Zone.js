var React = require('react');

var Zone = React.createClass({
  render: function () {
    return (
      <div className="row gc-message">
        <div className="col-md-offset-1 col-md-10">
          Joined zone: {this.props.data.zonehash}
        </div>
      </div>
    )
  }
})

module.exports = Zone
