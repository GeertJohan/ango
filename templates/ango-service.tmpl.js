angular.module('ango-{{.Service.Name}}fdsa', [])
	.provider('{{.Service.Name}}', function() {

		// static constant globals for this generated service
		var serviceName = "{{.Service.Name}}";
		var protocolVersion = "{{.ProtocolVersion}}";

		// state enum
		var stateInit = 0;
		var stateRunning = 1;
		var stateStopped = 2;

		// errors
		var errStateStopped = "AngoError: state == stateStopped";
		var errVersionMismatch = "AngoError: version mismatch";

		// exceptions
		var expMissingArgs = "AngoException: missing arguments";
		var expTooManyArgs = "AngoException: too many arguments";

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

		// settings concerning the location and path for the websocket
		var wsUriScheme = "ws://";
		this.setWsUriScheme = function(scheme) {
			wsUriScheme = scheme;
		}
		var wsUriHost = document.location.host;
		this.setWsUriHost = function (host) {
			wsHost = host;
		};
		var wsUriPath = "/websocket-ango-{{.Service.Name}}";
		this.setWsUriPath = function(path) {
			wsPath = path;
		};

		// debugging settings
		var debug = false;
		this.setDebug = function(d) {
			debug = d;
		};

		// SERVICE CREATOR
		this.$get = ['$rootScope', '$q', function($rootScope, $q) {
			if(debug) {
				console.log("Starting ango service "+serviceName+" with version "+protocolVersion);
			}
			var service = {};

			// keep all pending requests here until they get responses
			var callbacks = {};
			// create a unique callback ID to map requests to responses
			var currentCallbackID = 0;
			// queue to hold sends when socket isn't open
			var queue = [];
			// create our websocket object with the address to the websocket
			var ws = new WebSocket(wsUriScheme+wsUriHost+wsPath);

			var state = stateInit;
			
			ws.onopen = function(){
				if(debug) {
					console.log("websocket has been opened!");
				}

				// send version string
				ws.send(protocolVersion)

				//++ hook open listeners
			};
			
			ws.onmessage = function(message) {
				switch(state) {
				case stateRunning:
					handleMessage(JSON.parse(message.data));
					break;
				case stateInit:
					switch(message.data) {
					case "good":
						if(debug) {
							console.log("Connection initialized (version matches)");
						}
						// set state
						state = stateRunning;
						// send queue
						sendQueue();
						break;
					case "invalid":
						if(debug) {
							console.error("Cannot setup communication over websocket: invalid version string. (version mismatch between server and client?)");
						}
						// set state
						state = stateStopped;
						// error on all deferreds
						errQueue(errVersionMismatch);
						errCallbacks(errVersionMismatch);

						//++ hook faulty version
						break;
					}
					break
				}
			};

			ws.onerror = function(err) {
				if(debug) {
					console.error("Error on websocket: ", err)
				}
				//++ hook error listeners
			}

			ws.onclose = function() {
				if(debug) {
					console.error("ango websocket closed");
				}
				//++ hook close listeners
			}

			// getCallbackID creates a new callback ID for a request
			function getCallbackID() {
				currentCallbackID += 1;
				return currentCallbackID;
			}

			// sendQueue send all requests from queue
			function sendQueue() {
				if(queue.length > 0) {
					if(debug) {
						console.log("Going to send "+queue.length+" items from queue.");
					}
					for(var i = 0; i < queue.length; i++) {
						// send request
						ws.send(queue[i].requestJson);

						if(queue[i].hasOwnProperty('oneway_deferred')) {
							// is aparently a oneway request
							// resolve the deferred
							queue[i].oneway_deferred.resolve({});
						}
					}
					queue = [];
				}
			}

			// errQueue rejects all oneway deferreds from the queue and empties the queue.
			// this is done when the connection could not be set up or broke.
			function errQueue(err) {
				for(var i = 0; i < queue.length; i++) {
					if(queue[i].hasOwnProperty('oneway_deferred')) {
						// is aparently a oneway request
						// resolve the deferred
						queue[i].oneway_deferred.reject(err);
					}
				}
			}

			// errCallbacks rejects all deferred callbacks
			// this is done when the connection could not be set up or broke.
			function errCallbacks(err) {
				for(var cb_id in callbacks) {
					if(callbacks.hasOwnProperty(cb_id)) {
						callbacks[cb_id].deferred.reject(err);
						delete callbacks[cb_id];
					}
				}
			}

			// doRequest makes a new request
			// it's either sent directly, or placed on queue (during startup)
			function doRequest(name, oneway, data) {
				if(state == stateStopped) {
					var deferred = $q.defer();
					deferred.reject(errStateStopped);
					return deferred.promise;
				}

				// setup request
				var request = {
					type: "req",
					procedure: name,
				}

				// create deferred to return
				var deferred = $q.defer();

				if(!oneway) {
					// only send data on not-oneway requests
					request.data = data;

					// setup callback to resolve deferred
					var callbackID = getCallbackID();
					callbacks[callbackID] = {
						time: new Date(),
						deferred: deferred,
					};
					request.cb_id = callbackID;
				}

				if(debug) {
					console.log('Sending request', request);
				}

				var requestJson = JSON.stringify(request);
				if(ws.readyState == 1 && state == stateRunning && queue.length == 0) {
					// directly send when ws is live and queue was completely sent
					ws.send(requestJson);

					if(oneway) {
						// resolve oneway requests immediatly after sending
						deferred.resolve({});
					}
				} else {
					// ws is not ready or queue is not completely sent yet.
					// therefore, add item to queue
					queueItem = {
						requestJson: requestJson,
					}

					// for oneway requests: add oneway_deferred propertie on queueItem
					// when queue is being sent, this deferred will be resolved
					if(oneway) {
						queueItem.oneway_deferred = deferred;
					}

					queue.push(queueItem);
				}

				return defer.promise;
			}

			function handleMessage(messageObj) {
				console.log("Received data from websocket: ", messageObj);

				switch(messageObj.type) {
				case "res":
					handleResolveMessage(messageObj);
					break;
				case "req":
					handleRequestMessage(messageObj);
					break;
				default:
					console.error("message with unknown type: ", messageObj);
					break;
				}
			}

			// handleResolveMessage resolves an outgoing request
			function handleResolveMessage(messageObj) {
				// if an object exists with cb_id in our callbacks object, resolve the deferred
				if(callbacks.hasOwnProperty(messageObj.cb_id)) {
					//++ TODO: is this $rootScope.$apply proper way to do it?
					if(typeof(messageObj.error) == "object") {
						$rootScope.$apply(callbacks[messageObj.cb_id].deferred.reject(messageObj.error));
					} else {
						$rootScope.$apply(callbacks[messageObj.cb_id].deferred.resolve(messageObj.data));
					}

					delete callbacks[messageObj.cb_id];
					
					return
				}
				console.error("TODO: implement resolve() some more")
			}

			// handleRequestMessage handles an incomming request
			function handleRequestMessage() {
				//++ TODO: implement request()
				console.error("TODO: implement request()")
			}

			// PROCEDURES, as defined in .ango file
			{{range .Service.ServerProcedures}}
			service.{{.Name}} = function( {{.Args.CommaSeperatedString}} ) {
				if(arguments.length > {{len .Args}}) {
					throw expTooManyArgs;
				}
				if(arguments.length < {{len .Args}}) {
					throw expMissingArgs;
				}
				//++ arguments type checking
				var data = {
					{{range .Args}} "{{.Name}}": {{.Name}}, {{end}}
				};
				var promise = doRequest("{{.Name}}", {{.Oneway}}, data); 
				return promise;
			};
			{{end}}

			return service;
		}];
	});
