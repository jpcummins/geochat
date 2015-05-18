var React = require('react');

var ChatMap = React.createClass({

  componentDidMount: function () {
	  var mapOptions = {
	    zoom: 1,
	    center: new google.maps.LatLng({{ .boundary.SouthWestLat }}, {{ .boundary.SouthWestLong }}),
	    disableDefaultUI: true
	  };

	  map = new google.maps.Map(document.getElementById('map-canvas'), mapOptions);

	  var rectangle = new google.maps.Rectangle({
	    strokeColor: '#FF0000',
	    strokeOpacity: 0.8,
	    strokeWeight: 2,
	    fillColor: '#FF0000',
	    fillOpacity: 0.35,
	    map: map,
	    bounds: new google.maps.LatLngBounds(
	      new google.maps.LatLng({{ .boundary.SouthWestLat }}, {{ .boundary.SouthWestLong }}),
	      new google.maps.LatLng({{ .boundary.NorthEastLat }}, {{ .boundary.NorthEastLong }}))
	  });
  },

  componentWillRecieveProps: function (nextProps) {

  },

  render: function () {
    return (
    	<div id="map-canvas"></div>
    )
  }
})

module.exports = ChatMap