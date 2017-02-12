// Package timeseries implements encoder and decoder for time-series data point in similar to
// Facebook Gorilla time-series database.
// The spec is described in the paper "Gorilla: A Fast, Scalable, In-Memory Time Series Database"
// http://www.vldb.org/pvldb/vol8/p1816-teller.pdf
//
// This implementation is based on a third party Java implementation https://github.com/burmanm/gorilla-tsc/
// However this Go package removes the enhancements in this Java version.
//   - The precision of timestamps is one second (not a milisecond).
//   - The data point value type is flaot64 only.
//   - The first timestamp delta is sized at 14 bits. This size span a bit more than 4 hours (16,384 seconds).
package timeseries
