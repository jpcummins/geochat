var React = require('react'),
    stateTree = require('../stateTree');

var zoneCursor = stateTree.select('zone');

var ChatMap = React.createClass({

  updateMap: function (e) {
    var zone = e.data.data.data; // gross

	  var mapOptions = {
	    zoom: 1,
	    center: new google.maps.LatLng(zone.boundary.swlat, zone.boundary.swlong),
	    disableDefaultUI: true
	  };

	  var map = new google.maps.Map(this.getDOMNode(), mapOptions);

	  var rectangle = new google.maps.Rectangle({
	    strokeColor: '#FF0000',
	    strokeOpacity: 0.8,
	    strokeWeight: 2,
	    fillColor: '#FF0000',
	    fillOpacity: 0.35,
	    map: map,
	    bounds: new google.maps.LatLngBounds(
	      new google.maps.LatLng(zone.boundary.swlat, zone.boundary.swlong),
	      new google.maps.LatLng(zone.boundary.nelat, zone.boundary.nelong))
	  });

    this.setState({ map: map });
  },

  componentDidMount: function () {
    zoneCursor.on('update', this.updateMap);
  },

  render: function () {
    return (
    	<div ref="map" className="map-canvas"></div>
    )
  }
})

module.exports = ChatMap
