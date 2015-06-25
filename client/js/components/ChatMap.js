var React = require('react'),
    stateTree = require('../stateTree');

var zoneCursor = stateTree.select('zone');
var usersCursor = stateTree.select('users')

var seattle = {lat: 47.6235616, lng: -122.330341}

var ChatMap = React.createClass({

  mapOptions: {
    zoom: 1,
    center: seattle,
    disableDefaultUI: true
  },

  markers: [],

  updateMap: function (e) {
    var zone = e.data.data

    var centerLat = (zone.ne.lat + zone.sw.lat) / 2;
    var centerLng = (zone.sw.lng + zone.ne.lng) / 2;

    if (centerLat == 0 && centerLng == 0) {
      centerLat = seattle.lat
      centerLng = seattle.lng
      this.map.setCenter(new google.maps.LatLng(centerLat, centerLng));
    } else {
      this.map.setCenter(new google.maps.LatLng(centerLat, centerLng));
      this.map.panToBounds(new google.maps.LatLngBounds(
        new google.maps.LatLng(zone.sw.lat, zone.sw.lng),
        new google.maps.LatLng(zone.ne.lat, zone.ne.lng)));
    }

  },

  updateMarkers: function(e) {
    this.markers.map(function(marker) {
      marker.setMap(null)
    });
    this.markers = [];

    for (var userID in e.data.data) {
      var user = e.data.data[userID];
      this.markers.push(new google.maps.Marker({
        position: new google.maps.LatLng(user.location.lat, user.location.lng),
        map: this.map
      }));
    }
  },

  componentDidMount: function () {
    this.map = new google.maps.Map(this.getDOMNode(), this.mapOptions);
    zoneCursor.on('update', this.updateMap);
    usersCursor.on('update', this.updateMarkers);
  },

  render: function () {
    return (
    	<div ref="map" className="map-canvas"></div>
    )
  }
})

module.exports = ChatMap
