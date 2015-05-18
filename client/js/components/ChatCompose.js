var React = require('react');

var ChatCompose = React.createClass({

  handleKeypress: function (e) {
  	if (e.charCode == 13 || e.keyCode == 13) {
  		e.preventDefault();
  		var message = React.findDOMNode(this.refs.chatInput);

  		$.ajax({
  			url: '/s/' + this.props.subscription + '/message',
  			method: "POST",
  			data: {
  				text: message.value
  			}
  		});

  		message.value = "";
  	}
  },

  render: function () {
    return (
	    <div className="row gc-chat-textbox">
	      <div className="col-md-12">
	        <input type="text" className="form-control" placeholder="Say something…" ref="chatInput" autoComplete="off" autofocus onKeyPress={this.handleKeypress} />
	      </div>
	    </div>
    )
  }
})

module.exports = ChatCompose