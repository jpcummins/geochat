var React = require('react')

var Join = React.createClass({

  getInitialState: function () {
    return this.props.data
  },

  render: function () {
    return (
      <div className="row gc-message">
        <div className="col-md-offset-1 col-md-10">
          {this.state.name} joined
        </div>
      </div>
    )
  }
})

module.exports = Join
