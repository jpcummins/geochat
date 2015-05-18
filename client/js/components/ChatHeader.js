var React = require('react');

var ChatHeader = React.createClass({
  render: function () {
    return (
      <div className="row gc-header">
        <div className="col-md-12">
          <h4>Location: ?</h4>
          <span className="pull-right">
            <input type="checkbox" className="gc-notify">Notifications</input>
          </span>
        </div>
      </div>
    )
  }
})

module.exports = ChatHeader