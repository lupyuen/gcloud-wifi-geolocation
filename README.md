# gcloud-wifi-geolocation
Go web application for Blue Pill Geolocation based on Google Cloud App Engine Standard Environment

Check out the article:

https://medium.com/@ly.lee/connect-stm32-blue-pill-to-esp8266-with-apache-mynewt-7edceb9e3b8d

`gcloud-wifi-geolocation` is a Go web application hosted on Google Cloud Standard App Engine that
renders realtime temperature sensor data and geolocation on a map.  The map is rendered using Mapbox GL JS.

The sensor data and geolocation are pushed via HTTPS from thethings.io Cloud Code Trigger `forward_geolocate` and
Cloud Code Function `geolocate`:

https://github.com/lupyuen/thethingsio-wifi-geolocation

thethings.io receives WiFi Access Point MAC Addresses and Signal Strength scanned by STM32 Blue Pill, running Apache Mynewt connected to ESP8266:

https://github.com/lupyuen/stm32bluepill-mynewt-sensor/tree/esp8266

Note that this implementation only supports 1 server instance because the sensor and geolocation data are cached in
the RAM of the same server that's serving the web app. Since this app is written in Go, the single server instance 
will support a high volume of concurrent requests. With App Engine's Basic Scaling, the server instance
will be shut down automatically when there are no requests. The RAM will also be cleared.

For privacy, users are required to specify the Device ID when viewing the app. Adapted from:

https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine/go11x/helloworld/helloworld.go

https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine_flexible/pubsub/pubsub.go

https://docs.mapbox.com/mapbox-gl-js/example/3d-buildings/
