var React = require('react');

var ChatCompose = React.createClass({

  handleKeypress: function (e) {
  	if (e.charCode == 13 || e.keyCode == 13) {
  		e.preventDefault();
  		var inputBox = React.findDOMNode(this.refs.chatInput);
      var message = inputBox.value;

      if (/\//.test(message)) {
        message = message.replace(/\//, '').split(/\s(.+)?/)
        var command = message[0].trim()
        var args = message[1] ? message[1].trim() : ""

        $.ajax({
          url: '/chat/command',
          method: 'POST',
          data: {
            command: command,
            args: args
          }
        });
      } else {
    		$.ajax({
    			url: '/chat/message',
    			method: "POST",
          contentType: 'application/javascript',
    			data: JSON.stringify({
    				text: message
    			})
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
