angular.module('{{.ServiceName}}', [])
.provider('{{.ServiceName}}', function($q) {
	// SERVICE
	var service = function(wsPath) {

		// PROCEDURES, as defined in .ango file
		this.Add = function(a, b) {
			//++ return deferred
		}

		//++ do event handlers for incomming calls?
		//++ or, require provider or service to be set up with an object having functions for all handlers?
		//++ 	(like go's interface, but run-time checked if all handler methods are present)
		//++ this would mean it would be best to use the service provided here from another service that implements all the handlers
		// var handlers = {};
		// handlers.displayNotification(text) { alert(text); }
		// chatService.setHandlers(handlers)
		//
		//++ that would also mean it's not possible to have websocket running until handlers are set.
		//++ maybe set handlers on confiig phase? but that's not really possible... can't have a service in that phase.. so no handler..
		//++ this needs serious thinking... when can socket be started.. who handles messages? can't have a controller do it..
		//++ push notification from server needs to end up in client scope, most of the times.. so need event to get to right controller.??
		//++ are event registrations cleaned when controller is destroyed?
		//++ maybe provide simple handlerfunction that turns call into event?
		//++ maybe do start websocket right away, and cache all calls to non-registered handler, wait until handler is registered..
		//++ okay, might be bad idea, should probably implement handlers before websocket setup..
		//++ So can be done from controller (one that doesn't get destroyed). Or from a service (singleton).. Thats up to the user..
		//++ we can always add more helpers later-on.
		//++ service method: calcService.start(handlers) ?
		//++ handlers should return quite fast.. actually.. right? should they return answer as return? or using deferred?
		//++ automatically pickup returned deferred? (typeof deferred) and then wait? for resolve/reject? error on notify?
		//++ would be nice, makes implementing really smooth (choose wether to return answer directly, or as deferred)
		//++ answer should always be object with named fields. Whether it directly or using deferred.
		//++ when returning string, an error occurred. (in defered: reject("error message"))
		//++ when returning object, ok, and values are object. (in defered: resolve({field: "foo"}))

	};

	// PROVIDER VARIABLES AND SETTERS
	var wsPath = "/websocket-ango-{{.ServiceName}}";
	this.setWsPath = function(path) {
		wsPath = path;
	};

	// SERVICE CREATOR
	this.$get = function() {
		//++ open ws
		//++ setup queue
		return new service(wsPath);
	};
})
