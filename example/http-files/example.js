var chat = angular.module('chatExample', ['ango-chatservice']);

chat.config(function(chatserviceProvider) {
	chatserviceProvider.setDebug(true);

	chatserviceProvider.listenOnWsError(function(err) {
		console.error("ws error: " + err);
		console.log('setting reload timeout');
		setTimeout(function() {location.reload();}, 2000);
	});
	
	chatserviceProvider.listenOnWsClose(function() {
		console.log('ws closed, setting reload timeout');
		setTimeout(function() {location.reload();}, 2000);
	});
	
	chatserviceProvider.listenOnWrongVersion(function() {
		console.log('wrong version, setting reload timeout');
		setTimeout(function() {location.reload();}, 2000);
	});

	chatserviceProvider.setHandlers({
		displayNotification: function(subject, text) {
			alert(subject + ":\n\t" + text);
		},
		askQuestion: function(question) {
			return {
				answer: prompt(question),
			};
		},
	});
});

chat.controller('mainCtrl', function($scope, chatservice) {
	console.log(chatservice.getServiceName());
	console.log(chatservice.getProtocolVersion());

	chatservice.notify("hello ango server");

	chatservice.add8(100, 100).then(
		function(rets){
			console.log(rets.c);
		},
		function(err) {
			console.error(err);
		});

	$scope.calc = {
		a: 1,
		b: 2,
		add: function() {
			console.log('do add');
			chatservice.add(parseInt($scope.calc.a), parseInt($scope.calc.b)).then(
				function(retval) {
					$scope.calc.c = retval.c;
				}, function(err) {
					console.error(err);
				});
		}
	};

	$scope.name = {
		asked: false,
		answer: "my name",
	};

	$scope.foo = "controller is working";
});