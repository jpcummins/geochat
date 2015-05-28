var React = require('react');

var ChatCompose = React.createClass({

  handleKeypress: function (e) {
  	if (e.charCode == 13 || e.keyCode == 13) {
  		e.preventDefault();
  		var inputBox = React.findDOMNode(this.refs.chatInput);
      var message = inputBox.value;

      if (/\//.test(message)) {
        message = message.replace(/\//, '').trim()
        $.ajax({
          url: '/chat/command',
          method: 'POST',
          data: {
            command: message
          }
        });
      } else {
    		$.ajax({
    			url: '/chat/message',
    			method: "POST",
    			data: {
    				text: message
    			}
    		});
      }


      inputBox.value = "";
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
