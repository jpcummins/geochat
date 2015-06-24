var React = require('react'),
    stateTree = require('../stateTree');

var ChatMap = React.createClass({

  updateMap: function (e) {
    var zone = this.props.zone;

	  var mapOptions = {
	    zoom: 1,
	    center: new google.maps.LatLng(zone.sw.lat, zone.sw.lng),
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
	      new google.maps.LatLng(zone.sw.lat, zone.sw.lng),
	      new google.maps.LatLng(zone.ne.lat, zone.ne.lng))
	  });
  },

  componentDidUpdate: function () {
    this.updateMap();
  },

  render: function () {
    return (
    	<div ref="map" className="map-canvas"></div>
    )
  }
})

module.exports = ChatMap
