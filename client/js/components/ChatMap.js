var React = require('react'),
    stateTree = require('../stateTree');

var zoneCursor = stateTree.select('zone');

var ChatMap = React.createClass({

  updateMap: function (e) {
    var zone = e.data.data.data; // gross

	  var mapOptions = {
	    zoom: 1,
	    center: new google.maps.LatLng(zone.boundary.SouthWestLat, zone.boundary.SouthWestLong),
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
	      new google.maps.LatLng(zone.boundary.SouthWestLat, zone.boundary.SouthWestLong),
	      new google.maps.LatLng(zone.boundary.NorthEastLat, zone.boundary.NorthEastLong))
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
