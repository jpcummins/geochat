var React = require('react')

var Leave = React.createClass({

  getInitialState: function () {
    return this.props.data.user
  },

  render: function () {
    return (
      <div className="row gc-message">
        <div className="col-md-offset-1 col-md-10">
          {this.state.name} left
        </div>
      </div>
    )
  }
})

module.exports = Leave
