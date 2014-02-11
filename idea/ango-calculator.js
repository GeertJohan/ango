angular.module('ango-calculator', {
  setup: function() {
    // setup for ango-calculator 
  },
  teardown: function() {
    //teardown for ango-calculator
  }
}).factory('GongCalculator', [function () {
	var service = {};

	var server = {};
	server.add = function(a, b) {
		return c;
	}
	//++ more server services


	service.startCommunication = function(handler) {
		if(typeof(handler) != 'object') {
			throw "error, handler must be object";
		}
		if(typeof(handler.connected) != "function") {
			throw "error, expecting handler to have function connected(..)";
		}
		if(typeof(handler.disconnected) != "function") {
			throw "error, expecting handler to have function disconnected(..)";
		}
		if(typeof(handler.displayNotification) != "function") {
			throw "error, expecting handler to have function displayNotification(..)";
		}

		console.log('successfull argument `handler` for call to startCommunication(handler), continueing with setup comm');
		handler.connected(server);
		handler.disconnected(error);

		handler.displayNotification(..);
	}

	return service;
}])
