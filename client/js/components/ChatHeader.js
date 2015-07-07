var React = require('react'),
    stateTree = require('../stateTree');

var zoneCursor = stateTree.select('zone');

var ChatHeader = React.createClass({

  updateHeader: function () {

    var that = this;
    var zone = zoneCursor.get();
    var geohash = zone.id.split(':')[0]

    if (geohash.length >= 2) {
      var geocoder = new google.maps.Geocoder();
      geocoder.geocode({bounds: new google.maps.LatLngBounds(
        new google.maps.LatLng(zone.sw.lat, zone.sw.lng),
        new google.maps.LatLng(zone.ne.lat, zone.ne.lng))}, function(results, status) {

        if (results && results[0]) {
          var name = results[0].formatted_address
          $(that.getDOMNode()).find(".loc-name").text(name)
        }
      });
    }
    this.forceUpdate();
  },

  componentDidMount: function () {
    zoneCursor.on('update', this.updateHeader);
  },

  render: function () {
    return (
      <div className="row gc-header">
        <div className="col-md-12">
          <h4>
            Location: {zoneCursor.get().id}
            <span className="loc-name"></span>
          </h4>
        </div>
      </div>
    )
  }
})

module.exports = ChatHeader
