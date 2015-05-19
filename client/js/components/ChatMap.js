var React = require('react');

var ChatMap = React.createClass({

  componentDidUpdate: function (nextProps) {

    console.log("componentDidUpdate")

    // var zone = this.props.zone.data;
    //
	  // var mapOptions = {
	  //   zoom: 1,
	  //   center: new google.maps.LatLng(zone.boundary.SouthWestLat, zone.boundary.SouthWestLong),
	  //   disableDefaultUI: true
	  // };
    //
	  // map = new google.maps.Map(document.getElementById('map-canvas'), mapOptions);
    //
	  // var rectangle = new google.maps.Rectangle({
	  //   strokeColor: '#FF0000',
	  //   strokeOpacity: 0.8,
	  //   strokeWeight: 2,
	  //   fillColor: '#FF0000',
	  //   fillOpacity: 0.35,
	  //   map: map,
	  //   bounds: new google.maps.LatLngBounds(
	  //     new google.maps.LatLng(zone.boundary.SouthWestLat, zone.boundary.SouthWestLong),
	  //     new google.maps.LatLng(zone.boundary.NorthEastLat, zone.boundary.NorthEastLong))
	  // });
  },

  render: function () {
    return (
    	<div id="map-canvas"></div>
    )
  }
})

module.exports = ChatMap
