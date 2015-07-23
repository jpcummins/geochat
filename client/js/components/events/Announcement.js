var React = require('react'),
    stateTree = require('../../stateTree');

var Announcement = React.createClass({
  render: function () {
    return (
      <div className="row gc-message announcement">
        <div className="col-md-offset-1 col-md-10">
          {this.props.data.text}
        </div>
      </div>
    )
  }
})

module.exports = Announcement
