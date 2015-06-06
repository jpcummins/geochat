var React = require('react'),
    stateTree = require('../../stateTree');

var Leave = React.createClass({

  getInitialState: function () {
    return stateTree.select('users').get(this.props.data.user_id)
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
