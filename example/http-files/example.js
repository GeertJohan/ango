var chat = angular.module('chatExample', ['ango-chatService']);

chat.config(function(chatServiceProvider) {
	chatServiceProvider.setDebug(true);

	chatServiceProvider.listenOnWsError(function(err) {
		console.error("catched err! " + err);
	});
	
	chatServiceProvider.listenOnWsClose(function() {
		console.log('ws closed, setting reload timeouts');
		setTimeout(function() {location.reload();}, 2000);
	});
});

chat.controller('mainCtrl', function($scope, chatService) {
	console.log(chatService.getServiceName());
	console.log(chatService.getProtocolVersion());

	chatService.notify("hello ango server");

	chatService.add8(100, 100).then(
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
			chatService.add(parseInt($scope.calc.a), parseInt($scope.calc.b)).then(
				function(retval) {
					$scope.calc.c = retval.c;
				}, function(err) {
					console.error(err);
				});
		}
	};

	$scope.foo = "controller is working";
});