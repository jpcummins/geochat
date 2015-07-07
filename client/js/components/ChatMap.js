var React = require('react'),
    stateTree = require('../stateTree');

var zoneCursor = stateTree.select('zone');
var usersCursor = stateTree.select('users');
var userCursor = stateTree.select('user');

var ChatMap = React.createClass({

  mapOptions: {
    zoom: 1,
    disableDefaultUI: true
  },

  markers: [],

  updateMap: function (e) {
    var zone = e.data.data
    var user = userCursor.get()
    var center = new google.maps.LatLng(user.location.lat, user.location.lng)
    this.map.setCenter(center);

    if (Math.abs(zone.sw.lat) < 90 &&
        Math.abs(zone.ne.lat) < 90 &&
        Math.abs(zone.sw.lng) < 180 &&
        Math.abs(zone.ne.lng) < 180) {
      this.map.fitBounds(new google.maps.LatLngBounds(
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
